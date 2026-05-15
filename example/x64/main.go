// Package main 大漠插件64位单线程示例程序
//
// 本示例展示如何在64位环境下使用大漠插件Go绑定库
// 通过TCP+gob与32位helper进程通信，实现跨架构调用大漠DLL
//
// 功能演示：
//   - 基础信息获取（版本号、屏幕分辨率等）
//   - 窗口操作（查找、枚举窗口）
//   - 窗口绑定与截图（BMP和PNG格式）
//   - 取色功能（RGB/BGR/HSV）
//   - 鼠标键盘模拟操作
//   - 找图找色（含输出参数）
//   - OCR文字识别
//   - 剪贴板操作
//   - 内存读写操作
//
// 编译运行：
//   1. 先编译helper进程：GOARCH=386 go build -o dm_com_server.exe ../../cmd/dm_com_server/
//   2. 编译本程序：set GOARCH=amd64 && go build -o x64_example.exe .
//   3. 运行：x64_example.exe
//
// 注意事项：
//   - 必须以64位编译（amd64）
//   - 需要dm_com_server.exe在同级目录或../cmd/dm_com_server/目录
//   - 需要xd47243.dll（大漠DLL）和Go.dll（破解DLL）在同级目录
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	dmsoft "github.com/yuan71058/dm72424-go"
)

// 常量定义：DLL文件路径
const (
	DmPluginPath = "xd47243.dll" // 大漠插件DLL文件名
	CrackDllPath = "Go.dll"        // 破解DLL文件名（用于激活大漠插件）
)

// main 主函数：64位大漠插件使用示例入口
//
// 执行流程：
//   1. 检查编译架构（必须是amd64）
//   2. LoadDm - 加载大漠DLL（记录路径，定位helper）
//   3. CrackDm - 设置破解DLL路径
//   4. New + Init - 创建实例并启动helper进程，建立TCP连接
//   5. Reg - 注册大漠插件
//   6. 执行各项功能测试
//   7. Release - 释放资源（关闭helper进程）
//
// 控制台输出包含详细的步骤说明和原理讲解，便于理解64位调用机制
func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         大漠插件 64位 TCP+gob 跨进程调用示例            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if runtime.GOARCH != "amd64" {
		log.Fatalf("本示例仅支持 64 位编译，当前架构: %s\n请使用: set GOARCH=amd64 && go build", runtime.GOARCH)
	}
	fmt.Println("✓ 编译架构: amd64")
	fmt.Println("✓ 调用模式: TCP+gob 跨进程 (32位helper + 偏移量调用)")
	fmt.Println()

	absPluginPath, _ := filepath.Abs(DmPluginPath)
	absCrackPath, _ := filepath.Abs(CrackDllPath)

	if _, err := os.Stat(absPluginPath); os.IsNotExist(err) {
		log.Fatalf("大漠插件不存在: %s", absPluginPath)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第一步：LoadDm — 记录DLL路径 + 定位helper")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  原理：")
	fmt.Println("    64位进程无法直接加载32位DLL")
	fmt.Println("    通过启动32位helper进程加载dm.dll，TCP转发调用")
	fmt.Println()

	dmHModule, err := dmsoft.LoadDm(absPluginPath)
	if err != nil {
		if strings.Contains(err.Error(), "管理员权限") {
			fmt.Println("  ✗ 注册COM Surrogate失败：需要管理员权限！")
			fmt.Println("  请右键以管理员身份运行此程序")
		}
		log.Fatalf("LoadDm 失败: %v", err)
	}
	fmt.Printf("  ✓ LoadDm 成功 (返回值: %v)\n", dmHModule)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第二步：CrackDm — 设置破解DLL路径（必须在Init之前）")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  原理：")
	fmt.Println("    64位模式下，CrackDm仅记录破解DLL路径")
	fmt.Println("    实际破解在Init()启动的32位helper进程中执行")
	fmt.Println("    helper进程：加载dm.dll → 加载Go.dll调用Go() → 创建dm对象 → 监听TCP")
	fmt.Println()
	fmt.Println("  ⚠ 必须在Init()之前调用")
	fmt.Println()

	if _, err := os.Stat(absCrackPath); os.IsNotExist(err) {
		fmt.Printf("  ⊘ 破解DLL不存在，跳过: %s\n", absCrackPath)
		fmt.Println("  （如需使用正版注册码，可跳过此步骤）")
	} else {
		err := dmsoft.CrackDm(absCrackPath)
		if err != nil {
			fmt.Printf("  ✗ CrackDm 失败: %v\n", err)
		} else {
			fmt.Println("  ✓ CrackDm 路径已设置")
		}
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第三步：New + Init — 启动helper + TCP连接")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  原理：")
	fmt.Println("    Init() 启动32位helper进程(dm_com_server.exe)")
	fmt.Println("    helper进程内：加载dm.dll → 注入破解 → 偏移量创建dm对象 → 监听TCP端口")
	fmt.Println("    64位主进程通过TCP+gob序列化连接helper")
	fmt.Println("    后续所有方法调用通过TCP转发到helper执行")
	fmt.Println()

	dm := dmsoft.New()
	if dm == nil {
		log.Fatal("  ✗ New() 返回 nil")
	}
	fmt.Println("  ✓ New() 成功")

	err = dm.Init()
	if err != nil {
		log.Fatalf("  ✗ Init() 失败: %v\n  提示：请确认dm_com_server.exe已编译(GOARCH=386)且在同级目录", err)
	}
	fmt.Println("  ✓ Init() 成功 (helper进程已启动，TCP连接已建立)")
	defer dm.Release()
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第四步：Reg — 注册大漠插件")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	nret := dm.Reg("", "")
	if nret == 1 {
		fmt.Println("  ✓ 注册成功")
	} else {
		fmt.Printf("  ✗ 注册失败 (返回值: %d)\n", nret)
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  功能测试")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	testBasicInfo(dm)
	testWindowOps(dm)
	testBindAndCapture(dm)
	testColorOps(dm)
	testMouseKeyboard(dm)
	testFindPic(dm)
	testOCR(dm)
	testClipboard(dm)
	testMemory(dm)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  测试完成")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  64位调用时序总结：")
	fmt.Println("    LoadDm → CrackDm → New → Init → Reg → 使用功能 → Release")
	fmt.Println()
	fmt.Println("  常见问题排查：")
	fmt.Println("    1. LoadDm失败 → 检查dm.dll路径、dm_com_server.exe是否在同级目录")
	fmt.Println("    2. Init失败 → 检查dm_com_server.exe是否已编译(GOARCH=386)")
	fmt.Println("    3. CrackDm失败 → 检查Go.dll路径")
	fmt.Println("    4. 方法返回空/0 → 检查helper进程是否存活")
	fmt.Println("    5. 中文乱码 → TCP传输自动处理GBK↔UTF-8编码转换")
	fmt.Println()

	dm.UnBindWindow()
}

