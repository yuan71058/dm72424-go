// example/multithread/main.go - 多线程写入文本示例
//
// 版本: v1.4.0
// 更新日期: 2026-03-21
//
// 功能说明:
//   1. 首先注册破解大漠，创建大漠主对象
//   2. 再创建记事本
//   3. 枚举所有记事本
//   4. 查找每个记事本编辑框控件句柄
//   5. 再创建大漠子对象，绑定编辑框控件后台
//   6. 输入文本
//   7. 解绑释放资源
//
// 编译说明:
//   必须使用32位编译: GOARCH=386 go build -o multithread_test.exe
//
// 运行要求:
//   - 需要 xd47243.dll 和 Go.dll 在同一目录
//   - 需要安装 Notepad 或其他文本编辑器

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	dmsoft "github.com/yuan71058/dm72424-go"
)

const (
	DmPluginPath = "xd47243.dll"
	CrackDllPath = "Go.dll"
	WorkerCount  = 3
	TextCount    = 500
)

// TextWorker 文本写入工作线程
type TextWorker struct {
	ID       int
	Dm       *dmsoft.DmSoft
	MainHwnd int32
	EditHwnd int32
	Content  string
	Result   chan string
}

// NewTextWorker 创建新的工作线程
func NewTextWorker(id int, mainHwnd, editHwnd int32, content string, resultChan chan string) *TextWorker {
	return &TextWorker{
		ID:       id,
		MainHwnd: mainHwnd,
		EditHwnd: editHwnd,
		Content:  content,
		Result:   resultChan,
	}
}

// Init 初始化大漠子对象
func (w *TextWorker) Init() bool {
	w.Dm = dmsoft.New()
	if w.Dm == nil {
		w.Result <- fmt.Sprintf("[线程%d] 创建子对象失败", w.ID)
		return false
	}
	w.Dm.Init()
	w.Result <- fmt.Sprintf("[线程%d] 子对象初始化完成，地址: %p", w.ID, w.Dm)
	return true
}

