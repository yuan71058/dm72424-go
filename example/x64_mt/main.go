// Package main 大漠插件64位多线程示例程序
//
// 本示例展示如何在64位环境下使用大漠插件进行多线程并发操作
// 每个线程拥有独立的DmSoft实例和32位helper进程，实现真正的并行执行
//
// 功能演示：
//   - 多线程并发创建和初始化大漠对象
//   - 每个线程独立绑定不同的窗口
//   - 多线程同时向不同窗口写入文字
//   - 展示线程间完全隔离，互不干扰
//
// 应用场景：
//   - 同时控制多个游戏窗口（多开）
//   - 并行处理多个任务
//   - 提高自动化脚本效率
//
// 编译运行：
//   1. 先编译helper进程：GOARCH=386 go build -o dm_com_server.exe ../../cmd/dm_com_server/
//   2. 编译本程序：set GOARCH=amd64 && go build -o x64_mt_example.exe .
//   3. 运行：x64_mt_example.exe
//
// 核心原理：
//   - 64位模式下每个DmSoft.Init()会启动独立的32位helper进程
//   - 多个goroutine各自持有独立的DmSoft实例
//   - 每个实例通过独立的TCP连接与自己的helper通信
//   - 线程间天然隔离，无需额外的同步机制
//
// 注意事项：
//   - 每个线程会占用约5-10MB内存（独立helper进程）
//   - 建议线程数量不超过CPU核心数的2倍
//   - 必须正确Release()释放资源，避免helper进程残留
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	dmsoft "github.com/yuan71058/dm72424-go"
)

// 常量定义
const (
	DmPluginPath = "xd47243.dll" // 大漠插件DLL文件名
	CrackDllPath = "Go.dll"        // 破解DLL文件名
	WorkerCount  = 3               // 工作线程数量（同时操作的窗口数）
	TextCount    = 500             // 每个线程写入的字符数量
)

// TextWorker 文本工作线程结构体
// 每个TextWorker实例代表一个独立的工作线程，拥有自己的大漠对象和helper进程
// 用于向指定的Edit控件写入文字，展示多线程并发操作能力
type TextWorker struct {
	ID       int             // 线程ID（用于标识和日志输出）
	Dm       *dmsoft.DmSoft  // 大漠插件实例（每个线程独立）
	MainHwnd int32           // 主窗口句柄（记事本窗口）
	EditHwnd int32           // Edit编辑框控件句柄
	Content  string          // 要写入的文本内容
	Result   chan string     // 结果通道（用于输出状态信息到主线程）
}

// NewTextWorker 创建新的文本工作线程
// 参数:
//   - id: 线程ID
//   - mainHwnd: 主窗口句柄
//   - editHwnd: Edit控件句柄
//   - content: 要写入的文本
//   - resultChan: 结果通道（用于输出状态信息）
// 返回值: *TextWorker 工作线程实例指针
func NewTextWorker(id int, mainHwnd, editHwnd int32, content string, resultChan chan string) *TextWorker {
	return &TextWorker{
		ID:       id,
		MainHwnd: mainHwnd,
		EditHwnd: editHwnd,
		Content:  content,
		Result:   resultChan,
	}
}

// Init 初始化工作线程：创建大漠对象、启动helper进程、注册插件
// 返回值: bool 初始化是否成功
// 说明: 每个线程调用此方法时会启动独立的32位helper进程
func (w *TextWorker) Init() bool {
	w.Dm = dmsoft.New()
	if w.Dm == nil {
		w.Result <- fmt.Sprintf("[线程%d] 创建子对象失败", w.ID)
		return false
	}
	err := w.Dm.Init()
	if err != nil {
		w.Result <- fmt.Sprintf("[线程%d] 子对象初始化失败: %v", w.ID, err)
		return false
	}

	nret := w.Dm.Reg("", "")
	if nret != 1 {
		w.Result <- fmt.Sprintf("[线程%d] 子对象注册失败: %d", w.ID, nret)
		w.Dm.Release()
		return false
	}

	w.Result <- fmt.Sprintf("[线程%d] 子对象初始化完成(helper进程已启动, 已注册)", w.ID)
	return true
}

