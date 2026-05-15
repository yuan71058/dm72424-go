// Package dmsoft 提供大漠插件的 Go 语言绑定
// 支持两种运行模式：
//   - 32位模式: 直接通过syscall调用大漠DLL的COM接口
//   - 64位模式: 通过TCP+gob序列化与32位helper进程通信，实现跨进程调用
//
// 主要功能：
//   - 窗口操作：绑定、查找、枚举窗口
//   - 鼠标键盘模拟：点击、移动、按键等
//   - 图像识别：找图、找色、OCR文字识别
//   - 内存操作：读写内存、搜索数据
//   - 屏幕截图：捕获屏幕或窗口内容
//
// 使用流程：
//   1. LoadDm() 加载DLL
//   2. CrackDm() 破解（可选）
//   3. New() 创建实例
//   4. Init() 初始化COM对象
//   5. 调用各种API方法
//   6. Release() 释放资源
package dmsoft

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

// 全局变量定义
var (
	DmHModule     uintptr // 大漠DLL模块句柄（32位模式下使用）
	GoHModule     uintptr // 破解DLL模块句柄（32位模式下使用）
	comInitOnce   sync.Once // 确保COM初始化只执行一次
	is64bit       bool      // 当前是否为64位架构
	dmDllPath     string    // 大漠DLL文件路径
	crackDllPath  string    // 破解DLL文件路径
	helperExePath string    // 32位helper进程路径（dm_com_server.exe）
)

// uintptrSize 获取当前系统的指针大小
// 返回值: 4表示32位系统，8表示64位系统
// 用于判断当前运行架构，选择不同的调用方式
func uintptrSize() int {
	return int(unsafe.Sizeof(uintptr(0)))
}

// findHelperExe 查找32位helper进程可执行文件(dm_com_server.exe)
// 搜索顺序：
//   1. 当前可执行文件所在目录
//   2. ../cmd/dm_com_server/ 目录（开发环境）
//
// 返回值: 找到返回绝对路径，未找到返回空字符串
func findHelperExe() string {
	exePath, _ := os.Executable()
	dir := filepath.Dir(exePath)
	candidates := []string{
		filepath.Join(dir, "dm_com_server.exe"),
		filepath.Join(dir, "..", "cmd", "dm_com_server", "dm_com_server.exe"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}

// LoadDm 加载大漠插件DLL
// 参数:
//   - dmPath: 大漠DLL文件路径（如"dm.dll"、"xd47243.dll"等）
//
// 返回值:
//   - uintptr: 32位模式下返回DLL模块句柄，64位模式下返回1（表示成功）
//   - error: 错误信息
//
// 说明:
//   - 32位模式: 直接调用syscall.LoadLibrary加载DLL到当前进程
//   - 64位模式: 仅记录DLL路径并查找helper进程，实际加载在Init()时由helper完成
//   - 首次调用时会初始化COM库并检测当前架构（32/64位）
func LoadDm(dmPath string) (uintptr, error) {
	comInitOnce.Do(func() {
		ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
		is64bit = uintptrSize() == 8
	})

	dmDllPath = dmPath

	if is64bit {
		helperExePath = findHelperExe()
		if helperExePath == "" {
			return 0, fmt.Errorf("未找到dm_com_server.exe，请先编�?2位helper: GOARCH=386 go build -o dm_com_server.exe ./cmd/dm_com_server/")
		}
		return 1, nil
	}

	if DmHModule != 0 {
		return DmHModule, nil
	}

	module, err := syscall.LoadLibrary(dmPath)
	if err != nil {
		return 0, fmt.Errorf("加载大漠DLL失败: %v", err)
	}

	DmHModule = uintptr(module)
	return DmHModule, nil
}

// CrackDm 加载破解DLL并激活大漠插件
// 参数:
//   - crackPath: 破解DLL文件路径（如"Go.dll"等）
//
// 返回值: error 错误信息，成功时为nil
//
// 说明:
//   - 32位模式: 直接加载破解DLL并调用其导出函数Go()来激活大漠
//   - 64位模式: 仅记录破解DLL路径，实际破解在Init()时由helper进程完成
//   - 必须在LoadDm()之后、Init()之前调用
func CrackDm(crackPath string) error {
	if is64bit {
		crackDllPath = crackPath
		return nil
	}

	if DmHModule == 0 {
		return fmt.Errorf("大漠DLL未加�?请先调用 LoadDm 加载大漠插件")
	}

	if GoHModule != 0 {
		return nil
	}

	crackModule, err := syscall.LoadLibrary(crackPath)
	if err != nil {
		return fmt.Errorf("加载破解DLL失败: %v", err)
	}

	GoHModule = uintptr(crackModule)

	goFunAddr, err := syscall.GetProcAddress(syscall.Handle(GoHModule), "Go")
	if err != nil {
		syscall.FreeLibrary(syscall.Handle(GoHModule))
		GoHModule = 0
		return fmt.Errorf("获取破解函数地址失败: %v", err)
	}

	syscall.Syscall(uintptr(goFunAddr), 1, DmHModule, 0, 0)

	return nil
}

func FreeCrackDll() bool {
	if GoHModule == 0 {
		return true
	}

	err := syscall.FreeLibrary(syscall.Handle(GoHModule))
	if err != nil {
		return false
	}

	GoHModule = 0
	return true
}

func Load(path string) (uintptr, error) {
	return LoadDm(path)
}

func Free() bool {
	if is64bit {
		return true
	}

	if DmHModule == 0 {
		return true
	}

	err := syscall.FreeLibrary(syscall.Handle(DmHModule))
	if err != nil {
		return false
	}

	DmHModule = 0
	return true
}

// DmSoft 大漠插件实例结构体
// 封装了大漠COM对象的所有操作，支持32位和64位两种模式
type DmSoft struct {
	obj          uintptr       // COM对象指针（32位模式下使用）
	disp         *ole.IDispatch // COM IDispatch接口（预留，当前未使用）
	mu           sync.RWMutex  // 读写锁，保证线程安全
	dispCache    map[string]int32 // 调度ID缓存（预留）
	helperCmd    *exec.Cmd     // helper进程命令（64位模式下使用）
	helperStdin  io.WriteCloser // helper进程标准输入（用于发送退出信号）
	pipeConn     net.Conn      // TCP连接（64位模式下与helper通信）
	pipeEnc      *gob.Encoder  // gob编码器（用于发送请求到helper）
	pipeDec      *gob.Decoder  // gob解码器（用于接收helper的响应）
}

// New 创建大漠插件实例
// 返回值: *DmSoft 大漠实例指针，32位模式下如果DLL未加载则返回nil
//
// 说明:
//   - 每个DmSoft实例代表一个独立的大漠COM对象
//   - 64位模式下每个实例会启动独立的helper进程
//   - 多线程场景下，每个线程应创建独立的DmSoft实例
func New() *DmSoft {
	if is64bit {
		return &DmSoft{dispCache: make(map[string]int32)}
	}

	if DmHModule == 0 {
		return nil
	}

	return &DmSoft{dispCache: make(map[string]int32)}
}

// Init 初始化大漠对象，创建内部COM对象实例
// 返回值: error 错误信息
//
// 说明:
//   - 32位模式: 通过偏移量调用DLL中的创建函数，生成COM对象
//   - 64位模式: 启动32位helper进程(dm_com_server.exe)，建立TCP连接
//
// 64位模式的详细流程：
//   1. 启动helper进程，传入dm.dll和crack.dll路径作为参数
//   2. helper进程加载DLL、执行破解、创建dm对象
//   3. helper进程监听随机TCP端口并输出端口号到stdout
//   4. 主进程读取端口号，建立TCP连接
//   5. 创建gob编解码器用于后续通信
func (dm *DmSoft) Init() error {
	if is64bit {
		if helperExePath == "" {
			return fmt.Errorf("helper路径未初始化，请先调用LoadDm")
		}

		args := []string{dmDllPath}
		if crackDllPath != "" {
			args = append(args, crackDllPath)
		}

		cmd := exec.Command(helperExePath, args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("创建helper管道失败: %v", err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			stdin.Close()
			return fmt.Errorf("创建helper管道失败: %v", err)
		}
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("启动helper进程失败: %v", err)
		}

		scanner := bufio.NewScanner(stdout)
		if !scanner.Scan() {
			cmd.Process.Kill()
			cmd.Wait()
			return fmt.Errorf("helper进程无输出，请检查stderr")
		}
		line := scanner.Text()
		var port int
		if _, err := fmt.Sscanf(line, "READY %d", &port); err != nil || port == 0 {
			cmd.Process.Kill()
			cmd.Wait()
			return fmt.Errorf("helper进程输出异常: %s", line)
		}

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 5*time.Second)
		if err != nil {
			cmd.Process.Kill()
			cmd.Wait()
			return fmt.Errorf("连接helper TCP失败: %v", err)
		}

		dm.helperCmd = cmd
		dm.helperStdin = stdin
		dm.pipeConn = conn
		dm.pipeEnc = gob.NewEncoder(conn)
		dm.pipeDec = gob.NewDecoder(conn)
		return nil
	}

	if DmHModule == 0 {
		return fmt.Errorf("大漠DLL未加载，请先调用LoadDm")
	}
	createObjAddr := DmHModule + 98304
	dm.obj, _, _ = syscall.Syscall(createObjAddr, 0, 0, 0, 0)
	if dm.obj == 0 {
		return fmt.Errorf("创建大漠COM对象失败")
	}
	return nil
}

// Release 释放大漠对象和所有相关资源
//
// 说明:
//   - 64位模式: 关闭TCP连接、关闭helper进程的stdin（发送退出信号）、终止helper进程
//   - 32位模式: 调用DLL中的释放函数销毁COM对象
//   - 释放后该DmSoft实例不可再用，需要重新New()和Init()
func (dm *DmSoft) Release() {
	if dm.pipeConn != nil {
		dm.pipeConn.Close()
		dm.pipeConn = nil
		dm.pipeEnc = nil
		dm.pipeDec = nil
	}
	if dm.helperStdin != nil {
		dm.helperStdin.Close()
		dm.helperStdin = nil
	}
	if dm.helperCmd != nil {
		dm.helperCmd.Process.Kill()
		dm.helperCmd.Wait()
		dm.helperCmd = nil
	}

	if dm.obj == 0 {
		return
	}

	releaseObjAddr := DmHModule + 98400
	syscall.Syscall(releaseObjAddr, 1, dm.obj, 0, 0)
	dm.obj = 0
}

// GetDiskReversion 获取磁盘版本信息
// 参数: index - 索引(�?开�?
// 返回: 结果字符�?
func (dm *DmSoft) GetDiskReversion(index int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDiskReversion", index)
	}
	funAddr := DmHModule + 109040
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// LoadAiMemory 从内存加载AI模型
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LoadAiMemory(addr int32, size int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LoadAiMemory", addr, size)
	}
	funAddr := DmHModule + 108256
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(addr), uintptr(size))
	return int32(ret)
}

// FaqSend 发送FAQ请求到服务器
// 参数: server - 服务器地址
// 参数: handle - 句柄
// 参数: request_type - 请求类型
// 参数: time_out - 超时时间(毫秒)
// 返回: 结果字符�?
func (dm *DmSoft) FaqSend(server string, handle int32, request_type int32, time_out int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FaqSend", server, handle, request_type, time_out)
	}
	funAddr := DmHModule + 114016
	serverPtr := stringToBytePtr(server)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(serverPtr)), uintptr(handle), uintptr(request_type), uintptr(time_out), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindPicSimMem 在内存中查找图片(相似度模�?
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindPicSimMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindPicSimMem", []interface{}{x1, y1, x2, y2, pic_info, delta_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 121744
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall12(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0)
	return int32(ret)
}

// Ver 获取大漠插件版本�?// 返回: 结果字符�?
func (dm *DmSoft) Ver() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Ver")
	}
	funAddr := DmHModule + 100320
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetPath 设置资源文件路径
// 参数: path - 资源路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetPath(path string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetPath", path)
	}
	funAddr := DmHModule + 123808
	pathPtr := stringToBytePtr(path)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pathPtr)), 0)
	return int32(ret)
}

// SetShowAsmErrorMsg 设置是否显示汇编错误信息
// 参数: show - 显示标志(1:显示,0:隐藏)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetShowAsmErrorMsg(show int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetShowAsmErrorMsg", show)
	}
	funAddr := DmHModule + 101392
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(show), 0)
	return int32(ret)
}

// FindStrS 查找文字,返回找到的文字字符串
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStrWithOutVars("FindStrS", []interface{}{x1, y1, x2, y2, str, color, sim}, x, y)
	}
	funAddr := DmHModule + 116832
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetWordsNoDict 无字典获取区域内的所有文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 返回: 结果字符�?
func (dm *DmSoft) GetWordsNoDict(x1 int32, y1 int32, x2 int32, y2 int32, color string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWordsNoDict", x1, y1, x2, y2, color)
	}
	funAddr := DmHModule + 99024
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetOsBuildNumber 获取操作系统版本�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetOsBuildNumber() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetOsBuildNumber")
	}
	funAddr := DmHModule + 104240
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetID 获取大漠ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetID() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetID")
	}
	funAddr := DmHModule + 105184
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetMouseSpeed 设置鼠标移动速度
// 参数: speed - 移动速度(1-100)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMouseSpeed(speed int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMouseSpeed", speed)
	}
	funAddr := DmHModule + 124800
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(speed), 0)
	return int32(ret)
}

// FindData 查找数据
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: data - 数据
// 返回: 结果字符�?
func (dm *DmSoft) FindData(hwnd int32, addr_range string, data string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindData", hwnd, addr_range, data)
	}
	funAddr := DmHModule + 109760
	addr_rangePtr := stringToBytePtr(addr_range)
	dataPtr := stringToBytePtr(data)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(unsafe.Pointer(dataPtr)), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SendPaste 发送粘�?// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SendPaste(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SendPaste", hwnd)
	}
	funAddr := DmHModule + 122944
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// GetColor 获取指定坐标的颜�?// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) GetColor(x int32, y int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetColor", x, y)
	}
	funAddr := DmHModule + 117424
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(x), uintptr(y))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// LoadPicByte 从字节数据加载图片到内存
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 参数: name - 名称
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LoadPicByte(addr int32, size int32, name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LoadPicByte", addr, size, name)
	}
	funAddr := DmHModule + 121408
	namePtr := stringToBytePtr(name)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(addr), uintptr(size), uintptr(unsafe.Pointer(namePtr)), 0, 0)
	return int32(ret)
}

// WriteFloatAddr 写入浮点�?指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: float_value - 浮点数�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteFloatAddr(hwnd int32, addr int64, float_value float32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteFloatAddr", hwnd, addr, float64(float_value))
	}
	funAddr := DmHModule + 117312
	floatBits := uint32(math.Float32bits(float_value))
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(floatBits), 0)
	return int32(ret)
}

// SetWordLineHeight 设置文字识别行高
// 参数: line_height - 行高(像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWordLineHeight(line_height int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWordLineHeight", line_height)
	}
	funAddr := DmHModule + 101296
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(line_height), 0)
	return int32(ret)
}

// AsmCall 调用汇编代码
// 参数: hwnd - 窗口句柄
// 参数: mode - 模式
// 返回: 64位整数�?
func (dm *DmSoft) AsmCall(hwnd int32, mode int32) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("AsmCall", hwnd, mode)
	}
	funAddr := DmHModule + 114656
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(mode))
	return int64(ret)
}

// FindColorBlock 查找色块
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: count - 数量
// 参数: width - 宽度
// 参数: height - 高度
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindColorBlock(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindColorBlock", []interface{}{x1, y1, x2, y2, color, sim, count, width, height}, x, y)
	}
	funAddr := DmHModule + 113568
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 13, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(count), uintptr(width), uintptr(height), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0)
	return int32(ret)
}

// DisAssemble 反汇�?// 参数: asm_code - 汇编代码字符�?// 参数: base_addr - 模块基址
// 参数: is_64bit - 是否64�?// 返回: 结果字符�?
func (dm *DmSoft) DisAssemble(asm_code string, base_addr int64, is_64bit int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("DisAssemble", asm_code, base_addr, is_64bit)
	}
	funAddr := DmHModule + 112656
	asm_codePtr := stringToBytePtr(asm_code)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(asm_codePtr)), uintptr(base_addr), uintptr(base_addr>>32), uintptr(is_64bit), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// RegEx 扩展注册大漠插件
// 参数: code - 注册�?// 参数: ver - 版本�?// 参数: ip - IP地址
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RegEx(code string, ver string, ip string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RegEx", code, ver, ip)
	}
	funAddr := DmHModule + 98864
	codePtr := stringToBytePtr(code)
	verPtr := stringToBytePtr(ver)
	ipPtr := stringToBytePtr(ip)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(codePtr)), uintptr(unsafe.Pointer(verPtr)), uintptr(unsafe.Pointer(ipPtr)), 0, 0)
	return int32(ret)
}

// EncodeFile 加密文件
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EncodeFile(file string, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EncodeFile", file, pwd)
	}
	funAddr := DmHModule + 106528
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)))
	return int32(ret)
}

