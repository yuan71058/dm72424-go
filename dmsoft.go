// dmsoft.go - 大漠插件接口定义
// 对应原C++项目中的 obj.h
// 本文件由脚本自动生成
//
// 功能说明：
//   本文件定义了大漠插件(DmSoft)的所有接口方法，包括：
//   - 窗口操作：窗口查找、绑定、状态获取等
//   - 鼠标键盘：鼠标移动、点击、键盘按键等
//   - 图像处理：截图、找图、找色、OCR等
//   - 内存操作：进程内存读写、地址查找等
//   - 系统信息：获取硬件信息、系统时间等
//   - 文字识别：OCR识别、字库操作等
//   - 文件操作：文件读写、路径操作等
//
// 使用说明：
//   1. 首先调用 Load() 加载大漠插件DLL
//   2. 调用 New() 创建大漠对象
//   3. 调用 Init() 初始化对象
//   4. 调用 Reg() 或 RegEx() 进行注册
//   5. 调用 BindWindow() 绑定目标窗口
//   6. 使用各种功能方法
//   7. 最后调用 Release() 释放对象
//
// 注意事项：
//   - 大漠插件是32位DLL，需要使用32位编译 (GOARCH=386)
//   - 部分功能需要管理员权限
//   - 绑定窗口后才能进行截图、取色等操作
//   - 使用完毕后需要调用 Release() 释放资源

package dmsoft

