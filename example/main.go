// example/main.go - 大漠插件Go库使用示例
//
// 版本: v1.6.0
// 更新日期: 2026-03-22
//
// 功能说明:
//   本文件演示如何使用 dmsoft 库进行大漠插件开发。
//   包含完整的初始化、注册、窗口绑定和功能测试流程。
//
// 编译说明:
//   必须使用32位编译: GOARCH=386 go build -o dmsoft_test.exe
//
// 运行要求:
//   - 需要 xd47243.dll 和 Go.dll 在同一目录
//   - 部分功能需要管理员权限

package main

import (
	"fmt"
	"log"
	"time"

	dmsoft "github.com/yuan71058/dm72424-go"
)

const (
	DmPluginPath = "xd47243.dll"
	CrackDllPath = "Go.dll"
)

func main() {
	// ========== 第一步：加载大漠插件DLL ==========
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

	// ========== 第三步：创建大漠插件对象 ==========
	dm := dmsoft.New()
	if dm == nil {
		log.Fatal("创建大漠对象失败")
	}
	dm.Init()
	defer dm.Release()

	// ========== 第四步：注册大漠插件 ==========
	nret := dm.Reg("", "")
	if nret == 1 {
		fmt.Println("大漠注册成功")
	} else {
		fmt.Println("大漠注册失败")
	}

	fmt.Println("\n========== 常用函数测试 ==========")

	// ========== 测试: 获取插件版本 ==========
	version := dm.Ver()
	fmt.Printf("大漠插件版本: %s\n", version)

	// ========== 测试: 获取屏幕分辨率 ==========
	screenWidth := dm.GetScreenWidth()
	screenHeight := dm.GetScreenHeight()
	fmt.Printf("屏幕分辨率: %d x %d\n", screenWidth, screenHeight)

	// ========== 测试: 设置错误提示 ==========
	dm.SetShowErrorMsg(0)
	fmt.Println("关闭错误提示框")

	// ========== 测试: 获取当前路径 ==========
	path := dm.GetDir(0)
	fmt.Printf("当前路径: %s\n", path)

	// ========== 测试: 设置全局路径 ==========
	dm.SetPath("c:\\test")
	fmt.Println("设置全局路径: c:\\test")

	// ========== 测试: 获取当前时间 ==========
	currentTime := dm.GetTime()
	fmt.Printf("当前时间戳: %d\n", currentTime)

	// ========== 测试: 获取系统信息 ==========
	diskSerial := dm.GetDiskSerial(0)
	fmt.Printf("硬盘序列号: %s\n", diskSerial)

	cpuType := dm.GetCpuType()
	fmt.Printf("CPU类型: %d\n", cpuType)

	osType := dm.GetOsType()
	fmt.Printf("操作系统类型: %d\n", osType)

	// ========== 测试: 获取前台窗口信息 ==========
	foreHwnd := dm.GetForegroundWindow()
	fmt.Printf("前台窗口句柄: %d\n", foreHwnd)

	if foreHwnd != 0 {
		foreTitle := dm.GetWindowTitle(foreHwnd)
		fmt.Printf("前台窗口标题: %s\n", foreTitle)

		foreClass := dm.GetWindowClass(foreHwnd)
		fmt.Printf("前台窗口类名: %s\n", foreClass)
	}

	// ========== 测试: 绑定窗口 ==========
	fmt.Println("\n--- 窗口绑定测试 ---")
	var bindHwnd int32 = 0
	if foreHwnd != 0 {
		bindHwnd = foreHwnd
	} else {
		bindHwnd = dm.GetWindow(0, 0)
	}

	bindResult := dm.BindWindow(bindHwnd, "normal", "normal", "normal", 0)
	fmt.Printf("绑定窗口结果: %d (1=成功)\n", bindResult)

	if bindResult == 1 {
		fmt.Println("窗口绑定成功，可以进行后续操作")
		isBind := dm.IsBind(bindHwnd)
		fmt.Printf("窗口绑定状态: %d\n", isBind)
	}

	// ========== 测试: 屏幕取色 ==========
	fmt.Println("\n--- 屏幕取色测试 ---")
	if screenWidth > 0 && screenHeight > 0 {
		color := dm.GetColor(screenWidth/2, screenHeight/2)
		fmt.Printf("屏幕中心点颜色: %s\n", color)
	}

	pointColor := dm.GetColor(100, 100)
	fmt.Printf("屏幕(100,100)颜色: %s\n", pointColor)

	// ========== 测试: 鼠标操作 ==========
	fmt.Println("\n--- 鼠标操作测试 ---")
	mouseX, mouseY := int32(0), int32(0)
	dm.GetCursorPos(&mouseX, &mouseY)
	fmt.Printf("当前鼠标位置: (%d, %d)\n", mouseX, mouseY)

	centerX := screenWidth / 2
	centerY := screenHeight / 2
	moveResult := dm.MoveTo(centerX, centerY)
	fmt.Printf("移动鼠标到中心点(%d, %d): %d\n", centerX, centerY, moveResult)

	time.Sleep(100 * time.Millisecond)
	dm.MoveTo(mouseX, mouseY)
	fmt.Printf("恢复鼠标位置到(%d, %d)\n", mouseX, mouseY)

	// ========== 测试: 键盘操作 ==========
	fmt.Println("\n--- 键盘操作测试 ---")
	keyState := dm.GetKeyState(65)
	fmt.Printf("A键状态: %d (0=未按下, 1=已按下)\n", keyState)

	// ========== 测试: 屏幕截图 ==========
	fmt.Println("\n--- 截图测试 ---")
	captureResult := dm.Capture(0, 0, 200, 200, "test.bmp")
	fmt.Printf("截图测试结果: %d (1=成功, 0=失败)\n", captureResult)

	// ========== 测试: 图像相关 ==========
	fmt.Println("\n--- 图像相关测试 ---")
	fmt.Println("找图/找字/OCR/颜色查找测试跳过（需要字库或图片支持）")

	// ========== 测试: 内存操作 ==========
	fmt.Println("\n--- 内存操作测试 ---")
	processId := dm.GetWindowProcessId(bindHwnd)
	fmt.Printf("绑定窗口进程ID: %d\n", processId)
	fmt.Println("提示: 内存读写功能需要管理员权限")

	// ========== 测试: 设置功能 ==========
	fmt.Println("\n--- 设置功能测试 ---")
	dm.SetMouseDelay("normal", 100)
	fmt.Println("设置鼠标延迟: 100ms")

	dm.SetKeypadDelay("normal", 50)
	fmt.Println("设置键盘延迟: 50ms")

	dm.SetSimMode(0)
	fmt.Println("设置仿真模式: 0 (正常模式)")

	// ========== 测试: 窗口操作 ==========
	fmt.Println("\n--- 窗口操作测试 ---")
	desktopHwnd := dm.GetWindow(0, 0)
	fmt.Printf("桌面窗口句柄: %d\n", desktopHwnd)

	dcX1, dcY1, dcX2, dcY2 := int32(0), int32(0), int32(0), int32(0)
	getRectResult := dm.GetWindowRect(desktopHwnd, &dcX1, &dcY1, &dcX2, &dcY2)
	fmt.Printf("获取桌面窗口矩形: %d, (%d,%d)-(%d,%d)\n", getRectResult, dcX1, dcY1, dcX2, dcY2)

	// ========== 测试: 时间相关 ==========
	fmt.Println("\n--- 时间相关测试 ---")
	startTime := dm.GetTime()
	time.Sleep(100 * time.Millisecond)
	endTime := dm.GetTime()
	fmt.Printf("时间差: %d ms\n", endTime-startTime)

	// ========== 测试: 文件操作 ==========
	fmt.Println("\n--- 文件操作测试 ---")
	fileExists := dm.IsFileExist("test.bmp")
	fmt.Printf("文件是否存在(test.bmp): %d (1=存在, 0=不存在)\n", fileExists)

	// ========== 测试: 剪贴板操作 ==========
	fmt.Println("\n--- 剪贴板测试 ---")
	clipResult := dm.SetClipboard("大漠插件测试")
	fmt.Printf("设置剪贴板结果: %d (1=成功)\n", clipResult)

	clipContent := dm.GetClipboard()
	fmt.Printf("获取剪贴板内容: %s\n", clipContent)

	// ========== 测试: 枚举窗口 ==========
	fmt.Println("\n--- 枚举窗口测试 ---")
	hwndList := dm.EnumWindow(0, "", "", 1)
	fmt.Printf("枚举窗口结果长度: %d 字符\n", len(hwndList))

	// ========== 测试: 进程相关 ==========
	fmt.Println("\n--- 进程相关测试 ---")
	processPath := dm.GetWindowProcessPath(bindHwnd)
	fmt.Printf("绑定窗口进程路径: %s\n", processPath)

	pid := dm.GetWindowProcessId(bindHwnd)
	fmt.Printf("绑定窗口进程ID: %d\n", pid)

	// ========== 测试: 解绑窗口 ==========
	fmt.Println("\n--- 解绑窗口 ---")
	unbindResult := dm.UnBindWindow()
	fmt.Printf("解绑窗口结果: %d (1=成功)\n", unbindResult)

	fmt.Println("\n========== 测试完成 ==========")
	fmt.Println("所有常用函数测试执行完毕！")
	fmt.Println("\n使用建议:")
	fmt.Println("  1. 使用前确保已加载DLL并创建对象")
	fmt.Println("  2. 绑定窗口后才能使用截图、取色等功能")
	fmt.Println("  3. 使用完毕后调用Release()释放资源")
	fmt.Println("  4. 内存操作需要管理员权限")
}