// WriteString 写入字符�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteString(hwnd int32, addr string, type_ int32, v string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteString", hwnd, addr, type_, v)
	}
	funAddr := DmHModule + 115936
	addrPtr := stringToBytePtr(addr)
	vPtr := stringToBytePtr(v)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(type_), uintptr(unsafe.Pointer(vPtr)), 0)
	return int32(ret)
}

// FindStrFastEx 高级快速查找文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrFastEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrFastEx", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 122000
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// AsmCallEx 扩展调用汇编代码
// 参数: hwnd - 窗口句柄
// 参数: mode - 模式
// 参数: base_addr - 模块基址
// 返回: 64位整数�?
func (dm *DmSoft) AsmCallEx(hwnd int32, mode int32, base_addr string) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("AsmCallEx", hwnd, mode, base_addr)
	}
	funAddr := DmHModule + 99632
	base_addrPtr := stringToBytePtr(base_addr)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(mode), uintptr(unsafe.Pointer(base_addrPtr)), 0, 0)
	return int64(ret)
}

// FindDoubleEx 高级查找双精度浮点数
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: double_value_min - 最小双精度�?// 参数: double_value_max - 最大双精度�?// 参数: step - 步长
// 参数: multi_thread - 多线程数�?// 参数: mode - 模式
// 返回: 结果字符�?
func (dm *DmSoft) FindDoubleEx(hwnd int32, addr_range string, double_value_min float64, double_value_max float64, step int32, multi_thread int32, mode int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindDoubleEx", hwnd, addr_range, double_value_min, double_value_max, step, multi_thread, mode)
	}
	funAddr := DmHModule + 110416
	addr_rangePtr := stringToBytePtr(addr_range)
	minBits := math.Float64bits(double_value_min)
	maxBits := math.Float64bits(double_value_max)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(minBits), uintptr(minBits>>32), uintptr(maxBits), uintptr(maxBits>>32), uintptr(step), uintptr(multi_thread), uintptr(mode), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetFindPicMultithreadLimit 设置找图多线程限�?// 参数: limit - 限制数量
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetFindPicMultithreadLimit(limit int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetFindPicMultithreadLimit", limit)
	}
	funAddr := DmHModule + 107616
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(limit), 0)
	return int32(ret)
}

// SendString2 发送字符串2
// 参数: hwnd - 窗口句柄
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SendString2(hwnd int32, str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SendString2", hwnd, str)
	}
	funAddr := DmHModule + 99888
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(strPtr)))
	return int32(ret)
}

// DownCpu 降低CPU使用�?// 参数: type_ - 类型
// 参数: rate - 速率
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DownCpu(type_ int32, rate int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DownCpu", type_, rate)
	}
	funAddr := DmHModule + 112960
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(type_), uintptr(rate))
	return int32(ret)
}

// DmGuard 大漠守护
// 参数: enable - 启用标志(1:启用,0:禁用)
// 参数: type_ - 类型
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DmGuard(enable int32, type_ string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DmGuard", enable, type_)
	}
	funAddr := DmHModule + 103552
	type_Ptr := stringToBytePtr(type_)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(enable), uintptr(unsafe.Pointer(type_Ptr)))
	return int32(ret)
}

// SpeedNormalGraphic 速度正常图形
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SpeedNormalGraphic(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SpeedNormalGraphic", en)
	}
	funAddr := DmHModule + 101184
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// FindPicSim 查找图片相似�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindPicSim(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindPicSim", []interface{}{x1, y1, x2, y2, pic_name, delta_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 98768
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall12(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0)
	return int32(ret)
}

// WriteInt 写入整数
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteInt(hwnd int32, addr string, type_ int32, v int64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteInt", hwnd, addr, type_, v)
	}
	funAddr := DmHModule + 112416
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(type_), uintptr(v), 0)
	return int32(ret)
}

// SetMemoryHwndAsProcessId 设置内存操作使用进程ID
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMemoryHwndAsProcessId(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMemoryHwndAsProcessId", en)
	}
	funAddr := DmHModule + 107984
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// WriteDataFromBin 从二进制写入数据
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: data - 数据
// 参数: len - 长度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteDataFromBin(hwnd int32, addr string, data int32, len int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteDataFromBin", hwnd, addr, data, len)
	}
	funAddr := DmHModule + 118304
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(data), uintptr(len), 0)
	return int32(ret)
}

// SetMinColGap 设置最小列间距
// 参数: col_gap - 列间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMinColGap(col_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMinColGap", col_gap)
	}
	funAddr := DmHModule + 110560
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(col_gap), 0)
	return int32(ret)
}

// KeyPressStr 按键字符串序�?// 参数: key_str - 按键字符�?// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyPressStr(key_str string, delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyPressStr", key_str, delay)
	}
	funAddr := DmHModule + 102528
	key_strPtr := stringToBytePtr(key_str)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(key_strPtr)), uintptr(delay))
	return int32(ret)
}

// LockDisplay 锁定显示区域
// 参数: lock - 锁定标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LockDisplay(lock int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LockDisplay", lock)
	}
	funAddr := DmHModule + 108304
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(lock), 0)
	return int32(ret)
}

// FindStrWithFontE 指定字体查找文字,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: font_name - 字体名称
// 参数: font_size - 字体大小(像素)
// 参数: flag - 查找标志
// 返回: 结果字符�?
func (dm *DmSoft) FindStrWithFontE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrWithFontE", x1, y1, x2, y2, str, color, sim, font_name, font_size, flag)
	}
	funAddr := DmHModule + 112544
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	font_namePtr := stringToBytePtr(font_name)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(font_namePtr)), uintptr(font_size), uintptr(flag), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnumIniKey 枚举INI�?// 参数: section - INI节名
// 参数: file - 文件路径
// 返回: 结果字符�?
func (dm *DmSoft) EnumIniKey(section string, file string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumIniKey", section, file)
	}
	funAddr := DmHModule + 108032
	sectionPtr := stringToBytePtr(section)
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(filePtr)))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// MatchPicName 匹配图片名称
// 参数: pic_name - 图片名称(多个用|分隔)
// 返回: 结果字符�?
func (dm *DmSoft) MatchPicName(pic_name string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("MatchPicName", pic_name)
	}
	funAddr := DmHModule + 117984
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pic_namePtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableFakeActive 启用假激�?// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableFakeActive(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableFakeActive", en)
	}
	funAddr := DmHModule + 107888
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// FaqGetSize 获取FAQ数据大小
// 参数: handle - 句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqGetSize(handle int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqGetSize", handle)
	}
	funAddr := DmHModule + 103456
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(handle), 0)
	return int32(ret)
}

// ExecuteCmd 执行命令
// 参数: cmd - 命令字符�?// 参数: current_dir - 当前目录
// 参数: time_out - 超时时间(毫秒)
// 返回: 结果字符�?
func (dm *DmSoft) ExecuteCmd(cmd string, current_dir string, time_out int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ExecuteCmd", cmd, current_dir, time_out)
	}
	funAddr := DmHModule + 116928
	cmdPtr := stringToBytePtr(cmd)
	current_dirPtr := stringToBytePtr(current_dir)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(cmdPtr)), uintptr(unsafe.Pointer(current_dirPtr)), uintptr(time_out), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableRealKeypad 启用真实键盘
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableRealKeypad(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableRealKeypad", en)
	}
	funAddr := DmHModule + 105648
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// SetDisplayRefreshDelay 设置显示刷新延迟
// 参数: t - 时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDisplayRefreshDelay(t int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDisplayRefreshDelay", t)
	}
	funAddr := DmHModule + 111344
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(t), 0)
	return int32(ret)
}

// MiddleClick 鼠标中键单击
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MiddleClick() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MiddleClick")
	}
	funAddr := DmHModule + 108560
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// AiYoloSortsObjects YOLO检测结果排�?// 参数: objects - 检测到的对�?// 参数: height - 高度
// 返回: 结果字符�?
func (dm *DmSoft) AiYoloSortsObjects(objects string, height int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AiYoloSortsObjects", objects, height)
	}
	funAddr := DmHModule + 120480
	objectsPtr := stringToBytePtr(objects)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(objectsPtr)), uintptr(height))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// WriteDataAddr 写入数据(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: data - 数据
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteDataAddr(hwnd int32, addr int64, data string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteDataAddr", hwnd, addr, data)
	}
	funAddr := DmHModule + 105744
	dataPtr := stringToBytePtr(data)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(unsafe.Pointer(dataPtr)), 0)
	return int32(ret)
}

// RGB2BGR RGB颜色转BGR
// 参数: rgb_color - rgb_color
// 返回: 结果字符�?
func (dm *DmSoft) RGB2BGR(rgb_color string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("RGB2BGR", rgb_color)
	}
	funAddr := DmHModule + 115744
	rgb_colorPtr := stringToBytePtr(rgb_color)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(rgb_colorPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// DisablePowerSave 禁用节能模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DisablePowerSave() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DisablePowerSave")
	}
	funAddr := DmHModule + 121952
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetClientSize 获取客户区大�?// 参数: hwnd - 窗口句柄
// 参数: width - 宽度(输出参数)
// 参数: height - 高度(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetClientSize(hwnd int32, width *int32, height *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetClientSize", []interface{}{hwnd}, width, height)
	}
	funAddr := DmHModule + 103344
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(width)), uintptr(unsafe.Pointer(height)), 0, 0)
	return int32(ret)
}

// EnableMouseMsg 启用鼠标消息模拟
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableMouseMsg(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableMouseMsg", en)
	}
	funAddr := DmHModule + 101344
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// EnableKeypadMsg 启用键盘消息
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableKeypadMsg(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableKeypadMsg", en)
	}
	funAddr := DmHModule + 120864
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// GetFileLength 获取文件长度
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetFileLength(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetFileLength", file)
	}
	funAddr := DmHModule + 111296
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// GetRemoteApiAddress 获取远程API地址
// 参数: hwnd - 窗口句柄
// 参数: base_addr - 模块基址
// 参数: fun_name - 函数名称
// 返回: 64位整数�?
func (dm *DmSoft) GetRemoteApiAddress(hwnd int32, base_addr int64, fun_name string) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("GetRemoteApiAddress", hwnd, base_addr, fun_name)
	}
	funAddr := DmHModule + 122192
	fun_namePtr := stringToBytePtr(fun_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(base_addr), uintptr(base_addr>>32), uintptr(unsafe.Pointer(fun_namePtr)), 0)
	return int64(ret)
}

// DmGuardParams 大漠守护参数
// 参数: cmd - 命令字符�?// 参数: sub_cmd - 子命�?// 参数: param - 参数
// 返回: 结果字符�?
func (dm *DmSoft) DmGuardParams(cmd string, sub_cmd string, param string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("DmGuardParams", cmd, sub_cmd, param)
	}
	funAddr := DmHModule + 105472
	cmdPtr := stringToBytePtr(cmd)
	sub_cmdPtr := stringToBytePtr(sub_cmd)
	paramPtr := stringToBytePtr(param)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(cmdPtr)), uintptr(unsafe.Pointer(sub_cmdPtr)), uintptr(unsafe.Pointer(paramPtr)), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// DownloadFile 下载文件
// 参数: url - URL地址
// 参数: save_file - 保存文件路径
// 参数: timeout - 超时时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DownloadFile(url string, save_file string, timeout int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DownloadFile", url, save_file, timeout)
	}
	funAddr := DmHModule + 123648
	urlPtr := stringToBytePtr(url)
	save_filePtr := stringToBytePtr(save_file)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(urlPtr)), uintptr(unsafe.Pointer(save_filePtr)), uintptr(timeout), 0, 0)
	return int32(ret)
}

// WriteDoubleAddr 写入双精度浮点数(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: double_value - 双精度浮点数�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteDoubleAddr(hwnd int32, addr int64, double_value float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteDoubleAddr", hwnd, addr, double_value)
	}
	funAddr := DmHModule + 115232
	doubleBits := math.Float64bits(double_value)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(doubleBits), uintptr(doubleBits>>32))
	return int32(ret)
}

// EnableIme 启用输入�?// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableIme(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableIme", en)
	}
	funAddr := DmHModule + 120192
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// TerminateProcessTree 终止进程�?// 参数: pid - 进程ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) TerminateProcessTree(pid int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("TerminateProcessTree", pid)
	}
	funAddr := DmHModule + 114240
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(pid), 0)
	return int32(ret)
}

// FoobarClose 关闭Foobar窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarClose(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarClose", hwnd)
	}
	funAddr := DmHModule + 102480
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// FindNearestPos 查找最近位�?// 参数: all_pos - 所有位置字符串
// 参数: type_ - 类型
// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) FindNearestPos(all_pos string, type_ int32, x int32, y int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindNearestPos", all_pos, type_, x, y)
	}
	funAddr := DmHModule + 112480
	all_posPtr := stringToBytePtr(all_pos)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(all_posPtr)), uintptr(type_), uintptr(x), uintptr(y), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// CreateFoobarRect 创建矩形Foobar窗口
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CreateFoobarRect(hwnd int32, x int32, y int32, w int32, h int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CreateFoobarRect", hwnd, x, y, w, h)
	}
	funAddr := DmHModule + 119072
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h))
	return int32(ret)
}

// GetCursorPos 获取鼠标当前位置
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetCursorPos(x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetCursorPos", []interface{}{}, x, y)
	}
	funAddr := DmHModule + 121680
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)))
	return int32(ret)
}

// FindColorBlockEx 高级查找色块
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: count - 数量
// 参数: width - 宽度
// 参数: height - 高度
// 返回: 结果字符�?
func (dm *DmSoft) FindColorBlockEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindColorBlockEx", x1, y1, x2, y2, color, sim, count, width, height)
	}
	funAddr := DmHModule + 103840
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(count), uintptr(width), uintptr(height), 0, 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindFloat 查找浮点�?// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: float_value_min - 最小浮点�?// 参数: float_value_max - 最大浮点�?// 返回: 结果字符�?
func (dm *DmSoft) FindFloat(hwnd int32, addr_range string, float_value_min float32, float_value_max float32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindFloat", hwnd, addr_range, float64(float_value_min), float64(float_value_max))
	}
	funAddr := DmHModule + 103216
	addr_rangePtr := stringToBytePtr(addr_range)
	minBits := uint32(math.Float32bits(float_value_min))
	maxBits := uint32(math.Float32bits(float_value_max))
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(minBits), uintptr(maxBits), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetProcessInfo 获取进程信息
// 参数: pid - 进程ID
// 返回: 结果字符�?
func (dm *DmSoft) GetProcessInfo(pid int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetProcessInfo", pid)
	}
	funAddr := DmHModule + 119024
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(pid), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadFile 读取文件内容
// 参数: file - 文件路径
// 返回: 结果字符�?
func (dm *DmSoft) ReadFile(file string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadFile", file)
	}
	funAddr := DmHModule + 114464
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindShapeEx 高级查找形状
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindShapeEx(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindShapeEx", x1, y1, x2, y2, offset_color, sim, dir)
	}
	funAddr := DmHModule + 99792
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetWindowText SetWindowText
// 参数: hwnd - 窗口句柄
// 参数: text - 文本内容
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWindowText(hwnd int32, text string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWindowText", hwnd, text)
	}
	funAddr := DmHModule + 113008
	textPtr := stringToBytePtr(text)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(textPtr)))
	return int32(ret)
}

// ForceUnBindWindow 强制解绑窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ForceUnBindWindow(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ForceUnBindWindow", hwnd)
	}
	funAddr := DmHModule + 120144
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// ReadIntAddr 读取整数(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 返回: 64位整数�?
func (dm *DmSoft) ReadIntAddr(hwnd int32, addr int64, type_ int32) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("ReadIntAddr", hwnd, addr, type_)
	}
	funAddr := DmHModule + 99712
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(type_), 0)
	return int64(ret)
}

// FindShape 查找形状
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindShape(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindShape", []interface{}{x1, y1, x2, y2, offset_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 123856
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// GetRealPath 获取真实文件路径
// 参数: path - 资源路径
// 返回: 结果字符�?
func (dm *DmSoft) GetRealPath(path string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetRealPath", path)
	}
	funAddr := DmHModule + 105008
	pathPtr := stringToBytePtr(path)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pathPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableSpeedDx 启用速度DX
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableSpeedDx(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableSpeedDx", en)
	}
	funAddr := DmHModule + 115472
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// UnLoadDriver 卸载驱动
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) UnLoadDriver() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("UnLoadDriver")
	}
	funAddr := DmHModule + 105696
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetMemoryUsage 获取内存使用情况
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetMemoryUsage() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetMemoryUsage")
	}
	funAddr := DmHModule + 106064
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// MiddleDown 鼠标中键按下
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MiddleDown() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MiddleDown")
	}
	funAddr := DmHModule + 109872
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// EnumIniSection 枚举INI�?// 参数: file - 文件路径
// 返回: 结果字符�?
func (dm *DmSoft) EnumIniSection(file string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumIniSection", file)
	}
	funAddr := DmHModule + 117184
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// CheckUAC 检查UAC状�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CheckUAC() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CheckUAC")
	}
	funAddr := DmHModule + 123104
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// OpenProcess 打开进程
// 参数: pid - 进程ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) OpenProcess(pid int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("OpenProcess", pid)
	}
	funAddr := DmHModule + 124624
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(pid), 0)
	return int32(ret)
}

