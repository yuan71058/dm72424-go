// example/load_ai/main.go - 大漠插件 LoadAi 加载AI模型示例
//
// 版本: v1.6.0
// 更新日期: 2026-07-06
//
// 功能说明:
//   本文件演示如何使用 dmsoft 库加载大漠插件内置 AI 模型，
//   并使用 AI 找图系列函数进行图像识别。
//
// AI 功能相关函数:
//   - LoadAi(file string) int32
//       从文件加载 AI 模型。
//       file 参数为 AI 模型文件路径（支持绝对或相对路径），
//       大漠插件使用 .module 后缀的模型文件（本项目附带 ai.module，约 4.77 MB）。
//
//       返回值（官方定义）:
//         1   表示成功
//        -1   打开文件失败
//        -2   内存初始化失败（正式版本可联系作者解决）
//        -3   参数错误
//        -4   加载错误
//        -5   Ai 模块初始化失败
//        -6   内存分配失败
//
//   - LoadAiMemory(addr int32, size int32) int32
//       从内存加载 AI 模型（适用于模型已读入内存的场景）。
//       addr 为内存地址，size 为数据大小。
//
//   - AiFindPic(...) int32
//       AI 找图，返回图片索引，并通过 x,y 输出参数返回坐标。
//
//   - AiFindPicEx(...) string
//       AI 高级找图，返回所有匹配位置字符串（格式: "idx|x|y,idx|x|y,..."）。
//
//   - AiFindPicMem(...) int32
//       AI 内存找图，pic_info 为图片数据描述字符串。
//
//   - AiFindPicMemEx(...) string
//       AI 高级内存找图，返回所有匹配位置字符串。
//
// 使用说明:
//   1. LoadAi 只需在初始化后调用一次，多次调用会覆盖已加载模型。
//   2. AI 找图使用前必须成功调用 LoadAi 或 LoadAiMemory。
//   3. AI 找图不依赖 SetDict / SetPath，但 pic_name 仍受 SetPath 影响。
//   4. sim 参数取值范围 0.1 ~ 1.0（与 FindPic 一致），与 FindPicSim 的 0~100 不同。
//
// 编译说明:
//   必须使用32位编译: GOARCH=386 go build -o load_ai.exe
//
// 运行要求:
//   - 需要 xd47243.dll 和 Go.dll 在同一目录
//   - 需要 AI 模型文件（默认 model.dm，可由命令行参数指定）
//   - 部分功能需要管理员权限

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	dmsoft "github.com/yuan71058/dm72424-go"
)

// 常量定义：DLL 文件路径
const (
	DmPluginPath = "xd47243.dll" // 大漠插件 DLL 文件名
	CrackDllPath = "Go.dll"     // 破解 DLL 文件名（用于激活大漠插件）
	DefaultModel = "ai.module"  // 默认 AI 模型文件名（大漠专用 .module 格式）
	DefaultPic   = "target.bmp" // 默认查找图片文件名
)