import (
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// DmSoftInterface 定义大漠插件的所有接口方法
// 对应C++中的 dmsoft 类
//
// 接口分类：
//  1. 版本信息：Ver, GetID, GetLastError 等
//  2. 注册认证：Reg, RegEx, RegNoMac, RegExNoMac 等
//  3. 窗口操作：BindWindow, UnBindWindow, FindWindow, GetWindow 等
//  4. 鼠标操作：MoveTo, LeftClick, RightClick, GetCursorPos 等
//  5. 键盘操作：KeyPress, KeyDown, KeyUp, GetKeyState 等
//  6. 图像处理：Capture, FindPic, FindColor, GetColor 等
//  7. 文字识别：Ocr, FindStr, UseDict, AddDict 等
//  8. 内存操作：ReadInt, WriteInt, FindData, VirtualAllocEx 等
//  9. 系统信息：GetScreenWidth, GetScreenHeight, GetTime, GetOsType 等
//  10. 文件操作：ReadFile, WriteFile, CreateFolder, DeleteFile 等
type DmSoftInterface interface {
	// ==================== 磁盘与硬件信息 ====================

	// GetDiskReversion 获取磁盘版本信息
	// 参数: index - 磁盘索引
	// 返回: 磁盘版本字符串
	GetDiskReversion(index int32) string

	// GetDiskSerial 获取磁盘序列号
	// 参数: index - 磁盘索引
	// 返回: 磁盘序列号字符串
	GetDiskSerial(index int32) string

	// GetDiskModel 获取磁盘型号
	// 参数: index - 磁盘索引
	// 返回: 磁盘型号字符串
	GetDiskModel(index int32) string

	// ==================== AI与内存操作 ====================

	// LoadAiMemory 从内存加载AI模型
	// 参数: addr - 内存地址, size - 数据大小
	// 返回: 1成功, 0失败
	LoadAiMemory(addr int32, size int32) int32

	// LoadAi 从文件加载AI模型
	// 参数: file - AI模型文件路径
	// 返回: 1成功, 0失败
	LoadAi(file string) int32

	// ==================== FAQ相关 ====================

	// FaqSend 发送FAQ请求
	// 参数: server - 服务器地址, handle - 句柄, request_type - 请求类型, time_out - 超时时间
	// 返回: 结果字符串
	FaqSend(server string, handle int32, request_type int32, time_out int32) string

	// FaqGetSize 获取FAQ大小
	// 参数: handle - 句柄
	// 返回: 大小值
	FaqGetSize(handle int32) int32

	// FaqCapture FAQ截图
	// 参数: x1,y1,x2,y2 - 区域坐标, quality - 质量, delay - 延迟, time - 时间
	// 返回: 句柄
	FaqCapture(x1 int32, y1 int32, x2 int32, y2 int32, quality int32, delay int32, time int32) int32

	// FaqCaptureFromFile 从文件FAQ截图
	// 参数: x1,y1,x2,y2 - 区域坐标, file - 文件路径, quality - 质量
	// 返回: 句柄
	FaqCaptureFromFile(x1 int32, y1 int32, x2 int32, y2 int32, file string, quality int32) int32

	// FaqPost FAQ提交
	// 参数: server - 服务器地址, handle - 句柄, request_type - 请求类型, time_out - 超时时间
	// 返回: 1成功, 0失败
	FaqPost(server string, handle int32, request_type int32, time_out int32) int32

	// FaqIsPosted FAQ是否已提交
	// 返回: 1已提交, 0未提交
	FaqIsPosted() int32

	// FaqCancel 取消FAQ
	// 返回: 1成功, 0失败
	FaqCancel() int32

	// ==================== 图片查找 ====================

	// FindPic 查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: -1失败, 其他为找到的图片索引
	FindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32

	// FindPicEx 高级查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	FindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string

	// FindPicE 查找图片(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度, dir - 方向
	// 返回: "x,y" 格式字符串
	FindPicE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string

	// FindPicS 查找图片(返回图片索引)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 图片索引字符串
	FindPicS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) string

	// FindPicSim 查找图片相似度
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindPicSim(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32

	// FindPicSimE 查找图片相似度(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向
	// 返回: "索引|x|y" 格式字符串
	FindPicSimE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string

	// FindPicSimEx 高级查找图片相似度
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向
	// 返回: 所有找到的坐标字符串
	FindPicSimEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string

	// FindPicSimMem 在内存中查找图片(相似度)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_info - 图片信息, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindPicSimMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32

	// FindPicSimMemE 在内存中查找图片(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_info - 图片信息, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向
	// 返回: "x,y" 格式字符串
	FindPicSimMemE(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string

	// FindPicSimMemEx 高级在内存中查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_info - 图片信息, delta_color - 颜色偏差, sim - 相似度(0-100), dir - 方向
	// 返回: 所有找到的坐标字符串
	FindPicSimMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string

	// FindPicExS 高级查找图片(返回详细字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, delta_color - 颜色偏差, sim - 相似度, dir - 方向
	// 返回: 详细结果字符串
	FindPicExS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string

	// LoadPic 加载图片到内存
	// 参数: pic_name - 图片名称(支持多个,用|分隔)
	// 返回: 1成功, 0失败
	LoadPic(pic_name string) int32

	// LoadPicByte 从内存加载图片
	// 参数: addr - 内存地址, size - 数据大小, name - 图片名称
	// 返回: 1成功, 0失败
	LoadPicByte(addr int32, size int32, name string) int32

	// ==================== AI图片查找 ====================

	// AiFindPic AI查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	AiFindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32, x *int32, y *int32) int32

	// AiFindPicEx AI高级查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	AiFindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32) string

	// AiFindPicMem AI内存查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_info - 图片信息, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	AiFindPicMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32, x *int32, y *int32) int32

	// AiFindPicMemEx AI高级内存查找图片
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_info - 图片信息, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	AiFindPicMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32) string

	// AiEnableFindPicWindow 启用AI找图窗口
	// 参数: enable - 1启用, 0禁用
	// 返回: 1成功
	AiEnableFindPicWindow(enable int32) int32

	// ==================== YOLO相关 ====================

	// AiYoloSetVersion 设置YOLO版本
	// 参数: ver - 版本字符串
	// 返回: 1成功, 0失败
	AiYoloSetVersion(ver string) int32

	// AiYoloSetModel 设置YOLO模型
	// 参数: index - 模型索引, file - 模型文件路径, pwd - 密码
	// 返回: 1成功, 0失败
	AiYoloSetModel(index int32, file string, pwd string) int32

	// AiYoloSetModelMemory 从内存设置YOLO模型
	// 参数: index - 模型索引, addr - 内存地址, size - 数据大小, pwd - 密码
	// 返回: 1成功, 0失败
	AiYoloSetModelMemory(index int32, addr int32, size int32, pwd string) int32

	// AiYoloUseModel 使用YOLO模型
	// 参数: index - 模型索引
	// 返回: 1成功, 0失败
	AiYoloUseModel(index int32) int32

	// AiYoloDetectObjects YOLO检测对象
	// 参数: x1,y1,x2,y2 - 区域坐标, prob - 概率阈值, iou - IOU阈值
	// 返回: 检测结果字符串
	AiYoloDetectObjects(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32) string

	// AiYoloDetectObjectsToFile YOLO检测对象并保存到文件
	// 参数: x1,y1,x2,y2 - 区域坐标, prob - 概率阈值, iou - IOU阈值, file - 文件路径, mode - 模式
	// 返回: 1成功, 0失败
	AiYoloDetectObjectsToFile(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32, file string, mode int32) int32

	// AiYoloObjectsToString YOLO对象转字符串
	// 参数: objects - 对象字符串
	// 返回: 格式化后的字符串
	AiYoloObjectsToString(objects string) string

	// AiYoloSortsObjects YOLO排序对象
	// 参数: objects - 对象字符串, height - 高度
	// 返回: 排序后的字符串
	AiYoloSortsObjects(objects string, height int32) string

	// ==================== 版本与注册 ====================

	// Ver 获取大漠版本号
	// 返回: 版本字符串
	Ver() string

	// GetID 获取ID
	// 返回: ID值
	GetID() int32

	// GetMachineCode 获取机器码
	// 返回: 机器码字符串
	GetMachineCode() string

	// GetMachineCodeNoMac 获取机器码(不含MAC)
	// 返回: 机器码字符串
	GetMachineCodeNoMac() string

	// GetMac 获取MAC地址
	// 返回: MAC地址字符串
	GetMac() string

	// Reg 注册大漠
	// 参数: code - 注册码, ver - 版本
	// 返回: 1成功, 0失败
	Reg(code string, ver string) int32

	// RegEx 扩展注册
	// 参数: code - 注册码, ver - 版本, ip - IP地址
	// 返回: 1成功, 0失败
	RegEx(code string, ver string, ip string) int32

	// RegNoMac 注册(不含MAC)
	// 参数: code - 注册码, ver - 版本
	// 返回: 1成功, 0失败
	RegNoMac(code string, ver string) int32

	// GetLastError 获取最后错误码
	// 返回: 错误码
	GetLastError() int32

	// ==================== 路径与配置 ====================

	// SetPath 设置资源路径
	// 参数: path - 路径
	// 返回: 1成功, 0失败
	SetPath(path string) int32

	// GetPath 获取资源路径
	// 返回: 路径字符串
	GetPath() string

	// GetBasePath 获取基础路径
	// 返回: 基础路径字符串
	GetBasePath() string

	// GetRealPath 获取真实路径
	// 参数: path - 路径
	// 返回: 真实路径字符串
	GetRealPath(path string) string

	// ==================== 文字识别与OCR ====================

	// FindStrS 查找文字(返回字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度, x,y - 返回坐标
	// 返回: 找到的字符串
	FindStrS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) string

	// FindStr 查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度, x,y - 返回坐标
	// 返回: -1失败, 其他为找到的字符串索引
	FindStr(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32

	// FindStrEx 高级查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度
	// 返回: 所有找到的坐标字符串
	FindStrEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string

	// FindStrExS 高级查找文字(返回详细字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度
	// 返回: 详细结果字符串
	FindStrExS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string

	// FindStrFast 快速查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度, x,y - 返回坐标
	// 返回: -1失败, 其他为找到的字符串索引
	FindStrFast(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32

	// FindStrFastE 快速查找文字(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度
	// 返回: "x,y" 格式字符串
	FindStrFastE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string

	// FindStrFastEx 高级快速查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度
	// 返回: 所有找到的坐标字符串
	FindStrFastEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string

	// FindStrWithFontE 指定字体查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度, font_name - 字体名, font_size - 字号, flag - 标志
	// 返回: "x,y" 格式字符串
	FindStrWithFontE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string

	// FindStrWithFontEx 高级指定字体查找文字
	// 参数: x1,y1,x2,y2 - 区域坐标, str - 字符串, color - 颜色, sim - 相似度, font_name - 字体名, font_size - 字号, flag - 标志
	// 返回: 所有找到的坐标字符串
	FindStrWithFontEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string

	// Ocr OCR识别
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度
	// 返回: 识别结果字符串
	Ocr(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string

	// OcrExOne OCR识别单个
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度
	// 返回: 识别结果字符串
	OcrExOne(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string

	// OcrInFile 从文件OCR识别
	// 参数: x1,y1,x2,y2 - 区域坐标, pic_name - 图片名, color - 颜色, sim - 相似度
	// 返回: 识别结果字符串
	OcrInFile(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, color string, sim float64) string

	// GetWords 获取文字
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度
	// 返回: 文字字符串
	GetWords(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string

	// GetWordsNoDict 无字典获取文字
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色
	// 返回: 文字字符串
	GetWordsNoDict(x1 int32, y1 int32, x2 int32, y2 int32, color string) string

	// FetchWord 提取文字
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, word - 文字
	// 返回: 提取结果字符串
	FetchWord(x1 int32, y1 int32, x2 int32, y2 int32, color string, word string) string

	// GetWordResultCount 获取文字结果数量
	// 参数: str - 结果字符串
	// 返回: 数量
	GetWordResultCount(str string) int32

	// GetWordResultPos 获取文字结果位置
	// 参数: str - 结果字符串, index - 索引, x,y - 返回坐标
	// 返回: 1成功, 0失败
	GetWordResultPos(str string, index int32, x *int32, y *int32) int32

	// GetWordResultStr 获取文字结果字符串
	// 参数: str - 结果字符串, index - 索引
	// 返回: 文字字符串
	GetWordResultStr(str string, index int32) string

	// ==================== 字库操作 ====================

	// SetDict 设置字库
	// 参数: index - 字库索引, dict_name - 字库名
	// 返回: 1成功, 0失败
	SetDict(index int32, dict_name string) int32

	// UseDict 使用字库
	// 参数: index - 字库索引
	// 返回: 1成功, 0失败
	UseDict(index int32) int32

	// AddDict 添加字库
	// 参数: index - 字库索引, dict_info - 字库信息
	// 返回: 1成功, 0失败
	AddDict(index int32, dict_info string) int32

	// ClearDict 清除字库
	// 参数: index - 字库索引
	// 返回: 1成功, 0失败
	ClearDict(index int32) int32

	// SaveDict 保存字库
	// 参数: index - 字库索引, file - 文件路径
	// 返回: 1成功, 0失败
	SaveDict(index int32, file string) int32

	// GetDict 获取字库
	// 参数: index - 字库索引, font_index - 字体索引
	// 返回: 字库字符串
	GetDict(index int32, font_index int32) string

	// GetDictCount 获取字库数量
	// 参数: index - 字库索引
	// 返回: 数量
	GetDictCount(index int32) int32

	// GetNowDict 获取当前字库索引
	// 返回: 字库索引
	GetNowDict() int32

	// SetExportDict 设置导出字库
	// 参数: index - 字库索引, dict_name - 字库名
	// 返回: 1成功, 0失败
	SetExportDict(index int32, dict_name string) int32

	// EnableShareDict 启用共享字库
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableShareDict(en int32) int32

	// ==================== 字库参数设置 ====================

	// SetWordLineHeight 设置文字行高
	// 参数: line_height - 行高
	// 返回: 1成功
	SetWordLineHeight(line_height int32) int32

	// SetWordLineHeightNoDict 设置文字行高(无字典)
	// 参数: line_height - 行高
	// 返回: 1成功
	SetWordLineHeightNoDict(line_height int32) int32

	// SetMinColGap 设置最小列间距
	// 参数: col_gap - 列间距
	// 返回: 1成功
	SetMinColGap(col_gap int32) int32

	// SetColGapNoDict 设置列间距(无字典)
	// 参数: col_gap - 列间距
	// 返回: 1成功
	SetColGapNoDict(col_gap int32) int32

	// SetMinRowGap 设置最小行间距
	// 参数: row_gap - 行间距
	// 返回: 1成功
	SetMinRowGap(row_gap int32) int32

	// SetRowGapNoDict 设置行间距(无字典)
	// 参数: row_gap - 行间距
	// 返回: 1成功
	SetRowGapNoDict(row_gap int32) int32

	// SetExactOcr 设置精确OCR
	// 参数: exact_ocr - 1启用, 0禁用
	// 返回: 1成功
	SetExactOcr(exact_ocr int32) int32

	// ==================== 键盘操作 ====================

	// KeyPressChar 按键(字符)
	// 参数: key_str - 按键字符串
	// 返回: 1成功, 0失败
	KeyPressChar(key_str string) int32

	// KeyDownChar 按下按键(字符)
	// 参数: key_str - 按键字符串
	// 返回: 1成功, 0失败
	KeyDownChar(key_str string) int32

	// KeyUpChar 弹起按键(字符)
	// 参数: key_str - 按键字符串
	// 返回: 1成功, 0失败
	KeyUpChar(key_str string) int32

	// KeyPress 按键
	// 参数: vk - 虚拟键码
	// 返回: 1成功, 0失败
	KeyPress(vk int32) int32

	// KeyDown 按下按键
	// 参数: vk - 虚拟键码
	// 返回: 1成功, 0失败
	KeyDown(vk int32) int32

	// KeyUp 弹起按键
	// 参数: vk - 虚拟键码
	// 返回: 1成功, 0失败
	KeyUp(vk int32) int32

	// KeyPressStr 按键字符串
	// 参数: key_str - 按键字符串, delay - 延迟
	// 返回: 1成功, 0失败
	KeyPressStr(key_str string, delay int32) int32

	// WaitKey 等待按键
	// 参数: key_code - 键码, time_out - 超时
	// 返回: 1成功, 0失败
	WaitKey(key_code int32, time_out int32) int32

	// GetKeyState 获取按键状态
	// 参数: vk - 虚拟键码
	// 返回: 0弹起, 1按下
	GetKeyState(vk int32) int32

	// SetKeypadDelay 设置键盘延迟
	// 参数: type_ - 类型, delay - 延迟
	// 返回: 1成功
	SetKeypadDelay(type_ string, delay int32) int32

	// EnableKeypadSync 启用键盘同步
	// 参数: enable - 1启用, time_out - 超时
	// 返回: 1成功
	EnableKeypadSync(enable int32, time_out int32) int32

	// EnableKeypadMsg 启用键盘消息
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableKeypadMsg(en int32) int32

	// EnableRealKeypad 启用真实键盘
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableRealKeypad(en int32) int32

	// ==================== 鼠标操作 ====================

	// MoveTo 移动鼠标到指定位置
	// 参数: x - X坐标, y - Y坐标
	// 返回: 1成功, 0失败
	MoveTo(x int32, y int32) int32

	// MoveR 相对移动鼠标
	// 参数: rx - X偏移, ry - Y偏移
	// 返回: 1成功, 0失败
	MoveR(rx int32, ry int32) int32

	// LeftClick 左键单击
	// 返回: 1成功, 0失败
	LeftClick() int32

	// RightClick 右键单击
	// 返回: 1成功, 0失败
	RightClick() int32

	// MiddleClick 中键单击
	// 返回: 1成功, 0失败
	MiddleClick() int32

	// LeftDoubleClick 左键双击
	// 返回: 1成功, 0失败
	LeftDoubleClick() int32

	// LeftDown 左键按下
	// 返回: 1成功, 0失败
	LeftDown() int32

	// LeftUp 左键弹起
	// 返回: 1成功, 0失败
	LeftUp() int32

	// RightDown 右键按下
	// 返回: 1成功, 0失败
	RightDown() int32

	// RightUp 右键弹起
	// 返回: 1成功, 0失败
	RightUp() int32

	// MiddleDown 中键按下
	// 返回: 1成功, 0失败
	MiddleDown() int32

	// MiddleUp 中键弹起
	// 返回: 1成功, 0失败
	MiddleUp() int32

	// WheelUp 滚轮向上
	// 返回: 1成功, 0失败
	WheelUp() int32

	// WheelDown 滚轮向下
	// 返回: 1成功, 0失败
	WheelDown() int32

	// GetCursorPos 获取鼠标位置
	// 参数: x,y - 返回坐标
	// 返回: 1成功, 0失败
	GetCursorPos(x *int32, y *int32) int32

	// GetCursorShape 获取鼠标形状
	// 返回: 鼠标形状字符串
	GetCursorShape() string

	// SetMouseSpeed 设置鼠标速度
	// 参数: speed - 速度
	// 返回: 1成功
	SetMouseSpeed(speed int32) int32

	// EnableRealMouse 启用真实鼠标
	// 参数: en - 1启用, mousedelay - 延迟, mousestep - 步长
	// 返回: 1成功
	EnableRealMouse(en int32, mousedelay int32, mousestep int32) int32

	// EnableMouseSync 启用鼠标同步
	// 参数: enable - 1启用, time_out - 超时
	// 返回: 1成功
	EnableMouseSync(enable int32, time_out int32) int32

	// EnableMouseMsg 启用鼠标消息
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableMouseMsg(en int32) int32

	// EnableMouseAccuracy 启用鼠标精度
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableMouseAccuracy(en int32) int32

	// ==================== 颜色与取色 ====================

	// GetColor 获取颜色
	// 参数: x - X坐标, y - Y坐标
	// 返回: 颜色字符串(RGB)
	GetColor(x int32, y int32) string

	// GetColorBGR 获取颜色(BGR)
	// 参数: x - X坐标, y - Y坐标
	// 返回: 颜色字符串(BGR)
	GetColorBGR(x int32, y int32) string

	// GetColorHSV 获取颜色(HSV)
	// 参数: x - X坐标, y - Y坐标
	// 返回: 颜色字符串(HSV)
	GetColorHSV(x int32, y int32) string

	// CmpColor 比较颜色
	// 参数: x - X坐标, y - Y坐标, color - 颜色, sim - 相似度
	// 返回: 1匹配, 0不匹配
	CmpColor(x int32, y int32, color string, sim float64) int32

	// GetAveHSV 获取区域平均HSV
	// 参数: x1,y1,x2,y2 - 区域坐标
	// 返回: HSV字符串
	GetAveHSV(x1 int32, y1 int32, x2 int32, y2 int32) string

	// RGB2BGR RGB转BGR
	// 参数: rgb_color - RGB颜色
	// 返回: BGR颜色字符串
	RGB2BGR(rgb_color string) string

	// BGR2RGB BGR转RGB
	// 参数: bgr_color - BGR颜色
	// 返回: RGB颜色字符串
	BGR2RGB(bgr_color string) string

	// EnableGetColorByCapture 启用截图取色
	// 参数: enable - 1启用, 0禁用
	// 返回: 1成功
	EnableGetColorByCapture(enable int32) int32

	// ==================== 系统信息 ====================

	// GetScreenWidth 获取屏幕宽度
	// 返回: 屏幕宽度(像素)
	GetScreenWidth() int32

	// GetScreenHeight 获取屏幕高度
	// 返回: 屏幕高度(像素)
	GetScreenHeight() int32

	// GetScreenDepth 获取屏幕色深
	// 返回: 色深(位)
	GetScreenDepth() int32

	// GetDPI 获取DPI
	// 返回: DPI值
	GetDPI() int32

	// GetFps 获取FPS
	// 返回: FPS值
	GetFps() int32

	// GetOsType 获取操作系统类型
	// 返回: 系统类型
	GetOsType() int32

	// GetOsBuildNumber 获取系统版本号
	// 返回: 版本号
	GetOsBuildNumber() int32

	// GetTime 获取时间戳
	// 返回: 时间戳(毫秒)
	GetTime() int32

	// GetNetTime 获取网络时间
	// 返回: 网络时间字符串
	GetNetTime() string

	// GetNetTimeSafe 安全获取网络时间
	// 返回: 网络时间字符串
	GetNetTimeSafe() string

	// GetNetTimeByIp 通过IP获取网络时间
	// 参数: ip - IP地址
	// 返回: 网络时间字符串
	GetNetTimeByIp(ip string) string

	// GetCpuType 获取CPU类型
	// 返回: CPU类型
	GetCpuType() int32

	// GetCpuUsage 获取CPU使用率
	// 返回: 使用率(百分比)
	GetCpuUsage() int32

	// GetMemoryUsage 获取内存使用量
	// 返回: 使用量(MB)
	GetMemoryUsage() int32

	// Is64Bit 判断是否64位系统
	// 返回: 1是64位, 0是32位
	Is64Bit() int32

	// IsSurrpotVt 判断是否支持VT
	// 返回: 1支持, 0不支持
	IsSurrpotVt() int32

	// GetSystemInfo 获取系统信息
	// 参数: type_ - 类型, method - 方法
	// 返回: 系统信息字符串
	GetSystemInfo(type_ string, method int32) string

	// GetDisplayInfo 获取显示器信息
	// 返回: 显示器信息字符串
	GetDisplayInfo() string

	// ==================== 截图与图像 ====================

	// Capture 截图保存为BMP
	// 参数: x1,y1,x2,y2 - 区域坐标, file - 文件路径
	// 返回: 1成功, 0失败
	Capture(x1 int32, y1 int32, x2 int32, y2 int32, file string) int32

	// CapturePng 截图保存为PNG
	// 参数: x1,y1,x2,y2 - 区域坐标, file - 文件路径
	// 返回: 1成功, 0失败
	CapturePng(x1 int32, y1 int32, x2 int32, y2 int32, file string) int32

	// CaptureJpg 截图保存为JPG
	// 参数: x1,y1,x2,y2 - 区域坐标, file - 文件路径, quality - 质量
	// 返回: 1成功, 0失败
	CaptureJpg(x1 int32, y1 int32, x2 int32, y2 int32, file string, quality int32) int32

	// CaptureGif 截图保存为GIF
	// 参数: x1,y1,x2,y2 - 区域坐标, file - 文件路径, delay - 延迟, time - 时间
	// 返回: 1成功, 0失败
	CaptureGif(x1 int32, y1 int32, x2 int32, y2 int32, file string, delay int32, time int32) int32

	// CapturePre 预截图
	// 参数: file - 文件路径
	// 返回: 1成功, 0失败
	CapturePre(file string) int32

	// GetScreenData 获取屏幕数据
	// 参数: x1,y1,x2,y2 - 区域坐标
	// 返回: 数据句柄
	GetScreenData(x1 int32, y1 int32, x2 int32, y2 int32) int32

	// GetScreenDataBmp 获取屏幕BMP数据
	// 参数: x1,y1,x2,y2 - 区域坐标, data - 数据指针, size - 大小指针
	// 返回: 1成功, 0失败
	GetScreenDataBmp(x1 int32, y1 int32, x2 int32, y2 int32, data *int32, size *int32) int32

	// ImageToBmp 图片转BMP
	// 参数: pic_name - 图片名, bmp_name - BMP名
	// 返回: 1成功, 0失败
	ImageToBmp(pic_name string, bmp_name string) int32

	// GetPicSize 获取图片尺寸
	// 参数: pic_name - 图片名
	// 返回: "宽度,高度" 字符串
	GetPicSize(pic_name string) string

	// FreePic 释放图片
	// 参数: pic_name - 图片名
	// 返回: 1成功, 0失败
	FreePic(pic_name string) int32

	// EnablePicCache 启用图片缓存
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnablePicCache(en int32) int32

	// SetPicPwd 设置图片密码
	// 参数: pwd - 密码
	// 返回: 1成功
	SetPicPwd(pwd string) int32

	// AppendPicAddr 追加图片地址
	// 参数: pic_info - 图片信息, addr - 地址, size - 大小
	// 返回: 图片信息字符串
	AppendPicAddr(pic_info string, addr int32, size int32) string

	// ==================== 找色与形状 ====================

	// FindColor 查找颜色
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32, x *int32, y *int32) int32

	// FindColorE 查找颜色(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度, dir - 方向
	// 返回: "x,y" 格式字符串
	FindColorE(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string

	// FindColorEx 高级查找颜色
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	FindColorEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string

	// FindMultiColor 多点找色
	// 参数: x1,y1,x2,y2 - 区域坐标, first_color - 第一个颜色, offset_color - 偏移颜色, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindMultiColor(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32, x *int32, y *int32) int32

	// FindMultiColorE 多点找色(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, first_color - 第一个颜色, offset_color - 偏移颜色, sim - 相似度, dir - 方向
	// 返回: "x,y" 格式字符串
	FindMultiColorE(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string

	// FindMultiColorEx 高级多点找色
	// 参数: x1,y1,x2,y2 - 区域坐标, first_color - 第一个颜色, offset_color - 偏移颜色, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	FindMultiColorEx(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string

	// FindMulColor 查找多个颜色
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度
	// 返回: 1找到, 0未找到
	FindMulColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32

	// GetColorNum 获取颜色数量
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度
	// 返回: 颜色数量
	GetColorNum(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32

	// GetAveRGB 获取区域平均RGB
	// 参数: x1,y1,x2,y2 - 区域坐标
	// 返回: RGB字符串
	GetAveRGB(x1 int32, y1 int32, x2 int32, y2 int32) string

	// FindColorBlock 查找色块
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度, count - 数量, width - 宽度, height - 高度, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindColorBlock(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32, x *int32, y *int32) int32

	// FindColorBlockEx 高级查找色块
	// 参数: x1,y1,x2,y2 - 区域坐标, color - 颜色, sim - 相似度, count - 数量, width - 宽度, height - 高度
	// 返回: 所有找到的坐标字符串
	FindColorBlockEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32) string

	// FindShape 查找形状
	// 参数: x1,y1,x2,y2 - 区域坐标, offset_color - 偏移颜色, sim - 相似度, dir - 方向, x,y - 返回坐标
	// 返回: 1找到, 0未找到
	FindShape(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32, x *int32, y *int32) int32

	// FindShapeE 查找形状(返回坐标字符串)
	// 参数: x1,y1,x2,y2 - 区域坐标, offset_color - 偏移颜色, sim - 相似度, dir - 方向
	// 返回: "x,y" 格式字符串
	FindShapeE(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string

	// FindShapeEx 高级查找形状
	// 参数: x1,y1,x2,y2 - 区域坐标, offset_color - 偏移颜色, sim - 相似度, dir - 方向
	// 返回: 所有找到的坐标字符串
	FindShapeEx(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string

	// ==================== 窗口操作 ====================

	// BindWindow 绑定窗口
	// 参数: hwnd - 窗口句柄, display - 显示模式, mouse - 鼠标模式, keypad - 键盘模式, mode - 模式
	// 返回: 1成功, 0失败
	BindWindow(hwnd int32, display string, mouse string, keypad string, mode int32) int32

	// BindWindowEx 扩展绑定窗口
	// 参数: hwnd - 窗口句柄, display - 显示模式, mouse - 鼠标模式, keypad - 键盘模式, public_desc - 公共描述, mode - 模式
	// 返回: 1成功, 0失败
	BindWindowEx(hwnd int32, display string, mouse string, keypad string, public_desc string, mode int32) int32

	// UnBindWindow 解绑窗口
	// 返回: 1成功, 0失败
	UnBindWindow() int32

	// ForceUnBindWindow 强制解绑窗口
	// 参数: hwnd - 窗口句柄
	// 返回: 1成功, 0失败
	ForceUnBindWindow(hwnd int32) int32

	// SwitchBindWindow 切换绑定窗口
	// 参数: hwnd - 窗口句柄
	// 返回: 1成功, 0失败
	SwitchBindWindow(hwnd int32) int32

	// GetBindWindow 获取绑定窗口
	// 返回: 窗口句柄
	GetBindWindow() int32

	// IsBind 判断是否绑定
	// 参数: hwnd - 窗口句柄
	// 返回: 1已绑定, 0未绑定
	IsBind(hwnd int32) int32

	// EnableBind 启用绑定
	// 参数: en - 1启用, 0禁用
	// 返回: 1成功
	EnableBind(en int32) int32

	// FindWindow 查找窗口
	// 参数: class_name - 类名, title_name - 标题名
	// 返回: 窗口句柄
	FindWindow(class_name string, title_name string) int32

	// FindWindowEx 扩展查找窗口
	// 参数: parent - 父窗口, class_name - 类名, title_name - 标题名
	// 返回: 窗口句柄
	FindWindowEx(parent int32, class_name string, title_name string) int32

	// FindWindowByProcess 通过进程名查找窗口
	// 参数: process_name - 进程名, class_name - 类名, title_name - 标题名
	// 返回: 窗口句柄
	FindWindowByProcess(process_name string, class_name string, title_name string) int32

	// FindWindowByProcessId 通过进程ID查找窗口
	// 参数: process_id - 进程ID, class_name - 类名, title_name - 标题名
	// 返回: 窗口句柄
	FindWindowByProcessId(process_id int32, class_name string, title_name string) int32

	// FindWindowSuper 超级查找窗口
	// 参数: spec1 - 规格1, flag1 - 标志1, type1 - 类型1, spec2 - 规格2, flag2 - 标志2, type2 - 类型2
	// 返回: 窗口句柄
	FindWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32) int32

	// EnumWindow 枚举窗口
	// 参数: parent - 父窗口, title - 标题, class_name - 类名, filter - 过滤器
	// 返回: 窗口句柄列表字符串
	EnumWindow(parent int32, title string, class_name string, filter int32) string

	// EnumWindowByProcess 通过进程名枚举窗口
	// 参数: process_name - 进程名, title - 标题, class_name - 类名, filter - 过滤器
	// 返回: 窗口句柄列表字符串
	EnumWindowByProcess(process_name string, title string, class_name string, filter int32) string

	// EnumWindowByProcessId 通过进程ID枚举窗口
	// 参数: pid - 进程ID, title - 标题, class_name - 类名, filter - 过滤器
	// 返回: 窗口句柄列表字符串
	EnumWindowByProcessId(pid int32, title string, class_name string, filter int32) string

	// EnumWindowSuper 超级枚举窗口
	// 参数: spec1 - 规格1, flag1 - 标志1, type1 - 类型1, spec2 - 规格2, flag2 - 标志2, type2 - 类型2, sort - 排序
	// 返回: 窗口句柄列表字符串
	EnumWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32, sort int32) string

	// GetWindow 获取窗口
	// 参数: hwnd - 窗口句柄, flag - 标志
	// 返回: 窗口句柄
	GetWindow(hwnd int32, flag int32) int32

	// GetSpecialWindow 获取特殊窗口
	// 参数: flag - 标志
	// 返回: 窗口句柄
	GetSpecialWindow(flag int32) int32

	// GetWindowRect 获取窗口矩形
	// 参数: hwnd - 窗口句柄, x1,y1,x2,y2 - 返回坐标
	// 返回: 1成功, 0失败
	GetWindowRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32

	// GetClientRect 获取客户区矩形
	// 参数: hwnd - 窗口句柄, x1,y1,x2,y2 - 返回坐标
	// 返回: 1成功, 0失败
	GetClientRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32

	// GetClientSize 获取客户区大小
	// 参数: hwnd - 窗口句柄, width,height - 返回大小
	// 返回: 1成功, 0失败
	GetClientSize(hwnd int32, width *int32, height *int32) int32

	// SetClientSize 设置客户区大小
	// 参数: hwnd - 窗口句柄, width - 宽度, height - 高度
	// 返回: 1成功, 0失败
	SetClientSize(hwnd int32, width int32, height int32) int32

	// SetWindowSize 设置窗口大小
	// 参数: hwnd - 窗口句柄, width - 宽度, height - 高度
	// 返回: 1成功, 0失败
	SetWindowSize(hwnd int32, width int32, height int32) int32

	// MoveWindow 移动窗口
	// 参数: hwnd - 窗口句柄, x - X坐标, y - Y坐标
	// 返回: 1成功, 0失败
	MoveWindow(hwnd int32, x int32, y int32) int32

	// SetWindowState 设置窗口状态
	// 参数: hwnd - 窗口句柄, flag - 标志
	// 返回: 1成功, 0失败
	SetWindowState(hwnd int32, flag int32) int32

	// GetWindowState 获取窗口状态
	// 参数: hwnd - 窗口句柄, flag - 标志
	// 返回: 状态值
	GetWindowState(hwnd int32, flag int32) int32

	// GetWindowTitle 获取窗口标题
	// 参数: hwnd - 窗口句柄
	// 返回: 标题字符串
	GetWindowTitle(hwnd int32) string

	// SetWindowText 设置窗口标题
	// 参数: hwnd - 窗口句柄, text - 标题
	// 返回: 1成功, 0失败
	SetWindowText(hwnd int32, text string) int32

	// GetWindowClass 获取窗口类名
	// 参数: hwnd - 窗口句柄
	// 返回: 类名字符串
	GetWindowClass(hwnd int32) string

	// GetWindowProcessId 获取窗口进程ID
	// 参数: hwnd - 窗口句柄
	// 返回: 进程ID
	GetWindowProcessId(hwnd int32) int32

	// GetWindowThreadId 获取窗口线程ID
	// 参数: hwnd - 窗口句柄
	// 返回: 线程ID
	GetWindowThreadId(hwnd int32) int32

	// GetWindowProcessPath 获取窗口进程路径
	// 参数: hwnd - 窗口句柄
	// 返回: 进程路径字符串
	GetWindowProcessPath(hwnd int32) string

	// GetForegroundWindow 获取前台窗口
	// 返回: 窗口句柄
	GetForegroundWindow() int32

	// GetForegroundFocus 获取前台焦点窗口
	// 返回: 窗口句柄
	GetForegroundFocus() int32

	// GetMousePointWindow 获取鼠标指向窗口
	// 返回: 窗口句柄
	GetMousePointWindow() int32

	// GetPointWindow 获取指定坐标窗口
	// 参数: x - X坐标, y - Y坐标
	// 返回: 窗口句柄
	GetPointWindow(x int32, y int32) int32

	// ScreenToClient 屏幕坐标转客户区坐标
	// 参数: hwnd - 窗口句柄, x,y - 坐标
	// 返回: 1成功, 0失败
	ScreenToClient(hwnd int32, x *int32, y *int32) int32

	// ClientToScreen 客户区坐标转屏幕坐标
	// 参数: hwnd - 窗口句柄, x,y - 坐标
	// 返回: 1成功, 0失败
	ClientToScreen(hwnd int32, x *int32, y *int32) int32

	// SetWindowTransparent 设置窗口透明度
	// 参数: hwnd - 窗口句柄, v - 透明度(0-255)
	// 返回: 1成功, 0失败
	SetWindowTransparent(hwnd int32, v int32) int32

	// ShowTaskBarIcon 显示任务栏图标
	// 参数: hwnd - 窗口句柄, is_show - 1显示, 0隐藏
	// 返回: 1成功, 0失败
	ShowTaskBarIcon(hwnd int32, is_show int32) int32

	// SetEnumWindowDelay 设置枚举窗口延迟
	// 参数: delay - 延迟(毫秒)
	// 返回: 1成功
	SetEnumWindowDelay(delay int32) int32
}

// DmSoftBase 基础结构体，包含通用字段
// 用于存储大漠插件对象句柄和模块句柄
type DmSoftBase struct {
	obj     uintptr // 大漠插件COM对象指针
	hModule uintptr // DLL模块句柄
}

// utf8ToGbk 将UTF-8字符串转换为GBK编码的字节切片
// 参数: s - UTF-8编码的字符串
// 返回: GBK编码的字节切片
func utf8ToGbk(s string) []byte {
	data, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
	return data
}

// stringToBytePtr 将Go字符串转换为C风格字符串指针（GBK编码）
// 参数:
//   - s: Go字符串（UTF-8编码）
//
// 返回:
//   - C风格字符串指针（以\0结尾，GBK编码）
//
// 注意:
//   - 返回的指针指向的内存由Go垃圾回收器管理
//   - 仅用于传递给DLL函数，不要在Go代码中长期保存
//   - 自动将UTF-8编码转换为GBK编码，以兼容大漠插件
func stringToBytePtr(s string) *byte {
	gbkData := utf8ToGbk(s)
	gbkData = append(gbkData, 0)
	return &gbkData[0]
}

// int32Ptr 获取int32指针
// 参数:
//   - v: int32值
//
// 返回:
//   - int32指针
//
// 注意:
//   - 返回的指针指向的内存由Go垃圾回收器管理
func int32Ptr(v int32) *int32 {
	return &v
}

// gbkToUtf8 将GBK编码的字节切片转换为UTF-8字符串
// 参数: data - GBK编码的字节切片
// 返回: UTF-8编码的字符串
func gbkToUtf8(data []byte) string {
	result, _ := simplifiedchinese.GBK.NewDecoder().Bytes(data)
	return string(result)
}

// bytePtrToString 将C风格字符串指针转换为Go字符串
// 参数:
//   - ptr: C风格字符串指针（以\0结尾，GBK编码）
//
// 返回:
//   - Go字符串（UTF-8编码）
//
// 注意:
//   - 如果ptr为nil，返回空字符串
//   - 自动计算字符串长度直到遇到\0
//   - 自动将GBK编码转换为UTF-8编码，以兼容Go程序
func bytePtrToString(ptr *byte) string {
	if ptr == nil {
		return ""
	}
	n := 0
	for p := ptr; *p != 0; p = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1)) {
		n++
	}
	return gbkToUtf8(unsafe.Slice(ptr, n))
}