// IsDisplayDead 检测屏幕是否死�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: t - 时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) IsDisplayDead(x1 int32, y1 int32, x2 int32, y2 int32, t int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("IsDisplayDead", x1, y1, x2, y2, t)
	}
	funAddr := DmHModule + 114896
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(t))
	return int32(ret)
}

// WriteIniPwd 写入INI配置(带密�?
// 参数: section - INI节名
// 参数: key - 键名
// 参数: v - �?// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteIniPwd(section string, key string, v string, file string, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteIniPwd", section, key, v, file, pwd)
	}
	funAddr := DmHModule + 115872
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	vPtr := stringToBytePtr(v)
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(vPtr)), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)))
	return int32(ret)
}

// GetNetTime 获取网络时间
// 返回: 结果字符�?
func (dm *DmSoft) GetNetTime() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetNetTime")
	}
	funAddr := DmHModule + 107712
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadFloat 读取浮点�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 返回: 32位浮点数
func (dm *DmSoft) ReadFloat(hwnd int32, addr string) float32 {
	if dm.pipeConn != nil {
		return float32(dm.comCallFloat64("ReadFloat", hwnd, addr))
	}
	funAddr := DmHModule + 100976
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)))
	return float32(ret)
}

// DisableCloseDisplayAndSleep 禁用关闭显示器和睡眠
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DisableCloseDisplayAndSleep() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DisableCloseDisplayAndSleep")
	}
	funAddr := DmHModule + 114416
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetWindowTitle 获取窗口标题
// 参数: hwnd - 窗口句柄
// 返回: 结果字符�?
func (dm *DmSoft) GetWindowTitle(hwnd int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWindowTitle", hwnd)
	}
	funAddr := DmHModule + 110816
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// Assemble Assemble
// 参数: base_addr - 模块基址
// 参数: is_64bit - 是否64�?// 返回: 结果字符�?
func (dm *DmSoft) Assemble(base_addr int64, is_64bit int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Assemble", base_addr, is_64bit)
	}
	funAddr := DmHModule + 119584
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(base_addr), uintptr(base_addr>>32), uintptr(is_64bit), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetMousePointWindow 获取鼠标指向的窗�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetMousePointWindow() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetMousePointWindow")
	}
	funAddr := DmHModule + 105424
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetExportDict 设置导出字库
// 参数: index - 索引(�?开�?
// 参数: dict_name - 字库文件�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetExportDict(index int32, dict_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetExportDict", index, dict_name)
	}
	funAddr := DmHModule + 119392
	dict_namePtr := stringToBytePtr(dict_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(index), uintptr(unsafe.Pointer(dict_namePtr)))
	return int32(ret)
}

// Delay 延迟指定时间
// 参数: mis - mis
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Delay(mis int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Delay", mis)
	}
	funAddr := DmHModule + 106480
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(mis), 0)
	return int32(ret)
}

// Reg 注册大漠插件
// 参数: code - 注册�?// 参数: ver - 版本�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Reg(code string, ver string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Reg", code, ver)
	}
	funAddr := DmHModule + 121344
	codePtr := stringToBytePtr(code)
	verPtr := stringToBytePtr(ver)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(codePtr)), uintptr(unsafe.Pointer(verPtr)))
	return int32(ret)
}

// FoobarStopGif Foobar停止播放GIF
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarStopGif(hwnd int32, x int32, y int32, pic_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarStopGif", hwnd, x, y, pic_name)
	}
	funAddr := DmHModule + 108096
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(pic_namePtr)), 0)
	return int32(ret)
}

// ReadFileData 读取文件数据
// 参数: file - 文件路径
// 参数: start_pos - 起始位置
// 参数: end_pos - 结束位置
// 返回: 结果字符�?
func (dm *DmSoft) ReadFileData(file string, start_pos int32, end_pos int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadFileData", file, start_pos, end_pos)
	}
	funAddr := DmHModule + 115808
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(filePtr)), uintptr(start_pos), uintptr(end_pos), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindPicSimEx 高级查找图片相似�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicSimEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicSimEx", x1, y1, x2, y2, pic_name, delta_color, sim, dir)
	}
	funAddr := DmHModule + 113728
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// Capture 截取屏幕区域
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Capture(x1 int32, y1 int32, x2 int32, y2 int32, file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Capture", x1, y1, x2, y2, file)
	}
	funAddr := DmHModule + 119456
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(filePtr)))
	return int32(ret)
}

// GetScreenWidth 获取屏幕宽度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetScreenWidth() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetScreenWidth")
	}
	funAddr := DmHModule + 113920
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindStrWithFontEx 高级指定字体查找文字
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: font_name - 字体名称
// 参数: font_size - 字体大小(像素)
// 参数: flag - 查找标志
// 返回: 结果字符�?
func (dm *DmSoft) FindStrWithFontEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrWithFontEx", x1, y1, x2, y2, str, color, sim, font_name, font_size, flag)
	}
	funAddr := DmHModule + 118848
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	font_namePtr := stringToBytePtr(font_name)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(font_namePtr)), uintptr(font_size), uintptr(flag), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetLocale 设置区域
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetLocale() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetLocale")
	}
	funAddr := DmHModule + 100928
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// AsmAdd 添加汇编指令
// 参数: asm_ins - 汇编指令
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AsmAdd(asm_ins string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AsmAdd", asm_ins)
	}
	funAddr := DmHModule + 121232
	asm_insPtr := stringToBytePtr(asm_ins)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(asm_insPtr)), 0)
	return int32(ret)
}

// GetScreenHeight 获取屏幕高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetScreenHeight() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetScreenHeight")
	}
	funAddr := DmHModule + 117792
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// CaptureGif 截取屏幕区域为GIF
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: file - 文件路径
// 参数: delay - 延迟时间(毫秒)
// 参数: time - 时间
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CaptureGif(x1 int32, y1 int32, x2 int32, y2 int32, file string, delay int32, time int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CaptureGif", x1, y1, x2, y2, file, delay, time)
	}
	funAddr := DmHModule + 120912
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(filePtr)), uintptr(delay), uintptr(time), 0)
	return int32(ret)
}

// ReadDataAddrToBin 读取数据到二进制(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: len - 长度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ReadDataAddrToBin(hwnd int32, addr int64, len int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ReadDataAddrToBin", hwnd, addr, len)
	}
	funAddr := DmHModule + 111792
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(len), 0)
	return int32(ret)
}

// ReadDataToBin 读取数据到二进制
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: len - 长度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ReadDataToBin(hwnd int32, addr string, len int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ReadDataToBin", hwnd, addr, len)
	}
	funAddr := DmHModule + 104480
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(len), 0, 0)
	return int32(ret)
}

// FindPicS 查找图片,返回图片索引
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStrWithOutVars("FindPicS", []interface{}{x1, y1, x2, y2, pic_name, delta_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 101952
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindPic 在指定区域查找图�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindPic", []interface{}{x1, y1, x2, y2, pic_name, delta_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 104032
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0)
	return int32(ret)
}

// FindMultiColor 多点找色
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: first_color - 第一个颜�?// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindMultiColor(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindMultiColor", []interface{}{x1, y1, x2, y2, first_color, offset_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 109360
	first_colorPtr := stringToBytePtr(first_color)
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(first_colorPtr)), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0)
	return int32(ret)
}

// HackSpeed 加�?// 参数: rate - 速率
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) HackSpeed(rate float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("HackSpeed", rate)
	}
	funAddr := DmHModule + 104352
	rateBits := math.Float64bits(rate)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(rateBits), uintptr(rateBits>>32))
	return int32(ret)
}

// FindPicE 查找图片,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicE", x1, y1, x2, y2, pic_name, delta_color, sim, dir)
	}
	funAddr := DmHModule + 114144
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// MiddleUp 鼠标中键弹起
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MiddleUp() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MiddleUp")
	}
	funAddr := DmHModule + 115072
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetWindow 获取窗口
// 参数: hwnd - 窗口句柄
// 参数: flag - 查找标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWindow(hwnd int32, flag int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetWindow", hwnd, flag)
	}
	funAddr := DmHModule + 120752
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(flag))
	return int32(ret)
}

// SetUAC 设置UAC状�?// 参数: uac - UAC标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetUAC(uac int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetUAC", uac)
	}
	funAddr := DmHModule + 108608
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(uac), 0)
	return int32(ret)
}

// FoobarSetSave 设置Foobar保存
// 参数: hwnd - 窗口句柄
// 参数: file - 文件路径
// 参数: en - 启用标志(1:启用,0:禁用)
// 参数: header - header
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarSetSave(hwnd int32, file string, en int32, header string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarSetSave", hwnd, file, en, header)
	}
	funAddr := DmHModule + 124736
	filePtr := stringToBytePtr(file)
	headerPtr := stringToBytePtr(header)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(filePtr)), uintptr(en), uintptr(unsafe.Pointer(headerPtr)), 0)
	return int32(ret)
}

// WheelDown 鼠标滚轮向下
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WheelDown() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WheelDown")
	}
	funAddr := DmHModule + 112848
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FloatToData 浮点数转数据
// 参数: float_value - 浮点数�?// 返回: 结果字符�?
func (dm *DmSoft) FloatToData(float_value float32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FloatToData", float64(float_value))
	}
	funAddr := DmHModule + 100464
	floatBits := uint32(math.Float32bits(float_value))
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(floatBits), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableFindPicMultithread 启用找图多线�?// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableFindPicMultithread(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableFindPicMultithread", en)
	}
	funAddr := DmHModule + 118048
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// DisableScreenSave 禁用屏保
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DisableScreenSave() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DisableScreenSave")
	}
	funAddr := DmHModule + 112800
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// AiFindPicEx AI高级查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) AiFindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AiFindPicEx", x1, y1, x2, y2, pic_name, sim, dir)
	}
	funAddr := DmHModule + 119136
	pic_namePtr := stringToBytePtr(pic_name)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SendString 发送字符串
// 参数: hwnd - 窗口句柄
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SendString(hwnd int32, str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SendString", hwnd, str)
	}
	funAddr := DmHModule + 114832
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(strPtr)))
	return int32(ret)
}

// EnterCri 进入临界�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnterCri() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnterCri")
	}
	funAddr := DmHModule + 116336
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindPicSimMemE 在内存中查找图片,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicSimMemE(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicSimMemE", x1, y1, x2, y2, pic_info, delta_color, sim, dir)
	}
	funAddr := DmHModule + 113296
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// Delays 随机延迟
// 参数: min_s - min_s
// 参数: max_s - max_s
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Delays(min_s int32, max_s int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Delays", min_s, max_s)
	}
	funAddr := DmHModule + 123328
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(min_s), uintptr(max_s))
	return int32(ret)
}

// CreateFoobarCustom 创建自定义Foobar窗口
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: pic - 图片数据
// 参数: trans_color - 透明�?// 参数: sim - 相似�?0.1-1.0)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CreateFoobarCustom(hwnd int32, x int32, y int32, pic string, trans_color string, sim float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CreateFoobarCustom", hwnd, x, y, pic, trans_color, sim)
	}
	funAddr := DmHModule + 105872
	picPtr := stringToBytePtr(pic)
	trans_colorPtr := stringToBytePtr(trans_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(picPtr)), uintptr(unsafe.Pointer(trans_colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return int32(ret)
}

// FindStringEx 高级查找字符�?// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: string_value - 字符串�?// 参数: type_ - 类型
// 参数: step - 步长
// 参数: multi_thread - 多线程数�?// 参数: mode - 模式
// 返回: 结果字符�?
func (dm *DmSoft) FindStringEx(hwnd int32, addr_range string, string_value string, type_ int32, step int32, multi_thread int32, mode int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStringEx", hwnd, addr_range, string_value, type_, step, multi_thread, mode)
	}
	funAddr := DmHModule + 124384
	addr_rangePtr := stringToBytePtr(addr_range)
	string_valuePtr := stringToBytePtr(string_value)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(unsafe.Pointer(string_valuePtr)), uintptr(type_), uintptr(step), uintptr(multi_thread), uintptr(mode), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetClientRect 获取客户区矩�?// 参数: hwnd - 窗口句柄
// 参数: x1 - 左上角X坐标(输出参数)
// 参数: y1 - 左上角Y坐标(输出参数)
// 参数: x2 - 右下角X坐标(输出参数)
// 参数: y2 - 右下角Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetClientRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetClientRect", []interface{}{hwnd}, x1, y1, x2, y2)
	}
	funAddr := DmHModule + 105808
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(x1)), uintptr(unsafe.Pointer(y1)), uintptr(unsafe.Pointer(x2)), uintptr(unsafe.Pointer(y2)))
	return int32(ret)
}

// AiYoloSetModel 从文件加载YOLO模型
// 参数: index - 索引(�?开�?
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloSetModel(index int32, file string, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloSetModel", index, file, pwd)
	}
	funAddr := DmHModule + 104416
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(index), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)), 0, 0)
	return int32(ret)
}

// FoobarSetTrans 设置Foobar透明�?// 参数: hwnd - 窗口句柄
// 参数: trans - trans
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarSetTrans(hwnd int32, trans int32, color string, sim float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarSetTrans", hwnd, trans, color, sim)
	}
	funAddr := DmHModule + 117248
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(trans), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return int32(ret)
}

// GetForegroundFocus 获取前台焦点窗口
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetForegroundFocus() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetForegroundFocus")
	}
	funAddr := DmHModule + 108512
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetForegroundWindow 获取前台窗口
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetForegroundWindow() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetForegroundWindow")
	}
	funAddr := DmHModule + 115360
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetExcludeRegion 设置排除区域
// 参数: type_ - 类型
// 参数: info - 信息字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetExcludeRegion(type_ int32, info string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetExcludeRegion", type_, info)
	}
	funAddr := DmHModule + 104832
	infoPtr := stringToBytePtr(info)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(type_), uintptr(unsafe.Pointer(infoPtr)))
	return int32(ret)
}

// SendStringIme2 通过输入法发送字符串2
// 参数: hwnd - 窗口句柄
// 参数: str - 要查找的字符�?// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SendStringIme2(hwnd int32, str string, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SendStringIme2", hwnd, str, mode)
	}
	funAddr := DmHModule + 119520
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(strPtr)), uintptr(mode), 0, 0)
	return int32(ret)
}

// ActiveInputMethod 激活输入法
// 参数: hwnd - 窗口句柄
// 参数: id - 标识ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ActiveInputMethod(hwnd int32, id string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ActiveInputMethod", hwnd, id)
	}
	funAddr := DmHModule + 124320
	idPtr := stringToBytePtr(id)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(idPtr)))
	return int32(ret)
}

// FoobarDrawPic Foobar绘制图片
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: pic - 图片数据
// 参数: trans_color - 透明�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarDrawPic(hwnd int32, x int32, y int32, pic string, trans_color string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarDrawPic", hwnd, x, y, pic, trans_color)
	}
	funAddr := DmHModule + 114288
	picPtr := stringToBytePtr(pic)
	trans_colorPtr := stringToBytePtr(trans_color)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(picPtr)), uintptr(unsafe.Pointer(trans_colorPtr)))
	return int32(ret)
}

// AiYoloSetVersion 设置YOLO模型版本
// 参数: ver - 版本�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloSetVersion(ver string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloSetVersion", ver)
	}
	funAddr := DmHModule + 118496
	verPtr := stringToBytePtr(ver)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(verPtr)), 0)
	return int32(ret)
}

// FindColorE 查找颜色,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindColorE(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindColorE", x1, y1, x2, y2, color, sim, dir)
	}
	funAddr := DmHModule + 120384
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// LeftClick 鼠标左键单击
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LeftClick() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LeftClick")
	}
	funAddr := DmHModule + 118096
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// IsFileExist 判断文件是否存在
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) IsFileExist(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("IsFileExist", file)
	}
	funAddr := DmHModule + 113824
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// Is64Bit 判断是否64位系�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Is64Bit() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Is64Bit")
	}
	funAddr := DmHModule + 110512
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindShapeE 查找形状,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindShapeE(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindShapeE", x1, y1, x2, y2, offset_color, sim, dir)
	}
	funAddr := DmHModule + 120592
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetDisplayInfo 获取显示器信�?// 返回: 结果字符�?
func (dm *DmSoft) GetDisplayInfo() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDisplayInfo")
	}
	funAddr := DmHModule + 122992
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetEnumWindowDelay 设置枚举窗口延迟
// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetEnumWindowDelay(delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetEnumWindowDelay", delay)
	}
	funAddr := DmHModule + 114720
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(delay), 0)
	return int32(ret)
}

// RegNoMac 注册大漠插件(不含MAC)
// 参数: code - 注册�?// 参数: ver - 版本�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RegNoMac(code string, ver string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RegNoMac", code, ver)
	}
	funAddr := DmHModule + 118960
	codePtr := stringToBytePtr(code)
	verPtr := stringToBytePtr(ver)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(codePtr)), uintptr(unsafe.Pointer(verPtr)))
	return int32(ret)
}

