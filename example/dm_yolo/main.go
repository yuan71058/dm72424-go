// example/dm_yolo/main.go - 大漠插件原生 YOLO 目标检测 Demo
//
// 版本: v1.0.0
// 更新日期: 2026-07-10
//
// 功能说明:
//   本示例演示使用大漠插件内置的 YOLO 推理能力进行目标检测,
//   完全不依赖 Python, 纯 Go + 大漠 DLL 实现。
//
//   核心 API 调用链:
//     LoadAi            -> 加载 AI 模块 (ai.module)
//     AiYoloSetVersion  -> 设置 YOLO 版本 (目前仅支持 "v5-7.0")
//     AiYoloSetModel    -> 加载模型文件 (.onnx 或 .dmx)
//     AiYoloUseModel    -> 切换当前使用的模型
//     AiYoloDetectObjects     -> 检测目标, 返回 "类名,置信度,x,y,w,h|..."
//     AiYoloDetectObjectsToFile-> 检测并保存标注图到 BMP
//     AiYoloSortsObjects -> 排序检测结果 (从上到下, 从左到右)
//     AiYoloObjectsToString -> 提取类名连接字符串
//     AiYoloFreeModel   -> 释放模型
//
//   检测结果格式:
//     "类名,置信度,x,y,w,h|类名,置信度,x,y,w,h|..."
//     其中 x,y 为检测框左上角坐标, w,h 为宽高, 均为相对绑定窗口的坐标
//     中心点坐标 = (x + w/2, y + h/2)
//
// 模型文件说明:
//   大漠 YOLO 支持两种模型格式:
//   1. .onnx - 需在同目录放置同名 .class 文件 (类别名列表)
//      可用 yolov5-7.0/export.py 导出: python export.py --weights yolov5s.pt --simplify --include onnx
//      .class 文件内容为每行一个类名, 如 COCO 的 80 个类别
//   2. .dmx  - 大漠专有加密格式, 需提供密码, 由大漠 YOLO 综合工具生成
//
// 编译说明:
//   必须使用 32 位编译: GOARCH=386 go build -o dm_yolo.exe
//
// 运行要求:
//   - xd47243.dll 和 Go.dll 在同一目录
//   - ai.module (AI模块, 可从 example/load_ai/ 复制)
//   - YOLO 模型文件 (.onnx+.class 或 .dmx)
//
// 用法:
//   dm_yolo.exe                                        # 使用默认 onnx 模型全屏检测
//   dm_yolo.exe -model best.onnx                       # 指定 onnx 模型
//   dm_yolo.exe -model best.dmx -pwd 123456            # 指定加密 dmx 模型
//   dm_yolo.exe -x1 0 -y1 0 -x2 800 -y2 600            # 指定检测区域
//   dm_yolo.exe -save result.bmp                       # 检测并保存标注图
//   dm_yolo.exe -move                                  # 检测后移动鼠标到首个目标
//
// 模型准备 (仅需一次):
//   cd yolov5-7.0
//   python export.py --weights yolov5s.pt --simplify --include onnx   # 导出 onnx (可能生成 .onnx.data)
//   python merge_onnx.py                                              # 合并为单文件 onnx (大漠要求)
//   python gen_class.py                                               # 生成 .class 类名文件

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	dmsoft "github.com/yuan71058/dm72424-go"
	"golang.org/x/sys/windows"
)

// 常量定义
const (
	DmPluginPath   = "xd47243.dll"                         // 大漠插件 DLL
	CrackDllPath   = "Go.dll"                              // 破解 DLL
	DefaultAiMod   = "../load_ai/ai.module"                // 默认 AI 模块 (复用 load_ai 示例)
	DefaultVersion = "v5-7.0"                              // YOLO 版本 (目前仅支持此值)
	DefaultModel   = "yolov5-7.0/yolov5s_op12.onnx"     // 默认 ONNX 模型 (opset12/IR7, 大漠兼容)
)

// YoloObject 单个检测目标
type YoloObject struct {
	Class      string  // 类名 (或 class id)
	Confidence float64 // 置信度
	X          int     // 检测框左上角 X
	Y          int     // 检测框左上角 Y
	W          int     // 检测框宽度
	H          int     // 检测框高度
	CenterX    int     // 中心点 X = X + W/2
	CenterY    int     // 中心点 Y = Y + H/2
}

