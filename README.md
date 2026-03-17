# 大漠插件 Go 语言封装库

## 项目简介

本项目是大漠插件（dm.dll）的 Go 语言封装库，将原 C++ 版本的大漠插件接口完整翻译为 Go 语言，支持大漠插件 7.2424 版本的所有 428 个函数接口。

## 目录结构

```
E:\SRC\dm\72424_C++\go\
├── go.mod              # Go 模块定义文件
├── main.go             # 主程序入口，包含完整的使用示例
├── dmsoft.go           # 接口定义文件，定义 DmSoftInterface 接口
├── dmsoft_impl.go      # 接口实现文件，实现所有 428 个函数
├── dmsoft_test.exe     # 编译后的可执行文件（32位）
└── README.md           # 本文档
```

## 环境要求

- **操作系统**: Windows（32位或64位系统均可，但程序必须编译为32位）
- **Go 版本**: Go 1.16 或更高版本
- **必需文件**:
  - `xd47243.dll` - 大漠插件主文件
  - `Go.dll` - 破解补丁（可选，用于绕过注册验证）

## 编译说明

### 重要：必须编译为 32 位程序

大漠插件是 32 位 DLL，因此 Go 程序必须编译为 32 位才能正常调用。

```powershell
# 方法1：设置环境变量后编译
go env -w GOARCH=386
go build -o dmsoft_test.exe

# 方法2：临时设置编译
$env:GOARCH="386"; go build -o dmsoft_test.exe

# 方法3：命令行直接指定（Linux/Mac）
GOARCH=386 go build -o dmsoft_test.exe
```

## 快速开始

### 1. 基本使用流程

```go
package main

import (
    "fmt"
    "log"
    "syscall"
)

func main() {
    // 步骤1：加载大漠插件DLL
    dmHModule, err := LoadDm("xd47243.dll")
    if err != nil {
        log.Fatalf("加载大漠插件失败: %v", err)
    }
    fmt.Printf("大漠插件加载成功，模块句柄: %v\n", dmHModule)

    // 步骤2：加载破解补丁（可选）
    goHModule, err := syscall.LoadLibrary("Go.dll")
    if err != nil {
        log.Fatalf("加载破解DLL失败: %v", err)
    }
    defer syscall.FreeLibrary(goHModule)

    goFunAddr, _ := syscall.GetProcAddress(goHModule, "Go")
    syscall.Syscall(uintptr(goFunAddr), 1, dmHModule, 0, 0)

    // 步骤3：创建大漠对象
    dm := NewDmSoftImpl()
    if dm == nil {
        log.Fatal("创建大漠对象失败")
    }

    // 步骤4：初始化对象（重要！）
    dm.Init()
    defer dm.Release() // 确保程序退出时释放资源

    // 步骤5：注册插件
    ret := dm.Reg("", "")
    if ret == 1 {
        fmt.Println("大漠注册成功")
    } else {
        fmt.Println("大漠注册失败")
    }

    // 步骤6：使用各种功能
    version := dm.Ver()
    fmt.Printf("大漠版本: %s\n", version)

    // 绑定窗口后才能使用截图、取色等功能
    hwnd := dm.GetForegroundWindow()
    dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)

    // 获取屏幕颜色
    color := dm.GetColor(100, 100)
    fmt.Printf("屏幕(100,100)颜色: %s\n", color)

    // 解绑窗口
    dm.UnBindWindow()
}
```

### 2. 窗口绑定模式说明

```go
// 绑定窗口参数说明
// hwnd: 窗口句柄
// display: 显示模式
//   - "normal" - 正常模式，使用GDI截图
//   - "gdi"    - GDI模式，适合大多数情况
//   - "dx"     - DirectX模式，适合游戏
//   - "opengl" - OpenGL模式
// mouse: 鼠标模式
//   - "normal" - 正常模式
//   - "windows" - Windows消息模式
//   - "dx"     - DirectX模拟
// keypad: 键盘模式
//   - "normal" - 正常模式
//   - "windows" - Windows消息模式
// mode: 绑定模式（0-5）

ret := dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)
```

### 3. 常用功能示例

#### 截图功能

```go
// 设置截图保存路径
dm.SetPath("C:\\screenshots")

// 截取全屏
dm.Capture(0, 0, 1920, 1080, "screen.bmp")

// 截取指定区域
dm.Capture(100, 100, 500, 500, "region.bmp")

// 截图为JPG格式
dm.CaptureJpg(0, 0, 1920, 1080, "screen.jpg", 80)

// 截图为PNG格式
dm.CapturePng(0, 0, 1920, 1080, "screen.png")
```