// KeyUpChar 弹起按键(字符形式)
// 参数: key_str - 按键字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyUpChar(key_str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyUpChar", key_str)
	}
	funAddr := DmHModule + 121904
	key_strPtr := stringToBytePtr(key_str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(key_strPtr)), 0)
	return int32(ret)
}

// SetDisplayAcceler 设置显示加�?// 参数: level - 级别
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDisplayAcceler(level int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDisplayAcceler", level)
	}
	funAddr := DmHModule + 101088
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(level), 0)
	return int32(ret)
}

// SetRowGapNoDict 设置行间�?无字典模�?
// 参数: row_gap - 行间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetRowGapNoDict(row_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetRowGapNoDict", row_gap)
	}
	funAddr := DmHModule + 118256
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(row_gap), 0)
	return int32(ret)
}

// EnableMouseAccuracy 启用鼠标精度
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableMouseAccuracy(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableMouseAccuracy", en)
	}
	funAddr := DmHModule + 123760
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// MoveTo 移动鼠标到指定坐�?// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MoveTo(x int32, y int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MoveTo", x, y)
	}
	funAddr := DmHModule + 109088
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(x), uintptr(y))
	return int32(ret)
}

// KeyPressChar 按键(字符形式)
// 参数: key_str - 按键字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyPressChar(key_str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyPressChar", key_str)
	}
	funAddr := DmHModule + 116464
	key_strPtr := stringToBytePtr(key_str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(key_strPtr)), 0)
	return int32(ret)
}

// RightDown 鼠标右键按下
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RightDown() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RightDown")
	}
	funAddr := DmHModule + 124576
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// AiYoloSetModelMemory 从内存加载YOLO模型
// 参数: index - 索引(�?开�?
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloSetModelMemory(index int32, addr int32, size int32, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloSetModelMemory", index, addr, size, pwd)
	}
	funAddr := DmHModule + 117600
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(index), uintptr(addr), uintptr(size), uintptr(unsafe.Pointer(pwdPtr)), 0)
	return int32(ret)
}

// WriteIni 写入INI配置
// 参数: section - INI节名
// 参数: key - 键名
// 参数: v - �?// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteIni(section string, key string, v string, file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteIni", section, key, v, file)
	}
	funAddr := DmHModule + 101232
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	vPtr := stringToBytePtr(v)
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(vPtr)), uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// DmGuardLoadCustom 大漠守护加载自定义模�?// 参数: type_ - 类型
// 参数: path - 资源路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DmGuardLoadCustom(type_ string, path string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DmGuardLoadCustom", type_, path)
	}
	funAddr := DmHModule + 106896
	type_Ptr := stringToBytePtr(type_)
	pathPtr := stringToBytePtr(path)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(type_Ptr)), uintptr(unsafe.Pointer(pathPtr)))
	return int32(ret)
}

// CreateFolder 创建文件�?// 参数: folder_name - folder_name
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CreateFolder(folder_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CreateFolder", folder_name)
	}
	funAddr := DmHModule + 113120
	folder_namePtr := stringToBytePtr(folder_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(folder_namePtr)), 0)
	return int32(ret)
}

// EnableRealMouse 启用真实鼠标模拟
// 参数: en - 启用标志(1:启用,0:禁用)
// 参数: mousedelay - 鼠标延迟(毫秒)
// 参数: mousestep - 鼠标步长
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableRealMouse(en int32, mousedelay int32, mousestep int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableRealMouse", en, mousedelay, mousestep)
	}
	funAddr := DmHModule + 105952
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(en), uintptr(mousedelay), uintptr(mousestep), 0, 0)
	return int32(ret)
}

// GetBasePath 获取大漠基础路径
// 返回: 结果字符�?
func (dm *DmSoft) GetBasePath() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetBasePath")
	}
	funAddr := DmHModule + 107312
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetFps 获取FPS
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetFps() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetFps")
	}
	funAddr := DmHModule + 106016
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// EnableGetColorByCapture 启用截图取色模式
// 参数: enable - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableGetColorByCapture(enable int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableGetColorByCapture", enable)
	}
	funAddr := DmHModule + 109216
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(enable), 0)
	return int32(ret)
}

// SetDisplayInput 设置显示输入方式
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDisplayInput(mode string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDisplayInput", mode)
	}
	funAddr := DmHModule + 110944
	modePtr := stringToBytePtr(mode)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(modePtr)), 0)
	return int32(ret)
}

// Hex64 64位整数转十六进制字符�?// 参数: v - �?// 返回: 结果字符�?
func (dm *DmSoft) Hex64(v int64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Hex64", v)
	}
	funAddr := DmHModule + 105296
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(v), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ScreenToClient 屏幕坐标转客户区坐标
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ScreenToClient(hwnd int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("ScreenToClient", []interface{}{hwnd}, x, y)
	}
	funAddr := DmHModule + 111392
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0)
	return int32(ret)
}

// AiEnableFindPicWindow 启用AI找图窗口模式
// 参数: enable - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiEnableFindPicWindow(enable int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiEnableFindPicWindow", enable)
	}
	funAddr := DmHModule + 100064
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(enable), 0)
	return int32(ret)
}

// ReadIni 读取INI配置
// 参数: section - INI节名
// 参数: key - 键名
// 参数: file - 文件路径
// 返回: 结果字符�?
func (dm *DmSoft) ReadIni(section string, key string, file string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadIni", section, key, file)
	}
	funAddr := DmHModule + 102912
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(filePtr)), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ImageToBmp 图片转换为BMP格式
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: bmp_name - bmp_name
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ImageToBmp(pic_name string, bmp_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ImageToBmp", pic_name, bmp_name)
	}
	funAddr := DmHModule + 109152
	pic_namePtr := stringToBytePtr(pic_name)
	bmp_namePtr := stringToBytePtr(bmp_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(bmp_namePtr)))
	return int32(ret)
}

// SetDisplayDelay 设置显示延迟
// 参数: t - 时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDisplayDelay(t int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDisplayDelay", t)
	}
	funAddr := DmHModule + 122784
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(t), 0)
	return int32(ret)
}

// WheelUp 鼠标滚轮向上
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WheelUp() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WheelUp")
	}
	funAddr := DmHModule + 102688
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// CopyFile 复制文件
// 参数: src_file - 源文件路�?// 参数: dst_file - 目标文件路径
// 参数: over - 覆盖标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CopyFile(src_file string, dst_file string, over int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CopyFile", src_file, dst_file, over)
	}

	funAddr := DmHModule + 100688
	src_filePtr := stringToBytePtr(src_file)
	dst_filePtr := stringToBytePtr(dst_file)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(src_filePtr)), uintptr(unsafe.Pointer(dst_filePtr)), uintptr(over), 0, 0)
	return int32(ret)
}

// FindWindowEx 扩展查找窗口
// 参数: parent - 父窗口句�?// 参数: class_name - 窗口类名
// 参数: title_name - 窗口标题
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindWindowEx(parent int32, class_name string, title_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindWindowEx", parent, class_name, title_name)
	}
	funAddr := DmHModule + 115408
	class_namePtr := stringToBytePtr(class_name)
	title_namePtr := stringToBytePtr(title_name)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(parent), uintptr(unsafe.Pointer(class_namePtr)), uintptr(unsafe.Pointer(title_namePtr)), 0, 0)
	return int32(ret)
}

// SetFindPicMultithreadCount 设置找图多线程数�?// 参数: count - 数量
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetFindPicMultithreadCount(count int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetFindPicMultithreadCount", count)
	}
	funAddr := DmHModule + 106784
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(count), 0)
	return int32(ret)
}

// GetScreenDataBmp 获取屏幕BMP数据
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: data - 数据(输出参数)
// 参数: size - 大小(字节)(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetScreenDataBmp(x1 int32, y1 int32, x2 int32, y2 int32, data *int32, size *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetScreenDataBmp", []interface{}{x1, y1, x2, y2}, data, size)
	}
	funAddr := DmHModule + 107136
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(size)), 0, 0)
	return int32(ret)
}

// GetWordResultPos 获取文字识别结果位置
// 参数: str - 要查找的字符�?// 参数: index - 索引(�?开�?
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWordResultPos(str string, index int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetWordResultPos", []interface{}{str, index}, x, y)
	}
	funAddr := DmHModule + 114352
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(strPtr)), uintptr(index), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0)
	return int32(ret)
}

// LeftDoubleClick 鼠标左键双击
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LeftDoubleClick() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LeftDoubleClick")
	}
	funAddr := DmHModule + 101136
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// ReadStringAddr 读取字符�?指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: len - 长度
// 返回: 结果字符�?
func (dm *DmSoft) ReadStringAddr(hwnd int32, addr int64, type_ int32, len int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadStringAddr", hwnd, addr, type_, len)
	}

	funAddr := DmHModule + 118608
	ret, _, _ := syscall.Syscall9(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(type_), uintptr(len), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadData 读取数据
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: len - 长度
// 返回: 结果字符�?
func (dm *DmSoft) ReadData(hwnd int32, addr string, len int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadData", hwnd, addr, len)
	}
	funAddr := DmHModule + 111232
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(len), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// AddDict 添加字库条目
// 参数: index - 索引(�?开�?
// 参数: dict_info - 字库信息
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AddDict(index int32, dict_info string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AddDict", index, dict_info)
	}
	funAddr := DmHModule + 106336
	dict_infoPtr := stringToBytePtr(dict_info)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(index), uintptr(unsafe.Pointer(dict_infoPtr)))
	return int32(ret)
}

// SetInputDm 设置输入大漠
// 参数: input_dm - input_dm
// 参数: rx - 相对X偏移
// 参数: ry - 相对Y偏移
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetInputDm(input_dm int32, rx int32, ry int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetInputDm", input_dm, rx, ry)
	}
	funAddr := DmHModule + 108656
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(input_dm), uintptr(rx), uintptr(ry), 0, 0)
	return int32(ret)
}

// GetWindowProcessId 获取窗口进程ID
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWindowProcessId(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetWindowProcessId", hwnd)
	}
	funAddr := DmHModule + 124464
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// WriteDataAddrFromBin 从二进制写入数据(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: data - 数据
// 参数: len - 长度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteDataAddrFromBin(hwnd int32, addr int64, data int32, len int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteDataAddrFromBin", hwnd, addr, data, len)
	}
	funAddr := DmHModule + 121120
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(data), uintptr(len))
	return int32(ret)
}

// AiFindPicMemEx AI高级内存查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) AiFindPicMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AiFindPicMemEx", x1, y1, x2, y2, pic_info, sim, dir)
	}
	funAddr := DmHModule + 102976
	pic_infoPtr := stringToBytePtr(pic_info)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// TerminateProcess 终止进程
// 参数: pid - 进程ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) TerminateProcess(pid int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("TerminateProcess", pid)
	}
	funAddr := DmHModule + 112032
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(pid), 0)
	return int32(ret)
}

// VirtualQueryEx 查询目标进程内存信息
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: pmbi - 内存信息缓冲区指�?// 返回: 结果字符�?
func (dm *DmSoft) VirtualQueryEx(hwnd int32, addr int64, pmbi int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("VirtualQueryEx", hwnd, addr, pmbi)
	}
	funAddr := DmHModule + 101632
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(pmbi), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableKeypadSync 启用键盘同步
// 参数: enable - 启用标志(1:启用,0:禁用)
// 参数: time_out - 超时时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableKeypadSync(enable int32, time_out int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableKeypadSync", enable, time_out)
	}
	funAddr := DmHModule + 109968
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(enable), uintptr(time_out))
	return int32(ret)
}

// AiYoloUseModel 切换使用已加载的YOLO模型
// 参数: index - 索引(�?开�?
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloUseModel(index int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloUseModel", index)
	}
	funAddr := DmHModule + 110032
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return int32(ret)
}

// DeleteFile 删除文件
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DeleteFile(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DeleteFile", file)
	}
	funAddr := DmHModule + 99408
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// GetScreenDepth 获取屏幕色深
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetScreenDepth() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetScreenDepth")
	}
	funAddr := DmHModule + 102384
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindColor 在区域查找颜�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindColor", []interface{}{x1, y1, x2, y2, color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 106112
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// MoveR 相对移动鼠标
// 参数: rx - 相对X偏移
// 参数: ry - 相对Y偏移
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MoveR(rx int32, ry int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MoveR", rx, ry)
	}
	funAddr := DmHModule + 113504
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(rx), uintptr(ry))
	return int32(ret)
}

// LockInput 锁定输入
// 参数: lock - 锁定标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LockInput(lock int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LockInput", lock)
	}
	funAddr := DmHModule + 124272
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(lock), 0)
	return int32(ret)
}

// IntToData 整数转数�?// 参数: int_value - 整数�?// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) IntToData(int_value int64, type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("IntToData", int_value, type_)
	}
	funAddr := DmHModule + 122272
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(int_value), uintptr(type_))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FaqPost 异步发送FAQ请求
// 参数: server - 服务器地址
// 参数: handle - 句柄
// 参数: request_type - 请求类型
// 参数: time_out - 超时时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqPost(server string, handle int32, request_type int32, time_out int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqPost", server, handle, request_type, time_out)
	}
	funAddr := DmHModule + 107440
	serverPtr := stringToBytePtr(server)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(serverPtr)), uintptr(handle), uintptr(request_type), uintptr(time_out), 0)
	return int32(ret)
}

// GetColorHSV 获取指定坐标的HSV颜色
// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) GetColorHSV(x int32, y int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetColorHSV", x, y)
	}
	funAddr := DmHModule + 116192
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(x), uintptr(y))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindWindowSuper 超级查找窗口
// 参数: spec1 - 条件1字符�?// 参数: flag1 - 条件1标志
// 参数: type1 - 条件1类型
// 参数: spec2 - 条件2字符�?// 参数: flag2 - 条件2标志
// 参数: type2 - 条件2类型
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindWindowSuper", spec1, flag1, type1, spec2, flag2, type2)
	}
	funAddr := DmHModule + 108432
	spec1Ptr := stringToBytePtr(spec1)
	spec2Ptr := stringToBytePtr(spec2)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(unsafe.Pointer(spec1Ptr)), uintptr(flag1), uintptr(type1), uintptr(unsafe.Pointer(spec2Ptr)), uintptr(flag2), uintptr(type2), 0, 0)
	return int32(ret)
}

// EnableBind 启用/禁用绑定
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableBind(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableBind", en)
	}
	funAddr := DmHModule + 116576
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// SetAero 设置Aero效果
// 参数: enable - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetAero(enable int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetAero", enable)
	}
	funAddr := DmHModule + 102640
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(enable), 0)
	return int32(ret)
}

// DecodeFile 解密文件
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DecodeFile(file string, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DecodeFile", file, pwd)
	}
	funAddr := DmHModule + 122496
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)))
	return int32(ret)
}

// FindPicExS 高级查找图片,返回详细字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicExS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicExS", x1, y1, x2, y2, pic_name, delta_color, sim, dir)
	}
	funAddr := DmHModule + 100368
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// WriteStringAddr 写入字符�?指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteStringAddr(hwnd int32, addr int64, type_ int32, v string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteStringAddr", hwnd, addr, type_, v)
	}
	funAddr := DmHModule + 122720
	vPtr := stringToBytePtr(v)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(type_), uintptr(unsafe.Pointer(vPtr)), 0, 0, 0)
	return int32(ret)
}

// GetCommandLine 获取命令�?// 参数: hwnd - 窗口句柄
// 返回: 结果字符�?
func (dm *DmSoft) GetCommandLine(hwnd int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetCommandLine", hwnd)
	}
	funAddr := DmHModule + 100752
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SelectFile 选择文件对话�?// 返回: 结果字符�?
func (dm *DmSoft) SelectFile() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("SelectFile")
	}
	funAddr := DmHModule + 118144
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindPicSimMemEx 高级在内存中查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicSimMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicSimMemEx", x1, y1, x2, y2, pic_info, delta_color, sim, dir)
	}
	funAddr := DmHModule + 124912
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetWordResultStr 获取文字识别结果字符�?// 参数: str - 要查找的字符�?// 参数: index - 索引(�?开�?
// 返回: 结果字符�?
func (dm *DmSoft) GetWordResultStr(str string, index int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWordResultStr", str, index)
	}
	funAddr := DmHModule + 104768
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(strPtr)), uintptr(index))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnablePicCache 启用图片缓存
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnablePicCache(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnablePicCache", en)
	}
	funAddr := DmHModule + 99536
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// FindStrExS 高级查找文字,返回详细字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrExS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrExS", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 100528
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// LoadPic 预加载图片到内存
// 参数: pic_name - 图片名称(多个用|分隔)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LoadPic(pic_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LoadPic", pic_name)
	}
	funAddr := DmHModule + 124128
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pic_namePtr)), 0)
	return int32(ret)
}