// testBasicInfo 测试基础信息获取功能
// 演示：版本号、屏幕分辨率、色深、DPI、机器码等
func testBasicInfo(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 基础信息 ──")

	version := dm.Ver()
	fmt.Printf("  版本号:     %s\n", version)

	screenW := dm.GetScreenWidth()
	screenH := dm.GetScreenHeight()
	fmt.Printf("  屏幕分辨率: %d x %d\n", screenW, screenH)

	depth := dm.GetScreenDepth()
	fmt.Printf("  色深:       %d 位\n", depth)

	dpi := dm.GetDPI()
	fmt.Printf("  DPI:        %d\n", dpi)

	is64 := dm.Is64Bit()
	fmt.Printf("  64位系统:   %d (1=是)\n", is64)

	machineCode := dm.GetMachineCode()
	fmt.Printf("  机器码:     %s\n", machineCode)

	dm.SetShowErrorMsg(0)
	exePath, _ := os.Executable()
	workDir := filepath.Dir(exePath)
	dm.SetPath(workDir)
	fmt.Printf("  工作路径:   %s\n", workDir)
	fmt.Println()
}

// testWindowOps 测试窗口操作功能
// 演示：获取前台窗口、窗口标题、类名、进程信息、枚举窗口等
func testWindowOps(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 窗口操作 ──")

	foreHwnd := dm.GetForegroundWindow()
	fmt.Printf("  前台窗口句柄: %d\n", foreHwnd)

	if foreHwnd != 0 {
		title := dm.GetWindowTitle(foreHwnd)
		fmt.Printf("  窗口标题:     %s\n", title)

		className := dm.GetWindowClass(foreHwnd)
		fmt.Printf("  窗口类名:     %s\n", className)

		pid := dm.GetWindowProcessId(foreHwnd)
		fmt.Printf("  进程ID:       %d\n", pid)

		procPath := dm.GetWindowProcessPath(foreHwnd)
		fmt.Printf("  进程路径:     %s\n", procPath)

		x1, y1, x2, y2 := int32(0), int32(0), int32(0), int32(0)
		dm.GetWindowRect(foreHwnd, &x1, &y1, &x2, &y2)
		fmt.Printf("  窗口矩形:     (%d,%d)-(%d,%d)\n", x1, y1, x2, y2)
	}

	hwndList := dm.EnumWindow(0, "", "", 1)
	count := 0
	if hwndList != "" {
		count = len(strings.Split(hwndList, ","))
	}
	fmt.Printf("  枚举窗口数:   %d\n", count)
	fmt.Println()
}