func main() {
	// ========== 命令行参数 ==========
	modelFile := flag.String("model", DefaultModel, "YOLO 模型文件 (.onnx 需同名 .class, 或 .dmx)")
	modelPwd := flag.String("pwd", "", "模型密码 (仅 .dmx 格式需要)")
	aiModule := flag.String("aimod", DefaultAiMod, "AI 模块文件 (ai.module)")
	version := flag.String("ver", DefaultVersion, "YOLO 版本 (目前仅支持 v5-7.0)")
	prob := flag.Float64("prob", 0.5, "置信度阈值 (超过此值才检测)")
	iou := flag.Float64("iou", 0.45, "NMS IoU 阈值 (建议 0.4-0.6)")
	x1 := flag.Int("x1", 0, "检测区域左上角 X")
	y1 := flag.Int("y1", 0, "检测区域左上角 Y")
	x2 := flag.Int("x2", 0, "检测区域右下角 X (0=屏幕宽度)")
	y2 := flag.Int("y2", 0, "检测区域右下角 Y (0=屏幕高度)")
	saveBmp := flag.String("save", "", "保存标注图到指定 BMP 文件 (空=不保存)")
	moveMouse := flag.Bool("move", false, "检测后移动鼠标到首个目标中心")
	loop := flag.Bool("loop", false, "持续循环检测 (按 Ctrl+C 退出)")
	scanLines := flag.Bool("sort", false, "对检测结果排序 (从上到下,从左到右)")
	imageFile := flag.String("image", "", "对本地图片文件做检测 (用系统图片查看器打开后绑定窗口检测, 标注图保存到根目录)")
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║       大漠插件 原生 YOLO 目标检测 Demo                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Printf("  AI 模块:   %s\n", *aiModule)
	fmt.Printf("  YOLO 版本: %s\n", *version)
	fmt.Printf("  模型文件:  %s\n", strOrEmpty(*modelFile))
	fmt.Printf("  置信度:    %.2f\n", *prob)
	fmt.Printf("  IoU:       %.2f\n", *iou)
	fmt.Println()

	// 校验模型文件
	absModel, _ := filepath.Abs(*modelFile)
	if _, err := os.Stat(absModel); os.IsNotExist(err) {
		log.Fatalf("模型文件不存在: %s\n  获取 onnx: cd yolov5-7.0 && python export.py --weights yolov5s.pt --simplify --include onnx", absModel)
	}
	// onnx 格式: 大漠 7.2336+ 不再需要 .class 文件, 仅旧版本需要. 此处仅提示.
	if strings.HasSuffix(strings.ToLower(absModel), ".onnx") {
		classFile := absModel[:len(absModel)-5] + ".class"
		if _, err := os.Stat(classFile); os.IsNotExist(err) {
			fmt.Printf("⚠ 提示: 未找到同名 .class 文件 (%s)\n", filepath.Base(classFile))
			fmt.Println("  大漠 7.2336+ 的 onnx 格式不需要 .class; 旧版本可运行 gen_class.py 生成")
		}
	}

	// ========== 第一步: 初始化大漠插件 ==========
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第一步: 初始化大漠插件")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	dmHModule, err := dmsoft.LoadDm(DmPluginPath)
	if err != nil {
		log.Fatalf("加载大漠插件失败: %v", err)
	}
	fmt.Printf("✓ 大漠插件加载成功, 句柄: %v\n", dmHModule)

	if err := dmsoft.CrackDm(CrackDllPath); err != nil {
		log.Fatalf("破解大漠插件失败: %v", err)
	}
	fmt.Println("✓ 破解函数执行完成")

	dm := dmsoft.New()
	if dm == nil {
		log.Fatal("创建大漠对象失败")
	}
	if err := dm.Init(); err != nil {
		log.Fatalf("初始化大漠对象失败: %v", err)
	}
	defer dm.Release()
	fmt.Println("✓ 大漠对象创建并初始化完成")

	if dm.Reg("", "") != 1 {
		log.Fatal("大漠注册失败")
	}
	fmt.Println("✓ 大漠注册成功")
	fmt.Printf("✓ 大漠版本: %s\n", dm.Ver())

	dm.SetShowErrorMsg(1) // 开启错误弹窗, 便于调试 AiYoloSetModel 失败原因

	workDir, _ := os.Getwd()
	dm.SetPath(workDir)
	fmt.Printf("✓ 工作路径: %s\n", workDir)

	screenW := dm.GetScreenWidth()
	screenH := dm.GetScreenHeight()
	fmt.Printf("✓ 屏幕分辨率: %d x %d\n", screenW, screenH)

	// 计算检测区域
	rx2 := int32(*x2)
	ry2 := int32(*y2)
	if rx2 == 0 {
		rx2 = screenW
	}
	if ry2 == 0 {
		ry2 = screenH
	}
	fmt.Printf("✓ 检测区域: (%d,%d) - (%d,%d)\n", *x1, *y1, rx2, ry2)
	fmt.Println()

	// ========== 第二步: 加载 AI 模块 ==========
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第二步: LoadAi 加载 AI 模块")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	absAiMod, _ := filepath.Abs(*aiModule)
	if _, err := os.Stat(absAiMod); os.IsNotExist(err) {
		log.Fatalf("AI 模块不存在: %s\n  可从 example/load_ai/ 复制 ai.module", absAiMod)
	}

	loadRet := dm.LoadAi(absAiMod)
	if loadRet != 1 {
		log.Fatalf("LoadAi 失败, 返回值: %d\n  -1=打开失败 -2=内存初始化失败 -3=参数错误 -4=加载错误 -5=Ai模块初始化失败 -6=内存分配失败", loadRet)
	}
	fmt.Printf("✓ AI 模块加载成功: %s\n", absAiMod)
	fmt.Println()

	// ========== 第三步: 配置 YOLO ==========
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  第三步: AiYoloSetVersion + AiYoloSetModel + AiYoloUseModel")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 设置 YOLO 版本 (必须在 LoadAi 后第一时间调用)
	verRet := dm.AiYoloSetVersion(*version)
	if verRet != 1 {
		log.Fatalf("AiYoloSetVersion 失败, 返回值: %d (目前仅支持 v5-7.0)", verRet)
	}
	fmt.Printf("✓ AiYoloSetVersion(\"%s\") = 1\n", *version)

	// 加载模型 (index=0)
	// 大漠文档: AiYoloSetModel 的 file 参数是"模型文件名"(相对路径), 配合 SetPath 使用.
	// 大漠 7.2336+ 的 onnx 格式从模型 metadata 读取类名, 不再需要 .class 文件.
	// 因此导出 onnx 时必须写入 stride 和 names metadata (见 export_op12.py).
	modelDir := filepath.Dir(absModel)
	modelName := filepath.Base(absModel)
	dm.SetPath(modelDir)
	fmt.Printf("  → SetPath: %s\n", modelDir)
	fmt.Printf("  → 模型文件名: %s\n", modelName)

	setModelRet := dm.AiYoloSetModel(0, modelName, *modelPwd)
	if setModelRet != 1 {
		lastErr := dm.GetLastError()
		fmt.Printf("✗ AiYoloSetModel 失败, 返回值: %d, GetLastError: %d\n", setModelRet, lastErr)
		fmt.Println("  排查建议:")
		fmt.Println("    1. onnx 必须含 stride 和 names metadata (大漠 7.2336+ 从 metadata 读类名)")
		fmt.Println("    2. onnx 需 opset<=12, IR<=7 (大漠内置 onnxruntime 版本较低)")
		fmt.Println("    3. 路径不要含中文; 确保文件存在且未被占用")
		fmt.Println("    4. .dmx 需正确密码; 老版 dmx 需配合老版 ai.module")
		log.Fatalf("  AiYoloSetModel 无法加载模型")
	}
	fmt.Printf("✓ AiYoloSetModel(0, %s) = 1\n", modelName)

	// 切换使用模型 0
	useRet := dm.AiYoloUseModel(0)
	if useRet != 1 {
		log.Fatalf("AiYoloUseModel 失败, 返回值: %d", useRet)
	}
	fmt.Println("✓ AiYoloUseModel(0) = 1")
	fmt.Println()

	// 确保退出时释放模型
	defer func() {
		freeRet := dm.AiYoloFreeModel(0)
		fmt.Printf("\n✓ AiYoloFreeModel(0) = %d (1=成功)\n", freeRet)
	}()

	// ========== 图片文件检测模式 ==========
	// 大漠 YOLO 只能对屏幕区域检测, 没有"从图片文件检测"的接口.
	// 这里用系统图片查看器打开图片 -> 大漠绑定窗口 -> 对客户区检测 -> 保存标注图到 demo 根目录.
	if *imageFile != "" {
		detectImageFile(dm, *imageFile, *prob, *iou, *scanLines)
		return
	}

	// ========== 第四步: YOLO 检测 ==========
	detectOnce := func() []YoloObject {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("  第四步: AiYoloDetectObjects 目标检测")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		start := time.Now()
		objectsStr := dm.AiYoloDetectObjects(int32(*x1), int32(*y1), rx2, ry2, float32(*prob), float32(*iou))
		elapsed := time.Since(start)

		if objectsStr == "" {
			fmt.Printf("  未检测到任何目标 (耗时 %v)\n", elapsed)
			return nil
		}

		objects := parseYoloObjects(objectsStr)
		fmt.Printf("✓ 检测到 %d 个目标 (耗时 %v)\n\n", len(objects), elapsed)

		// 可选: 排序
		if *scanLines {
			sortedStr := dm.AiYoloSortsObjects(objectsStr, 30)
			objects = parseYoloObjects(sortedStr)
			fmt.Println("  (已按从上到下、从左到右排序)")
		}

		// 打印检测结果
		fmt.Printf("  %-4s %-16s %-8s %-18s %-12s\n", "序号", "类名", "置信度", "边界框(x,y,w,h)", "中心点")
		fmt.Println("  " + strings.Repeat("-", 70))
		for i, o := range objects {
			bbox := fmt.Sprintf("(%d,%d,%d,%d)", o.X, o.Y, o.W, o.H)
			center := fmt.Sprintf("(%d,%d)", o.CenterX, o.CenterY)
			fmt.Printf("  %-4d %-16s %-8.4f %-18s %-12s\n",
				i+1, o.Class, o.Confidence, bbox, center)
		}
		fmt.Println()

		// 打印类名连接字符串
		classStr := dm.AiYoloObjectsToString(objectsStr)
		if classStr != "" {
			fmt.Printf("  类名摘要: %s\n\n", classStr)
		}

		// 可选: 保存标注图
		if *saveBmp != "" {
			saveRet := dm.AiYoloDetectObjectsToFile(int32(*x1), int32(*y1), rx2, ry2, float32(*prob), float32(*iou), *saveBmp, 0)
			if saveRet == 1 {
				absSave, _ := filepath.Abs(*saveBmp)
				fmt.Printf("✓ 标注图已保存: %s\n", absSave)
			} else {
				fmt.Printf("✗ 保存标注图失败, 返回值: %d\n", saveRet)
			}
		}

		return objects
	}

	objects := detectOnce()

	// ========== 第五步: 联动操作 ==========
	if *moveMouse && len(objects) > 0 {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("  第五步: 移动鼠标到首个目标中心")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		o := objects[0]
		fmt.Printf("  目标: %s (置信度 %.4f)\n", o.Class, o.Confidence)
		fmt.Printf("  中心点: (%d, %d)\n", o.CenterX, o.CenterY)

		origX, origY := int32(0), int32(0)
		dm.GetCursorPos(&origX, &origY)
		fmt.Printf("  原鼠标位置: (%d, %d)\n", origX, origY)

		fmt.Printf("  → 移动鼠标到 (%d, %d) ...\n", o.CenterX, o.CenterY)
		if dm.MoveTo(int32(o.CenterX), int32(o.CenterY)) == 1 {
			fmt.Println("  ✓ 鼠标移动成功")
		} else {
			fmt.Println("  ✗ 鼠标移动失败")
		}
		time.Sleep(1500 * time.Millisecond)
		dm.MoveTo(origX, origY)
		fmt.Printf("  ✓ 已恢复鼠标到原位置\n", origX, origY)
	}

	// ========== 循环检测模式 ==========
	if *loop {
		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("  循环检测模式 (每秒一次, Ctrl+C 退出)")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for {
			time.Sleep(time.Second)
			objectsStr := dm.AiYoloDetectObjects(int32(*x1), int32(*y1), rx2, ry2, float32(*prob), float32(*iou))
			ts := time.Now().Format("15:04:05")
			if objectsStr == "" {
				fmt.Printf("  [%s] 未检测到目标\n", ts)
				continue
			}
			objs := parseYoloObjects(objectsStr)
			fmt.Printf("  [%s] 检测到 %d 个目标: ", ts, len(objs))
			for i, o := range objs {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s(%.2f)@(%d,%d)", o.Class, o.Confidence, o.CenterX, o.CenterY)
			}
			fmt.Println()
		}
	}

	fmt.Println("\n========== Demo 执行结束 ==========")
}