#### 找图功能

```go
// 预加载图片到内存（提高查找速度）
dm.LoadPic("target.bmp")

// 在屏幕上查找图片
var x, y int32
ret := dm.FindPic(0, 0, 1920, 1080, "target.bmp", "000000", 0.9, 0, &x, &y)
if ret != -1 {
    fmt.Printf("找到图片，位置: (%d, %d)\n", x, y)
}

// 高级找图，返回所有匹配位置
result := dm.FindPicEx(0, 0, 1920, 1080, "target.bmp", "000000", 0.9, 0)
fmt.Printf("所有匹配位置: %s\n", result)
```

#### 找色功能

```go
// 在屏幕上查找颜色
var x, y int32
ret := dm.FindColor(0, 0, 1920, 1080, "FF0000", "000000", 1.0, 0, &x, &y)
if ret == 1 {
    fmt.Printf("找到颜色，位置: (%d, %d)\n", x, y)
}

// 多点找色
ret := dm.FindMultiColor(0, 0, 1920, 1080, "FF0000", "5|0|00FF00|10|5|0000FF", "000000", 1.0, 0, &x, &y)
```

#### 文字识别（OCR）

```go
// 设置字库
dm.SetDict(0, "dict.txt")

// OCR识别文字
text := dm.Ocr(0, 0, 500, 100, "FFFFFF-000000", 1.0)
fmt.Printf("识别到的文字: %s\n", text)

// 查找文字
var x, y int32
ret := dm.FindStr(0, 0, 1920, 1080, "登录", "FFFFFF-000000", 1.0, &x, &y)
if ret == 1 {
    fmt.Printf("找到文字，位置: (%d, %d)\n", x, y)
}
```

#### 鼠标操作

```go
// 移动鼠标
dm.MoveTo(500, 300)

// 相对移动
dm.MoveR(10, 10)

// 点击
dm.LeftClick()
dm.RightClick()
dm.LeftDoubleClick()

// 按住和释放
dm.LeftDown()
// ... 执行拖动操作 ...
dm.LeftUp()

// 获取鼠标位置
var x, y int32
dm.GetCursorPos(&x, &y)

// 滚轮
dm.WheelDown()  // 向下滚动
dm.WheelUp()    // 向上滚动
```

#### 键盘操作

```go
// 按键（字符形式）
dm.KeyPressChar("a")
dm.KeyDownChar("ctrl")
dm.KeyUpChar("ctrl")

// 按键（虚拟键码）
dm.KeyPress(65)  // A键

// 发送字符串
dm.SendString(hwnd, "Hello World")

// 按键组合
dm.KeyDownChar("ctrl")
dm.KeyPressChar("c")
dm.KeyUpChar("ctrl")
```

#### 内存操作

```go
// 注意：内存操作需要管理员权限

// 打开进程
hProcess := dm.OpenProcess(0x1F0FFF, pid)

// 读取整数
value := dm.ReadInt(hwnd, 0x12345678, 0)

// 写入整数
dm.WriteInt(hwnd, 0x12345678, 0, 12345)

// 读取字符串
str := dm.ReadString(hwnd, 0x12345678, 0, 100)

// 查找特征码
addr := dm.FindData(hwnd, 0x400000, 0x500000, "FF ?? 00 ??")
```

## 注意事项

### 1. 初始化顺序（重要！）

必须按照以下顺序初始化，否则会导致注册失败或功能异常：

```go
// 正确顺序：
// 1. 加载大漠插件DLL
// 2. 加载并执行破解DLL
// 3. 创建大漠对象
// 4. 调用Init()初始化对象 ← 重要！
// 5. 注册插件
// 6. 绑定窗口（如需要）
// 7. 使用各种功能

dm := NewDmSoftImpl()
dm.Init()  // 必须调用！否则dm.obj为0，所有函数调用都会失败
```

### 2. 窗口绑定

大部分屏幕操作（截图、取色、找图等）需要先绑定窗口：

```go
// 获取窗口句柄
hwnd := dm.GetForegroundWindow()

// 绑定窗口
ret := dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)
if ret != 1 {
    fmt.Println("绑定窗口失败")
    return
}

// 现在可以进行截图、取色等操作
color := dm.GetColor(100, 100)

// 使用完毕后解绑
dm.UnBindWindow()
```

### 3. 资源释放

使用完毕后必须释放资源，避免内存泄漏：

```go
dm := NewDmSoftImpl()
dm.Init()
defer dm.Release()  // 确保程序退出时释放
```