// FindStrFast 快速查找文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindStrFast(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindStrFast", []interface{}{x1, y1, x2, y2, str, color, sim}, x, y)
	}
	funAddr := DmHModule + 115584
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// FindDouble 查找双精度浮点数
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: double_value_min - 最小双精度�?// 参数: double_value_max - 最大双精度�?// 返回: 结果字符�?
func (dm *DmSoft) FindDouble(hwnd int32, addr_range string, double_value_min float64, double_value_max float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindDouble", hwnd, addr_range, double_value_min, double_value_max)
	}
	funAddr := DmHModule + 102192
	addr_rangePtr := stringToBytePtr(addr_range)
	minBits := math.Float64bits(double_value_min)
	maxBits := math.Float64bits(double_value_max)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(minBits), uintptr(minBits>>32), uintptr(maxBits), uintptr(maxBits>>32), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetParam64ToPointer 设置64位参数转指针
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetParam64ToPointer() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetParam64ToPointer")
	}
	funAddr := DmHModule + 99952
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetMemoryFindResultToFile 设置内存查找结果输出到文�?// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMemoryFindResultToFile(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMemoryFindResultToFile", file)
	}
	funAddr := DmHModule + 110704
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// WaitKey 等待按键
// 参数: key_code - 键码
// 参数: time_out - 超时时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WaitKey(key_code int32, time_out int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WaitKey", key_code, time_out)
	}
	funAddr := DmHModule + 114528
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(key_code), uintptr(time_out))
	return int32(ret)
}

// CreateFoobarEllipse 创建椭圆Foobar窗口
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CreateFoobarEllipse(hwnd int32, x int32, y int32, w int32, h int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CreateFoobarEllipse", hwnd, x, y, w, h)
	}
	funAddr := DmHModule + 114592
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h))
	return int32(ret)
}

// MoveFile 移动文件
// 参数: src_file - 源文件路�?// 参数: dst_file - 目标文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MoveFile(src_file string, dst_file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MoveFile", src_file, dst_file)
	}
	funAddr := DmHModule + 102272
	src_filePtr := stringToBytePtr(src_file)
	dst_filePtr := stringToBytePtr(dst_file)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(src_filePtr)), uintptr(unsafe.Pointer(dst_filePtr)))
	return int32(ret)
}

// Stop 停止
// 参数: id - 标识ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Stop(id int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Stop", id)
	}
	funAddr := DmHModule + 100880
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(id), 0)
	return int32(ret)
}

// ReleaseRef 释放引用
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ReleaseRef() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ReleaseRef")
	}
	funAddr := DmHModule + 111072
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetColorBGR 获取指定坐标的BGR颜色
// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) GetColorBGR(x int32, y int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetColorBGR", x, y)
	}
	funAddr := DmHModule + 100000
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(x), uintptr(y))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnumIniKeyPwd 枚举INI�?带密�?
// 参数: section - INI节名
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 结果字符�?
func (dm *DmSoft) EnumIniKeyPwd(section string, file string, pwd string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumIniKeyPwd", section, file, pwd)
	}
	funAddr := DmHModule + 116768
	sectionPtr := stringToBytePtr(section)
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetMac 获取本机MAC地址
// 返回: 结果字符�?
func (dm *DmSoft) GetMac() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetMac")
	}
	funAddr := DmHModule + 123536
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// UseDict 切换使用指定索引的字�?// 参数: index - 索引(�?开�?
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) UseDict(index int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("UseDict", index)
	}
	funAddr := DmHModule + 104656
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return int32(ret)
}

// FindDataEx 高级查找数据
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: data - 数据
// 参数: step - 步长
// 参数: multi_thread - 多线程数�?// 参数: mode - 模式
// 返回: 结果字符�?
func (dm *DmSoft) FindDataEx(hwnd int32, addr_range string, data string, step int32, multi_thread int32, mode int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindDataEx", hwnd, addr_range, data, step, multi_thread, mode)
	}
	funAddr := DmHModule + 123200
	addr_rangePtr := stringToBytePtr(addr_range)
	dataPtr := stringToBytePtr(data)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(unsafe.Pointer(dataPtr)), uintptr(step), uintptr(multi_thread), uintptr(mode), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// Md5 计算MD5�?// 参数: str - 要查找的字符�?// 返回: 结果字符�?
func (dm *DmSoft) Md5(str string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Md5", str)
	}
	funAddr := DmHModule + 117376
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(strPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// BGR2RGB BGR颜色转RGB
// 参数: bgr_color - bgr_color
// 返回: 结果字符�?
func (dm *DmSoft) BGR2RGB(bgr_color string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("BGR2RGB", bgr_color)
	}
	funAddr := DmHModule + 118736
	bgr_colorPtr := stringToBytePtr(bgr_color)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(bgr_colorPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindColorEx 高级查找颜色
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindColorEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindColorEx", x1, y1, x2, y2, color, sim, dir)
	}
	funAddr := DmHModule + 103600
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// OcrExOne OCR识别单个区域
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) OcrExOne(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("OcrExOne", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 112080
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// CmpColor 比较颜色
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CmpColor(x int32, y int32, color string, sim float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CmpColor", x, y, color, sim)
	}
	funAddr := DmHModule + 109648
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(x), uintptr(y), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32))
	return int32(ret)
}

// OcrInFile 从文件中识别文字

// OcrInFile 从图片文件进行OCR识别
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) OcrInFile(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("OcrInFile", x1, y1, x2, y2, pic_name, color, sim)
	}
	funAddr := DmHModule + 110608
	pic_namePtr := stringToBytePtr(pic_name)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// CheckInputMethod 检查输入法状�?// 参数: hwnd - 窗口句柄
// 参数: id - 标识ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CheckInputMethod(hwnd int32, id string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CheckInputMethod", hwnd, id)
	}
	funAddr := DmHModule + 101792
	idPtr := stringToBytePtr(id)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(idPtr)))
	return int32(ret)
}

// MoveWindow 移动窗口
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MoveWindow(hwnd int32, x int32, y int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MoveWindow", hwnd, x, y)
	}
	funAddr := DmHModule + 119648
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), 0, 0)
	return int32(ret)
}

// GetClipboard 获取剪贴板内�?// 返回: 结果字符�?
func (dm *DmSoft) GetClipboard() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetClipboard")
	}
	funAddr := DmHModule + 116624
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindStr 在指定区域查找文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindStr(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindStr", []interface{}{x1, y1, x2, y2, str, color, sim}, x, y)
	}
	funAddr := DmHModule + 110320
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// FoobarClearText Foobar清除文字
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarClearText(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarClearText", hwnd)
	}
	funAddr := DmHModule + 113072
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// ClientToScreen 客户区坐标转屏幕坐标
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ClientToScreen(hwnd int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("ClientToScreen", []interface{}{hwnd}, x, y)
	}
	funAddr := DmHModule + 116512
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0)
	return int32(ret)
}

// GetCursorShape 获取鼠标形状
// 返回: 结果字符�?
func (dm *DmSoft) GetCursorShape() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetCursorShape")
	}
	funAddr := DmHModule + 111984
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetWordResultCount 获取文字识别结果数量
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWordResultCount(str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetWordResultCount", str)
	}
	funAddr := DmHModule + 103984
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(strPtr)), 0)
	return int32(ret)
}

// SelectDirectory 选择目录对话�?// 返回: 结果字符�?
func (dm *DmSoft) SelectDirectory() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("SelectDirectory")
	}
	funAddr := DmHModule + 116000
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// CapturePng 截取屏幕区域为PNG
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CapturePng(x1 int32, y1 int32, x2 int32, y2 int32, file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CapturePng", x1, y1, x2, y2, file)
	}
	funAddr := DmHModule + 114080
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(filePtr)))
	return int32(ret)
}

// KeyDownChar 按下按键(字符形式)
// 参数: key_str - 按键字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyDownChar(key_str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyDownChar", key_str)
	}
	funAddr := DmHModule + 105600
	key_strPtr := stringToBytePtr(key_str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(key_strPtr)), 0)
	return int32(ret)
}

// CaptureJpg 截取屏幕区域为JPG
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: file - 文件路径
// 参数: quality - 图片质量(1-100)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CaptureJpg(x1 int32, y1 int32, x2 int32, y2 int32, file string, quality int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CaptureJpg", x1, y1, x2, y2, file, quality)
	}
	funAddr := DmHModule + 106400
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(filePtr)), uintptr(quality), 0, 0)
	return int32(ret)
}

// FindStrEx 高级查找文字,返回所有匹配位�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrEx", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 106640
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FaqCapture FAQ截图并保存到缓存
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: quality - 图片质量(1-100)
// 参数: delay - 延迟时间(毫秒)
// 参数: time - 时间
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqCapture(x1 int32, y1 int32, x2 int32, y2 int32, quality int32, delay int32, time int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqCapture", x1, y1, x2, y2, quality, delay, time)
	}
	funAddr := DmHModule + 118416
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(quality), uintptr(delay), uintptr(time), 0)
	return int32(ret)
}

// ShowScrMsg 显示屏幕消息
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: msg - 消息内容
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ShowScrMsg(x1 int32, y1 int32, x2 int32, y2 int32, msg string, color string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ShowScrMsg", x1, y1, x2, y2, msg, color)
	}
	funAddr := DmHModule + 112208
	msgPtr := stringToBytePtr(msg)
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(colorPtr)), 0, 0)
	return int32(ret)
}

// SetKeypadDelay 设置键盘按键延迟
// 参数: type_ - 类型
// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetKeypadDelay(type_ string, delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetKeypadDelay", type_, delay)
	}
	funAddr := DmHModule + 110256
	type_Ptr := stringToBytePtr(type_)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(type_Ptr)), uintptr(delay))
	return int32(ret)
}

// SetScreen 设置屏幕参数
// 参数: width - 宽度
// 参数: height - 高度
// 参数: depth - depth
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetScreen(width int32, height int32, depth int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetScreen", width, height, depth)
	}
	funAddr := DmHModule + 115168
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(width), uintptr(height), uintptr(depth), 0, 0)
	return int32(ret)
}

// Play 播放声音文件
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Play(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Play", file)
	}
	funAddr := DmHModule + 105072
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// FindWindowByProcessId 通过进程ID查找窗口
// 参数: process_id - 进程ID
// 参数: class_name - 窗口类名
// 参数: title_name - 窗口标题
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindWindowByProcessId(process_id int32, class_name string, title_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindWindowByProcessId", process_id, class_name, title_name)
	}
	funAddr := DmHModule + 104176
	class_namePtr := stringToBytePtr(class_name)
	title_namePtr := stringToBytePtr(title_name)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(process_id), uintptr(unsafe.Pointer(class_namePtr)), uintptr(unsafe.Pointer(title_namePtr)), 0, 0)
	return int32(ret)
}

// WriteDouble 写入双精度浮点数
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: double_value - 双精度浮点数�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteDouble(hwnd int32, addr string, double_value float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteDouble", hwnd, addr, double_value)
	}
	funAddr := DmHModule + 116048
	addrPtr := stringToBytePtr(addr)
	doubleBits := math.Float64bits(double_value)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(doubleBits), uintptr(doubleBits>>32), 0)
	return int32(ret)
}

// GetWindowThreadId 获取窗口线程ID
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWindowThreadId(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetWindowThreadId", hwnd)
	}
	funAddr := DmHModule + 107504
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// GetBindWindow 获取绑定的窗口句�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetBindWindow() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetBindWindow")
	}
	funAddr := DmHModule + 109712
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindWindow 查找窗口
// 参数: class_name - 窗口类名
// 参数: title_name - 窗口标题
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindWindow(class_name string, title_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindWindow", class_name, title_name)
	}
	funAddr := DmHModule + 104288
	class_namePtr := stringToBytePtr(class_name)
	title_namePtr := stringToBytePtr(title_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(class_namePtr)), uintptr(unsafe.Pointer(title_namePtr)))
	return int32(ret)
}

// AiFindPic 使用AI模型查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiFindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("AiFindPic", []interface{}{x1, y1, x2, y2, pic_name, sim, dir}, x, y)
	}
	funAddr := DmHModule + 121536
	pic_namePtr := stringToBytePtr(pic_name)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// FindInt 查找整数
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: int_value_min - 最小整数�?// 参数: int_value_max - 最大整数�?// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) FindInt(hwnd int32, addr_range string, int_value_min int64, int_value_max int64, type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindInt", hwnd, addr_range, int_value_min, int_value_max, type_)
	}
	funAddr := DmHModule + 106256
	addr_rangePtr := stringToBytePtr(addr_range)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(int_value_min), uintptr(int_value_max), uintptr(type_))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// IsBind 判断是否已绑定窗�?// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) IsBind(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("IsBind", hwnd)
	}
	funAddr := DmHModule + 119232
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// SetSimMode 设置模拟模式
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetSimMode(mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetSimMode", mode)
	}
	funAddr := DmHModule + 122896
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(mode), 0)
	return int32(ret)
}

// GetNowDict 获取当前使用的字库索�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetNowDict() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetNowDict")
	}
	funAddr := DmHModule + 101584
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetNetTimeSafe 安全获取网络时间
// 返回: 结果字符�?
func (dm *DmSoft) GetNetTimeSafe() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetNetTimeSafe")
	}
	funAddr := DmHModule + 107760
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetMachineCode 获取机器�?// 返回: 结果字符�?
func (dm *DmSoft) GetMachineCode() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetMachineCode")
	}
	funAddr := DmHModule + 113456
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// VirtualAllocEx 在目标进程分配内�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 参数: type_ - 类型
// 返回: 64位整数�?
func (dm *DmSoft) VirtualAllocEx(hwnd int32, addr int64, size int32, type_ int32) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("VirtualAllocEx", hwnd, addr, size, type_)
	}
	funAddr := DmHModule + 99104
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(size), uintptr(type_))
	return int64(ret)
}

// GetPath 获取当前资源路径
// 返回: 结果字符�?
func (dm *DmSoft) GetPath() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetPath")
	}
	funAddr := DmHModule + 109600
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnumWindowSuper 超级枚举窗口
// 参数: spec1 - 条件1字符�?// 参数: flag1 - 条件1标志
// 参数: type1 - 条件1类型
// 参数: spec2 - 条件2字符�?// 参数: flag2 - 条件2标志
// 参数: type2 - 条件2类型
// 参数: sort - 排序方式
// 返回: 结果字符�?
func (dm *DmSoft) EnumWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32, sort int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumWindowSuper", spec1, flag1, type1, spec2, flag2, type2, sort)
	}
	funAddr := DmHModule + 107360
	spec1Ptr := stringToBytePtr(spec1)
	spec2Ptr := stringToBytePtr(spec2)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(unsafe.Pointer(spec1Ptr)), uintptr(flag1), uintptr(type1), uintptr(unsafe.Pointer(spec2Ptr)), uintptr(flag2), uintptr(type2), uintptr(sort), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetModuleBaseAddr 获取模块基址
// 参数: hwnd - 窗口句柄
// 参数: module_name - 模块名称
// 返回: 64位整数�?
func (dm *DmSoft) GetModuleBaseAddr(hwnd int32, module_name string) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("GetModuleBaseAddr", hwnd, module_name)
	}
	funAddr := DmHModule + 108848
	module_namePtr := stringToBytePtr(module_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(module_namePtr)))
	return int64(ret)
}

// EnumWindowByProcessId 通过进程ID枚举窗口
// 参数: pid - 进程ID
// 参数: title - 窗口标题
// 参数: class_name - 窗口类名
// 参数: filter - 过滤条件
// 返回: 结果字符�?
func (dm *DmSoft) EnumWindowByProcessId(pid int32, title string, class_name string, filter int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumWindowByProcessId", pid, title, class_name, filter)
	}
	funAddr := DmHModule + 124672
	titlePtr := stringToBytePtr(title)
	class_namePtr := stringToBytePtr(class_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(pid), uintptr(unsafe.Pointer(titlePtr)), uintptr(unsafe.Pointer(class_namePtr)), uintptr(filter), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// UnBindWindow 解绑窗口
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) UnBindWindow() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("UnBindWindow")
	}
	funAddr := DmHModule + 101904
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetLastError 获取最后一次错误码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetLastError() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetLastError")
	}
	funAddr := DmHModule + 107936
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FoobarDrawText Foobar绘制文字
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 参数: text - 文本内容
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: align - align
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarDrawText(hwnd int32, x int32, y int32, w int32, h int32, text string, color string, align int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarDrawText", hwnd, x, y, w, h, text, color, align)
	}
	funAddr := DmHModule + 119712
	textPtr := stringToBytePtr(text)
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(align))
	return int32(ret)
}

// SetMinRowGap 设置最小行间距
// 参数: row_gap - 行间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMinRowGap(row_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMinRowGap", row_gap)
	}
	funAddr := DmHModule + 122144
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(row_gap), 0)
	return int32(ret)
}

// LeftUp 鼠标左键弹起
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LeftUp() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LeftUp")
	}
	funAddr := DmHModule + 113680
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// WriteFile 写入文件
// 参数: file - 文件路径
// 参数: content - content
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteFile(file string, content string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteFile", file, content)
	}
	funAddr := DmHModule + 105536
	filePtr := stringToBytePtr(file)
	contentPtr := stringToBytePtr(content)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(contentPtr)))
	return int32(ret)
}