// parseYoloObjects 解析大漠 YOLO 返回的字符串
// 格式: "类名,置信度,x,y,w,h|类名,置信度,x,y,w,h|..."
func parseYoloObjects(s string) []YoloObject {
	if s == "" {
		return nil
	}
	items := strings.Split(s, "|")
	objects := make([]YoloObject, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		parts := strings.Split(item, ",")
		if len(parts) < 6 {
			continue
		}
		conf, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		x, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
		w, _ := strconv.Atoi(strings.TrimSpace(parts[4]))
		h, _ := strconv.Atoi(strings.TrimSpace(parts[5]))
		objects = append(objects, YoloObject{
			Class:      strings.TrimSpace(parts[0]),
			Confidence: conf,
			X:          x,
			Y:          y,
			W:          w,
			H:          h,
			CenterX:    x + w/2,
			CenterY:    y + h/2,
		})
	}
	return objects
}

// strOrEmpty 空字符串显示占位符
func strOrEmpty(s string) string {
	if s == "" {
		return "(未指定)"
	}
	return s
}

// getClientRect 直接调用 Windows API GetClientRect 获取窗口客户区大小.
// 用于大漠 GetClientSize 对某些窗口返回 0 时的兜底方案.
// 返回值: 宽度, 高度, 是否成功.
var (
	modUser32     = syscall.NewLazyDLL("user32.dll")
	procGetClient = modUser32.NewProc("GetClientRect")
	procShowWin   = modUser32.NewProc("ShowWindow")
)

