# 大漠插件 64位调用方案

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://www.microsoft.com/windows)
[![Arch](https://img.shields.io/badge/Arch-amd64%20%7C%20386-orange?style=flat)]()

> 大漠插件 7.2424 的 64 位 Go 调用方案，通过 TCP+gob 跨进程通信实现 64 位程序调用 32 位 DLL。

---

## 方案原理

大漠插件（dm.dll）是 32 位 DLL，无法在 64 位进程中直接加载。本方案通过 **TCP+gob 跨进程调用** 解决此问题：

```
┌─────────────────────┐         TCP+gob         ┌─────────────────────┐
│   64位主进程          │ ◄──────────────────────► │   32位helper进程     │
│   (GOARCH=amd64)     │    127.0.0.1:PORT       │   (GOARCH=386)      │
│                      │                          │                     │
│   DmSoft对象          │   CallRequest/Response   │   dm.dll + Go.dll   │
│   pipeCall() ────────┤ ───────────────────────► │   handleCall()      │
│   gob编码/解码        │                          │   偏移量直接调用      │
└─────────────────────┘                          └─────────────────────┘
```

### 工作流程

1. **LoadDm** — 记录 DLL 路径，定位 `dm_com_server.exe`
2. **CrackDm** — 记录破解 DLL 路径（实际破解在 helper 进程中执行）
3. **New + Init** — 启动 32 位 helper 进程，建立 TCP 连接
4. **方法调用** — 通过 TCP+gob 序列化转发到 helper 执行
5. **Release** — 关闭 TCP 连接，终止 helper 进程

### 参数类型协议

| Type | 类型 | 32位栈槽数 | 说明 |
|------|------|-----------|------|
| 0 | int32 | 1 | 整数参数 |
| 1 | string | 1 | GBK编码C字符串指针 |
| 2 | float64 | 2 | double，低位在前高位在后 |
| 3 | *int32 | 1 | 输出参数指针 |
| 4 | float32 | 1 | 单精度浮点 |
| 5 | int64 | 2 | 64位整数，低位在前高位在后 |

---

## 快速开始

### 编译

```powershell
# 1. 编译32位helper进程
cd cmd\dm_com_server
$env:GOARCH="386"; $env:CGO_ENABLED="0"; go build -o dm_com_server.exe .
cd ..\..

# 2. 编译64位示例
cd example\x64
$env:GOARCH="amd64"; go build -o x64_demo.exe .
```

### 运行

确保以下文件在同一目录：
- `x64_demo.exe` — 64位主程序
- `dm_com_server.exe` — 32位helper进程
- `xd47243.dll` — 大漠插件DLL
- `Go.dll` — 破解DLL

```powershell
# 以管理员身份运行
x64_demo.exe
```

### 代码示例

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    dmsoft "github.com/yuan71058/dm72424-go"
)

func main() {
    absPluginPath, _ := filepath.Abs("xd47243.dll")
    absCrackPath, _ := filepath.Abs("Go.dll")

    // 1. 加载大漠插件（记录路径，不实际加载）
    _, err := dmsoft.LoadDm(absPluginPath)
    if err != nil {
        log.Fatal(err)
    }

    // 2. 设置破解DLL路径（必须在Init之前）
    err = dmsoft.CrackDm(absCrackPath)
    if err != nil {
        log.Fatal(err)
    }

    // 3. 创建对象并初始化（启动helper + TCP连接）
    dm := dmsoft.New()
    err = dm.Init()
    if err != nil {
        log.Fatal(err)
    }
    defer dm.Release()

    // 4. 注册
    if dm.Reg("", "") == 1 {
        fmt.Println("注册成功")
    }

    // 5. 正常使用，API与32位完全一致
    fmt.Printf("版本: %s\n", dm.Ver())
    fmt.Printf("分辨率: %d x %d\n", dm.GetScreenWidth(), dm.GetScreenHeight())

    // 找图（含float64参数和输出参数）
    var x, y int32
    ret := dm.FindPic(0, 0, 1920, 1080, "test.bmp", "000000", 0.9, 0, &x, &y)

    // 剪贴板（自动GBK↔UTF-8转换）
    dm.SetClipboard("Hello世界")
    text := dm.GetClipboard()
}
```

---

## 多线程

64位模式下，每个 `DmSoft` 对象启动**独立的 helper 进程**，互不干扰，真正并行。

```go
// 每个goroutine创建独立对象
for i := 0; i < 3; i++ {
    go func(id int) {
        dm := dmsoft.New()
        dm.Init()
        dm.Reg("", "")  // 每个helper进程必须单独注册！
        defer dm.Release()

        dm.BindWindow(hwnd, "gdi", "windows", "windows", 0)
        // ... 操作 ...
        dm.UnBindWindow()
    }(i)
}
```

> **重要**：64位多线程中，每个子对象必须调用 `Reg()`。与32位不同，64位每个 helper 是独立进程，不共享 dm.dll 全局状态。

完整示例：`example/x64_mt/main.go`

---

## 项目结构

```
go-dm72424/
├── dmsoft.go              # DmSoftInterface + DmSoftBase + GBK工具函数
├── dmsoft_impl.go         # DmSoft实现 + LoadDm/CrackDm/Init/Release
├── dm_x64_pipe.go         # TCP客户端 + 426个方法偏移量表 + pipeCall
├── dm_x64_helpers.go      # comCall→pipeCall 桥接层
├── cmd/
│   └── dm_com_server/     # 32位helper进程
│       └── main.go        # TCP服务器 + 偏移量调用dm.dll
├── example/
│   ├── main.go            # 32位基础示例
│   ├── multithread/       # 32位多线程示例
│   ├── x64/               # 64位基础示例
│   └── x64_mt/            # 64位多线程示例
```

---

## 32位 vs 64位对比

| | 32位模式 | 64位模式 |
|---|---|---|
| GOARCH | 386 | amd64 |
| DLL加载 | 直接 `syscall.LoadLibrary` | helper进程加载 |
| 方法调用 | 偏移量直接 `syscall.Syscall` | TCP+gob 跨进程转发 |
| 编码转换 | 进程内 GBK↔UTF-8 | TCP传输自动处理 |
| 多线程 | 共享进程，子对象无需Reg | 独立helper，每个必须Reg |
| 输出参数 | 直接指针 | gob序列化回传 |
| 性能 | ~微秒级 | ~毫秒级（TCP开销） |
| 资源占用 | 1个进程 | N+1个进程（N=DmSoft对象数） |

---

## 常见问题

### Q1: Init 失败 "未找到dm_com_server.exe"

确保 `dm_com_server.exe` 已编译（GOARCH=386）且在以下位置之一：
- 与主程序同目录
- `../cmd/dm_com_server/dm_com_server.exe`

### Q2: Init 失败 "helper进程无输出"

- 检查 `dm_com_server.exe` 是否为32位编译（`GOARCH=386`）
- 检查 `xd47243.dll` 和 `Go.dll` 路径是否正确
- 以管理员身份运行

### Q3: 方法调用崩溃

- 确保每个 DmSoft 对象都调用了 `Reg()`
- 检查 helper 进程是否存活
- 查看 helper 进程的 stderr 输出

### Q4: 中文乱码

TCP 传输自动处理 GBK↔UTF-8 编码转换，无需手动处理。如仍有问题，请检查：
- 字库文件编码
- `SetDict` 加载的字库是否为 GBK 编码

### Q5: 性能优化

- 减少不必要的调用次数
- 使用 `Capture` 缓存截图，避免重复截图
- 多线程场景下，合理控制 helper 进程数量
- 对于高频调用，考虑使用32位模式直接调用

---

## 许可证

本项目仅供学习交流使用，请勿用于商业用途。

大漠插件版权归大漠插件作者所有。
