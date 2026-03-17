// main.go - 程序入口
// 对应原C++项目中的 72424_C++.cpp
//
// 功能说明：
//   本文件是大漠插件Go版本的示例程序和测试入口。
//   演示了如何加载大漠插件、创建对象、注册、绑定窗口以及使用各种功能。
//
// 程序流程：
//   1. 加载大漠插件DLL (xd47243.dll)
//   2. 加载破解DLL (Go.dll) 并执行破解函数
//   3. 创建大漠插件对象
//   4. 注册大漠插件
//   5. 获取系统信息和屏幕分辨率
//   6. 获取前台窗口并绑定
//   7. 执行各种功能测试（取色、鼠标、键盘、截图等）
//   8. 解绑窗口并释放资源
//
// 编译说明：
//   - 必须使用32位编译，因为大漠插件是32位DLL
//   - 编译命令: GOARCH=386 go build -o dmsoft.exe
//   - 或者: go env -w GOARCH=386 && go build -o dmsoft.exe
//
// 运行要求：
//   - 需要xd47243.dll（大漠插件）和Go.dll（破解）在同一目录
//   - 部分功能需要管理员权限
//   - 绑定窗口后才能进行截图、取色等操作
//
// 测试内容：
//   - 基础信息：版本、屏幕分辨率、系统信息
//   - 窗口操作：查找、绑定、获取信息
//   - 鼠标操作：获取位置、移动
//   - 键盘操作：获取按键状态
//   - 屏幕操作：取色、截图
//   - 剪贴板：设置和获取内容
//   - 文件操作：检查文件存在
//   - 进程信息：获取进程ID和路径

package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

const (
	// DmPluginPath 大漠插件DLL文件路径
	// 这是大漠插件的核心DLL文件，提供所有自动化功能
	DmPluginPath = "xd47243.dll"

	// CrackDllPath 破解DLL文件路径
	// 用于破解大漠插件注册限制
	CrackDllPath = "Go.dll"
)