// testBindAndCapture 测试窗口绑定与截图功能
// 演示：绑定窗口（GDI模式）、获取客户区大小、截图保存为BMP和PNG
func testBindAndCapture(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 窗口绑定与截图 ──")

	hwnd := dm.GetForegroundWindow()
	if hwnd == 0 {
		hwnd = dm.GetWindow(0, 0)
	}

	bindResult := dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)
	fmt.Printf("  BindWindow:   %d (1=成功)\n", bindResult)

	if bindResult == 1 {
		var w, h int32
		dm.GetClientSize(hwnd, &w, &h)
		fmt.Printf("  客户区大小:   %d x %d\n", w, h)

		captureResult := dm.Capture(0, 0, 200, 200, "test_x64.bmp")
		fmt.Printf("  Capture:      %d (1=成功)\n", captureResult)

		capturePng := dm.CapturePng(0, 0, 200, 200, "test_x64.png")
		fmt.Printf("  CapturePng:   %d (1=成功)\n", capturePng)
	} else {
		fmt.Println("  绑定失败，跳过截图测试")
	}
	fmt.Println()
}

// testColorOps 测试取色功能
// 演示：获取指定坐标的颜色（RGB/BGR/HSV格式）、区域平均颜色
func testColorOps(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 取色测试 ──")

	color := dm.GetColor(100, 100)
	fmt.Printf("  (100,100) RGB:  %s\n", color)

	colorBGR := dm.GetColorBGR(100, 100)
	fmt.Printf("  (100,100) BGR:  %s\n", colorBGR)

	colorHSV := dm.GetColorHSV(100, 100)
	fmt.Printf("  (100,100) HSV:  %s\n", colorHSV)

	aveRGB := dm.GetAveRGB(0, 0, 100, 100)
	fmt.Printf("  区域平均RGB:    %s\n", aveRGB)

	aveHSV := dm.GetAveHSV(0, 0, 100, 100)
	fmt.Printf("  区域平均HSV:    %s\n", aveHSV)
	fmt.Println()
}

// testMouseKeyboard 测试鼠标键盘操作
// 演示：获取鼠标位置、移动鼠标、获取鼠标形状、检测按键状态
func testMouseKeyboard(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 鼠标键盘 ──")

	mouseX, mouseY := int32(0), int32(0)
	dm.GetCursorPos(&mouseX, &mouseY)
	fmt.Printf("  鼠标位置:    (%d, %d)\n", mouseX, mouseY)

	moveResult := dm.MoveTo(100, 100)
	fmt.Printf("  MoveTo:      %d (1=成功)\n", moveResult)

	time.Sleep(50 * time.Millisecond)
	dm.MoveTo(mouseX, mouseY)
	fmt.Printf("  恢复鼠标:    (%d, %d)\n", mouseX, mouseY)

	cursorShape := dm.GetCursorShape()
	fmt.Printf("  鼠标形状:    %s\n", cursorShape)

	keyState := dm.GetKeyState(65)
	fmt.Printf("  A键状态:     %d (0=弹起, 1=按下)\n", keyState)
	fmt.Println()
}

