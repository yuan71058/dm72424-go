# 大漠插件 Go 语言封装库

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![Arch](https://img.shields.io/badge/Arch-386%20%7C%20amd64-orange?style=flat)]()
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

> 大漠插件 7.2424 破解版本的 Go 语言封装库，支持 428 个函数接口，**完整支持32位和64位架构**，开箱即用！

---

## 特性

- **完整封装** - 支持大漠插件 7.2424 版本全部 428 个函数
- **双架构支持** - 完整支持32位(x86)和64位(x64/amd64)架构
- **开箱即用** - 简单导入即可使用，无需复杂配置
- **详细注释** - 所有函数都有完整的中文注释
- **类型安全** - 完整的类型定义，编译时检查
- **示例丰富** - 提供完整的32位和64位使用示例
- **自动编码转换** - 内置 UTF-8 到 GBK 自动转换，中文参数无需手动处理

---

## 目录

- [架构选择](#架构选择)
- [环境要求](#环境要求)
- [快速开始](#快速开始)
  - [32位模式](#32位模式)
  - [64位模式](#64位模式)
- [安装方法](#安装方法)
- [核心函数](#核心函数)
- [使用示例](#使用示例)
- [多线程操作](#多线程操作)
- [注意事项](#注意事项)
- [常见问题](#常见问题)
- [更新日志](#更新日志)

---

## 架构选择

本库支持两种编译架构，请根据需求选择：

| 架构 | 编译目标 | 适用场景 | 性能 | 文档 |
|------|---------|---------|------|------|
| **32位** | `GOARCH=386` | 传统应用、高性能需求 | 微秒级 | 本文档 |
| **64位** | `GOARCH=amd64` | 现代应用、大内存需求 | 毫秒级 | [README_x64.md](README_x64.md) |

### 如何选择？

- 如果你的主程序是32位的，或者对性能要求极高 → 选择 **32位模式**
- 如果你的主程序是64位的，或者需要访问超过4GB内存 → 选择 **64位模式**

> **重要**：大漠插件DLL本身是32位的。64位模式通过启动独立的32位helper进程来调用DLL，通过TCP+gob进行跨进程通信。

---

## 环境要求

- **操作系统**: Windows（32位或64位均可）
- **Go 版本**: Go 1.16 或更高版本
- **编译要求**: 
  - 32位模式：必须编译为 32 位程序 (`GOARCH=386`)
  - 64位模式：必须编译为 64 位程序 (`GOARCH=amd64`)，并需要helper进程

---

## 快速开始

### 32位模式

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
}
```

**编译命令**：
```powershell
$env:GOARCH="386"; go build -o myapp.exe
```

### 64位模式

64位模式需要先编译helper进程，然后编译主程序：

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    dmsoft "github.com/yuan71058/dm72424-go"
)

func main() {
    absPluginPath, _ := filepath.Abs("xd47243.dll")
    absCrackPath, _ := filepath.Abs("Go.dll")

    // 1. 加载大漠插件（记录路径）
    _, err := dmsoft.LoadDm(absPluginPath)
    if err != nil {
        log.Fatal(err)
    }

    // 2. 设置破解DLL路径
    err = dmsoft.CrackDm(absCrackPath)
    if err != nil {
        log.Fatal(err)
    }

    // 3. 创建对象并初始化（启动helper进程）
    dm := dmsoft.New()
    err = dm.Init()
    if err != nil {
        log.Fatal(err)
    }
    defer dm.Release()

    // 4. 注册（每个helper进程都需要注册）
    if dm.Reg("", "") == 1 {
        fmt.Println("注册成功")
    }

    // 5. 正常使用，API与32位完全一致
    fmt.Printf("版本: %s\n", dm.Ver())
}
```

**编译命令**：
```powershell
# 1. 编译32位helper进程
cd cmd\dm_com_server
$env:GOARCH="386"; $env:CGO_ENABLED="0"; go build -o dm_com_server.exe .
cd ..\..

# 2. 编译64位主程序
$env:GOARCH="amd64"; go build -o myapp.exe
```

> 📖 **完整64位文档**：请参考 [README_x64.md](README_x64.md) 获取详细的64位使用说明、架构原理和多线程示例。

---

## 安装方法

### 方法一：Go Modules（推荐）

```bash
go get github.com/yuan71058/dm72424-go
```

### 方法二：本地引用

```bash
git clone https://github.com/yuan71058/dm72424-go.git
```

在 go.mod 中添加：
```
replace github.com/yuan71058/dm72424-go => ./dm72424-go
```

### 编译说明

**32位模式**：大漠插件是 32 位 DLL，必须编译为 32 位程序！

```powershell
go env -w GOARCH=386
go build -o myapp.exe
```

**64位模式**：需要先编译helper进程，再编译主程序

```powershell
# 编译helper进程（32位）
cd cmd\dm_com_server
$env:GOARCH="386"; $env:CGO_ENABLED="0"; go build -o dm_com_server.exe .
cd ..\..

# 编译主程序（64位）
$env:GOARCH="amd64"; go build -o myapp.exe
```

> 📖 64位模式详细说明请参考 [README_x64.md](README_x64.md)

---

## 核心函数

### 插件管理

| 函数 | 说明 |
|------|------|
| `LoadDm(path)` | 加载大漠插件 DLL |
| `CrackDm(path)` | 破解大漠插件 |
| `FreeCrackDll()` | 释放破解 DLL |
| `Free()` | 释放大漠插件 |

### 对象管理

| 函数 | 说明 |
|------|------|
| `New()` | 创建大漠对象 |
| `Init()` | 初始化对象 |
| `Release()` | 释放对象 |

---

## 使用示例

### 窗口操作

```go
// 查找窗口（支持中文标题）
hwnd := dm.FindWindow("Qt51514QWindowIcon", "朋友圈")
if hwnd > 0 {
    fmt.Printf("找到窗口，句柄: %d\n", hwnd)
    
    // 绑定窗口
    ret := dm.BindWindow(hwnd, "gdi", "windows3", "windows", 0)
    if ret == 1 {
        fmt.Println("绑定成功")
        dm.UnBindWindow()
    }
}
```

### 截图功能

```go
// 设置保存路径
dm.SetPath("C:\\screenshots")

// 截取全屏
dm.Capture(0, 0, 1920, 1080, "screen.bmp")

// JPG格式截图
dm.CaptureJpg(0, 0, 1920, 1080, "screen.jpg", 80)
```

### 找图功能

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

### 找色功能

```go
// 查找颜色
var x, y int32
ret := dm.FindColor(0, 0, 1920, 1080, "FF0000", "000000", 1.0, 0, &x, &y)

// 多点找色
ret = dm.FindMultiColor(0, 0, 1920, 1080, "FF0000", "5|0|00FF00", "000000", 1.0, 0, &x, &y)
```

### 文字识别（OCR）

```go
// 设置字库
dm.SetDict(0, "dict.txt")

// OCR识别
text := dm.Ocr(0, 0, 500, 100, "FFFFFF-000000", 1.0)

// 查找文字
var x, y int32
ret := dm.FindStr(0, 0, 1920, 1080, "登录", "FFFFFF-000000", 1.0, &x, &y)
```

### 鼠标操作

```go
// 移动鼠标
dm.MoveTo(500, 300)

// 点击
dm.LeftClick()
dm.RightClick()

// 获取位置
var x, y int32
dm.GetCursorPos(&x, &y)
```

### 键盘操作

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

### 内存操作

```go
// 读取内存
value := dm.ReadInt(hwnd, 0x12345678, 0)

// 写入内存
dm.WriteInt(hwnd, 0x12345678, 0, 12345)

// 查找特征码
addr := dm.FindData(hwnd, 0x400000, 0x500000, "FF ?? 00 ??")
```

---

## 多线程操作

### 32位模式

- **主对象**：全局唯一，负责注册，最后释放
- **子对象**：每线程独立创建，各自 Init/Release
- **注册**：只需在主对象中注册一次
- **绑定**：每个子对象独立绑定自己的目标窗口

### 64位模式

- **独立进程**：每个 `DmSoft` 对象启动独立的 helper 进程
- **并行执行**：多个 helper 进程互不干扰，真正并行
- **注册要求**：每个 helper 进程必须单独调用 `Reg()`
- **资源占用**：每个 helper 进程约占用 5-10MB 内存

### 32位多线程示例

```go
// 主线程：创建主对象并注册
mainDm := dmsoft.New()
mainDm.Init()
mainDm.Reg("", "")  // 只需注册一次

// 子线程：每个线程创建独立对象
for i := 0; i < 3; i++ {
    go func() {
        subDm := dmsoft.New()
        subDm.Init()  // 独立初始化
        subDm.BindWindow(hwnd, "gdi", "windows", "windows", 0)
        // ... 操作 ...
        subDm.UnBindWindow()
        subDm.Release()  // 各自释放
    }()
}

defer mainDm.Release()  // 最后释放主对象
```

### 64位多线程示例

```go
// 每个goroutine创建独立对象和helper进程
for i := 0; i < 3; i++ {
    go func(id int) {
        dm := dmsoft.New()
        dm.Init()
        dm.Reg("", "")  // 64位模式下每个helper必须单独注册！
        defer dm.Release()

        dm.BindWindow(hwnd, "gdi", "windows", "windows", 0)
        // ... 操作 ...
        dm.UnBindWindow()
    }(i)
}
```

完整示例请参考：
- 32位：`example/multithread/main.go`
- 64位：`example/x64_mt/main.go`

---

## 函数分类

| 分类 | 数量 | 主要函数 |
|------|------|----------|
| 窗口操作 | ~50 | BindWindow, FindWindow, GetWindowRect |
| 鼠标操作 | ~20 | MoveTo, LeftClick, GetCursorPos |
| 键盘操作 | ~15 | KeyPress, KeyDown, SendString |
| 图像处理 | ~30 | Capture, FindPic, LoadPic |
| 颜色操作 | ~15 | GetColor, FindColor, CmpColor |
| OCR识别 | ~20 | Ocr, FindStr, SetDict |
| 内存操作 | ~40 | ReadInt, WriteInt, FindData |
| 系统信息 | ~20 | Ver, GetOsType, GetTime |
| 文件操作 | ~15 | ReadFile, WriteFile, IsFileExist |
| AI功能 | ~10 | LoadAi, FindPicAi |

---

## 项目结构

```
go-dm72424/
├── dmsoft.go              # DmSoftInterface + DmSoftBase + GBK工具函数
├── dmsoft_impl.go         # DmSoft实现 + LoadDm/CrackDm/Init/Release
├── dm_x64_pipe.go         # TCP客户端 + 方法偏移量表 + pipeCall
├── dm_x64_helpers.go      # comCall→pipeCall 桥接层
├── cmd/
│   └── dm_com_server/     # 32位helper进程（64位模式必需）
│       └── main.go        # TCP服务器 + 偏移量调用dm.dll
├── example/
│   ├── main.go            # 32位基础示例
│   ├── multithread/       # 32位多线程示例
│   ├── x64/               # 64位基础示例
│   └── x64_mt/            # 64位多线程示例
├── README.md              # 本文档（32位模式为主）
└── README_x64.md          # 64位模式详细文档
```

---

## 注意事项

### 1. 初始化顺序

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

本库已内置 UTF-8 到 GBK 的自动编码转换：

```go
// 直接使用中文参数，无需手动转换
hwnd := dm.FindWindow("Qt51514QWindowIcon", "朋友圈")
dm.SetPath("C:\\测试目录")
```

---

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

### Q4: 编译报错 "not a valid Win32 application"？

- 确保编译为32位程序：`go env -w GOARCH=386`
- 大漠插件是32位DLL，不能在64位程序中调用

### Q5: 如何在64位程序中使用大漠插件？

- 使用本库的64位支持方案，通过TCP+gob跨进程调用
- 先编译32位helper进程：`GOARCH=386 go build -o dm_com_server.exe cmd/dm_com_server/`
- 再编译64位主程序：`GOARCH=amd64 go build`
- 详细说明请参考 [README_x64.md](README_x64.md)

### Q6: 64位模式性能如何？

- 64位模式通过TCP跨进程调用，延迟约毫秒级（32位为微秒级）
- 对于大多数自动化场景（找图、找色、OCR等），性能差异可忽略
- 高频调用场景（如每秒数千次）建议使用32位模式

---

## 更新日志

### v1.6.0 (2026-03-22)

- 修复返回字符串的编码转换问题
- 新增 GBK → UTF-8 转换函数
- 验证所有 428 个函数的编码处理

### v1.5.0 (2026-03-21)

- 修复 32 位 int64 参数传递问题
- 修复 19 个函数的参数传递
- 完成所有函数参数验证

### v1.4.0 (2026-03-21)

- 修复 32 位 float64/float32 参数传递问题
- 修复 10 个函数的浮点参数传递

### v1.3.0 (2026-03-20)

- 新增自动编码转换
- 新增 GetDmHModule() 函数
- 新增 GetObjPtr() 方法

### v1.2.0 (2026-03-19)

- 新增 LoadDm() 函数
- 新增 CrackDm() 函数
- 新增 FreeCrackDll() 函数

### v1.1.0 (2026-03-18)

- 添加多线程示例
- 修复 Release 函数问题

### v1.0.0 (2026-03-18)

- 初始版本发布
- 完成 428 个函数的翻译

---

## 许可证

本项目仅供学习交流使用，请勿用于商业用途。

大漠插件版权归大漠插件作者所有。

---

## 联系方式

GitHub: [https://github.com/yuan71058/dm72424-go](https://github.com/yuan71058/dm72424-go)

---

如果这个项目对你有帮助，请给一个 ⭐️ Star 支持一下！
