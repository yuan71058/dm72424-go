// example/find_window/main.go - 大漠插件查找窗口示例
//
// 版本: v1.6.0
// 更新日期: 2026-03-22
//
// 功能说明:
//   本文件演示如何使用 dmsoft 库查找指定类名和标题的窗口，并进行找图测试。
//   示例：查找类名为 Qt51514QWindowIcon，标题为 朋友圈 的窗口。
//   库已内置 UTF-8 到 GBK 自动编码转换，中文参数可直接使用。
//
// 找图方法演示:
//   基本找图系列 (sim参数0.1-1.0):
//   - FindPic: 基本找图，返回图片索引
//   - FindPicE: 找图返回坐标字符串
//   - FindPicEx: 高级找图，返回所有匹配位置
//   - FindPicS: 找图返回图片索引字符串
//   - FindPicExS: 高级找图返回详细字符串
//
//   相似度找图系列 (sim参数0-100):
//   - FindPicSim: 相似度找图
//   - FindPicSimE: 相似度找图返回坐标字符串
//   - FindPicSimEx: 相似度高级找图
//
//   内存找图系列:
//   - FindPicSimMem: 内存相似度找图
//   - FindPicSimMemE: 内存相似度找图返回坐标字符串
//   - FindPicSimMemEx: 内存相似度高级找图
//
//   AI找图系列:
//   - AiFindPic: AI找图
//   - AiFindPicEx: AI高级找图
//   - AiFindPicMem: AI内存找图
//   - AiFindPicMemEx: AI高级内存找图
//
// 编译说明:
//   必须使用32位编译: GOARCH=386 go build -o find_window.exe
//
// 运行要求:
//   - 需要 xd47243.dll 和 Go.dll 在同一目录
//   - 部分功能需要管理员权限

package main

import (
	"fmt"
	"log"

	dmsoft "github.com/yuan71058/dm72424-go"
)