func getClientRect(hwnd int32) (int32, int32, bool) {
	var rect windows.Rect
	ret, _, _ := procGetClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return 0, 0, false
	}
	w := rect.Right - rect.Left
	h := rect.Bottom - rect.Top
	if w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return int32(w), int32(h), true
}

// detectImageFile 对本地图片文件做 YOLO 检测
// 大漠 YOLO 只能对屏幕区域检测, 这里用系统图片查看器打开图片 -> 绑定窗口 -> 检测客户区 -> 保存标注图.
func detectImageFile(dm *dmsoft.DmSoft, imageFile string, prob, iou float64, sortObjs bool) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  图片文件检测模式")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	absImage, _ := filepath.Abs(imageFile)
	if _, err := os.Stat(absImage); os.IsNotExist(err) {
		log.Fatalf("图片文件不存在: %s", absImage)
	}
	fmt.Printf("  图片路径: %s\n", absImage)

	// 1. 用系统图片查看器 (Windows 照片查看器) 打开图片
	//    通过 rundll32 调用 shimgvw.dll, 以全屏图片查看器方式打开.
	//    窗口类名固定为 "Photo_Lightweight_Viewer", 标题为文件名.
	imgBase := filepath.Base(absImage)
	fmt.Printf("  → 启动系统图片查看器打开图片 ...\n")
	cmd := exec.Command("rundll32.exe", "shimgvw.dll,ImageView_Fullscreen", absImage)
	if err := cmd.Start(); err != nil {
		log.Fatalf("启动图片查看器失败: %v", err)
	}
	viewerProc := cmd.Process
	defer func() {
		if viewerProc != nil {
			viewerProc.Kill()
		}
	}()

	// 等待图片查看器完全加载和渲染
	fmt.Println("  → 等待图片查看器加载 (3秒) ...")
	time.Sleep(3 * time.Second)

	// 2. 等待图片查看器窗口出现并查找它 (优先按窗口类名查找, 兜底按标题包含文件名)
	const viewerClass = "Photo_Lightweight_Viewer"
	var hwnd int32 = -1
	for i := 0; i < 40; i++ { // 最多等 4 秒
		time.Sleep(100 * time.Millisecond)
		// 先按类名查找 (图片查看器窗口类名固定)
		hwnd = dm.FindWindow(viewerClass, "")
		if hwnd != 0 {
			break
		}
		// 兜底: 按标题包含文件名查找
		hwnd = dm.FindWindow("", imgBase)
		if hwnd != 0 {
			break
		}
	}
	if hwnd == 0 {
		log.Fatalf("未找到图片查看器窗口 (类名 %q 或标题含 %q)", viewerClass, imgBase)
	}
	fmt.Printf("✓ 找到图片查看器窗口, 句柄: %d\n", hwnd)

	// 3. 最大化图片查看器窗口, 确保图片以最大尺寸显示, 提升小目标检测率
	//    大漠 SetWindowState 参数: 3=最小化, 4=最大化
	dm.SetWindowState(hwnd, 4)
	fmt.Println("✓ 已最大化图片查看器窗口")

	// 4. 获取窗口客户区大小 (最大化后需等待窗口稳定, 循环直到非零)
	//    优先用大漠 GetClientSize; 若该接口对当前窗口返回 0 (部分系统/窗口会失败),
	//    则兜底直接调用 Windows API (GetClientRect) 获取客户区.
	var cliW, cliH int32
	for i := 0; i < 20; i++ {
		time.Sleep(200 * time.Millisecond)
		if dm.GetClientSize(hwnd, &cliW, &cliH) == 1 && cliW > 0 && cliH > 0 {
			break
		}
		// 兜底: 直接调用 Windows API 获取客户区大小
		if fw, fh, ok := getClientRect(hwnd); ok {
			cliW, cliH = fw, fh
			break
		}
	}
	if cliW <= 0 || cliH <= 0 {
		log.Fatalf("获取客户区大小失败: %d x %d (窗口句柄 %d)", cliW, cliH, hwnd)
	}
	fmt.Printf("✓ 客户区大小: %d x %d\n", cliW, cliH)

	// 5. 绑定窗口 (display=gdi 避免抢占前台; mouse/keypad 禁用避免干扰)
	if dm.BindWindow(hwnd, "gdi", "windows", "windows", 0) != 1 {
		log.Fatal("BindWindow 失败")
	}
	fmt.Println("✓ 已绑定图片查看器窗口 (gdi 模式)")
	defer dm.UnBindWindow()

	// 等待窗口画面渲染稳定
	time.Sleep(500 * time.Millisecond)

	// 6. 对客户区做 YOLO 检测
	start := time.Now()
	objectsStr := dm.AiYoloDetectObjects(0, 0, cliW, cliH, float32(prob), float32(iou))
	elapsed := time.Since(start)

	if objectsStr == "" {
		fmt.Printf("  未检测到任何目标 (耗时 %v)\n", elapsed)
		return
	}

	objects := parseYoloObjects(objectsStr)
	fmt.Printf("✓ 检测到 %d 个目标 (耗时 %v)\n\n", len(objects), elapsed)

	if sortObjs {
		objects = parseYoloObjects(dm.AiYoloSortsObjects(objectsStr, 30))
		fmt.Println("  (已按从上到下、从左到右排序)")
	}

	fmt.Printf("  %-4s %-16s %-8s %-18s %-12s\n", "序号", "类名", "置信度", "边界框(x,y,w,h)", "中心点")
	fmt.Println("  " + strings.Repeat("-", 70))
	for i, o := range objects {
		bbox := fmt.Sprintf("(%d,%d,%d,%d)", o.X, o.Y, o.W, o.H)
		center := fmt.Sprintf("(%d,%d)", o.CenterX, o.CenterY)
		fmt.Printf("  %-4d %-16s %-8.4f %-18s %-12s\n",
			i+1, o.Class, o.Confidence, bbox, center)
	}
	fmt.Println()

	classStr := dm.AiYoloObjectsToString(objectsStr)
	if classStr != "" {
		fmt.Printf("  类名摘要: %s\n\n", classStr)
	}

	// 6. 保存标注图到 demo 根目录
	demoRoot, _ := os.Getwd()
	outBmp := filepath.Join(demoRoot, "yolo_result.bmp")
	saveRet := dm.AiYoloDetectObjectsToFile(0, 0, cliW, cliH, float32(prob), float32(iou), outBmp, 0)
	if saveRet == 1 {
		fmt.Printf("✓ 标注图已保存到 demo 根目录: %s\n", outBmp)
	} else {
		fmt.Printf("✗ 保存标注图失败, 返回值: %d\n", saveRet)
	}

	fmt.Println("\n========== 图片检测结束 ==========")
}
