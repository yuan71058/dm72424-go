package main

import (
	"fmt"
	"log"
	"syscall"

	dmsoft "gitee.com/yuan71058/dm72424-go/dmsoft"
)

func main() {
	fmt.Println("========== 大漠插件 Go 语言封装库示例 ==========")

	// ========== 第一步：加载大漠插件DLL ==========
	// LoadDm加载大漠插件DLL文件
	// 参数: DLL文件路径(相对或绝对路径)
	// 返回: 模块句柄和错误信息
	dmHModule, err := dmsoft.LoadDm("xd47243.dll")
	if err != nil {
		log.Fatalf("加载大漠插件失败: %v", err)
	}
	fmt.Printf("大漠插件加载成功，模块句柄: %v\n", dmHModule)

	// ========== 第二步：加载破解补丁（可选）==========
	// 如果需要绕过注册验证，可以加载破解DLL
	goHModule, err := syscall.LoadLibrary("Go.dll")
	if err != nil {
		log.Printf("加载破解DLL失败(可忽略): %v", err)
	} else {
		defer syscall.FreeLibrary(goHModule)
		goFunAddr, _ := syscall.GetProcAddress(goHModule, "Go")
		syscall.Syscall(uintptr(goFunAddr), 1, dmHModule, 0, 0)
		fmt.Println("破解函数执行完成")
	}

	// ========== 第三步：创建大漠插件对象 ==========
	// NewDmSoftImpl创建大漠插件实例，需要调用Init()初始化对象
	dm := dmsoft.NewDmSoftImpl()
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

	// ========== 第五步：获取基本信息 ==========
	fmt.Println("\n========== 基本信息 ==========")
	fmt.Printf("大漠插件版本: %s\n", dm.Ver())
	fmt.Printf("屏幕分辨率: %d x %d\n", dm.GetScreenWidth(), dm.GetScreenHeight())

	// 关闭错误提示框
	dm.SetShowErrorMsg(0)

	// ========== 第六步：绑定窗口 ==========
	fmt.Println("\n========== 窗口操作 ==========")
	hwnd := dm.GetForegroundWindow()
	fmt.Printf("前台窗口句柄: %d\n", hwnd)

	// 绑定窗口后才能使用截图、取色等功能
	ret := dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)
	if ret == 1 {
		fmt.Println("窗口绑定成功")

		// ========== 第七步：屏幕取色 ==========
		fmt.Println("\n========== 屏幕取色 ==========")
		color := dm.GetColor(100, 100)
		fmt.Printf("屏幕(100,100)颜色: %s\n", color)

		// ========== 第八步：解绑窗口 ==========
		dm.UnBindWindow()
		fmt.Println("窗口解绑成功")
	} else {
		fmt.Println("窗口绑定失败")
	}

	// ========== 使用建议 ==========
	fmt.Println("\n========== 使用建议 ==========")
	fmt.Println("1. 使用前确保已加载DLL并创建对象")
	fmt.Println("2. 绑定窗口后才能使用截图、取色等功能")
	fmt.Println("3. 使用完毕后调用Release()释放资源")
	fmt.Println("4. 内存操作需要管理员权限")

	fmt.Println("\n========== 示例完成 ==========")
}