// BindWindow 绑定窗口
func (w *TextWorker) BindWindow() bool {
	w.Result <- fmt.Sprintf("[线程%d] 准备绑定Edit控件: %d", w.ID, w.EditHwnd)

	// GDI模式绑定Edit编辑框
	fmt.Printf("地址: %p, EditHwnd: %d\n", w.Dm, w.EditHwnd)
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

// WriteText 写入文字
func (w *TextWorker) WriteText() {
	w.Result <- fmt.Sprintf("[线程%d] 开始写入文字...", w.ID)

	// 使用SendString2发送文字
	ret := w.Dm.SendString2(w.EditHwnd, w.Content)
	if ret != 1 {
		lastError := w.Dm.GetLastError()
		w.Result <- fmt.Sprintf("[线程%d] SendString2失败: %d, 错误码: %d", w.ID, ret, lastError)
		return
	}

	// 延时，让文字显示更清晰
	time.Sleep(500 * time.Millisecond)

	w.Result <- fmt.Sprintf("[线程%d] 写入完成，共%d个字符", w.ID, len(w.Content))
}

// UnbindWindow 解绑窗口
func (w *TextWorker) UnbindWindow() {
	if w.EditHwnd != 0 {
		w.Dm.UnBindWindow()
		w.Result <- fmt.Sprintf("[线程%d] 解绑窗口", w.ID)
	}
}

// Release 释放资源
func (w *TextWorker) Release() {
	if w.Dm != nil {
		w.Dm.Release()
	}
}

// Run 运行工作线程
func (w *TextWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	w.Result <- fmt.Sprintf("[线程%d] 子对象地址: %p", w.ID, w.Dm)

	// 绑定窗口
	if !w.BindWindow() {
		return
	}

	// 写入文字
	w.WriteText()

	// 解绑窗口
	w.UnbindWindow()

	// 释放资源
	w.Release()

	w.Result <- fmt.Sprintf("[线程%d] 完成", w.ID)
}

func main() {
	fmt.Println("========== 多线程写入文本示例 ==========")
	fmt.Println()

	// ========== 第一步：加载大漠插件DLL ==========
	fmt.Println("========== 第一步：加载大漠插件 ==========")
	dmHModule, err := dmsoft.LoadDm(DmPluginPath)
	if err != nil {
		log.Fatalf("加载大漠插件失败: %v", err)
	}
	fmt.Printf("大漠插件加载成功，模块句柄: %v\n", dmHModule)

	// ========== 第二步：破解大漠插件 ==========
	err = dmsoft.CrackDm(CrackDllPath)
	if err != nil {
		log.Fatalf("破解大漠插件失败: %v", err)
	}
	fmt.Println("破解函数执行完成")

	// ========== 第三步：创建主对象并注册 ==========
	fmt.Println("\n========== 第二步：创建主对象并注册 ==========")
	mainDm := dmsoft.New()
	if mainDm == nil {
		log.Fatal("创建主对象失败")
	}
	mainDm.Init()

	nret := mainDm.Reg("", "")
	if nret == 1 {
		fmt.Println("大漠注册成功")
	} else {
		log.Fatal("大漠注册失败")
	}

	version := mainDm.Ver()
	fmt.Printf("大漠插件版本: %s\n", version)

	// ========== 第四步：创建3个空文本文件 ==========
	fmt.Println("\n========== 第三步：创建文本文件 ==========")
	textFiles := make([]string, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		textFiles[i] = filepath.Join(os.TempDir(), fmt.Sprintf("test_%d.txt", i+1))
		file, err := os.Create(textFiles[i])
		if err != nil {
			log.Fatalf("创建文件失败: %v", err)
		}
		file.Close()
		fmt.Printf("创建文件: %s\n", textFiles[i])
	}

	// ========== 第五步：用记事本打开3个文本文件 ==========
	fmt.Println("\n========== 第四步：打开文本窗口 ==========")
	processes := make([]*os.Process, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		cmd := exec.Command("notepad.exe", textFiles[i])
		cmd.Start()
		processes[i] = cmd.Process
		fmt.Printf("打开窗口: %s (PID: %d)\n", textFiles[i], cmd.Process.Pid)
	}

	// 等待窗口打开
	fmt.Println("\n等待窗口打开...")
	time.Sleep(2 * time.Second)
	fmt.Println()

	// ========== 第六步：枚举所有记事本窗口 ==========
	fmt.Println("========== 第六步：枚举所有记事本窗口 ==========")

	// 先尝试按类名枚举
	hwndList := mainDm.EnumWindow(0, "", "Notepad", 2)
	fmt.Printf("枚举结果(按类名): %s\n", hwndList)

	// 解析窗口句柄列表
	hwnds := parseHwndList(hwndList)

	// 如果按类名枚举失败，尝试按标题枚举
	if len(hwnds) < WorkerCount {
		hwndList = mainDm.EnumWindow(0, "记事本", "", 1)
		fmt.Printf("枚举结果(按标题): %s\n", hwndList)
		hwnds = parseHwndList(hwndList)
	}

	fmt.Printf("找到 %d 个记事本窗口\n", len(hwnds))

	if len(hwnds) < WorkerCount {
		log.Fatalf("记事本窗口数量不足，需要 %d 个，找到 %d 个", WorkerCount, len(hwnds))
	}

	// ========== 第七步：智能排列窗口 ==========
	fmt.Println("\n========== 第七步：智能排列窗口 ==========")

	// 计算窗口位置
	windowWidth := 800
	windowHeight := 600
	margin := 20

	for i := 0; i < WorkerCount; i++ {
		// 计算窗口位置（网格排列）
		row := i / 2
		col := i % 2
		x := int32(col * (windowWidth + margin))
		y := int32(row * (windowHeight + margin))

		// 移动窗口
		ret := mainDm.MoveWindow(hwnds[i], x, y)
		if ret == 1 {
			fmt.Printf("窗口%d 移动到 (%d, %d)\n", i+1, x, y)
		} else {
			fmt.Printf("窗口%d 移动失败\n", i+1)
		}
	}

	// 等待窗口移动完成
	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	// ========== 第七步：查找每个记事本的Edit编辑框控件 ==========
	fmt.Println("\n========== 第六步：查找Edit编辑框控件 ==========")
	editHwnds := make([]int32, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		editHwnd := mainDm.FindWindowEx(hwnds[i], "Edit", "")
		if editHwnd == 0 {
			log.Fatalf("找不到记事本 %d 的Edit编辑框", i+1)
		}
		editHwnds[i] = editHwnd
		fmt.Printf("记事本 %d: 主窗口=%d, Edit控件=%d\n", i+1, hwnds[i], editHwnd)
	}

	// ========== 第八步：创建大漠子对象并绑定编辑框控件 ==========
	fmt.Println("\n========== 第七步：创建大漠子对象并绑定编辑框控件 ==========")
	var wg sync.WaitGroup
	resultChan := make(chan string, WorkerCount*20)

	// 准备不同的文字内容
	contents := make([]string, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		contents[i] = generateTextForThread(i+1, TextCount)
		fmt.Printf("线程%d 文本长度: %d\n", i+1, len(contents[i]))
	}

	// 创建工作线程（大漠子对象）
	workers := make([]*TextWorker, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = NewTextWorker(i+1, hwnds[i], editHwnds[i], contents[i], resultChan)
		if !workers[i].Init() {
			log.Fatalf("线程%d初始化失败", i+1)
		}
		fmt.Printf("创建大漠子对象 %d, 准备绑定Edit控件: %d\n", i+1, editHwnds[i])
	}

	// 启动结果收集器
	go func() {
		for result := range resultChan {
			fmt.Println(result)
		}
	}()

	// ========== 第九步：并发执行写入 ==========
	fmt.Println("\n========== 第八步：开始多线程写入 ==========")
	startTime := time.Now()

	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go workers[i].Run(&wg)
	}

	// 等待所有线程完成
	wg.Wait()
	close(resultChan)

	elapsed := time.Since(startTime)
	fmt.Println()
	fmt.Printf("========== 所有线程完成 ==========\n")
	fmt.Printf("总耗时: %v\n", elapsed)
	fmt.Println()

	// 延时5秒，让用户看到效果
	fmt.Println("等待5秒，查看效果...")
	time.Sleep(5 * time.Second)

	// ========== 第十步：关闭记事本窗口 ==========
	fmt.Println("========== 第九步：关闭窗口 ==========")
	for i, proc := range processes {
		proc.Kill()
		fmt.Printf("关闭窗口: %s\n", textFiles[i])
	}

	// ========== 第十一步：清理临时文件 ==========
	fmt.Println("\n========== 第十步：清理临时文件 ==========")
	for i := 0; i < WorkerCount; i++ {
		os.Remove(textFiles[i])
		fmt.Printf("删除文件: %s\n", textFiles[i])
	}

	// ========== 第十二步：释放大漠插件 ==========
	fmt.Println("\n========== 第十一步：释放资源 ==========")

	// 先释放主对象
	mainDm.Release()
	fmt.Println("主对象已释放")

	// 再释放大漠插件
	dmsoft.Free()
	fmt.Println("大漠插件已释放")
	fmt.Println("\n========== 示例完成 ==========")
}

// parseHwndList 解析窗口句柄列表字符串
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

// generateText 生成指定长度的文字
func generateText(prefix string, count int) string {
	text := ""
	for i := 0; i < count; i++ {
		text += fmt.Sprintf("%s%d ", prefix, i+1)
		if (i+1)%10 == 0 {
			text += "\n"
		}
	}
	return text
}

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