const (
	DmPluginPath = "E:\\SRC\\dm\\72424_C++\\go\\example\\find_window\\xd47243.dll"
	CrackDllPath = "E:\\SRC\\dm\\72424_C++\\go\\example\\find_window\\Go.dll"
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

	// ========== 第五步：设置错误提示 ==========
	dm.SetShowErrorMsg(0)
	fmt.Println("关闭错误提示框")

	// ========== 第六步：查找窗口 ==========
	fmt.Println("\n========== 查找窗口 ==========")

	className := "Qt51514QWindowIcon"
	windowTitle := "朋友圈"

	// 直接使用中文参数，库内部自动进行GBK编码转换
	hwnd := dm.FindWindow(className, windowTitle)

	if hwnd > 0 {
		fmt.Printf("找到窗口！\n")
		fmt.Printf("  窗口句柄: %d\n", hwnd)
		fmt.Printf("  窗口类名: %s\n", className)
		fmt.Printf("  窗口标题: %s\n", windowTitle)

		// 获取窗口详细信息
		x1, y1, x2, y2 := int32(0), int32(0), int32(0), int32(0)
		dm.GetWindowRect(hwnd, &x1, &y1, &x2, &y2)
		fmt.Printf("  窗口位置: (%d, %d) - (%d, %d)\n", x1, y1, x2, y2)
		fmt.Printf("  窗口大小: %d x %d\n", x2-x1, y2-y1)

		// 获取进程信息
		pid := dm.GetWindowProcessId(hwnd)
		fmt.Printf("  进程ID: %d\n", pid)

		processPath := dm.GetWindowProcessPath(hwnd)
		fmt.Printf("  进程路径: %s\n", processPath)

		// ========== 第七步：绑定窗口 ==========
		fmt.Println("\n========== 绑定窗口 ==========")

		// 绑定模式说明:
		//   display: gdi - GDI模式，适用于大多数窗口
		//   mouse: windows3 - Windows3鼠标模式，后台鼠标模拟
		//   keypad: windows - Windows键盘模式，后台键盘模拟
		//   mode: 0 - 普通模式
		bindResult := dm.BindWindow(hwnd, "gdi", "windows3", "windows", 0)
		fmt.Printf("绑定窗口结果: %d (1=成功)\n", bindResult)

		if bindResult == 1 {
			fmt.Println("窗口绑定成功！")

			// 检查绑定状态
			isBind := dm.IsBind(hwnd)
			fmt.Printf("窗口绑定状态: %d (1=已绑定)\n", isBind)

			// ========== 第八步：截图测试 ==========
			fmt.Println("\n========== 截图测试 ==========")

			// 截取窗口区域
			captureResult := dm.Capture(x1, y1, x2, y2, "window_capture.bmp")
			fmt.Printf("截图结果: %d (1=成功)\n", captureResult)

			// 截取窗口左上角区域作为找图模板
			templateX1, templateY1 := x1, y1
			templateX2, templateY2 := x1+100, y1+100
			captureTemplate := dm.Capture(templateX1, templateY1, templateX2, templateY2, "template.bmp")
			fmt.Printf("截取模板图片结果: %d (1=成功)\n", captureTemplate)

			// 获取窗口中心点颜色
			centerX := (x1 + x2) / 2
			centerY := (y1 + y2) / 2
			color := dm.GetColor(centerX, centerY)
			fmt.Printf("窗口中心点(%d, %d)颜色: %s\n", centerX, centerY, color)

			// ========== 第九步：找图测试 ==========
			fmt.Println("\n========== 找图测试 ==========")

			picPath := "E:\\SRC\\biaoge\\ai-excel-assistant\\build\\bin\\pic\\FABIAO.bmp"
			findX, findY := int32(0), int32(0)

			// ==================== 基本找图系列 (sim参数0.1-1.0) ====================
			fmt.Println("\n========== 基本找图系列 (sim: 0.1-1.0) ==========")

			// 方法1: FindPic - 基本找图
			fmt.Println("\n--- 方法1: FindPic 基本找图 ---")
			findResult := dm.FindPic(x1, y1, x2, y2, picPath, "000000", 0.8, 0, &findX, &findY)
			if findResult >= 0 {
				fmt.Printf("FindPic 找到图片！\n")
				fmt.Printf("  图片索引: %d\n", findResult)
				fmt.Printf("  找到位置: (%d, %d)\n", findX, findY)
			} else {
				fmt.Printf("FindPic 未找到图片，返回值: %d\n", findResult)
			}

			// 方法2: FindPicE - 找图返回坐标字符串
			fmt.Println("\n--- 方法2: FindPicE 返回坐标字符串 ---")
			findStr := dm.FindPicE(x1, y1, x2, y2, picPath, "000000", 0.8, 0)
			if findStr != "" && findStr != "-1|-1|-1" {
				fmt.Printf("FindPicE 找到图片！\n")
				fmt.Printf("  结果: %s (格式: 索引|x|y)\n", findStr)
			} else {
				fmt.Printf("FindPicE 未找到图片，结果: %s\n", findStr)
			}

			// 方法3: FindPicEx - 高级找图，返回所有匹配位置
			fmt.Println("\n--- 方法3: FindPicEx 高级找图 ---")
			findExStr := dm.FindPicEx(x1, y1, x2, y2, picPath, "000000", 0.8, 0)
			if findExStr != "" {
				fmt.Printf("FindPicEx 结果:\n")
				fmt.Printf("  %s\n", findExStr)
			} else {
				fmt.Printf("FindPicEx 未找到图片\n")
			}

			// 方法4: FindPicS - 找图返回图片索引字符串
			fmt.Println("\n--- 方法4: FindPicS 返回索引字符串 ---")
			findSResult := dm.FindPicS(x1, y1, x2, y2, picPath, "000000", 0.8, 0, &findX, &findY)
			fmt.Printf("FindPicS 结果: %s\n", findSResult)
			fmt.Printf("  坐标: (%d, %d)\n", findX, findY)

			// 方法5: FindPicExS - 高级找图返回详细字符串
			fmt.Println("\n--- 方法5: FindPicExS 高级找图详细 ---")
			findExSStr := dm.FindPicExS(x1, y1, x2, y2, picPath, "000000", 0.8, 0)
			if findExSStr != "" {
				fmt.Printf("FindPicExS 结果:\n")
				fmt.Printf("  %s\n", findExSStr)
			} else {
				fmt.Printf("FindPicExS 未找到图片\n")
			}

			// ==================== 相似度找图系列 (sim参数0-100) ====================
			fmt.Println("\n========== 相似度找图系列 (sim: 0-100) ==========")

			// 方法6: FindPicSim - 相似度找图
			fmt.Println("\n--- 方法6: FindPicSim 相似度找图 ---")
			findSimResult := dm.FindPicSim(x1, y1, x2, y2, picPath, "000000", 80, 0, &findX, &findY)
			if findSimResult >= 0 {
				fmt.Printf("FindPicSim 找到图片！\n")
				fmt.Printf("  图片索引: %d\n", findSimResult)
				fmt.Printf("  坐标: (%d, %d)\n", findX, findY)
			} else {
				fmt.Printf("FindPicSim 未找到图片，返回值: %d\n", findSimResult)
			}

			// 方法7: FindPicSimE - 相似度找图返回坐标字符串
			fmt.Println("\n--- 方法7: FindPicSimE 返回坐标字符串 ---")
			findSimE := dm.FindPicSimE(x1, y1, x2, y2, picPath, "000000", 80, 0)
			fmt.Printf("FindPicSimE 结果: %s\n", findSimE)

			// 方法8: FindPicSimEx - 相似度高级找图
			fmt.Println("\n--- 方法8: FindPicSimEx 高级找图 ---")
			findSimEx := dm.FindPicSimEx(x1, y1, x2, y2, picPath, "000000", 80, 0)
			fmt.Printf("FindPicSimEx 结果: %s\n", findSimEx)

			// ==================== 多图片查找 ====================
			fmt.Println("\n========== 多图片查找 (用|分隔多个图片) ==========")
			multiResult := dm.FindPicE(x1, y1, x2, y2, picPath+"|template.bmp", "000000", 0.8, 0)
			fmt.Printf("多图片查找结果: %s\n", multiResult)

			// 释放图片资源
			fmt.Println("\n--- 释放图片资源 ---")
			freeResult := dm.FreePic(picPath + "|template.bmp")
			fmt.Printf("释放图片资源结果: %d (1=成功)\n", freeResult)

			// ========== 第十步：解绑窗口 ==========
			fmt.Println("\n========== 解绑窗口 ==========")
			unbindResult := dm.UnBindWindow()
			fmt.Printf("解绑窗口结果: %d (1=成功)\n", unbindResult)
		} else {
			fmt.Println("窗口绑定失败！")
			lastError := dm.GetLastError()
			fmt.Printf("错误码: %d\n", lastError)
		}
	} else {
		fmt.Printf("未找到窗口！\n")
		fmt.Printf("  查找条件:\n")
		fmt.Printf("    窗口类名: %s\n", className)
		fmt.Printf("    窗口标题: %s\n", windowTitle)
		fmt.Println("\n请确保目标窗口已打开！")
	}

	fmt.Println("\n========== 查找完成 ==========")
}