// main 主函数：LoadAi 加载 AI 模型示例入口
//
// 执行流程:
//   1. 解析命令行参数（模型文件、查找图片）
//   2. LoadDm  - 加载大漠插件 DLL
//   3. CrackDm - 破解大漠插件
//   4. New + Init - 创建对象并初始化
//   5. Reg     - 注册大漠插件
//   6. LoadAi  - 加载 AI 模型（核心步骤）
//   7. AiFindPic / AiFindPicEx - 执行 AI 找图测试
//   8. Release - 释放资源
func main() {
	// ========== 命令行参数解析 ==========
	// 支持通过命令行指定 AI 模型和查找图片，便于灵活测试
	modelPath := flag.String("model", DefaultModel, "AI 模型文件路径（大漠专用格式）")
	picPath := flag.String("pic", DefaultPic, "查找图片文件路径")
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         大漠插件 LoadAi 加载AI模型 示例                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Printf("  AI 模型: %s\n", *modelPath)
	fmt.Printf("  查找图片: %s\n", *picPath)
	fmt.Println()

	// ========== 第一步：加载大漠插件 DLL ==========
	// LoadDm 内部通过 LoadLibrary 加载 xd47243.dll，返回模块句柄
	dmHModule, err := dmsoft.LoadDm(DmPluginPath)
	if err != nil {
		log.Fatalf("加载大漠插件失败: %v", err)
	}
	fmt.Printf("✓ 大漠插件加载成功，模块句柄: %v\n", dmHModule)

	// ========== 第二步：破解大漠插件 ==========
	// CrackDm 加载 Go.dll，通过补丁方式激活大漠插件
	err = dmsoft.CrackDm(CrackDllPath)
	if err != nil {
		log.Fatalf("破解大漠插件失败: %v", err)
	}
	fmt.Println("✓ 破解函数执行完成")

	// ========== 第三步：创建大漠插件对象 ==========
	dm := dmsoft.New()
	if dm == nil {
		log.Fatal("创建大漠对象失败")
	}
	err = dm.Init()
	if err != nil {
		log.Fatalf("初始化大漠对象失败: %v", err)
	}
	defer dm.Release() // 确保程序退出时释放资源
	fmt.Println("✓ 大漠对象创建并初始化完成")

	// ========== 第四步：注册大漠插件 ==========
	// Reg("", "") 使用免费注册码进行注册，返回 1 表示成功
	nret := dm.Reg("", "")
	if nret == 1 {
		fmt.Println("✓ 大漠注册成功")
	} else {
		log.Fatalf("大漠注册失败，返回值: %d", nret)
	}

	// ========== 第五步：关闭错误提示框 ==========
	// 避免 AI 找图失败时弹出对话框阻塞程序
	dm.SetShowErrorMsg(0)
	fmt.Println("✓ 已关闭错误提示框")

	// ========== 第六步：设置全局路径 ==========
	// 设置图片、字库等资源的查找目录（使用绝对路径避免歧义）
	absPath, _ := os.Getwd()
	dm.SetPath(absPath)
	fmt.Printf("✓ 已设置全局路径: %s\n", absPath)

	// ========== 第七步：加载 AI 模型 ==========
	// 【核心步骤】LoadAi 从文件加载大漠 AI 模型
	//
	// 官方返回值定义:
	//    1   成功
	//   -1   打开文件失败
	//   -2   内存初始化失败
	//   -3   参数错误
	//   -4   加载错误
	//   -5   Ai 模块初始化失败
	//   -6   内存分配失败
	//
	// 加载成功后，AiFindPic 系列函数方可正常工作
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第七步：LoadAi — 加载 AI 模型")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 校验模型文件存在性，避免无意义调用
	absModelPath, _ := filepath.Abs(*modelPath)
	if _, err := os.Stat(absModelPath); os.IsNotExist(err) {
		fmt.Printf("⚠ AI 模型文件不存在: %s\n", absModelPath)
		fmt.Println("  跳过 AI 找图测试（请放置模型文件后重试）")
		fmt.Println("\n========== 示例执行结束 ==========")
		return
	}

	loadRet := dm.LoadAi(absModelPath)
	// 输出 LoadAi 返回值及对应含义，便于调试与排查
	fmt.Printf("LoadAi 返回值: %d\n", loadRet)
	switch loadRet {
	case 1:
		fmt.Printf("✓ AI 模型加载成功: %s\n", absModelPath)
	case -1:
		fmt.Println("✗ 失败: 打开文件失败 (-1)")
		fmt.Println("  排查: 检查文件路径、权限、是否被占用")
	case -2:
		fmt.Println("✗ 失败: 内存初始化失败 (-2)")
		fmt.Println("  排查: 正式版本出现此错误可联系大漠作者解决")
	case -3:
		fmt.Println("✗ 失败: 参数错误 (-3)")
		fmt.Println("  排查: file 参数为空或类型不正确")
	case -4:
		fmt.Println("✗ 失败: 加载错误 (-4)")
		fmt.Println("  排查: 模型文件已损坏或格式不符")
	case -5:
		fmt.Println("✗ 失败: Ai 模块初始化失败 (-5)")
		fmt.Println("  排查: 模型与插件版本不匹配")
	case -6:
		fmt.Println("✗ 失败: 内存分配失败 (-6)")
		fmt.Println("  排查: 系统内存不足，关闭其他程序后重试")
	default:
		fmt.Printf("✗ 未知返回值: %d\n", loadRet)
	}

	if loadRet != 1 {
		fmt.Println("\n========== 示例执行结束 ==========")
		return
	}

	// ========== 第八步：AI 找图测试 ==========
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第八步：AiFindPic / AiFindPicEx — AI 找图测试")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 校验查找图片存在性
	absPicPath, _ := filepath.Abs(*picPath)
	if _, err := os.Stat(absPicPath); os.IsNotExist(err) {
		fmt.Printf("⚠ 查找图片不存在: %s\n", absPicPath)
		fmt.Println("  请放置目标图片后重试")
		fmt.Println("\n========== 示例执行结束 ==========")
		return
	}

	// 使用全屏区域作为查找范围（可根据实际需求调整）
	x1, y1 := int32(0), int32(0)
	x2 := dm.GetScreenWidth()
	y2 := dm.GetScreenHeight()
	fmt.Printf("查找区域: (%d, %d) - (%d, %d)\n", x1, y1, x2, y2)
	fmt.Printf("查找图片: %s\n", absPicPath)
	fmt.Printf("相似度:   %.2f\n", 0.8)

	// ---------- 方法1: AiFindPic ----------
	// AiFindPic 返回图片索引（>=0 表示找到，-1 表示未找到）
	// 通过 x, y 输出参数返回匹配位置坐标
	fmt.Println("\n--- 方法1: AiFindPic 基本 AI 找图 ---")
	findX, findY := int32(0), int32(0)
	aiRet := dm.AiFindPic(x1, y1, x2, y2, absPicPath, 0.8, 0, &findX, &findY)
	if aiRet >= 0 {
		fmt.Printf("✓ AiFindPic 找到图片！\n")
		fmt.Printf("  图片索引: %d\n", aiRet)
		fmt.Printf("  匹配位置: (%d, %d)\n", findX, findY)
	} else {
		fmt.Printf("AiFindPic 未找到图片，返回值: %d\n", aiRet)
	}

	// ---------- 方法2: AiFindPicEx ----------
	// AiFindPicEx 返回字符串，包含所有匹配位置
	// 格式: "idx|x|y,idx|x|y,..."（多个匹配以逗号分隔）
	fmt.Println("\n--- 方法2: AiFindPicEx 高级 AI 找图 ---")
	findExStr := dm.AiFindPicEx(x1, y1, x2, y2, absPicPath, 0.8, 0)
	if findExStr != "" {
		fmt.Printf("✓ AiFindPicEx 找到匹配项:\n")
		fmt.Printf("  %s\n", findExStr)
		fmt.Println("  说明: 格式为 idx|x|y（多个以逗号分隔）")
	} else {
		fmt.Printf("AiFindPicEx 未找到图片\n")
	}

	// ---------- 方法3: 多图查找（pic_name 用 | 分隔）----------
	// AiFindPic/AiFindPicEx 支持多图查找，pic_name 用 "|" 分隔多个图片
	// 返回值为匹配到的图片索引（从0开始）
	fmt.Println("\n--- 方法3: AiFindPic 多图查找 ---")
	multiPic := absPicPath + "|target2.bmp"
	multiRet := dm.AiFindPic(x1, y1, x2, y2, multiPic, 0.8, 0, &findX, &findY)
	if multiRet >= 0 {
		fmt.Printf("✓ 多图查找命中，索引: %d, 位置: (%d, %d)\n", multiRet, findX, findY)
	} else {
		fmt.Printf("多图查找未命中，返回值: %d\n", multiRet)
	}

	// ========== 第九步：释放图片资源 ==========
	// FindPic / AiFindPic 系列函数会缓存图片，使用完毕后建议释放
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第九步：FreePic — 释放图片资源")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	freeRet := dm.FreePic(absPicPath)
	fmt.Printf("✓ FreePic 返回: %d (1=成功)\n", freeRet)

	fmt.Println("\n========== 示例执行结束 ==========")
	fmt.Println("使用建议:")
	fmt.Println("  1. LoadAi 只需调用一次，多次调用会覆盖已加载模型")
	fmt.Println("  2. AI 找图 sim 参数范围 0.1~1.0，与 FindPic 一致")
	fmt.Println("  3. 多图查找使用 '|' 分隔图片路径，返回值为匹配索引")
	fmt.Println("  4. 使用完毕后调用 FreePic 释放图片缓存")
	fmt.Println("  5. 默认模型文件 ai.module 随项目附带，可直接运行测试")
}
