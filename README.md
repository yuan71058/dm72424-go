# 🚀 大漠插件 Go 语言封装库

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuan71058/dm72424-go)](https://goreportcard.com/report/github.com/yuan71058/dm72424-go)

> 大漠插件 7.2424 版本的 Go 语言封装库，支持 428 个函数接口，开箱即用！

---

## 📖 目录

- [项目简介](#-项目简介)
- [快速开始](#-快速开始)
- [安装使用](#-安装使用)
- [API 文档](#-api-文档)
- [函数分类](#-函数分类)
- [注意事项](#-注意事项)
- [常见问题](#-常见问题)
- [更新日志](#-更新日志)

---

## 🎯 项目简介

本项目是大漠插件（dm.dll）的 Go 语言封装库，将原 C++ 版本的大漠插件接口完整翻译为 Go 语言。

### ✨ 特性

- 📦 **完整封装** - 支持大漠插件 7.2424 版本全部 428 个函数
- 🎯 **开箱即用** - 简单导入即可使用，无需复杂配置
- 📝 **详细注释** - 所有函数都有完整的中文注释
- 🔧 **类型安全** - 完整的类型定义，编译时检查
- 📚 **示例丰富** - 提供完整的使用示例
- 🌐 **自动编码转换** - 内置 UTF-8 到 GBK 自动转换，中文参数无需手动处理

### 📁 目录结构

```
dm72424-go/
├── dmsoft.go           # 接口定义文件
├── dmsoft_impl.go      # 接口实现文件（428个函数）
├── go.mod              # Go 模块定义
├── Go.dll              # 破解补丁
├── xd47243.dll         # 大漠插件主文件
├── example/            # 使用示例
│   ├── main.go         # 基础示例
│   ├── go.mod          # 示例模块
│   ├── go.sum          # 示例依赖
│   ├── multithread/    # 多线程示例
│   │   ├── main.go     # 多线程写入文本示例
│   │   └── go.mod      # 多线程示例模块
│   └── find_window/    # 查找窗口示例
│       └── main.go     # 查找中文窗口、绑定、截图示例
└── README.md           # 本文档
```

---

## 🚀 快速开始

### 环境要求

- **操作系统**: Windows（32位或64位均可）
- **Go 版本**: Go 1.16 或更高版本
- **编译要求**: 必须编译为 32 位程序

### 30秒上手

```go
package main

import (
    "fmt"
    "log"
    
    dmsoft "github.com/yuan71058/dm72424-go"
)

func main() {
    // 1. 加载大漠插件
    _, err := dmsoft.LoadDm("xd47243.dll")
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. 破解大漠插件
    err = dmsoft.CrackDm("Go.dll")
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. 创建对象并初始化
    dm := dmsoft.New()
    dm.Init()
    defer dm.Release()
    
    // 4. 注册
    if dm.Reg("", "") == 1 {
        fmt.Println("注册成功！")
    }
    
    // 5. 开始使用
    fmt.Printf("版本: %s\n", dm.Ver())
    fmt.Printf("分辨率: %d x %d\n", dm.GetScreenWidth(), dm.GetScreenHeight())
}
```

---

## 📥 安装使用

### 方法一：Go Modules（推荐）

```bash
# 在项目目录下执行
go get github.com/yuan71058/dm72424-go
```

### 方法二：本地引用

```bash
# 克隆项目
git clone https://github.com/yuan71058/dm72424-go.git

# 在你的 go.mod 中添加 replace
replace github.com/yuan71058/dm72424-go => ./dm72424-go
```

### 编译说明

⚠️ **重要**：大漠插件是 32 位 DLL，必须编译为 32 位程序！

```powershell
# 设置环境变量
go env -w GOARCH=386

# 编译
go build -o myapp.exe
```

---

## 📚 API 文档

### 核心函数

| 函数 | 说明 |
|------|------|
| `LoadDm(dmPath string) (uintptr, error)` | 加载大漠插件 DLL |
| `CrackDm(crackDllPath string) error` | 破解大漠插件 |
| `FreeCrackDll() bool` | 释放破解 DLL |
| `Load(path string) (uintptr, error)` | 加载大漠插件 DLL（旧版） |
| `Free() bool` | 释放大漠插件 |
| `New() *DmSoft` | 创建大漠对象 |
| `Init()` | 初始化对象 |
| `Release()` | 释放对象 |

### 常用功能示例

#### 🪟 窗口操作

```go
// 查找窗口（支持中文标题）
hwnd := dm.FindWindow("Qt51514QWindowIcon", "朋友圈")
if hwnd > 0 {
    fmt.Printf("找到窗口，句柄: %d\n", hwnd)
    
    // 获取窗口信息
    var x1, y1, x2, y2 int32
    dm.GetWindowRect(hwnd, &x1, &y1, &x2, &y2)
    fmt.Printf("窗口位置: (%d,%d) - (%d,%d)\n", x1, y1, x2, y2)
    
    // 绑定窗口（gdi图像模式 + windows3鼠标模式 + windows键盘模式）
    ret := dm.BindWindow(hwnd, "gdi", "windows3", "windows", 0)
    if ret == 1 {
        fmt.Println("绑定成功")
        // ... 执行操作 ...
        dm.UnBindWindow()
    }
}

// 枚举窗口
hwndList := dm.EnumWindow(0, "", "Notepad", 2)
```

#### 🖼️ 截图功能

```go
// 设置保存路径
dm.SetPath("C:\\screenshots")

// 截取全屏
dm.Capture(0, 0, 1920, 1080, "screen.bmp")

// 截取指定区域
dm.Capture(100, 100, 500, 500, "region.bmp")

// JPG格式截图
dm.CaptureJpg(0, 0, 1920, 1080, "screen.jpg", 80)
```

#### 🔍 找图功能

```go
// 预加载图片
dm.LoadPic("target.bmp")

// 查找图片
var x, y int32
ret := dm.FindPic(0, 0, 1920, 1080, "target.bmp", "000000", 0.9, 0, &x, &y)
if ret != -1 {
    fmt.Printf("找到图片: (%d, %d)\n", x, y)
}
```

#### 🎨 找色功能

```go
// 查找颜色
var x, y int32
ret := dm.FindColor(0, 0, 1920, 1080, "FF0000", "000000", 1.0, 0, &x, &y)

// 多点找色
ret = dm.FindMultiColor(0, 0, 1920, 1080, "FF0000", "5|0|00FF00", "000000", 1.0, 0, &x, &y)
```

#### 📝 文字识别（OCR）

```go
// 设置字库
dm.SetDict(0, "dict.txt")

// OCR识别
text := dm.Ocr(0, 0, 500, 100, "FFFFFF-000000", 1.0)

// 查找文字
var x, y int32
ret := dm.FindStr(0, 0, 1920, 1080, "登录", "FFFFFF-000000", 1.0, &x, &y)
```

#### 🖱️ 鼠标操作

```go
// 移动鼠标
dm.MoveTo(500, 300)

// 点击
dm.LeftClick()
dm.RightClick()
dm.LeftDoubleClick()

// 获取位置
var x, y int32
dm.GetCursorPos(&x, &y)
```

#### ⌨️ 键盘操作

```go
// 按键
dm.KeyPressChar("a")
dm.KeyPress(65)  // A键

// 组合键
dm.KeyDownChar("ctrl")
dm.KeyPressChar("c")
dm.KeyUpChar("ctrl")

// 发送字符串
dm.SendString(hwnd, "Hello World")
```

#### 💾 内存操作

```go
// 读取内存
value := dm.ReadInt(hwnd, 0x12345678, 0)

// 写入内存
dm.WriteInt(hwnd, 0x12345678, 0, 12345)

// 查找特征码
addr := dm.FindData(hwnd, 0x400000, 0x500000, "FF ?? 00 ??")
```

#### 🧵 多线程操作

```go
// 多线程示例：同时向多个窗口写入文本
// 完整示例请参考 example/multithread/main.go

// ========== 1. 定义工作线程结构体 ==========
type TextWorker struct {
    ID       int
    Dm       *dmsoft.DmSoft
    EditHwnd int32
    Content  string
    Result   chan string
}

// Init 初始化大漠子对象
func (w *TextWorker) Init() bool {
    w.Dm = dmsoft.New()
    if w.Dm == nil {
        return false
    }
    w.Dm.Init()  // 每个线程独立初始化
    return true
}

// Run 运行工作线程
func (w *TextWorker) Run(wg *sync.WaitGroup) {
    defer wg.Done()
    
    // 绑定窗口
    ret := w.Dm.BindWindow(w.EditHwnd, "gdi", "windows", "windows", 0)
    if ret != 1 {
        w.Result <- fmt.Sprintf("[线程%d] 绑定失败", w.ID)
        return
    }
    
    // 写入文字
    w.Dm.SendString2(w.EditHwnd, w.Content)
    
    // 解绑并释放
    w.Dm.UnBindWindow()
    w.Dm.Release()
    
    w.Result <- fmt.Sprintf("[线程%d] 完成", w.ID)
}

// ========== 2. 主函数 ==========
func main() {
    // 加载并破解大漠
    dmsoft.LoadDm("xd47243.dll")
    dmsoft.CrackDm("Go.dll")
    
    // 创建主对象并注册（只需一次）
    mainDm := dmsoft.New()
    mainDm.Init()
    mainDm.Reg("", "")
    
    // 枚举窗口、查找Edit控件...
    hwndList := mainDm.EnumWindow(0, "", "Notepad", 2)
    // ... 省略窗口枚举和控件查找代码
    
    // 创建工作线程
    var wg sync.WaitGroup
    resultChan := make(chan string, 100)
    workers := make([]*TextWorker, 3)
    
    for i := 0; i < 3; i++ {
        workers[i] = &TextWorker{
            ID:       i + 1,
            EditHwnd: editHwnds[i],
            Content:  fmt.Sprintf("线程%d的内容...", i+1),
            Result:   resultChan,
        }
        workers[i].Init()
    }
    
    // 并发执行
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go workers[i].Run(&wg)
    }
    
    wg.Wait()
    close(resultChan)
    
    // 最后释放主对象
    mainDm.Release()
}
```

**多线程要点**：

| 要点 | 说明 |
|------|------|
| 主对象 | 全局唯一，负责注册，最后释放 |
| 子对象 | 每线程独立创建，各自 Init/Release |
| 注册 | 只需在主对象中注册一次，子对象无需注册 |
| 绑定 | 每个子对象独立绑定自己的目标窗口 |
| 释放顺序 | 先释放所有子对象，最后释放主对象 |

---

## 📋 函数分类

| 分类 | 函数数量 | 主要函数 |
|------|----------|----------|
| 🪟 窗口操作 | ~50 | `BindWindow`, `FindWindow`, `GetWindowRect` |
| 🖱️ 鼠标操作 | ~20 | `MoveTo`, `LeftClick`, `GetCursorPos` |
| ⌨️ 键盘操作 | ~15 | `KeyPress`, `KeyDown`, `SendString` |
| 🖼️ 图像处理 | ~30 | `Capture`, `FindPic`, `LoadPic` |
| 🎨 颜色操作 | ~15 | `GetColor`, `FindColor`, `CmpColor` |
| 📝 OCR识别 | ~20 | `Ocr`, `FindStr`, `SetDict` |
| 💾 内存操作 | ~40 | `ReadInt`, `WriteInt`, `FindData` |
| 💻 系统信息 | ~20 | `Ver`, `GetOsType`, `GetTime` |
| 📁 文件操作 | ~15 | `ReadFile`, `WriteFile`, `IsFileExist` |
| 🤖 AI功能 | ~10 | `LoadAi`, `FindPicAi` |

---

## ⚠️ 注意事项

### 1. 初始化顺序

必须按照以下顺序初始化：

```go
dm := dmsoft.New()
dm.Init()  // 必须调用！
defer dm.Release()
```

### 2. 窗口绑定

大部分屏幕操作需要先绑定窗口：

```go
hwnd := dm.GetForegroundWindow()
dm.BindWindow(hwnd, "gdi", "normal", "normal", 0)
// ... 操作 ...
dm.UnBindWindow()
```

### 3. 管理员权限

以下功能需要管理员权限：
- 内存读写操作
- 某些窗口绑定模式（dx模式）
- 进程操作

### 4. 编码问题

本库已内置 UTF-8 到 GBK 的自动编码转换，所有字符串参数函数都会自动处理中文编码：

```go
// 直接使用中文参数，无需手动转换
hwnd := dm.FindWindow("Qt51514QWindowIcon", "朋友圈")
dm.SetPath("C:\\测试目录")
text := dm.Ocr(0, 0, 500, 100, "FFFFFF-000000", 1.0)
```

**大漠插件返回的字符串可能是 GBK 编码**，如果显示乱码，可手动转换：

```go
import "golang.org/x/text/encoding/simplifiedchinese"

func gbkToUtf8(s string) string {
    data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(s))
    return string(data)
}
```

### 5. Init/Release 对象生命周期

`Init()` 和 `Release()` 管理大漠 COM 对象实例的生命周期：

```go
dm := dmsoft.New()    // 创建实例(仅分配结构体指针,不分配COM对象)
dm.Init()             // 初始化(必须调用,创建COM对象)
defer dm.Release()    // 程序结束时释放(COM对象销毁后不可再用)
```

**调用场景**：

| 场景 | Init() 调用次数 | 说明 |
|------|----------------|------|
| 单线程/全局使用 | **1次** | 全局仅创建一个 DmSoft 实例，初始化一次即可 |
| 多线程场景 | **每线程1次** | 大漠 COM 对象线程相关，每个线程需独立创建 DmSoft 实例并调用 Init() |

**多线程示例**：

```go
// 主线程：创建主对象并注册
mainDm := dmsoft.New()
mainDm.Init()
mainDm.Reg("", "")  // 只需在主对象中注册一次

// 子线程：每个线程创建独立对象(无需再次注册)
for i := 0; i < 3; i++ {
    go func(threadID int) {
        subDm := dmsoft.New()
        subDm.Init()  // 各自独立初始化
        subDm.BindWindow(hwnd, "gdi", "windows", "windows", 0)
        // ... 操作 ...
        subDm.UnBindWindow()
        subDm.Release()  // 各自释放
    }(i)
}

defer mainDm.Release()  // 最后释放主对象
```

### 6. DLL 偏移地址说明

本库直接调用大漠 DLL 内部函数，地址偏移量固定：

| 功能 | 偏移地址 | 说明 |
|------|----------|------|
| 创建对象 | `DmHModule + 98304` | Init() 调用，偏移 0x18000 |
| 释放对象 | `DmHModule + 98400` | Release() 调用，偏移 0x18090 |
| 其他函数 | 各不相同 |详见 dmsoft_impl.go 中的函数定义 |

---

## ❓ 常见问题

<details>
<summary><b>Q1: 注册失败怎么办？</b></summary>

- 确保破解DLL已正确加载和执行
- 检查破解DLL版本是否与大漠插件版本匹配
- 尝试以管理员权限运行
</details>

<details>
<summary><b>Q2: 截图返回黑色图像？</b></summary>

- 确保已正确绑定窗口
- 尝试不同的绑定模式（gdi/dx/opengl）
- 检查窗口是否最小化或被遮挡
</details>

<details>
<summary><b>Q3: 找图/找色失败？</b></summary>

- 检查图片路径是否正确
- 确保已设置正确的资源路径（SetPath）
- 调整相似度参数
- 检查颜色格式是否正确
</details>

<details>
<summary><b>Q4: 编译报错 "not a valid Win32 application"？</b></summary>

- 确保编译为32位程序：`go env -w GOARCH=386`
- 大漠插件是32位DLL，不能在64位程序中调用
</details>

---

## 📝 更新日志

### v1.3.0 (2026-03-20)

- 🌐 **新增自动编码转换** - 所有字符串参数函数自动将 UTF-8 转换为 GBK
- ✨ 新增 `GetDmHModule()` 函数 - 获取大漠插件模块句柄
- ✨ 新增 `GetObjPtr()` 方法 - 获取大漠对象指针
- 📦 添加 `golang.org/x/text` 依赖用于编码转换
- 📝 更新 README 文档，说明编码自动转换特性
- 🔧 优化 309 处字符串参数调用，统一使用 GBK 编码
- 📚 新增 `example/find_window` 示例 - 演示查找中文窗口、绑定窗口、截图功能
- 🏷️ 添加版本标签 v1.3.0，支持 Go Modules 版本管理

### v1.2.0 (2026-03-19)

- ✨ 新增 `LoadDm()` 函数 - 加载大漠插件DLL
- ✨ 新增 `CrackDm()` 函数 - 破解大漠插件
- ✨ 新增 `FreeCrackDll()` 函数 - 释放破解DLL
- 📝 更新示例代码，使用新的封装函数
- 📝 更新README文档

### v1.1.0 (2026-03-18)

- ✨ 添加多线程示例（example/multithread）
- 🐛 修复 Release 函数释放后未清零指针的问题
- 📝 添加多线程使用文档和注意事项
- 🎯 智能窗口排列功能
- 🔧 优化示例代码，添加详细调试信息

### v1.0.0 (2026-03-18)

- ✨ 初始版本发布
- 📦 完成428个函数的翻译
- 📝 添加详细中文注释
- 🐛 修复偏移地址错误
- 📚 添加完整使用示例

---

## 📄 许可证

本项目仅供学习交流使用，请勿用于商业用途。

大漠插件版权归大漠插件作者所有。

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

## 📮 联系方式

- GitHub: [https://github.com/yuan71058/dm72424-go](https://github.com/yuan71058/dm72424-go)

---

<p align="center">
  如果这个项目对你有帮助，请给一个 ⭐️ Star 支持一下！
</p>