// testFindPic 测试找图找色功能（含输出参数）
// 演示：FindPic/FindPicE/FindPicEx、FindColor、FindMultiColor
// 特别说明：64位TCP模式下，输出参数(*int32)通过gob序列化跨进程回传
func testFindPic(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 找图测试（含输出参数） ──")
	fmt.Println("  说明：64位TCP模式下，输出参数(*int32)通过gob序列化跨进程回传")

	findX, findY := int32(0), int32(0)
	findResult := dm.FindPic(0, 0, 200, 200, "test_x64.bmp", "000000", 0.9, 0, &findX, &findY)
	if findResult >= 0 {
		fmt.Printf("  FindPic:     找到！索引=%d, 坐标=(%d,%d)\n", findResult, findX, findY)
	} else {
		fmt.Printf("  FindPic:     未找到 (返回值=%d)\n", findResult)
	}

	findE := dm.FindPicE(0, 0, 200, 200, "test_x64.bmp", "000000", 0.9, 0)
	fmt.Printf("  FindPicE:    %s\n", findE)

	findEx := dm.FindPicEx(0, 0, 200, 200, "test_x64.bmp", "000000", 0.9, 0)
	if findEx != "" {
		fmt.Printf("  FindPicEx:   %s\n", findEx)
	}

	colorX, colorY := int32(0), int32(0)
	colorResult := dm.FindColor(0, 0, 200, 200, "FFFFFF", 1.0, 0, &colorX, &colorY)
	fmt.Printf("  FindColor:   %d, 坐标=(%d,%d)\n", colorResult, colorX, colorY)

	multiX, multiY := int32(0), int32(0)
	multiResult := dm.FindMultiColor(0, 0, 200, 200, "FFFFFF", "1|0|000000", 1.0, 0, &multiX, &multiY)
	fmt.Printf("  FindMulti:   %d, 坐标=(%d,%d)\n", multiResult, multiX, multiY)
	fmt.Println()
}

// testOCR 测试OCR文字识别功能
// 演示：Ocr文字识别、FindStr文字查找（含输出参数）
// 说明：TCP传输自动处理UTF-8↔GBK编码转换
func testOCR(dm *dmsoft.DmSoft) {
	fmt.Println("  ── OCR文字识别 ──")
	fmt.Println("  说明：TCP传输自动处理UTF-8↔GBK编码转换")

	ocrResult := dm.Ocr(0, 0, 200, 200, "ffffff-000000", 1.0)
	fmt.Printf("  Ocr:         %s\n", truncate(ocrResult, 50))

	findStrX, findStrY := int32(0), int32(0)
	findStrResult := dm.FindStr(0, 0, 500, 500, "测试", "ffffff-000000", 1.0, &findStrX, &findStrY)
	fmt.Printf("  FindStr:     %d, 坐标=(%d,%d)\n", findStrResult, findStrX, findStrY)

	findStrE := dm.FindStrE(0, 0, 500, 500, "测试", "ffffff-000000", 1.0)
	fmt.Printf("  FindStrE:    %s\n", findStrE)
	fmt.Println()
}

// testClipboard 测试剪贴板操作
// 演示：设置和获取剪贴板内容（包含中文）
func testClipboard(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 剪贴板 ──")

	testStr := "大漠64位COM测试Hello世界"
	setResult := dm.SetClipboard(testStr)
	fmt.Printf("  SetClipboard: %d (1=成功)\n", setResult)

	getContent := dm.GetClipboard()
	fmt.Printf("  GetClipboard: %s\n", getContent)
	fmt.Println()
}

// testMemory 测试内存操作功能
// 演示：获取进程ID、打开进程、获取模块基址
func testMemory(dm *dmsoft.DmSoft) {
	fmt.Println("  ── 内存操作 ──")

	hwnd := dm.GetForegroundWindow()
	pid := dm.GetWindowProcessId(hwnd)
	fmt.Printf("  进程ID:      %d\n", pid)

	if pid != 0 {
		openResult := dm.OpenProcess(pid)
		fmt.Printf("  OpenProcess: %d (非0=成功)\n", openResult)

		baseAddr := dm.GetModuleBaseAddr(pid, "")
		fmt.Printf("  模块基址:    0x%X\n", baseAddr)
	}
	fmt.Println()
}

// truncate 截断字符串到指定长度，超出部分用"..."代替
// 参数:
//   - s: 原始字符串
//   - maxLen: 最大长度
// 返回值: 截断后的字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
