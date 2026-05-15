// Package main 实现32位helper进程（dm_com_server.exe）
//
// 功能说明：
//   本程序是64位模式下的关键组件，运行在32位进程中，负责：
//   1. 加载大漠DLL（dm.dll/xd47243.dll等）
//   2. 可选：加载破解DLL（Go.dll）并执行破解
//   3. 创建大漠COM对象
//   4. 监听TCP端口，接收64位主进程的调用请求
//   5. 通过偏移量直接调用大漠API函数
//   6. 将结果通过gob序列化返回给主进程
//
// 工作流程：
//   1. 启动时加载命令行指定的DLL
//   2. 执行破解（如果提供了crack.dll）
//   3. 创建大漠对象（通过固定偏移量0x18000）
//   4. 监听随机TCP端口，将端口号输出到stdout
//   5. 等待64位主进程连接
//   6. 循环处理调用请求直到连接断开
//
// 编译要求：
//   必须使用32位编译：GOARCH=386 go build -o dm_com_server.exe ./cmd/dm_com_server/
//
// 使用方式：
//   由64位主进程自动启动，无需手动运行
package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// CallRequest 从64位主进程接收的调用请求
type CallRequest struct {
	Offset  uint32    // 方法在DLL中的偏移量（对应methodOffsets表）
	RetType uint8     // 返回值类型: 0=int32, 1=string, 2=float64, 3=int64
	NOut    uint8     // 输出参数个数
	Args    []CallArg // 参数列表
}

// CallArg 调用参数，支持多种数据类型
type CallArg struct {
	Type  uint8   // 参数类型: 0=int32, 1=string, 2=float64, 3=输出参数指针, 4=float32, 5=int64
	IVal  int32   // 整数值
	SVal  string  // 字符串值
	FVal  float64 // 浮点数值
	I64Val int64  // 64位整数值
}

// CallResponse 返回给64位主进程的响应
type CallResponse struct {
	RetType uint8    // 返回值类型
	IRet    int32    // int32类型返回值
	SRet    string   // string类型返回值
	FRet    float64  // float64类型返回值
	IRet64  int64    // int64类型返回值
	OutVals []int32  // 输出参数值列表
	Err     string   // 错误信息
}

// dmObj 大漠COM对象句柄，全局唯一
var dmObj uintptr

// keepAlive 用于保持字符串参数在内存中不被GC回收
var keepAlive [][]byte

// utf8ToGBK 将UTF-8字符串转换为GBK编码的字节切片
// 参数: s - UTF-8编码的字符串
// 返回值: GBK编码的字节切片和错误信息
// 说明: 大漠DLL使用GBK编码，需要从UTF-8转换
func utf8ToGBK(s string) ([]byte, error) {
	return simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
}

// gbkToUTF8 将GBK编码的字节切片转换为UTF-8字符串
// 参数: b - GBK编码的字节切片
// 返回值: UTF-8编码的字节切片和错误信息
// 说明: 用于将大漠DLL返回的GBK字符串转换为Go的UTF-8格式
func gbkToUTF8(b []byte) ([]byte, error) {
	return simplifiedchinese.GBK.NewDecoder().Bytes(b)
}

// readCString 读取C风格字符串（以\0结尾）并转换为Go字符串
// 参数: ptr - 字符串指针地址
// 返回值: Go字符串（UTF-8编码）
//
// 说明:
//   1. 从指定地址开始读取字节，直到遇到\0结束符
//   2. 将GBK编码的字节转换为UTF-8
//   3. 如果转换失败，返回原始字节的字符串表示
func readCString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	var n int
	for *(*byte)(unsafe.Pointer(ptr + uintptr(n))) != 0 {
		n++
	}
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
	}
	out, err := gbkToUTF8(buf)
	if err != nil {
		return string(buf)
	}
	return string(out)
}

// dmModule 大漠DLL模块句柄
var dmModule syscall.Handle