// BindWindow 绑定窗口到当前线程的大漠对象
// 返回值: bool 绑定是否成功
// 说明: 使用"gdi"+"windows"+"windows"模式绑定Edit控件
func (w *TextWorker) BindWindow() bool {
	w.Result <- fmt.Sprintf("[线程%d] 准备绑定Edit控件: %d", w.ID, w.EditHwnd)

	ret := w.Dm.BindWindow(w.EditHwnd, "gdi", "windows", "windows", 0)
	w.Result <- fmt.Sprintf("[线程%d] BindWindow返回值: %d", w.ID, ret)

	if ret != 1 {
		lastError := w.Dm.GetLastError()
		w.Result <- fmt.Sprintf("[线程%d] 绑定窗口失败: %d, 错误码: %d", w.ID, ret, lastError)
		return false
	}

	w.Result <- fmt.Sprintf("[线程%d] 绑定Edit控件成功: %d", w.ID, w.EditHwnd)
	return true
}

// WriteText 向绑定的Edit控件写入文本内容
// 使用SendString2方法发送文字，模拟键盘输入
func (w *TextWorker) WriteText() {
	w.Result <- fmt.Sprintf("[线程%d] 开始写入文字...", w.ID)

	ret := w.Dm.SendString2(w.EditHwnd, w.Content)
	if ret != 1 {
		lastError := w.Dm.GetLastError()
		w.Result <- fmt.Sprintf("[线程%d] SendString2失败: %d, 错误码: %d", w.ID, ret, lastError)
		return
	}

	time.Sleep(500 * time.Millisecond)

	w.Result <- fmt.Sprintf("[线程%d] 写入完成，共%d个字符", w.ID, len(w.Content))
}

// UnbindWindow 解绑窗口，释放绑定资源
func (w *TextWorker) UnbindWindow() {
	if w.EditHwnd != 0 {
		w.Dm.UnBindWindow()
		w.Result <- fmt.Sprintf("[线程%d] 解绑窗口", w.ID)
	}
}

// Release 释放大漠对象和helper进程资源
// 说明: 调用Release后会关闭TCP连接并终止该线程的32位helper进程
func (w *TextWorker) Release() {
	if w.Dm != nil {
		w.Dm.Release()
		w.Result <- fmt.Sprintf("[线程%d] 释放子对象(helper进程已关闭)", w.ID)
	}
}

// Run 工作线程的主执行函数
// 参数: wg - WaitGroup用于等待所有工作线程完成
// 执行流程: 绑定窗口 -> 写入文字 -> 解绑 -> 释放资源
func (w *TextWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	if !w.BindWindow() {
		w.Release()
		return
	}

	w.WriteText()

	w.UnbindWindow()

	w.Release()

	w.Result <- fmt.Sprintf("[线程%d] 完成", w.ID)
}