// SetWindowSize 设置窗口大小
// 参数: hwnd - 窗口句柄
// 参数: width - 宽度
// 参数: height - 高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWindowSize(hwnd int32, width int32, height int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWindowSize", hwnd, width, height)
	}
	funAddr := DmHModule + 98560
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(width), uintptr(height), 0, 0)
	return int32(ret)
}

// FaqCaptureFromFile 从文件加载图片到FAQ缓存
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: file - 文件路径
// 参数: quality - 图片质量(1-100)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqCaptureFromFile(x1 int32, y1 int32, x2 int32, y2 int32, file string, quality int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqCaptureFromFile", x1, y1, x2, y2, file, quality)
	}
	funAddr := DmHModule + 116256
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(filePtr)), uintptr(quality), 0, 0)
	return int32(ret)
}

// ReadDataAddr 读取数据(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: length - 长度
// 返回: 结果字符�?
func (dm *DmSoft) ReadDataAddr(hwnd int32, addr int64, length int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadDataAddr", hwnd, addr, length)
	}
	funAddr := DmHModule + 123584
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(length), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// IsSurrpotVt 判断是否支持VT虚拟�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) IsSurrpotVt() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("IsSurrpotVt")
	}
	funAddr := DmHModule + 106992
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetWindowProcessPath 获取窗口进程路径
// 参数: hwnd - 窗口句柄
// 返回: 结果字符�?
func (dm *DmSoft) GetWindowProcessPath(hwnd int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWindowProcessPath", hwnd)
	}
	funAddr := DmHModule + 105232
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ClearDict 清除指定索引的字�?// 参数: index - 索引(�?开�?
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ClearDict(index int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ClearDict", index)
	}
	funAddr := DmHModule + 123152
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return int32(ret)
}

// SaveDict 保存字库到文�?// 参数: index - 索引(�?开�?
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SaveDict(index int32, file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SaveDict", index, file)
	}
	funAddr := DmHModule + 115520
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(index), uintptr(unsafe.Pointer(filePtr)))
	return int32(ret)
}

// ShowTaskBarIcon 显示/隐藏任务栏图�?// 参数: hwnd - 窗口句柄
// 参数: is_show - is_show
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ShowTaskBarIcon(hwnd int32, is_show int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ShowTaskBarIcon", hwnd, is_show)
	}
	funAddr := DmHModule + 119328
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(is_show))
	return int32(ret)
}

// GetAveHSV 获取区域平均HSV�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) GetAveHSV(x1 int32, y1 int32, x2 int32, y2 int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetAveHSV", x1, y1, x2, y2)
	}
	funAddr := DmHModule + 100176
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadIniPwd 读取INI配置(带密�?
// 参数: section - INI节名
// 参数: key - 键名
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 结果字符�?
func (dm *DmSoft) ReadIniPwd(section string, key string, file string, pwd string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadIniPwd", section, key, file, pwd)
	}
	funAddr := DmHModule + 102064
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FaqIsPosted 检查FAQ请求是否完成
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqIsPosted() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqIsPosted")
	}
	funAddr := DmHModule + 102864
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// LeftDown 鼠标左键按下
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LeftDown() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LeftDown")
	}
	funAddr := DmHModule + 106736
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// DmGuardExtract 大漠守护解压
// 参数: type_ - 类型
// 参数: path - 资源路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DmGuardExtract(type_ string, path string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DmGuardExtract", type_, path)
	}
	funAddr := DmHModule + 112160
	type_Ptr := stringToBytePtr(type_)
	pathPtr := stringToBytePtr(path)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(type_Ptr)), uintptr(unsafe.Pointer(pathPtr)))
	return int32(ret)
}

// ExitOs 退出系�?// 参数: type_ - 类型
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) ExitOs(type_ int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("ExitOs", type_)
	}
	funAddr := DmHModule + 115024
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(type_), 0)
	return int32(ret)
}

// FetchWord 提取文字
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: word - 要查找的文字
// 返回: 结果字符�?
func (dm *DmSoft) FetchWord(x1 int32, y1 int32, x2 int32, y2 int32, color string, word string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FetchWord", x1, y1, x2, y2, color, word)
	}
	funAddr := DmHModule + 117840
	colorPtr := stringToBytePtr(color)
	wordPtr := stringToBytePtr(word)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(unsafe.Pointer(wordPtr)), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetDiskSerial 获取磁盘序列�?// 参数: index - 索引(�?开�?
// 返回: 结果字符�?
func (dm *DmSoft) GetDiskSerial(index int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDiskSerial", index)
	}
	funAddr := DmHModule + 112352
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetDictCount 获取字库条目数量
// 参数: index - 索引(�?开�?
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetDictCount(index int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetDictCount", index)
	}
	funAddr := DmHModule + 99584
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return int32(ret)
}

// GetDict 获取字库内容
// 参数: index - 索引(�?开�?
// 参数: font_index - 字体索引
// 返回: 结果字符�?
func (dm *DmSoft) GetDict(index int32, font_index int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDict", index, font_index)
	}
	funAddr := DmHModule + 99184
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(index), uintptr(font_index))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetDict 设置字库文件
// 参数: index - 索引(�?开�?
// 参数: dict_name - 字库文件�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDict(index int32, dict_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDict", index, dict_name)
	}
	funAddr := DmHModule + 121280
	dict_namePtr := stringToBytePtr(dict_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(index), uintptr(unsafe.Pointer(dict_namePtr)))
	return int32(ret)
}

// AiYoloObjectsToString YOLO检测结果转字符�?// 参数: objects - 检测到的对�?// 返回: 结果字符�?
func (dm *DmSoft) AiYoloObjectsToString(objects string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AiYoloObjectsToString", objects)
	}
	funAddr := DmHModule + 111456
	objectsPtr := stringToBytePtr(objects)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(objectsPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetKeyState 获取按键状�?// 参数: vk - 虚拟键码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetKeyState(vk int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetKeyState", vk)
	}
	funAddr := DmHModule + 103296
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(vk), 0)
	return int32(ret)
}

// RightClick 鼠标右键单击
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RightClick() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RightClick")
	}
	funAddr := DmHModule + 101040
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// EnumWindowByProcess 通过进程名枚举窗�?// 参数: process_name - 进程名称
// 参数: title - 窗口标题
// 参数: class_name - 窗口类名
// 参数: filter - 过滤条件
// 返回: 结果字符�?
func (dm *DmSoft) EnumWindowByProcess(process_name string, title string, class_name string, filter int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumWindowByProcess", process_name, title, class_name, filter)
	}
	funAddr := DmHModule + 110192
	process_namePtr := stringToBytePtr(process_name)
	titlePtr := stringToBytePtr(title)
	class_namePtr := stringToBytePtr(class_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(process_namePtr)), uintptr(unsafe.Pointer(titlePtr)), uintptr(unsafe.Pointer(class_namePtr)), uintptr(filter), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetDiskModel 获取磁盘型号
// 参数: index - 索引(�?开�?
// 返回: 结果字符�?
func (dm *DmSoft) GetDiskModel(index int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDiskModel", index)
	}
	funAddr := DmHModule + 102128
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SendStringIme 通过输入法发送字符串
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SendStringIme(str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SendStringIme", str)
	}
	funAddr := DmHModule + 124000
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(strPtr)), 0)
	return int32(ret)
}

// AppendPicAddr 追加图片地址
// 参数: pic_info - 图片信息
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 返回: 结果字符�?
func (dm *DmSoft) AppendPicAddr(pic_info string, addr int32, size int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AppendPicAddr", pic_info, addr, size)
	}
	funAddr := DmHModule + 106832
	pic_infoPtr := stringToBytePtr(pic_info)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(addr), uintptr(size), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// DeleteFolder 删除文件�?// 参数: folder_name - folder_name
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DeleteFolder(folder_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DeleteFolder", folder_name)
	}
	funAddr := DmHModule + 118800
	folder_namePtr := stringToBytePtr(folder_name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(folder_namePtr)), 0)
	return int32(ret)
}

// GetDPI 获取系统DPI
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetDPI() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetDPI")
	}
	funAddr := DmHModule + 107664
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetCpuType 获取CPU类型
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetCpuType() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetCpuType")
	}
	funAddr := DmHModule + 102432
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// WriteIntAddr 写入整数(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteIntAddr(hwnd int32, addr int64, type_ int32, v int64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteIntAddr", hwnd, addr, type_, v)
	}
	funAddr := DmHModule + 100240
	ret, _, _ := syscall.Syscall12(funAddr, 8, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(type_), uintptr(v), uintptr(v>>32), 0, 0, 0, 0, 0)
	return int32(ret)
}

// GetSpecialWindow GetSpecialWindow
// 参数: flag - 查找标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetSpecialWindow(flag int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetSpecialWindow", flag)
	}
	funAddr := DmHModule + 102336
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(flag), 0)
	return int32(ret)
}

// EnumProcess 枚举进程
// 参数: name - 名称
// 返回: 结果字符�?
func (dm *DmSoft) EnumProcess(name string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumProcess", name)
	}
	funAddr := DmHModule + 112288
	namePtr := stringToBytePtr(name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(namePtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// AsmClear 清除汇编代码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AsmClear() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AsmClear")
	}
	funAddr := DmHModule + 119968
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetWindowState 获取窗口状�?// 参数: hwnd - 窗口句柄
// 参数: flag - 查找标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWindowState(hwnd int32, flag int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetWindowState", hwnd, flag)
	}
	funAddr := DmHModule + 100112
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(flag))
	return int32(ret)
}

// FindStrFastE 快速查找文�?返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrFastE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrFastE", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 120288
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetColGapNoDict 设置列间�?无字典模�?
// 参数: col_gap - 列间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetColGapNoDict(col_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetColGapNoDict", col_gap)
	}
	funAddr := DmHModule + 102592
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(col_gap), 0)
	return int32(ret)
}

// AiYoloDetectObjects YOLO目标检�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: prob - 概率阈�?// 参数: iou - IOU阈�?// 返回: 结果字符�?
func (dm *DmSoft) AiYoloDetectObjects(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("AiYoloDetectObjects", x1, y1, x2, y2, float64(prob), float64(iou))
	}
	funAddr := DmHModule + 116112
	probBits := uint32(math.Float32bits(prob))
	iouBits := uint32(math.Float32bits(iou))
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(probBits), uintptr(iouBits), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// RunApp 运行程序
// 参数: path - 资源路径
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RunApp(path string, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RunApp", path, mode)
	}
	funAddr := DmHModule + 122832
	pathPtr := stringToBytePtr(path)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(pathPtr)), uintptr(mode))
	return int32(ret)
}

// FindString 查找字符�?// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: string_value - 字符串�?// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) FindString(hwnd int32, addr_range string, string_value string, type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindString", hwnd, addr_range, string_value, type_)
	}
	funAddr := DmHModule + 110752
	addr_rangePtr := stringToBytePtr(addr_range)
	string_valuePtr := stringToBytePtr(string_value)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(unsafe.Pointer(string_valuePtr)), uintptr(type_), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetOsType 获取操作系统类型
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetOsType() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetOsType")
	}
	funAddr := DmHModule + 121632
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// Ocr OCR识别指定区域文字
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) Ocr(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Ocr", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 110992
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadString 读取字符�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 参数: length - 长度
// 返回: 结果字符�?
func (dm *DmSoft) ReadString(hwnd int32, addr string, type_ int32, length int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ReadString", hwnd, addr, type_, length)
	}
	funAddr := DmHModule + 121472
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(type_), uintptr(length), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// ReadFloatAddr 读取浮点�?指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 返回: 32位浮点数
func (dm *DmSoft) ReadFloatAddr(hwnd int32, addr int64) float32 {
	if dm.pipeConn != nil {
		return float32(dm.comCallFloat64("ReadFloatAddr", hwnd, addr))
	}
	funAddr := DmHModule + 100816
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), 0, 0)
	return float32(ret)
}

// Beep 蜂鸣�?// 参数: fre - 频率(Hz)
// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Beep(fre int32, delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Beep", fre, delay)
	}
	funAddr := DmHModule + 104544
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(fre), uintptr(delay))
	return int32(ret)
}

// LoadAi 从文件加载AI模型
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LoadAi(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LoadAi", file)
	}
	funAddr := DmHModule + 106944
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// GetCpuUsage 获取CPU使用�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetCpuUsage() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetCpuUsage")
	}
	funAddr := DmHModule + 121072
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// EnableShareDict 启用共享字库
// 参数: en - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableShareDict(en int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableShareDict", en)
	}
	funAddr := DmHModule + 108992
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(en), 0)
	return int32(ret)
}

// AiYoloDetectObjectsToFile YOLO目标检测并保存到文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: prob - 概率阈�?// 参数: iou - IOU阈�?// 参数: file - 文件路径
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloDetectObjectsToFile(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32, file string, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloDetectObjectsToFile", x1, y1, x2, y2, float64(prob), float64(iou), file, mode)
	}
	funAddr := DmHModule + 109504
	filePtr := stringToBytePtr(file)
	probBits := uint32(math.Float32bits(prob))
	iouBits := uint32(math.Float32bits(iou))
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(probBits), uintptr(iouBits), uintptr(unsafe.Pointer(filePtr)), uintptr(mode))
	return int32(ret)
}

// FoobarUnlock 解锁Foobar窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarUnlock(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarUnlock", hwnd)
	}
	funAddr := DmHModule + 123952
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// GetSystemInfo 获取系统信息
// 参数: type_ - 类型
// 参数: method - method
// 返回: 结果字符�?
func (dm *DmSoft) GetSystemInfo(type_ string, method int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetSystemInfo", type_, method)
	}
	funAddr := DmHModule + 115680
	type_Ptr := stringToBytePtr(type_)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(type_Ptr)), uintptr(method))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetResultCount GetResultCount
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetResultCount(str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetResultCount", str)
	}
	funAddr := DmHModule + 116720
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(strPtr)), 0)
	return int32(ret)
}

// EnumWindow 枚举窗口
// 参数: parent - 父窗口句�?// 参数: title - 窗口标题
// 参数: class_name - 窗口类名
// 参数: filter - 过滤条件
// 返回: 结果字符�?
func (dm *DmSoft) EnumWindow(parent int32, title string, class_name string, filter int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumWindow", parent, title, class_name, filter)
	}
	funAddr := DmHModule + 115296
	titlePtr := stringToBytePtr(title)
	class_namePtr := stringToBytePtr(class_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(parent), uintptr(unsafe.Pointer(titlePtr)), uintptr(unsafe.Pointer(class_namePtr)), uintptr(filter), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetResultPos GetResultPos
// 参数: str - 要查找的字符�?// 参数: index - 索引(�?开�?
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetResultPos(str string, index int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetResultPos", []interface{}{str, index}, x, y)
	}
	funAddr := DmHModule + 102800
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(strPtr)), uintptr(index), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0)
	return int32(ret)
}

// KeyDown 按下按键(虚拟键码)
// 参数: vk - 虚拟键码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyDown(vk int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyDown", vk)
	}
	funAddr := DmHModule + 115120
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(vk), 0)
	return int32(ret)
}

// SetWordLineHeightNoDict 设置文字识别行高(无字典模�?
// 参数: line_height - 行高(像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWordLineHeightNoDict(line_height int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWordLineHeightNoDict", line_height)
	}
	funAddr := DmHModule + 103792
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(line_height), 0)
	return int32(ret)
}

// AiFindPicMem AI内存查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiFindPicMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("AiFindPicMem", []interface{}{x1, y1, x2, y2, pic_info, sim, dir}, x, y)
	}
	funAddr := DmHModule + 111696
	pic_infoPtr := stringToBytePtr(pic_info)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return int32(ret)
}

// FoobarTextRect 设置Foobar文字区域
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarTextRect(hwnd int32, x int32, y int32, w int32, h int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarTextRect", hwnd, x, y, w, h)
	}
	funAddr := DmHModule + 108784
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h))
	return int32(ret)
}

// GetPointWindow 获取指定坐标的窗�?// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetPointWindow(x int32, y int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetPointWindow", x, y)
	}
	funAddr := DmHModule + 118544
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(x), uintptr(y))
	return int32(ret)
}

// FindMultiColorEx 高级多点找色
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: first_color - 第一个颜�?// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindMultiColorEx(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindMultiColorEx", x1, y1, x2, y2, first_color, offset_color, sim, dir)
	}
	funAddr := DmHModule + 122560
	first_colorPtr := stringToBytePtr(first_color)
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(first_colorPtr)), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FreeProcessMemory 释放进程内存
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FreeProcessMemory(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FreeProcessMemory", hwnd)
	}
	funAddr := DmHModule + 111120
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// GetMachineCodeNoMac 获取机器�?不含MAC地址)
// 返回: 结果字符�?
func (dm *DmSoft) GetMachineCodeNoMac() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetMachineCodeNoMac")
	}
	funAddr := DmHModule + 120544
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindWindowByProcess 通过进程名查找窗�?// 参数: process_name - 进程名称
// 参数: class_name - 窗口类名
// 参数: title_name - 窗口标题
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindWindowByProcess(process_name string, class_name string, title_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindWindowByProcess", process_name, class_name, title_name)
	}
	funAddr := DmHModule + 122336
	process_namePtr := stringToBytePtr(process_name)
	class_namePtr := stringToBytePtr(class_name)
	title_namePtr := stringToBytePtr(title_name)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(process_namePtr)), uintptr(unsafe.Pointer(class_namePtr)), uintptr(unsafe.Pointer(title_namePtr)), 0, 0)
	return int32(ret)
}