// handleCall 处理来自64位主进程的调用请求
// 参数: req - 调用请求，包含方法偏移量、参数列表等
// 返回值: CallResponse 调用结果或错误信息
//
// 详细流程：
//   1. 根据偏移量计算函数地址：fnAddr = dmModule + offset
//   2. 解析参数列表，将不同类型参数转换为uintptr数组
//   3. 字符串参数需要UTF-8转GBK，并保持引用防止GC回收
//   4. 根据参数数量选择合适的syscall函数（Syscall/Syscall6/Syscall9等）
//   5. 调用DLL函数并获取返回值
//   6. 根据返回值类型提取相应的返回数据
//   7. 收集输出参数的值
func handleCall(req CallRequest) CallResponse {
	resp := CallResponse{
		RetType: req.RetType,
	}

	fnAddr := uintptr(dmModule) + uintptr(req.Offset)

	keepAlive = nil

	outCount := 0
	for _, arg := range req.Args {
		if arg.Type == 3 {
			outCount++
		}
	}
	outVars := make([]int32, outCount)
	outIdx := 0

	const maxSlots = 16
	args := make([]uintptr, maxSlots)
	args[0] = dmObj
	idx := 1

	for _, arg := range req.Args {
		switch arg.Type {
		case 0:
			args[idx] = uintptr(arg.IVal)
			idx++
		case 1:
			gbk, err := utf8ToGBK(arg.SVal)
			if err != nil {
				resp.Err = fmt.Sprintf("utf8ToGBK: %v", err)
				return resp
			}
			cstr := make([]byte, len(gbk)+1)
			copy(cstr, gbk)
			keepAlive = append(keepAlive, cstr)
			args[idx] = uintptr(unsafe.Pointer(&cstr[0]))
			idx++
		case 2:
			bits := math.Float64bits(arg.FVal)
			args[idx] = uintptr(bits)
			args[idx+1] = uintptr(bits >> 32)
			idx += 2
		case 3:
			args[idx] = uintptr(unsafe.Pointer(&outVars[outIdx]))
			outIdx++
			idx++
		case 4:
			bits := math.Float32bits(float32(arg.FVal))
			args[idx] = uintptr(bits)
			idx++
		case 5:
			args[idx] = uintptr(arg.I64Val)
			args[idx+1] = uintptr(arg.I64Val >> 32)
			idx += 2
		default:
			args[idx] = 0
			idx++
		}
	}

	nargs := idx

	var r1, r2 uintptr
	switch {
	case nargs <= 3:
		r1, r2, _ = syscall.Syscall(fnAddr, uintptr(nargs), args[0], args[1], args[2])
	case nargs <= 6:
		r1, r2, _ = syscall.Syscall6(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5])
	case nargs <= 9:
		r1, r2, _ = syscall.Syscall9(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	case nargs <= 12:
		r1, r2, _ = syscall.Syscall12(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11])
	case nargs <= 15:
		r1, r2, _ = syscall.Syscall15(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14])
	default:
		resp.Err = fmt.Sprintf("too many args: %d", nargs)
		return resp
	}

	switch req.RetType {
	case 0:
		resp.IRet = int32(r1)
	case 1:
		resp.SRet = readCString(r1)
	case 2:
		resp.FRet = math.Float64frombits(uint64(r1))
	case 3:
		resp.IRet64 = int64(r2)<<32 | int64(r1)
	}

	resp.OutVals = make([]int32, outCount)
	copy(resp.OutVals, outVars)

	return resp
}

// main helper进程入口函数
//
// 执行流程：
//   1. 锁定OS线程（COM要求）
//   2. 解析命令行参数（dm.dll路径和可选的crack.dll路径）
//   3. 加载大漠DLL
//   4. 如果提供了crack.dll，加载并执行破解
//   5. 通过偏移量0x18000创建大漠对象
//   6. 监听随机TCP端口
//   7. 输出"READY <port>"到stdout，通知主进程
//   8. 启动goroutine监听stdin，用于接收退出信号
//   9. 等待主进程连接
//   10. 循环处理调用请求直到连接断开或stdin关闭
func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: dm_com_server <dm.dll> [crack.dll]\n")
		os.Exit(1)
	}
	dmPath := os.Args[1]
	var crackPath string
	if len(os.Args) >= 3 {
		crackPath = os.Args[2]
	}

	module, err := syscall.LoadLibrary(dmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadLibrary dm.dll failed: %v\n", err)
		os.Exit(1)
	}
	dmModule = module
	defer syscall.FreeLibrary(dmModule)

	if crackPath != "" {
		crackModule, err := syscall.LoadLibrary(crackPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadLibrary crack failed: %v\n", err)
			os.Exit(1)
		}
		defer syscall.FreeLibrary(crackModule)

		goProc, err := syscall.GetProcAddress(crackModule, "Go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "GetProcAddress Go failed: %v\n", err)
			os.Exit(1)
		}
		syscall.Syscall(uintptr(goProc), 1, uintptr(dmModule), 0, 0)
	}

	createObjAddr := uintptr(dmModule) + 98304
	dmObj, _, _ = syscall.Syscall(createObjAddr, 0, 0, 0, 0)
	if dmObj == 0 {
		fmt.Fprintf(os.Stderr, "CreateDmObject(offset 0x18000) failed\n")
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Listen failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("READY %d\n", ln.Addr().(*net.TCPAddr).Port)
	os.Stdout.Sync()

	go func() {
		var buf [1]byte
		for {
			_, err := os.Stdin.Read(buf[:])
			if err != nil {
				os.Exit(0)
			}
		}
	}()

	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	for {
		var req CallRequest
		if err := dec.Decode(&req); err != nil {
			return
		}
		resp := handleCall(req)
		if err := enc.Encode(resp); err != nil {
			return
		}
	}
}
