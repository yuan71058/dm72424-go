# 🚀 大漠插件 Go 语言封装库

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/yuan71058/dm72424-go)](https://goreportcard.com/report/gitee.com/yuan71058/dm72424-go)

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

### 📁 目录结构

```
dm72424-go/
├── dmsoft.go           # 接口定义文件
├── dmsoft_impl.go      # 接口实现文件（428个函数）
├── go.mod              # Go 模块定义
├── Go.dll              # 破解补丁
├── xd47243.dll         # 大漠插件主文件
├── example/            # 使用示例
│   ├── main.go         # 示例代码
│   └── go.mod          # 示例模块
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
    "syscall"
    
    dmsoft "gitee.com/yuan71058/dm72424-go"
)

func main() {
    // 1. 加载大漠插件
    dmHModule, err := dmsoft.Load("xd47243.dll")
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. 执行破解（可选）
    goHModule, _ := syscall.LoadLibrary("Go.dll")
    goFunAddr, _ := syscall.GetProcAddress(goHModule, "Go")
    syscall.SyscallN(uintptr(goFunAddr), dmHModule)
    defer syscall.FreeLibrary(goHModule)
    
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
go get gitee.com/yuan71058/dm72424-go
```

### 方法二：本地引用

```bash
# 克隆项目
git clone https://gitee.com/yuan71058/dm72424-go.git

# 在你的 go.mod 中添加 replace
replace gitee.com/yuan71058/dm72424-go => ./dm72424-go
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
| `Load(path string) (uintptr, error)` | 加载大漠插件 DLL |
| `Free() bool` | 释放大漠插件 |
| `New() *DmSoft` | 创建大漠对象 |
| `Init()` | 初始化对象 |
| `Release()` | 释放对象 |

### 常用功能示例

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

大漠插件返回的字符串可能是 GBK 编码：

```go
import "golang.org/x/text/encoding/simplifiedchinese"

func gbkToUtf8(s string) string {
    data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(s))
    return string(data)
}
```

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

### v1.0.0 (2024-03-18)

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

- Gitee: [https://gitee.com/yuan71058/dm72424-go](https://gitee.com/yuan71058/dm72424-go)

---

<p align="center">
  如果这个项目对你有帮助，请给一个 ⭐️ Star 支持一下！
</p>