// GetColorNum 获取区域中指定颜色的数量
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetColorNum(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetColorNum", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 124048
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0, 0)
	return int32(ret)
}

// SetWindowState 设置窗口状�?// 参数: hwnd - 窗口句柄
// 参数: flag - 查找标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWindowState(hwnd int32, flag int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWindowState", hwnd, flag)
	}
	funAddr := DmHModule + 102736
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(flag))
	return int32(ret)
}

// CheckFontSmooth 检查字体平�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CheckFontSmooth() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CheckFontSmooth")
	}
	funAddr := DmHModule + 117552
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// IsFolderExist 判断文件夹是否存�?// 参数: folder - folder
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) IsFolderExist(folder string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("IsFolderExist", folder)
	}
	funAddr := DmHModule + 121184
	folderPtr := stringToBytePtr(folder)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(folderPtr)), 0)
	return int32(ret)
}

// FaqCancel 取消FAQ请求
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqCancel() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqCancel")
	}
	funAddr := DmHModule + 113968
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetWindowTransparent 设置窗口透明�?// 参数: hwnd - 窗口句柄
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWindowTransparent(hwnd int32, v int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWindowTransparent", hwnd, v)
	}
	funAddr := DmHModule + 112896
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(v))
	return int32(ret)
}

// SwitchBindWindow 切换绑定窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SwitchBindWindow(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SwitchBindWindow", hwnd)
	}
	funAddr := DmHModule + 109920
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// EnableFontSmooth 启用字体平滑
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableFontSmooth() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableFontSmooth")
	}
	funAddr := DmHModule + 103936
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// StringToData 字符串转数据
// 参数: string_value - 字符串�?// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) StringToData(string_value string, type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("StringToData", string_value, type_)
	}
	funAddr := DmHModule + 114768
	string_valuePtr := stringToBytePtr(string_value)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(string_valuePtr)), uintptr(type_))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetWindowRect 获取窗口矩形
// 参数: hwnd - 窗口句柄
// 参数: x1 - 左上角X坐标(输出参数)
// 参数: y1 - 左上角Y坐标(输出参数)
// 参数: x2 - 右下角X坐标(输出参数)
// 参数: y2 - 右下角Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetWindowRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("GetWindowRect", []interface{}{hwnd}, x1, y1, x2, y2)
	}
	funAddr := DmHModule + 122656
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(x1)), uintptr(unsafe.Pointer(y1)), uintptr(unsafe.Pointer(x2)), uintptr(unsafe.Pointer(y2)))
	return int32(ret)
}

// FindPicEx 高级查找图片,返回所有匹配位�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicEx", x1, y1, x2, y2, pic_name, delta_color, sim, dir)
	}
	funAddr := DmHModule + 108160
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetWords 获取区域内的所有文�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) GetWords(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWords", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 107808
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetExactOcr 设置精确OCR模式
// 参数: exact_ocr - 精确OCR标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetExactOcr(exact_ocr int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetExactOcr", exact_ocr)
	}
	funAddr := DmHModule + 123280
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(exact_ocr), 0)
	return int32(ret)
}

// EnableMouseSync EnableMouseSync
// 参数: enable - 启用标志(1:启用,0:禁用)
// 参数: time_out - 超时时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableMouseSync(enable int32, time_out int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableMouseSync", enable, time_out)
	}
	funAddr := DmHModule + 98496
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(enable), uintptr(time_out))
	return int32(ret)
}

// CapturePre CapturePre
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CapturePre(file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CapturePre", file)
	}
	funAddr := DmHModule + 109456
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(filePtr)), 0)
	return int32(ret)
}

// BindWindowEx 扩展绑定窗口
// 参数: hwnd - 窗口句柄
// 参数: display - 显示模式
// 参数: mouse - 鼠标模式
// 参数: keypad - 键盘模式
// 参数: public_desc - 公共描述
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) BindWindowEx(hwnd int32, display string, mouse string, keypad string, public_desc string, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("BindWindowEx", hwnd, display, mouse, keypad, public_desc, mode)
	}
	funAddr := DmHModule + 99456
	displayPtr := stringToBytePtr(display)
	mousePtr := stringToBytePtr(mouse)
	keypadPtr := stringToBytePtr(keypad)
	public_descPtr := stringToBytePtr(public_desc)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(displayPtr)), uintptr(unsafe.Pointer(mousePtr)), uintptr(unsafe.Pointer(keypadPtr)), uintptr(unsafe.Pointer(public_descPtr)), uintptr(mode), 0, 0)
	return int32(ret)
}

// FaqCaptureString 将字符串添加到FAQ缓存
// 参数: str - 要查找的字符�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FaqCaptureString(str string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FaqCaptureString", str)
	}
	funAddr := DmHModule + 106208
	strPtr := stringToBytePtr(str)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(strPtr)), 0)
	return int32(ret)
}

// FoobarTextLineGap 设置Foobar文字行间�?// 参数: hwnd - 窗口句柄
// 参数: gap - gap
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarTextLineGap(hwnd int32, gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarTextLineGap", hwnd, gap)
	}
	funAddr := DmHModule + 124848
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(gap))
	return int32(ret)
}

// FoobarDrawLine Foobar绘制线条
// 参数: hwnd - 窗口句柄
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: style - style
// 参数: width - 宽度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarDrawLine(hwnd int32, x1 int32, y1 int32, x2 int32, y2 int32, color string, style int32, width int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarDrawLine", hwnd, x1, y1, x2, y2, color, style, width)
	}
	funAddr := DmHModule + 116384
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(hwnd), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(style), uintptr(width))
	return int32(ret)
}

// FindInputMethod 查找输入�?// 参数: id - 标识ID
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindInputMethod(id string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindInputMethod", id)
	}
	funAddr := DmHModule + 113872
	idPtr := stringToBytePtr(id)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(idPtr)), 0)
	return int32(ret)
}

// SetPicPwd SetPicPwd
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetPicPwd(pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetPicPwd", pwd)
	}
	funAddr := DmHModule + 123712
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pwdPtr)), 0)
	return int32(ret)
}

// GetCursorSpot 获取鼠标光点位置
// 返回: 结果字符�?
func (dm *DmSoft) GetCursorSpot() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetCursorSpot")
	}
	funAddr := DmHModule + 125056
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// InitCri 初始化临界区
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) InitCri() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("InitCri")
	}
	funAddr := DmHModule + 120240
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindPicMemE 在内存中查找图片,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicMemE(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicMemE", x1, y1, x2, y2, pic_info, delta_color, sim, dir)
	}
	funAddr := DmHModule + 109264
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindStrFastS 快速查找文�?返回字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrFastS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStrWithOutVars("FindStrFastS", []interface{}{x1, y1, x2, y2, str, color, sim}, x, y)
	}
	funAddr := DmHModule + 98672
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 11, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// DeleteIniPwd 删除INI配置�?带密�?
// 参数: section - INI节名
// 参数: key - 键名
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DeleteIniPwd(section string, key string, file string, pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DeleteIniPwd", section, key, file, pwd)
	}
	funAddr := DmHModule + 99344
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)), 0)
	return int32(ret)
}

// AiYoloDetectObjectsToDataBmp YOLO目标检测到BMP数据
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: prob - 概率阈�?// 参数: iou - IOU阈�?// 参数: data - 数据(输出参数)
// 参数: size - 大小(字节)(输出参数)
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloDetectObjectsToDataBmp(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32, data *int32, size *int32, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("AiYoloDetectObjectsToDataBmp", []interface{}{x1, y1, x2, y2, float64(prob), float64(iou), mode}, data, size)
	}
	funAddr := DmHModule + 98928
	probBits := uint32(math.Float32bits(prob))
	iouBits := uint32(math.Float32bits(iou))
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(probBits), uintptr(iouBits), uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(size)), uintptr(mode), 0, 0)
	return int32(ret)
}

// AiYoloFreeModel 释放YOLO模型
// 参数: index - 索引(�?开�?
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AiYoloFreeModel(index int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AiYoloFreeModel", index)
	}
	funAddr := DmHModule + 106592
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(index), 0)
	return int32(ret)
}

// DisableFontSmooth 禁用字体平滑
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DisableFontSmooth() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DisableFontSmooth")
	}
	funAddr := DmHModule + 118368
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// SetExitThread 设置退出线�?// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetExitThread(mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetExitThread", mode)
	}
	funAddr := DmHModule + 101536
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(mode), 0)
	return int32(ret)
}

// FindPicMemEx 高级在内存中查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicMemEx", x1, y1, x2, y2, pic_info, delta_color, sim, dir)
	}
	funAddr := DmHModule + 101440
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetDmCount 获取大漠对象计数
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetDmCount() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetDmCount")
	}
	funAddr := DmHModule + 125008
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FindMulColor FindMulColor
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindMulColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FindMulColor", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 111552
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return int32(ret)
}

// FaqFetch 获取FAQ返回结果
// 返回: 结果字符�?
func (dm *DmSoft) FaqFetch() string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FaqFetch")
	}
	funAddr := DmHModule + 117744
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// RegExNoMac 扩展注册大漠插件(不含MAC)
// 参数: code - 注册�?// 参数: ver - 版本�?// 参数: ip - IP地址
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RegExNoMac(code string, ver string, ip string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RegExNoMac", code, ver, ip)
	}
	funAddr := DmHModule + 107552
	codePtr := stringToBytePtr(code)
	verPtr := stringToBytePtr(ver)
	ipPtr := stringToBytePtr(ip)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(codePtr)), uintptr(unsafe.Pointer(verPtr)), uintptr(unsafe.Pointer(ipPtr)), 0, 0)
	return int32(ret)
}

// FoobarUpdate 更新Foobar窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarUpdate(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarUpdate", hwnd)
	}
	funAddr := DmHModule + 119280
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// ReadDouble 读取双精度浮点数
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 返回: 64位浮点数
func (dm *DmSoft) ReadDouble(hwnd int32, addr string) float64 {
	if dm.pipeConn != nil {
		return dm.comCallFloat64("ReadDouble", hwnd, addr)
	}
	funAddr := DmHModule + 110128
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)))
	return float64(ret)
}

// GetCursorShapeEx 获取鼠标形状(扩展)
// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) GetCursorShapeEx(type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetCursorShapeEx", type_)
	}
	funAddr := DmHModule + 117488
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(type_), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// DoubleToData 双精度转数据
// 参数: double_value - 双精度浮点数�?// 返回: 结果字符�?
func (dm *DmSoft) DoubleToData(double_value float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("DoubleToData", double_value)
	}
	funAddr := DmHModule + 111856
	doubleBits := math.Float64bits(double_value)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(doubleBits), uintptr(doubleBits>>32))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetWordGapNoDict 设置字间�?无字典模�?
// 参数: word_gap - 字间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWordGapNoDict(word_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWordGapNoDict", word_gap)
	}
	funAddr := DmHModule + 123392
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(word_gap), 0)
	return int32(ret)
}

// ReadDoubleAddr 读取双精度浮点数(指定地址)
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 返回: 64位浮点数
func (dm *DmSoft) ReadDoubleAddr(hwnd int32, addr int64) float64 {
	if dm.pipeConn != nil {
		return dm.comCallFloat64("ReadDoubleAddr", hwnd, addr)
	}
	funAddr := DmHModule + 113392
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), 0, 0)
	return float64(ret)
}

// FoobarLock 锁定Foobar窗口
// 参数: hwnd - 窗口句柄
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarLock(hwnd int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarLock", hwnd)
	}
	funAddr := DmHModule + 109824
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return int32(ret)
}

// FindStrFastExS 高级快速查找文�?返回详细字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrFastExS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrFastExS", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 124176
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FindStrWithFont 指定字体查找文字
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 参数: font_name - 字体名称
// 参数: font_size - 字体大小(像素)
// 参数: flag - 查找标志
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindStrWithFont(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindStrWithFont", []interface{}{x1, y1, x2, y2, str, color, sim, font_name, font_size, flag}, x, y)
	}
	funAddr := DmHModule + 119856
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	font_namePtr := stringToBytePtr(font_name)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 14, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(unsafe.Pointer(font_namePtr)), uintptr(font_size), uintptr(flag), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0)
	return int32(ret)
}

// VirtualProtectEx 修改目标进程内存保护属�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 参数: type_ - 类型
// 参数: old_protect - 旧保护属性指�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) VirtualProtectEx(hwnd int32, addr int64, size int32, type_ int32, old_protect int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("VirtualProtectEx", hwnd, addr, size, type_, old_protect)
	}
	funAddr := DmHModule + 108912
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), uintptr(size), uintptr(type_), uintptr(old_protect), 0, 0)
	return int32(ret)
}

// GetWindowClass 获取窗口类名
// 参数: hwnd - 窗口句柄
// 返回: 结果字符�?
func (dm *DmSoft) GetWindowClass(hwnd int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetWindowClass", hwnd)
	}
	funAddr := DmHModule + 117056
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(hwnd), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetMouseDelay 设置鼠标操作延迟
// 参数: type_ - 类型
// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetMouseDelay(type_ string, delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetMouseDelay", type_, delay)
	}
	funAddr := DmHModule + 104592
	type_Ptr := stringToBytePtr(type_)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(type_Ptr)), uintptr(delay))
	return int32(ret)
}

// ReadInt 读取整数
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: type_ - 类型
// 返回: 64位整数�?
func (dm *DmSoft) ReadInt(hwnd int32, addr string, type_ int32) int64 {
	if dm.pipeConn != nil {
		return dm.comCallInt64("ReadInt", hwnd, addr, type_)
	}
	funAddr := DmHModule + 112720
	addrPtr := stringToBytePtr(addr)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(type_), 0, 0)
	return int64(ret)
}

// GetAveRGB GetAveRGB
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) GetAveRGB(x1 int32, y1 int32, x2 int32, y2 int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetAveRGB", x1, y1, x2, y2)
	}
	funAddr := DmHModule + 118192
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// GetScreenData 获取屏幕数据
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetScreenData(x1 int32, y1 int32, x2 int32, y2 int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetScreenData", x1, y1, x2, y2)
	}
	funAddr := DmHModule + 125104
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), 0)
	return int32(ret)
}

// GetMouseSpeed 获取鼠标移动速度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetMouseSpeed() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetMouseSpeed")
	}
	funAddr := DmHModule + 99248
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// Int64ToInt32 int64转int32
// 参数: v - �?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) Int64ToInt32(v int64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("Int64ToInt32", v)
	}
	funAddr := DmHModule + 110880
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(v), 0)
	return int32(ret)
}

// FindFloatEx 高级查找浮点�?// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: float_value_min - 最小浮点�?// 参数: float_value_max - 最大浮点�?// 参数: step - 步长
// 参数: multi_thread - 多线程数�?// 参数: mode - 模式
// 返回: 结果字符�?
func (dm *DmSoft) FindFloatEx(hwnd int32, addr_range string, float_value_min float32, float_value_max float32, step int32, multi_thread int32, mode int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindFloatEx", hwnd, addr_range, float64(float_value_min), float64(float_value_max), step, multi_thread, mode)
	}
	funAddr := DmHModule + 107040
	addr_rangePtr := stringToBytePtr(addr_range)
	minBits := uint32(math.Float32bits(float_value_min))
	maxBits := uint32(math.Float32bits(float_value_max))
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(minBits), uintptr(maxBits), uintptr(step), uintptr(multi_thread), uintptr(mode), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FoobarPrintText Foobar打印文字
// 参数: hwnd - 窗口句柄
// 参数: text - 文本内容
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarPrintText(hwnd int32, text string, color string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarPrintText", hwnd, text, color)
	}
	funAddr := DmHModule + 108720
	textPtr := stringToBytePtr(text)
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(colorPtr)), 0, 0)
	return int32(ret)
}

// OcrEx 高级OCR识别
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) OcrEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("OcrEx", x1, y1, x2, y2, color, sim)
	}
	funAddr := DmHModule + 113168
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// FreePic FreePic
// 参数: pic_name - 图片名称(多个用|分隔)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FreePic(pic_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FreePic", pic_name)
	}
	funAddr := DmHModule + 103408
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pic_namePtr)), 0)
	return int32(ret)
}

// WriteData 写入数据
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: data - 数据
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteData(hwnd int32, addr string, data string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteData", hwnd, addr, data)
	}
	funAddr := DmHModule + 123040
	addrPtr := stringToBytePtr(addr)
	dataPtr := stringToBytePtr(data)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(unsafe.Pointer(dataPtr)), 0, 0)
	return int32(ret)
}

// MoveDD DD驱动移动鼠标
// 参数: dx - dx
// 参数: dy - dy
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) MoveDD(dx int32, dy int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("MoveDD", dx, dy)
	}
	funAddr := DmHModule + 121840
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(dx), uintptr(dy))
	return int32(ret)
}