// main 主函数：64位多线程示例入口
//
// 执行流程：
//   1. 检查编译架构（必须是amd64）
//   2. LoadDm + CrackDm - 初始化DLL路径
//   3. 创建主对象并注册（用于窗口枚举等全局操作）
//   4. 创建多个记事本窗口（模拟多开场景）
//   5. 枚举并定位每个记事本的Edit控件
//   6. 为每个窗口创建独立的TextWorker（每个启动自己的helper进程）
//   7. 并发执行所有工作线程
//   8. 等待完成，统计耗时
//   9. 清理资源（关闭窗口、删除临时文件）
//
// 核心演示：
//   - 展示64位模式下真正的多线程并行能力
//   - 每个线程完全独立，互不干扰
//   - 通过控制台输出实时显示各线程状态
func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║     大漠插件 64位 多线程 TCP+gob 跨进程调用示例        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if runtime.GOARCH != "amd64" {
		log.Fatalf("本示例仅支持 64 位编译，当前架构: %s\n请使用: set GOARCH=amd64 && go build", runtime.GOARCH)
	}
	fmt.Println("✓ 编译架构: amd64")
	fmt.Println("✓ 调用模式: TCP+gob 跨进程 (32位helper + 偏移量调用)")
	fmt.Println()
	fmt.Println("  多线程原理：")
	fmt.Println("    64位模式下，每个DmSoft对象启动独立的32位helper进程")
	fmt.Println("    多个goroutine各自持有独立DmSoft对象，互不干扰，真正并行")
	fmt.Println("    流程：LoadDm → CrackDm → 主对象注册 → 创建记事本 →")
	fmt.Println("          枚举窗口 → 创建子对象 → 绑定 → 写入 → 解绑 → 释放")
	fmt.Println()

	absPluginPath, _ := filepath.Abs(DmPluginPath)
	absCrackPath, _ := filepath.Abs(CrackDllPath)

	if _, err := os.Stat(absPluginPath); os.IsNotExist(err) {
		log.Fatalf("大漠插件不存在: %s", absPluginPath)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第一步：LoadDm + CrackDm")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	_, err := dmsoft.LoadDm(absPluginPath)
	if err != nil {
		log.Fatalf("LoadDm 失败: %v", err)
	}
	fmt.Println("  ✓ LoadDm 成功")

	if _, err := os.Stat(absCrackPath); os.IsNotExist(err) {
		fmt.Printf("  ⊘ 破解DLL不存在，跳过: %s\n", absCrackPath)
	} else {
		err := dmsoft.CrackDm(absCrackPath)
		if err != nil {
			log.Fatalf("CrackDm 失败: %v", err)
		}
		fmt.Println("  ✓ CrackDm 成功")
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第二步：创建主对象并注册")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mainDm := dmsoft.New()
	if mainDm == nil {
		log.Fatal("  ✗ New() 返回 nil")
	}

	err = mainDm.Init()
	if err != nil {
		log.Fatalf("  ✗ Init() 失败: %v", err)
	}
	fmt.Println("  ✓ 主对象 Init() 成功")

	nret := mainDm.Reg("", "")
	if nret == 1 {
		fmt.Println("  ✓ 注册成功")
	} else {
		log.Fatalf("  ✗ 注册失败 (返回值: %d)", nret)
	}

	version := mainDm.Ver()
	fmt.Printf("  ✓ 大漠插件版本: %s\n", version)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第三步：创建记事本窗口")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	textFiles := make([]string, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		textFiles[i] = filepath.Join(os.TempDir(), fmt.Sprintf("dm64_test_%d.txt", i+1))
		file, err := os.Create(textFiles[i])
		if err != nil {
			log.Fatalf("创建文件失败: %v", err)
		}
		file.Close()
		fmt.Printf("  创建文件: %s\n", textFiles[i])
	}

	processes := make([]*os.Process, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		cmd := exec.Command("notepad.exe", textFiles[i])
		cmd.Start()
		processes[i] = cmd.Process
		fmt.Printf("  打开记事本: PID=%d\n", cmd.Process.Pid)
	}

	fmt.Println("\n  等待窗口打开...")
	time.Sleep(2 * time.Second)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第四步：枚举记事本窗口")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	hwndList := mainDm.EnumWindow(0, "", "Notepad", 2)
	fmt.Printf("  枚举结果(按类名): %s\n", hwndList)

	hwnds := parseHwndList(hwndList)

	if len(hwnds) < WorkerCount {
		hwndList = mainDm.EnumWindow(0, "记事本", "", 1)
		fmt.Printf("  枚举结果(按标题): %s\n", hwndList)
		hwnds = parseHwndList(hwndList)
	}

	fmt.Printf("  找到 %d 个记事本窗口\n", len(hwnds))

	if len(hwnds) < WorkerCount {
		log.Fatalf("记事本窗口数量不足，需要 %d 个，找到 %d 个", WorkerCount, len(hwnds))
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第五步：排列窗口")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	windowWidth := 800
	windowHeight := 600
	margin := 20

	for i := 0; i < WorkerCount; i++ {
		row := i / 2
		col := i % 2
		x := int32(col * (windowWidth + margin))
		y := int32(row * (windowHeight + margin))

		ret := mainDm.MoveWindow(hwnds[i], x, y)
		if ret == 1 {
			fmt.Printf("  窗口%d 移动到 (%d, %d)\n", i+1, x, y)
		} else {
			fmt.Printf("  窗口%d 移动失败\n", i+1)
		}
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第六步：查找Edit编辑框控件")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	editHwnds := make([]int32, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		editHwnd := mainDm.FindWindowEx(hwnds[i], "Edit", "")
		if editHwnd == 0 {
			log.Fatalf("找不到记事本 %d 的Edit编辑框", i+1)
		}
		editHwnds[i] = editHwnd
		fmt.Printf("  记事本 %d: 主窗口=%d, Edit控件=%d\n", i+1, hwnds[i], editHwnd)
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第七步：创建大漠子对象（每个线程独立helper进程）")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  说明：64位模式下，每个DmSoft.Init()会启动独立的")
	fmt.Println("       32位helper进程(dm_com_server.exe)，互不干扰")
	fmt.Println()

	var wg sync.WaitGroup
	resultChan := make(chan string, WorkerCount*20)

	contents := make([]string, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		contents[i] = generateTextForThread(i+1, TextCount)
		fmt.Printf("  线程%d 文本长度: %d\n", i+1, len(contents[i]))
	}

	workers := make([]*TextWorker, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = NewTextWorker(i+1, hwnds[i], editHwnds[i], contents[i], resultChan)
		if !workers[i].Init() {
			log.Fatalf("线程%d初始化失败", i+1)
		}
		fmt.Printf("  ✓ 线程%d 子对象创建完成，准备绑定Edit控件: %d\n", i+1, editHwnds[i])
	}

	go func() {
		for result := range resultChan {
			fmt.Printf("    %s\n", result)
		}
	}()

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第八步：多线程并发写入")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	startTime := time.Now()

	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go workers[i].Run(&wg)
	}

	wg.Wait()
	close(resultChan)

	elapsed := time.Since(startTime)
	fmt.Println()
	fmt.Printf("  ✓ 所有线程完成，总耗时: %v\n", elapsed)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第九步：释放主对象")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mainDm.Release()
	fmt.Println("  ✓ 主对象已释放")
	fmt.Println()

	fmt.Println("  等待5秒，查看记事本效果...")
	time.Sleep(5 * time.Second)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第十步：关闭窗口并清理")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, proc := range processes {
		proc.Kill()
		fmt.Printf("  关闭窗口: %s\n", textFiles[i])
	}

	for i := 0; i < WorkerCount; i++ {
		os.Remove(textFiles[i])
		fmt.Printf("  删除文件: %s\n", textFiles[i])
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  测试完成")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  64位多线程调用时序：")
	fmt.Println("    LoadDm → CrackDm → 主对象New+Init+Reg")
	fmt.Println("    → 创建记事本 → 枚举窗口 → 查找Edit控件")
	fmt.Println("    → 每线程: 子对象New+Init+Reg → BindWindow → SendString2 → UnBind → Release")
	fmt.Println()
	fmt.Println("  关键点：")
	fmt.Println("    1. 每个DmSoft对象有独立的helper进程，互不干扰")
	fmt.Println("    2. 多个goroutine可真正并行操作不同窗口")
	fmt.Println("    3. TCP+gob自动处理GBK↔UTF-8编码和参数序列化")
	fmt.Println()
}

// parseHwndList 解析窗口句柄列表字符串
// 参数: hwndList - 逗号分隔的窗口句柄字符串（如"123,456,789"）
// 返回值: []int32 窗口句柄切片
func parseHwndList(hwndList string) []int32 {
	if hwndList == "" {
		return nil
	}

	parts := strings.Split(hwndList, ",")
	hwnds := make([]int32, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		hwnd, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			continue
		}
		hwnds = append(hwnds, int32(hwnd))
	}

	return hwnds
}

// generateTextForThread 为指定线程生成测试文本
// 参数:
//   - threadID: 线程ID（用于标识文本来源）
//   - count: 生成字符数量
// 返回值: string 格式化的文本字符串（每10个字符换行）
func generateTextForThread(threadID int, count int) string {
	text := ""
	for i := 0; i < count; i++ {
		text += fmt.Sprintf("%d ", threadID)
		if (i+1)%10 == 0 {
			text += "\n"
		}
	}
	return text
}