// main 程序入口函数
// 功能: 演示大漠插件的完整使用流程和各种功能测试
func main() {
	// ========== 第一步：加载大漠插件DLL ==========
	// LoadDm函数会加载指定路径的DLL文件，并返回模块句柄
	dmHModule, err := LoadDm(DmPluginPath)
	if err != nil {
		log.Fatalf("加载大漠插件失败: %v", err)
		return
	}
	fmt.Printf("大漠插件加载成功，模块句柄: %v\n", dmHModule)

	// ========== 第二步：加载并执行破解DLL ==========
	// 破解DLL用于绕过注册验证，必须先于大漠对象创建执行
	goHModule, err := syscall.LoadLibrary(CrackDllPath)
	if err != nil {
		log.Fatalf("加载破解DLL失败: %v", err)
		return
	}
	defer syscall.FreeLibrary(goHModule) // 确保程序退出时释放DLL

	// 获取破解函数地址
	goFunAddr, err := syscall.GetProcAddress(goHModule, "Go")
	if err != nil {
		log.Fatalf("获取Go函数地址失败: %v", err)
		return
	}

	// 执行破解函数 - 使用SyscallN直接调用，不使用回调
	// 这会修改大漠插件的内存，绕过注册检查
	syscall.SyscallN(uintptr(goFunAddr), dmHModule)
	fmt.Println("破解函数执行完成")

	// ========== 第三步：创建大漠插件对象 ==========
	// NewDmSoftImpl创建大漠插件实例，需要调用Init()初始化对象
	dm := NewDmSoftImpl()
	if dm == nil {
		log.Fatal("创建大漠对象失败")
		return
	}
	dm.Init() // 初始化大漠对象，创建内部COM对象
	defer dm.Release() // 确保程序退出时释放对象，避免内存泄漏

	// ========== 第四步：注册大漠插件 ==========
	// Reg函数用于注册插件，参数为注册码和版本号
	// 由于已经执行破解，传入空字符串即可
	nret := dm.Reg("", "")
	if nret == 1 {
		fmt.Println("大漠注册成功")
	} else {
		fmt.Println("大漠注册失败")
	}

	fmt.Println("\n========== 常用函数测试 ==========")

	// ========== 测试1: 获取插件版本 ==========
	// Ver函数返回大漠插件的版本号
	version := dm.Ver()
	fmt.Printf("大漠插件版本: %s\n", version)

	// ========== 测试2: 获取屏幕分辨率 ==========
	// GetScreenWidth和GetScreenHeight返回主屏幕的宽度和高度
	screenWidth := dm.GetScreenWidth()
	screenHeight := dm.GetScreenHeight()
	fmt.Printf("屏幕分辨率: %d x %d\n", screenWidth, screenHeight)

	// ========== 测试3: 设置错误提示 ==========
	// SetShowErrorMsg(0)关闭错误提示框，避免弹出窗口干扰
	dm.SetShowErrorMsg(0)
	fmt.Println("关闭错误提示框")

	// ========== 测试4: 获取当前路径 ==========
	// GetDir(0)获取当前工作目录
	path := dm.GetDir(0)
	fmt.Printf("当前路径: %s\n", path)

	// ========== 测试5: 设置全局路径 ==========
	// SetPath设置大漠插件的全局工作路径，用于保存截图等文件
	dm.SetPath("c:\\test")
	fmt.Println("设置全局路径: c:\\test")

	// ========== 测试6: 获取当前时间 ==========
	// GetTime返回当前系统时间戳（毫秒）
	currentTime := dm.GetTime()
	fmt.Printf("当前时间戳: %d\n", currentTime)

	// ========== 测试7: 获取系统信息 ==========
	// 获取硬盘序列号、CPU类型、操作系统类型等硬件信息
	diskSerial := dm.GetDiskSerial(0) // 0表示主硬盘
	fmt.Printf("硬盘序列号: %s\n", diskSerial)

	cpuType := dm.GetCpuType() // 1=Intel, 2=AMD
	fmt.Printf("CPU类型: %d\n", cpuType)

	osType := dm.GetOsType() // 操作系统类型代码
	fmt.Printf("操作系统类型: %d\n", osType)

	// ========== 测试8: 获取前台窗口信息 ==========
	// GetForegroundWindow获取当前活动窗口的句柄
	foreHwnd := dm.GetForegroundWindow()
	fmt.Printf("前台窗口句柄: %d\n", foreHwnd)

	if foreHwnd != 0 {
		// 获取窗口标题和类名
		foreTitle := dm.GetWindowTitle(foreHwnd)
		fmt.Printf("前台窗口标题: %s\n", foreTitle)

		foreClass := dm.GetWindowClass(foreHwnd)
		fmt.Printf("前台窗口类名: %s\n", foreClass)
	}

	// ========== 测试9: 绑定窗口 ==========
	// 绑定窗口是使用大漠插件功能的关键步骤
	// 绑定后才能进行截图、取色、鼠标键盘操作等
	fmt.Println("\n--- 窗口绑定测试 ---")
	var bindHwnd int32 = 0
	if foreHwnd != 0 {
		bindHwnd = foreHwnd // 使用前台窗口
	} else {
		bindHwnd = dm.GetWindow(0, 0) // 如果没有前台窗口，使用桌面窗口
	}

	// BindWindow参数说明:
	//   hwnd: 窗口句柄
	//   display: 显示模式 ("normal"=正常, "gdi"=GDI, "dx"=DirectX)
	//   mouse: 鼠标模式 ("normal"=正常, "windows"=Windows消息)
	//   keypad: 键盘模式 ("normal"=正常, "windows"=Windows消息)
	//   mode: 附加模式 (0=默认)
	bindResult := dm.BindWindow(bindHwnd, "normal", "normal", "normal", 0)
	fmt.Printf("绑定窗口结果: %d (1=成功)\n", bindResult)

	if bindResult != 1 {
		fmt.Println("窗口绑定失败，部分功能可能无法使用")
	} else {
		fmt.Println("窗口绑定成功，可以进行后续操作")

		// 检查绑定状态
		isBind := dm.IsBind(bindHwnd)
		fmt.Printf("窗口绑定状态: %d\n", isBind)
	}

	// ========== 测试10: 屏幕取色 ==========
	// GetColor获取指定屏幕坐标点的颜色值（16进制字符串）
	// 必须在绑定窗口后才能正常工作
	fmt.Println("\n--- 屏幕取色测试 ---")
	if screenWidth > 0 && screenHeight > 0 {
		color := dm.GetColor(screenWidth/2, screenHeight/2)
		fmt.Printf("屏幕中心点颜色: %s\n", color)
	}

	pointColor := dm.GetColor(100, 100)
	fmt.Printf("屏幕(100,100)颜色: %s\n", pointColor)

	// ========== 测试11: 鼠标操作 ==========
	// 鼠标相关操作：获取位置、移动鼠标
	// 需要在绑定窗口后才能正常工作
	fmt.Println("\n--- 鼠标操作测试 ---")
	mouseX, mouseY := int32(0), int32(0)
	dm.GetCursorPos(&mouseX, &mouseY)
	fmt.Printf("当前鼠标位置: (%d, %d)\n", mouseX, mouseY)

	// 移动鼠标到屏幕中心
	centerX := screenWidth / 2
	centerY := screenHeight / 2
	moveResult := dm.MoveTo(centerX, centerY)
	fmt.Printf("移动鼠标到中心点(%d, %d): %d\n", centerX, centerY, moveResult)

	// 延迟一下，让用户看到鼠标移动
	time.Sleep(100 * time.Millisecond)

	// 移动回原位置
	dm.MoveTo(mouseX, mouseY)
	fmt.Printf("恢复鼠标位置到(%d, %d)\n", mouseX, mouseY)

	// ========== 测试12: 键盘操作 ==========
	// GetKeyState获取指定按键的状态
	// 参数是虚拟键码，65对应A键
	fmt.Println("\n--- 键盘操作测试 ---")
	keyState := dm.GetKeyState(65) // 65 = A键的虚拟键码
	fmt.Printf("A键状态: %d (0=未按下, 1=已按下)\n", keyState)

	// ========== 测试13: 屏幕截图 ==========
	// Capture截取指定区域的屏幕图像并保存为BMP文件
	// 参数: x1, y1, x2, y2, 文件名
	fmt.Println("\n--- 截图测试 ---")
	captureResult := dm.Capture(0, 0, 200, 200, "test.bmp")
	fmt.Printf("截图测试结果: %d (1=成功, 0=失败)\n", captureResult)

	// ========== 测试14-18: 图像相关功能 ==========
	// 这些功能需要字库文件或预存图片，当前跳过
	fmt.Println("\n--- 图像相关测试 ---")
	fmt.Println("找图/找字/OCR/颜色查找测试跳过（需要字库或图片支持）")
	fmt.Println("  - FindPic: 在屏幕上查找图片")
	fmt.Println("  - FindStr: 在屏幕上查找文字")
	fmt.Println("  - Ocr: 识别屏幕上的文字")
	fmt.Println("  - FindColor: 在屏幕上查找颜色")

	// ========== 测试19: 内存操作 ==========
	// 获取绑定窗口的进程ID，用于后续的内存操作
	// 注意：内存读写功能需要管理员权限
	fmt.Println("\n--- 内存操作测试 ---")
	processId := dm.GetWindowProcessId(bindHwnd)
	fmt.Printf("绑定窗口进程ID: %d\n", processId)
	fmt.Println("提示: 内存读写功能需要管理员权限")

	// ========== 测试20: 设置功能 ==========
	// 设置大漠插件的各种参数
	fmt.Println("\n--- 设置功能测试 ---")

	// SetMouseDelay设置鼠标操作延迟，防止操作过快被检测
	dm.SetMouseDelay("normal", 100)
	fmt.Println("设置鼠标延迟: 100ms")

	// SetKeypadDelay设置键盘操作延迟
	dm.SetKeypadDelay("normal", 50)
	fmt.Println("设置键盘延迟: 50ms")

	// SetSimMode设置仿真模式
	dm.SetSimMode(0)
	fmt.Println("设置仿真模式: 0 (正常模式)")

	// ========== 测试21: 窗口操作 ==========
	// 获取窗口信息和位置
	fmt.Println("\n--- 窗口操作测试 ---")

	// GetWindow(0, 0)获取桌面窗口句柄
	desktopHwnd := dm.GetWindow(0, 0)
	fmt.Printf("桌面窗口句柄: %d\n", desktopHwnd)

	// GetWindowRect获取窗口在屏幕上的位置和大小（包含标题栏和边框）
	dcX1, dcY1, dcX2, dcY2 := int32(0), int32(0), int32(0), int32(0)
	getRectResult := dm.GetWindowRect(desktopHwnd, &dcX1, &dcY1, &dcX2, &dcY2)
	fmt.Printf("获取桌面窗口矩形: %d, (%d,%d)-(%d,%d)\n", getRectResult, dcX1, dcY1, dcX2, dcY2)

	// GetClientRect获取窗口客户区的位置和大小（不包含标题栏和边框）
	clientX1, clientY1, clientX2, clientY2 := int32(0), int32(0), int32(0), int32(0)
	getClientResult := dm.GetClientRect(desktopHwnd, &clientX1, &clientY1, &clientX2, &clientY2)
	fmt.Printf("获取桌面客户区矩形: %d, (%d,%d)-(%d,%d)\n", getClientResult, clientX1, clientY1, clientX2, clientY2)

	// ========== 测试22: 时间相关 ==========
	// GetTime返回系统运行时间（毫秒），可用于计时
	fmt.Println("\n--- 时间相关测试 ---")
	startTime := dm.GetTime()
	time.Sleep(100 * time.Millisecond) // 休眠100毫秒
	endTime := dm.GetTime()
	fmt.Printf("时间差: %d ms\n", endTime-startTime)

	// ========== 测试23: 随机数 ==========
	// 大漠插件的随机数功能有限，建议使用Go标准库math/rand
	fmt.Println("\n--- 随机数测试 ---")
	fmt.Println("随机数测试跳过（建议使用标准库math/rand）")

	// ========== 测试24: 字符串操作 ==========
	// 字符串转换建议使用Go标准库strconv
	fmt.Println("\n--- 字符串操作测试 ---")
	fmt.Println("字符串转数字测试跳过（建议使用标准库strconv）")

	// ========== 测试25: 文件操作 ==========
	// IsFileExist检查文件是否存在
	fmt.Println("\n--- 文件操作测试 ---")
	fileExists := dm.IsFileExist("test.bmp")
	fmt.Printf("文件是否存在(test.bmp): %d (1=存在, 0=不存在)\n", fileExists)

	// ========== 测试26: 剪贴板操作 ==========
	// SetClipboard设置剪贴板内容，GetClipboard获取剪贴板内容
	fmt.Println("\n--- 剪贴板测试 ---")
	clipResult := dm.SetClipboard("大漠插件测试")
	fmt.Printf("设置剪贴板结果: %d (1=成功)\n", clipResult)

	clipContent := dm.GetClipboard()
	fmt.Printf("获取剪贴板内容: %s\n", clipContent)

	// ========== 测试27: 枚举窗口 ==========
	// EnumWindow枚举系统中所有窗口，返回窗口句柄列表字符串
	fmt.Println("\n--- 枚举窗口测试 ---")
	hwndList := dm.EnumWindow(0, "", "", 1)
	fmt.Printf("枚举窗口结果长度: %d 字符\n", len(hwndList))
	fmt.Println("提示: 可以使用strings.Split(hwndList, \",\")解析句柄列表")

	// ========== 测试28: 进程相关 ==========
	// 获取绑定窗口的进程信息
	fmt.Println("\n--- 进程相关测试 ---")
	processPath := dm.GetWindowProcessPath(bindHwnd)
	fmt.Printf("绑定窗口进程路径: %s\n", processPath)

	pid := dm.GetWindowProcessId(bindHwnd)
	fmt.Printf("绑定窗口进程ID: %d\n", pid)

	// ========== 测试29: 解绑窗口 ==========
	// UnBindWindow解除当前绑定的窗口
	// 解绑后可以进行新的窗口绑定
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