### 4. 字库和图片路径

使用OCR和找图功能前，需要设置正确的路径：

```go
// 设置全局路径（图片、字库等文件的查找路径）
dm.SetPath("C:\\my_project\\resources")

// 设置字库
dm.SetDict(0, "my_dict.txt")  // 文件在C:\my_project\resources\my_dict.txt
```

### 5. 错误处理

建议关闭错误提示框，避免弹窗干扰：

```go
// 关闭错误提示框
dm.SetShowErrorMsg(0)

// 获取最后错误码
errCode := dm.GetLastError()
```

### 6. 管理员权限

以下功能需要管理员权限：
- 内存读写操作
- 某些窗口绑定模式（dx模式）
- 进程操作

### 7. 多线程安全

大漠插件对象不是线程安全的，如果需要在多线程环境中使用，应该：
- 每个线程创建独立的大漠对象
- 或者使用互斥锁保护共享对象

### 8. 编码问题

大漠插件返回的字符串可能是 GBK 编码，在 Go 中可能需要转换：

```go
import "golang.org/x/text/encoding/simplifiedchinese"

// GBK 转 UTF-8
func gbkToUtf8(s string) string {
    data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(s))
    return string(data)
}

title := dm.GetWindowTitle(hwnd)
title = gbkToUtf8(title)
```

## 函数分类

### 窗口操作（约50个函数）
- `BindWindow`, `UnBindWindow`, `FindWindow`, `EnumWindow`
- `GetWindowRect`, `GetWindowTitle`, `SetWindowTitle`
- `MoveWindow`, `SetWindowSize`, `SetWindowTransparent`

### 鼠标操作（约20个函数）
- `MoveTo`, `MoveR`, `LeftClick`, `RightClick`
- `GetCursorPos`, `SetMouseSpeed`, `LockMouseRect`

### 键盘操作（约15个函数）
- `KeyPress`, `KeyDown`, `KeyUp`
- `KeyPressChar`, `KeyDownChar`, `KeyUpChar`
- `SendString`, `WaitKey`

### 图像处理（约30个函数）
- `FindPic`, `FindPicEx`, `FindPicMem`
- `Capture`, `CaptureJpg`, `CapturePng`
- `LoadPic`, `FreePic`

### 颜色操作（约15个函数）
- `GetColor`, `FindColor`, `FindColorEx`
- `FindMultiColor`, `CmpColor`
- `RGB2BGR`, `BGR2RGB`

### OCR文字识别（约20个函数）
- `Ocr`, `FindStr`, `FindStrEx`
- `SetDict`, `GetWords`, `FetchWord`

### 内存操作（约40个函数）
- `ReadInt`, `WriteInt`, `ReadFloat`, `WriteFloat`
- `FindData`, `VirtualAllocEx`, `OpenProcess`

### 系统信息（约20个函数）
- `Ver`, `GetOsType`, `GetCpuType`
- `GetScreenWidth`, `GetScreenHeight`
- `GetTime`, `GetNetTime`

### 文件操作（约15个函数）
- `ReadFile`, `WriteFile`, `CopyFile`
- `IsFileExist`, `CreateFolder`

### 其他功能
- AI图像识别（YOLO）
- Foobar窗口
- 汇编调用
- 进程管理

## 常见问题

### Q1: 注册失败怎么办？
- 确保破解DLL已正确加载和执行
- 检查破解DLL版本是否与大漠插件版本匹配
- 尝试以管理员权限运行

### Q2: 截图返回黑色图像？
- 确保已正确绑定窗口
- 尝试不同的绑定模式（gdi/dx/opengl）
- 检查窗口是否最小化或被遮挡

### Q3: 找图/找色失败？
- 检查图片路径是否正确
- 确保已设置正确的资源路径（SetPath）
- 调整相似度参数
- 检查颜色格式是否正确

### Q4: 内存读写失败？
- 确保以管理员权限运行
- 检查目标进程是否存在
- 确认内存地址是否正确

### Q5: 编译报错 "not a valid Win32 application"？
- 确保编译为32位程序：`go env -w GOARCH=386`
- 大漠插件是32位DLL，不能在64位程序中调用

## 版本历史

- **v1.0** - 初始版本，完成428个函数的翻译
  - 从C++版本完整移植
  - 添加详细的中文注释
  - 修复偏移地址错误
  - 添加完整的使用示例

## 许可证

本项目仅供学习交流使用，请勿用于商业用途。

大漠插件版权归大漠插件作者所有。