// SetShowErrorMsg 设置是否显示错误信息
// 参数: show - 显示标志(1:显示,0:隐藏)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetShowErrorMsg(show int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetShowErrorMsg", show)
	}
	funAddr := DmHModule + 101856
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(show), 0)
	return int32(ret)
}

// SetDictMem 从内存设置字�?// 参数: index - 索引(�?开�?
// 参数: addr - 内存地址
// 参数: size - 大小(字节)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDictMem(index int32, addr int32, size int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDictMem", index, addr, size)
	}
	funAddr := DmHModule + 104704
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(index), uintptr(addr), uintptr(size), 0, 0)
	return int32(ret)
}

// SetClipboard 设置剪贴板内�?// 参数: data - 数据
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetClipboard(data string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetClipboard", data)
	}
	funAddr := DmHModule + 104960
	dataPtr := stringToBytePtr(data)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(dataPtr)), 0)
	return int32(ret)
}

// FindPicMem 在内存中查找图片
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_info - 图片信息
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 参数: x - X坐标(输出参数)
// 参数: y - Y坐标(输出参数)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FindPicMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallWithOutVars("FindPicMem", []interface{}{x1, y1, x2, y2, pic_info, delta_color, sim, dir}, x, y)
	}
	funAddr := DmHModule + 103696
	pic_infoPtr := stringToBytePtr(pic_info)
	delta_colorPtr := stringToBytePtr(delta_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall15(funAddr, 12, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_infoPtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), uintptr(unsafe.Pointer(x)), uintptr(unsafe.Pointer(y)), 0, 0, 0)
	return int32(ret)
}

// CreateFoobarRoundRect 创建圆角矩形Foobar窗口
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 参数: rw - rw
// 参数: rh - rh
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) CreateFoobarRoundRect(hwnd int32, x int32, y int32, w int32, h int32, rw int32, rh int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("CreateFoobarRoundRect", hwnd, x, y, w, h, rw, rh)
	}
	funAddr := DmHModule + 108352
	ret, _, _ := syscall.Syscall9(funAddr, 8, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(rw), uintptr(rh), 0)
	return int32(ret)
}

// WriteFloat 写入浮点�?// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 参数: float_value - 浮点数�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) WriteFloat(hwnd int32, addr string, float_value float32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("WriteFloat", hwnd, addr, float64(float_value))
	}
	funAddr := DmHModule + 111920
	addrPtr := stringToBytePtr(addr)
	floatBits := uint32(math.Float32bits(float_value))
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addrPtr)), uintptr(floatBits), 0, 0)
	return int32(ret)
}

// VirtualFreeEx VirtualFreeEx
// 参数: hwnd - 窗口句柄
// 参数: addr - 内存地址
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) VirtualFreeEx(hwnd int32, addr int64) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("VirtualFreeEx", hwnd, addr)
	}
	funAddr := DmHModule + 105120
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(addr), uintptr(addr>>32), 0, 0)
	return int32(ret)
}

// GetDictInfo 获取字库信息
// 参数: str - 要查找的字符�?// 参数: font_name - 字体名称
// 参数: font_size - 字体大小(像素)
// 参数: flag - 查找标志
// 返回: 结果字符�?
func (dm *DmSoft) GetDictInfo(str string, font_name string, font_size int32, flag int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDictInfo", str, font_name, font_size, flag)
	}
	funAddr := DmHModule + 100624
	strPtr := stringToBytePtr(str)
	font_namePtr := stringToBytePtr(font_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(font_namePtr)), uintptr(font_size), uintptr(flag), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// KeyPress 按键(虚拟键码)
// 参数: vk - 虚拟键码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyPress(vk int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyPress", vk)
	}
	funAddr := DmHModule + 118688
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(vk), 0)
	return int32(ret)
}

// SetClientSize SetClientSize
// 参数: hwnd - 窗口句柄
// 参数: width - 宽度
// 参数: height - 高度
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetClientSize(hwnd int32, width int32, height int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetClientSize", hwnd, width, height)
	}
	funAddr := DmHModule + 104896
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(hwnd), uintptr(width), uintptr(height), 0, 0)
	return int32(ret)
}

// ExcludePos 排除位置
// 参数: all_pos - 所有位置字符串
// 参数: type_ - 类型
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) ExcludePos(all_pos string, type_ int32, x1 int32, y1 int32, x2 int32, y2 int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("ExcludePos", all_pos, type_, x1, y1, x2, y2)
	}
	funAddr := DmHModule + 120992
	all_posPtr := stringToBytePtr(all_pos)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(unsafe.Pointer(all_posPtr)), uintptr(type_), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// MoveToEx 扩展移动鼠标,支持随机偏移
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: w - 宽度
// 参数: h - 高度
// 返回: 结果字符�?
func (dm *DmSoft) MoveToEx(x int32, y int32, w int32, h int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("MoveToEx", x, y, w, h)
	}
	funAddr := DmHModule + 120688
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetDictPwd 设置字库密码
// 参数: pwd - 密码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetDictPwd(pwd string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetDictPwd", pwd)
	}
	funAddr := DmHModule + 104128
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pwdPtr)), 0)
	return int32(ret)
}

// FoobarSetFont 设置Foobar字体
// 参数: hwnd - 窗口句柄
// 参数: font_name - 字体名称
// 参数: size - 大小(字节)
// 参数: flag - 查找标志
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarSetFont(hwnd int32, font_name string, size int32, flag int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarSetFont", hwnd, font_name, size, flag)
	}
	funAddr := DmHModule + 111632
	font_namePtr := stringToBytePtr(font_name)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(font_namePtr)), uintptr(size), uintptr(flag), 0)
	return int32(ret)
}

// GetNetTimeByIp GetNetTimeByIp
// 参数: ip - IP地址
// 返回: 结果字符�?
func (dm *DmSoft) GetNetTimeByIp(ip string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetNetTimeByIp", ip)
	}
	funAddr := DmHModule + 105360
	ipPtr := stringToBytePtr(ip)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(ipPtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableKeypadPatch 启用键盘补丁
// 参数: enable - 启用标志(1:启用,0:禁用)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableKeypadPatch(enable int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableKeypadPatch", enable)
	}
	funAddr := DmHModule + 116672
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(enable), 0)
	return int32(ret)
}

// FoobarStartGif Foobar开始播放GIF
// 参数: hwnd - 窗口句柄
// 参数: x - X坐标
// 参数: y - Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: repeat_limit - repeat_limit
// 参数: delay - 延迟时间(毫秒)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarStartGif(hwnd int32, x int32, y int32, pic_name string, repeat_limit int32, delay int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarStartGif", hwnd, x, y, pic_name, repeat_limit, delay)
	}
	funAddr := DmHModule + 117664
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(repeat_limit), uintptr(delay), 0, 0)
	return int32(ret)
}

// FindMultiColorE FindMultiColorE
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: first_color - 第一个颜�?// 参数: offset_color - 偏移颜色
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindMultiColorE(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindMultiColorE", x1, y1, x2, y2, first_color, offset_color, sim, dir)
	}
	funAddr := DmHModule + 101696
	first_colorPtr := stringToBytePtr(first_color)
	offset_colorPtr := stringToBytePtr(offset_color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 10, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(first_colorPtr)), uintptr(unsafe.Pointer(offset_colorPtr)), uintptr(simBits), uintptr(simBits>>32), uintptr(dir), 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// SetWordGap 设置字间�?// 参数: word_gap - 字间�?像素)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) SetWordGap(word_gap int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("SetWordGap", word_gap)
	}
	funAddr := DmHModule + 98624
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(word_gap), 0)
	return int32(ret)
}

// GetLocale 获取区域
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetLocale() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetLocale")
	}
	funAddr := DmHModule + 122096
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetModuleSize 获取模块大小
// 参数: hwnd - 窗口句柄
// 参数: module_name - 模块名称
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetModuleSize(hwnd int32, module_name string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetModuleSize", hwnd, module_name)
	}
	funAddr := DmHModule + 120016
	module_namePtr := stringToBytePtr(module_name)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(module_namePtr)))
	return int32(ret)
}

// FindStrE 查找文字,返回坐标字符�?// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: str - 要查找的字符�?// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 参数: sim - 相似�?0.1-1.0)
// 返回: 结果字符�?
func (dm *DmSoft) FindStrE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindStrE", x1, y1, x2, y2, str, color, sim)
	}
	funAddr := DmHModule + 122400
	strPtr := stringToBytePtr(str)
	colorPtr := stringToBytePtr(color)
	simBits := math.Float64bits(sim)
	ret, _, _ := syscall.Syscall12(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(strPtr)), uintptr(unsafe.Pointer(colorPtr)), uintptr(simBits), uintptr(simBits>>32), 0, 0, 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// KeyUp 弹起按键(虚拟键码)
// 参数: vk - 虚拟键码
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) KeyUp(vk int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("KeyUp", vk)
	}
	funAddr := DmHModule + 113248
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(vk), 0)
	return int32(ret)
}

// SortPosDistance 按距离排序位�?// 参数: all_pos - 所有位置字符串
// 参数: type_ - 类型
// 参数: x - X坐标
// 参数: y - Y坐标
// 返回: 结果字符�?
func (dm *DmSoft) SortPosDistance(all_pos string, type_ int32, x int32, y int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("SortPosDistance", all_pos, type_, x, y)
	}
	funAddr := DmHModule + 117120
	all_posPtr := stringToBytePtr(all_pos)
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(unsafe.Pointer(all_posPtr)), uintptr(type_), uintptr(x), uintptr(y), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnableDisplayDebug 启用显示调试
// 参数: enable_debug - enable_debug
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) EnableDisplayDebug(enable_debug int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("EnableDisplayDebug", enable_debug)
	}
	funAddr := DmHModule + 99296
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(enable_debug), 0)
	return int32(ret)
}

// DeleteIni 删除INI配置�?// 参数: section - INI节名
// 参数: key - 键名
// 参数: file - 文件路径
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) DeleteIni(section string, key string, file string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("DeleteIni", section, key, file)
	}
	funAddr := DmHModule + 111168
	sectionPtr := stringToBytePtr(section)
	keyPtr := stringToBytePtr(key)
	filePtr := stringToBytePtr(file)
	ret, _, _ := syscall.Syscall6(funAddr, 4, dm.obj, uintptr(unsafe.Pointer(sectionPtr)), uintptr(unsafe.Pointer(keyPtr)), uintptr(unsafe.Pointer(filePtr)), 0, 0)
	return int32(ret)
}

// FindIntEx 高级查找整数
// 参数: hwnd - 窗口句柄
// 参数: addr_range - 地址范围
// 参数: int_value_min - 最小整数�?// 参数: int_value_max - 最大整数�?// 参数: type_ - 类型
// 参数: step - 步长
// 参数: multi_thread - 多线程数�?// 参数: mode - 模式
// 返回: 结果字符�?
func (dm *DmSoft) FindIntEx(hwnd int32, addr_range string, int_value_min int64, int_value_max int64, type_ int32, step int32, multi_thread int32, mode int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindIntEx", hwnd, addr_range, int_value_min, int_value_max, type_, step, multi_thread, mode)
	}
	funAddr := DmHModule + 107216
	addr_rangePtr := stringToBytePtr(addr_range)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(addr_rangePtr)), uintptr(int_value_min), uintptr(int_value_max), uintptr(type_), uintptr(step), uintptr(multi_thread), uintptr(mode))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// BindWindow 绑定窗口
// 参数: hwnd - 窗口句柄
// 参数: display - 显示模式
// 参数: mouse - 鼠标模式
// 参数: keypad - 键盘模式
// 参数: mode - 模式
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) BindWindow(hwnd int32, display string, mouse string, keypad string, mode int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("BindWindow", hwnd, display, mouse, keypad, mode)
	}
	funAddr := DmHModule + 120080
	displayPtr := stringToBytePtr(display)
	mousePtr := stringToBytePtr(mouse)
	keypadPtr := stringToBytePtr(keypad)
	ret, _, _ := syscall.Syscall6(funAddr, 6, dm.obj, uintptr(hwnd), uintptr(unsafe.Pointer(displayPtr)), uintptr(unsafe.Pointer(mousePtr)), uintptr(unsafe.Pointer(keypadPtr)), uintptr(mode))
	return int32(ret)
}

// GetPicSize GetPicSize
// 参数: pic_name - 图片名称(多个用|分隔)
// 返回: 结果字符�?
func (dm *DmSoft) GetPicSize(pic_name string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetPicSize", pic_name)
	}
	funAddr := DmHModule + 114960
	pic_namePtr := stringToBytePtr(pic_name)
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(unsafe.Pointer(pic_namePtr)), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// AsmSetTimeout 设置汇编执行超时
// 参数: time_out - 超时时间(毫秒)
// 参数: param - 参数
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) AsmSetTimeout(time_out int32, param int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("AsmSetTimeout", time_out, param)
	}
	funAddr := DmHModule + 117920
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(time_out), uintptr(param))
	return int32(ret)
}

// LockMouseRect 锁定鼠标移动区域
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LockMouseRect(x1 int32, y1 int32, x2 int32, y2 int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LockMouseRect", x1, y1, x2, y2)
	}
	funAddr := DmHModule + 119792
	ret, _, _ := syscall.Syscall6(funAddr, 5, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), 0)
	return int32(ret)
}

// FindPicSimE FindPicSimE
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: pic_name - 图片名称(多个用|分隔)
// 参数: delta_color - 颜色偏差
// 参数: sim - 相似�?0.1-1.0)
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 结果字符�?
func (dm *DmSoft) FindPicSimE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("FindPicSimE", x1, y1, x2, y2, pic_name, delta_color, sim, dir)
	}
	funAddr := DmHModule + 123440
	pic_namePtr := stringToBytePtr(pic_name)
	delta_colorPtr := stringToBytePtr(delta_color)
	ret, _, _ := syscall.Syscall9(funAddr, 9, dm.obj, uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(pic_namePtr)), uintptr(unsafe.Pointer(delta_colorPtr)), uintptr(sim), uintptr(dir))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// EnumIniSectionPwd 枚举INI�?带密�?
// 参数: file - 文件路径
// 参数: pwd - 密码
// 返回: 结果字符�?
func (dm *DmSoft) EnumIniSectionPwd(file string, pwd string) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("EnumIniSectionPwd", file, pwd)
	}
	funAddr := DmHModule + 116992
	filePtr := stringToBytePtr(file)
	pwdPtr := stringToBytePtr(pwd)
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(unsafe.Pointer(filePtr)), uintptr(unsafe.Pointer(pwdPtr)))
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// RightUp 鼠标右键弹起
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) RightUp() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("RightUp")
	}
	funAddr := DmHModule + 111504
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FoobarTextPrintDir 设置Foobar文字打印方向
// 参数: hwnd - 窗口句柄
// 参数: dir - 查找方向(0:从左到右,1:从右到左,2:从上到下,3:从下到上)
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarTextPrintDir(hwnd int32, dir int32) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarTextPrintDir", hwnd, dir)
	}
	funAddr := DmHModule + 103072
	ret, _, _ := syscall.Syscall(funAddr, 3, dm.obj, uintptr(hwnd), uintptr(dir))
	return int32(ret)
}

// GetDir 获取特殊目录路径
// 参数: type_ - 类型
// 返回: 结果字符�?
func (dm *DmSoft) GetDir(type_ int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("GetDir", type_)
	}
	funAddr := DmHModule + 124512
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(type_), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// Hex32 32位整数转十六进制字符�?// 参数: v - �?// 返回: 结果字符�?
func (dm *DmSoft) Hex32(v int32) string {
	if dm.pipeConn != nil {
		return dm.comCallStr("Hex32", v)
	}
	funAddr := DmHModule + 110080
	ret, _, _ := syscall.Syscall(funAddr, 2, dm.obj, uintptr(v), 0)
	return bytePtrToString((*byte)(unsafe.Pointer(ret)))
}

// LeaveCri 离开临界�?// 返回: 成功返回1,失败返回0
func (dm *DmSoft) LeaveCri() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("LeaveCri")
	}
	funAddr := DmHModule + 120816
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// GetTime 获取系统时间
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) GetTime() int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("GetTime")
	}
	funAddr := DmHModule + 103504
	ret, _, _ := syscall.Syscall(funAddr, 1, dm.obj, 0, 0)
	return int32(ret)
}

// FoobarFillRect Foobar填充矩形
// 参数: hwnd - 窗口句柄
// 参数: x1 - 左上角X坐标
// 参数: y1 - 左上角Y坐标
// 参数: x2 - 右下角X坐标
// 参数: y2 - 右下角Y坐标
// 参数: color - 颜色(格式: "RRGGBB-DRDGDB")
// 返回: 成功返回1,失败返回0
func (dm *DmSoft) FoobarFillRect(hwnd int32, x1 int32, y1 int32, x2 int32, y2 int32, color string) int32 {
	if dm.pipeConn != nil {
		return dm.comCallInt32("FoobarFillRect", hwnd, x1, y1, x2, y2, color)
	}
	funAddr := DmHModule + 103136
	colorPtr := stringToBytePtr(color)
	ret, _, _ := syscall.Syscall9(funAddr, 7, dm.obj, uintptr(hwnd), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(colorPtr)), 0, 0)
	return int32(ret)
}
