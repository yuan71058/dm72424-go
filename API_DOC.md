# 大漠插件 Go 绑定 API 文档

> 版本: v1.6.0  
> 更新日期: 2026-03-22  
> 作者: dm72424-go 项目组

---

## 目录

1. [初始化与基础函数](#1-初始化与基础函数)
2. [窗口操作函数](#2-窗口操作函数)
3. [鼠标操作函数](#3-鼠标操作函数)
4. [键盘操作函数](#4-键盘操作函数)
5. [找图函数](#5-找图函数)
6. [找色函数](#6-找色函数)
7. [找字函数](#7-找字函数)
8. [OCR文字识别函数](#8-ocr文字识别函数)
9. [内存操作函数](#9-内存操作函数)
10. [文件操作函数](#10-文件操作函数)
11. [进程操作函数](#11-进程操作函数)
12. [屏幕截图函数](#12-屏幕截图函数)
13. [Foobar绘图函数](#13-foobar绘图函数)
14. [AI相关函数](#14-ai相关函数)
15. [汇编相关函数](#15-汇编相关函数)
16. [网络相关函数](#16-网络相关函数)
17. [INI配置文件函数](#17-ini配置文件函数)
18. [字库相关函数](#18-字库相关函数)
19. [杂项函数](#19-杂项函数)
20. [64位(x64)支持说明](#20-64位x64支持说明)

---

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    dmsoft "github.com/yuan71058/dm72424-go"
)

func main() {
    // 1. 加载大漠插件DLL
    dmHModule, err := dmsoft.LoadDm("xd47243.dll")
    if err != nil {
        log.Fatalf("加载大漠插件失败: %v", err)
    }

    // 2. 破解大漠插件
    err = dmsoft.CrackDm("Go.dll")
    if err != nil {
        log.Fatalf("破解大漠插件失败: %v", err)
    }

    // 3. 创建大漠对象
    dm := dmsoft.New()
    dm.Init()           // 必须调用！初始化COM对象
    defer dm.Release()  // 程序结束时释放

    // 4. 使用大漠功能
    var x, y int32
    ret := dm.GetCursorPos(&x, &y)
    fmt.Printf("鼠标位置: (%d, %d), 返回值: %d\n", x, y, ret)
}
```

---

## 1. 初始化与基础函数

### LoadDm

**功能说明**: 加载大漠插件DLL

**函数签名**:
```go
func LoadDm(dmPath string) (uintptr, error)
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| dmPath | string | 大漠DLL文件路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| uintptr | DLL模块句柄 |
| error | 错误信息，成功时为nil |

**示例**:
```go
module, err := dmsoft.LoadDm("dm.dll")
if err != nil {
    fmt.Printf("加载大漠DLL失败: %v\n", err)
    return
}
```

---

### CrackDm

**功能说明**: 加载破解DLL并激活大漠插件

**函数签名**:
```go
func CrackDm(crackDllPath string) error
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| crackDllPath | string | 破解DLL文件路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| error | 错误信息，成功时为nil |

**示例**:
```go
err := dmsoft.CrackDm("dm_crack.dll")
if err != nil {
    fmt.Printf("破解失败: %v\n", err)
    return
}
```

---

### FreeCrackDll

**功能说明**: 释放破解DLL

**函数签名**:
```go
func FreeCrackDll() bool
```

**返回值**:
| 类型 | 说明 |
|------|------|
| bool | 成功返回true，失败返回false |

---

### Load

**功能说明**: 加载DLL（通用方法）

**函数签名**:
```go
func Load(path string) (uintptr, error)
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | DLL文件路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| uintptr | DLL模块句柄 |
| error | 错误信息，成功时为nil |

---

### Free

**功能说明**: 释放DLL

**函数签名**:
```go
func Free() bool
```

**返回值**:
| 类型 | 说明 |
|------|------|
| bool | 成功返回true，失败返回false |

---

### New

**功能说明**: 创建大漠插件实例（非线程安全）

**函数签名**:
```go
func New() *DmSoft
```

**返回值**:
| 类型 | 说明 |
|------|------|
| *DmSoft | 大漠插件实例指针，失败返回nil |

**注意事项**:
- 单线程场景下全局创建一次即可
- 多线程场景下每个线程需独立创建并调用Init()

**示例**:
```go
dm := dmsoft.New()
if dm == nil {
    fmt.Println("创建大漠实例失败")
    return
}
```

---

### Init

**功能说明**: 初始化大漠对象，创建内部COM对象实例

**函数签名**:
```go
func (dm *DmSoft) Init()
```

**调用场景说明**:
- 单线程/全局使用: 全局只需调用一次
- 多线程场景: 每个线程需独立创建DmSoft实例并各自调用Init()

**示例**:
```go
dm := dmsoft.New()
dm.Init()
defer dm.Release()
```

---

### Release

**功能说明**: 释放大漠对象，销毁内部COM对象实例

**函数签名**:
```go
func (dm *DmSoft) Release()
```

**注意事项**:
- 释放后该DmSoft实例不可再用
- 需重新New()和Init()创建新实例

---

### Ver

**功能说明**: 获取大漠插件版本号

**函数签名**:
```go
func (dm *DmSoft) Ver() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 版本号字符串 |

**示例**:
```go
version := dm.Ver()
fmt.Printf("大漠版本: %s\n", version)
```

---

### GetID

**功能说明**: 获取大漠ID

**函数签名**:
```go
func (dm *DmSoft) GetID() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | ID值 |

---

### GetLastError

**功能说明**: 获取最后一次错误码

**函数签名**:
```go
func (dm *DmSoft) GetLastError() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 错误码 |

---

### GetMachineCode

**功能说明**: 获取机器码

**函数签名**:
```go
func (dm *DmSoft) GetMachineCode() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 机器码字符串 |

---

### GetMachineCodeNoMac

**功能说明**: 获取机器码（不含MAC地址）

**函数签名**:
```go
func (dm *DmSoft) GetMachineCodeNoMac() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 机器码字符串 |

---

### GetMac

**功能说明**: 获取本机MAC地址

**函数签名**:
```go
func (dm *DmSoft) GetMac() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | MAC地址字符串 |

---

### Reg

**功能说明**: 注册大漠插件

**函数签名**:
```go
func (dm *DmSoft) Reg(code string, ver string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | string | 注册码 |
| ver | string | 版本号 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
ret := dm.Reg("注册码", "版本号")
if ret == 1 {
    fmt.Println("注册成功")
}
```

---

### RegEx

**功能说明**: 扩展注册大漠插件

**函数签名**:
```go
func (dm *DmSoft) RegEx(code string, ver string, ip string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | string | 注册码 |
| ver | string | 版本号 |
| ip | string | IP地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetDmCount

**功能说明**: 获取大漠对象数量

**函数签名**:
```go
func (dm *DmSoft) GetDmCount() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 大漠对象数量 |

---

### SetDisplayInput

**功能说明**: 设置显示输入

**函数签名**:
```go
func (dm *DmSoft) SetDisplayInput(display_input int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| display_input | int32 | 显示输入类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetShowErrorMsg

**功能说明**: 设置显示错误信息

**函数签名**:
```go
func (dm *DmSoft) SetShowErrorMsg(show int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| show | int32 | 1显示错误信息，0不显示 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SpeedNormalGraphic

**功能说明**: 加速普通图形

**函数签名**:
```go
func (dm *DmSoft) SpeedNormalGraphic(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 5. 找图函数

### FindPic

**功能说明**: 在指定区域查找图片

**函数签名**:
```go
func (dm *DmSoft) FindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名，多个用\|分隔 |
| delta_color | string | 颜色偏差，如"203040" |
| sim | float64 | 相似度(0.1-1.0) |
| dir | int32 | 查找方向: 0从左到右从上到下, 1从左到右从下到上, 2从右到左从上到下, 3从右到左从下到上 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引(从0开始)，失败返回-1 |

**示例**:
```go
var x, y int32
idx := dm.FindPic(0, 0, 800, 600, "test.bmp|test2.bmp", "000000", 0.9, 0, &x, &y)
if idx >= 0 {
    fmt.Printf("找到第%d张图片，坐标: (%d, %d)\n", idx, x, y)
}
```

---

### FindPicEx

**功能说明**: 扩展找图，返回所有找到的坐标

**函数签名**:
```go
func (dm *DmSoft) FindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 所有找到的坐标，格式: "id,x,y\|id,x,y..." |

---

### FindPicE

**功能说明**: 找图并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindPicE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y"，失败返回"-1" |

---

### FindPicS

**功能说明**: 找图并返回相似度

**函数签名**:
```go
func (dm *DmSoft) FindPicS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 相似度字符串 |

---

### FindPicExS

**功能说明**: 扩展找图并返回相似度

**函数签名**:
```go
func (dm *DmSoft) FindPicExS(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y,sim\|id,x,y,sim..." |

---

### FindPicMem

**功能说明**: 从内存中找图

**函数签名**:
```go
func (dm *DmSoft) FindPicMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据（Base64编码） |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引，失败返回-1 |

---

### FindPicMemEx

**功能说明**: 从内存中找图（扩展）

**函数签名**:
```go
func (dm *DmSoft) FindPicMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 所有找到的坐标 |

---

### FindPicMemE

**功能说明**: 从内存中找图并返回字符串

**函数签名**:
```go
func (dm *DmSoft) FindPicMemE(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y" |

---

### FindPicSim

**功能说明**: 相似度找图

**函数签名**:
```go
func (dm *DmSoft) FindPicSim(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率(0-100) |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引，失败返回-1 |

---

### FindPicSimMem

**功能说明**: 从内存中相似度找图

**函数签名**:
```go
func (dm *DmSoft) FindPicSimMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引，失败返回-1 |

---

### FindPicSimE

**功能说明**: 相似度找图并返回字符串

**函数签名**:
```go
func (dm *DmSoft) FindPicSimE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y" |

---

### FindPicSimEx

**功能说明**: 相似度找图并返回所有坐标

**函数签名**:
```go
func (dm *DmSoft) FindPicSimEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim int32, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y\|id,x,y..." |

---

### FindMultiPic

**功能说明**: 查找多张图片

**函数签名**:
```go
func (dm *DmSoft) FindMultiPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名，多个用\|分隔 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引 |

---

### FindMultiPicE

**功能说明**: 查找多张图片并返回字符串

**函数签名**:
```go
func (dm *DmSoft) FindMultiPicE(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, delta_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| delta_color | string | 颜色偏差 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y" |

---

### SetPicPwd

**功能说明**: 设置图片密码

**函数签名**:
```go
func (dm *DmSoft) SetPicPwd(pwd string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### GetPicSize

**功能说明**: 获取图片大小

**函数签名**:
```go
func (dm *DmSoft) GetPicSize(pic_name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_name | string | 图片名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "width,height" |

---

### LoadPic

**功能说明**: 预加载图片到内存

**函数签名**:
```go
func (dm *DmSoft) LoadPic(pic_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_name | string | 图片名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FreePic

**功能说明**: 释放预加载的图片

**函数签名**:
```go
func (dm *DmSoft) FreePic(pic_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_name | string | 图片名，空字符串释放所有 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### LoadPicByte

**功能说明**: 从内存加载图片

**函数签名**:
```go
func (dm *DmSoft) LoadPicByte(pic_info string, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_info | string | 图片数据（Base64编码） |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AppendPicAddr

**功能说明**: 追加图片地址

**函数签名**:
```go
func (dm *DmSoft) AppendPicAddr(pic_info string, addr int32, size int32, size1 int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_info | string | 图片数据 |
| addr | int32 | 内存地址 |
| size | int32 | 大小 |
| size1 | int32 | 大小2 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnablePicCache

**功能说明**: 启用图片缓存

**函数签名**:
```go
func (dm *DmSoft) EnablePicCache(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FindPicSimMemE

**功能说明**: 从内存中相似度找图并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindPicSimMemE(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率(0-100) |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y" |

---

### FindPicSimMemEx

**功能说明**: 从内存中相似度找图（扩展）

**函数签名**:
```go
func (dm *DmSoft) FindPicSimMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, delta_color string, sim int32, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| delta_color | string | 颜色偏差 |
| sim | int32 | 最小百分比相似率 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "id,x,y|id,x,y..." |

---

### EnableDisplayDebug

**功能说明**: 启用显示调试

**函数签名**:
```go
func (dm *DmSoft) EnableDisplayDebug(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableFindPicMultithread

**功能说明**: 启用多线程找图

**函数签名**:
```go
func (dm *DmSoft) EnableFindPicMultithread(enable int32, count int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |
| count | int32 | 线程数量 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### IsDisplayDead

**功能说明**: 判断显示器是否失效

**函数签名**:
```go
func (dm *DmSoft) IsDisplayDead() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 1失效，0正常 |

---

### MatchPicName

**功能说明**: 匹配图片名称

**函数签名**:
```go
func (dm *DmSoft) MatchPicName(pic_name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_name | string | 图片名称 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 匹配到的图片名称 |

---

### SetExcludeRegion

**功能说明**: 设置排除区域

**函数签名**:
```go
func (dm *DmSoft) SetExcludeRegion(x1 int32, y1 int32, x2 int32, y2 int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetFindPicMultithreadCount

**功能说明**: 设置多线程找图数量

**函数签名**:
```go
func (dm *DmSoft) SetFindPicMultithreadCount(count int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| count | int32 | 线程数量 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetFindPicMultithreadLimit

**功能说明**: 设置多线程找图限制

**函数签名**:
```go
func (dm *DmSoft) SetFindPicMultithreadLimit(limit int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| limit | int32 | 限制数量 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 6. 找色函数

### GetColor

**功能说明**: 获取指定坐标颜色

**函数签名**:
```go
func (dm *DmSoft) GetColor(x int32, y int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 颜色值，格式: "RRGGBB" |

**示例**:
```go
color := dm.GetColor(100, 200)
fmt.Printf("颜色: %s\n", color)
```

---

### GetColorBGR

**功能说明**: 获取指定坐标颜色(BGR格式)

**函数签名**:
```go
func (dm *DmSoft) GetColorBGR(x int32, y int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 颜色值，格式: "BBGGRR" |

---

### FindColor

**功能说明**: 在指定区域查找颜色

**函数签名**:
```go
func (dm *DmSoft) FindColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值，格式: "RRGGBB-DRDGDB" |
| sim | float64 | 相似度(0.1-1.0) |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FindColorEx

**功能说明**: 扩展找色，返回所有找到的坐标

**函数签名**:
```go
func (dm *DmSoft) FindColorEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y\|x,y..." |

---

### FindColorE

**功能说明**: 找色并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindColorE(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y"，失败返回"-1" |

---

### FindMultiColor

**功能说明**: 多点找色

**函数签名**:
```go
func (dm *DmSoft) FindMultiColor(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| first_color | string | 第一个颜色 |
| offset_color | string | 偏移颜色，格式: "x1\|color1,x2\|color2..." |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FindMultiColorEx

**功能说明**: 扩展多点找色

**函数签名**:
```go
func (dm *DmSoft) FindMultiColorEx(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| first_color | string | 第一个颜色 |
| offset_color | string | 偏移颜色 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y\|x,y..." |

---

### FindMultiColorE

**功能说明**: 多点找色并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindMultiColorE(x1 int32, y1 int32, x2 int32, y2 int32, first_color string, offset_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| first_color | string | 第一个颜色 |
| offset_color | string | 偏移颜色 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y" |

---

### FindColorBlock

**功能说明**: 查找颜色块

**函数签名**:
```go
func (dm *DmSoft) FindColorBlock(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |
| count | int32 | 最小像素数 |
| width | int32 | 宽度 |
| height | int32 | 高度 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FindColorBlockEx

**功能说明**: 扩展查找颜色块

**函数签名**:
```go
func (dm *DmSoft) FindColorBlockEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64, count int32, width int32, height int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |
| count | int32 | 最小像素数 |
| width | int32 | 宽度 |
| height | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y\|x,y..." |

---

### CmpColor

**功能说明**: 比较颜色

**函数签名**:
```go
func (dm *DmSoft) CmpColor(x int32, y int32, color string, sim float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 相等返回1，不相等返回0 |

---

### SetFindColorMode

**功能说明**: 设置找色模式

**函数签名**:
```go
func (dm *DmSoft) SetFindColorMode(mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| mode | int32 | 模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### GetColorHSV

**功能说明**: 获取指定坐标颜色(HSV格式)

**函数签名**:
```go
func (dm *DmSoft) GetColorHSV(x int32, y int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | HSV颜色值 |

---

### GetAveHSV

**功能说明**: 获取区域平均HSV值

**函数签名**:
```go
func (dm *DmSoft) GetAveHSV(x1 int32, y1 int32, x2 int32, y2 int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 平均HSV值 |

---

### GetAveRGB

**功能说明**: 获取区域平均RGB值

**函数签名**:
```go
func (dm *DmSoft) GetAveRGB(x1 int32, y1 int32, x2 int32, y2 int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 平均RGB值 |

---

### GetColorNum

**功能说明**: 获取指定区域内的颜色数量

**函数签名**:
```go
func (dm *DmSoft) GetColorNum(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 颜色数量 |

---

### FindMulColor

**功能说明**: 查找多个颜色

**函数签名**:
```go
func (dm *DmSoft) FindMulColor(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色值，多个用|分隔 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到返回1，失败返回0 |

---

### RGB2BGR

**功能说明**: RGB颜色格式转BGR

**函数签名**:
```go
func (dm *DmSoft) RGB2BGR(rgb string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| rgb | string | RGB颜色值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | BGR颜色值 |

---

### BGR2RGB

**功能说明**: BGR颜色格式转RGB

**函数签名**:
```go
func (dm *DmSoft) BGR2RGB(bgr string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| bgr | string | BGR颜色值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | RGB颜色值 |

---

### EnableGetColorByCapture

**功能说明**: 启用截图取色

**函数签名**:
```go
func (dm *DmSoft) EnableGetColorByCapture(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FindShape

**功能说明**: 查找指定形状

**函数签名**:
```go
func (dm *DmSoft) FindShape(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| offset_color | string | 偏移颜色 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FindShapeE

**功能说明**: 查找形状并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindShapeE(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| offset_color | string | 偏移颜色 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y" |

---

### FindShapeEx

**功能说明**: 扩展查找形状

**函数签名**:
```go
func (dm *DmSoft) FindShapeEx(x1 int32, y1 int32, x2 int32, y2 int32, offset_color string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| offset_color | string | 偏移颜色 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "x,y|x,y..." |

---

## 7. 找字函数

### FindStr

**功能说明**: 在指定区域查找文字

**函数签名**:
```go
func (dm *DmSoft) FindStr(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的文字索引，失败返回-1 |

---

### FindStrEx

**功能说明**: 扩展找字，返回所有找到的坐标

**函数签名**:
```go
func (dm *DmSoft) FindStrEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y\|str,x,y..." |

---

### FindStrE

**功能说明**: 找字并返回坐标字符串

**函数签名**:
```go
func (dm *DmSoft) FindStrE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y" |

---

### FindStrS

**功能说明**: 找字并返回相似度

**函数签名**:
```go
func (dm *DmSoft) FindStrS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 相似度字符串 |

---

### FindStrFast

**功能说明**: 快速找字

**函数签名**:
```go
func (dm *DmSoft) FindStrFast(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的文字索引 |

---

### FindStrFastEx

**功能说明**: 扩展快速找字

**函数签名**:
```go
func (dm *DmSoft) FindStrFastEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y\|str,x,y..." |

---

### FindStrFastE

**功能说明**: 快速找字并返回字符串

**函数签名**:
```go
func (dm *DmSoft) FindStrFastE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y" |

---

### FindStrFastS

**功能说明**: 快速找字并返回相似度

**函数签名**:
```go
func (dm *DmSoft) FindStrFastS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, x *int32, y *int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 相似度字符串 |

---

### FindStrWithFont

**功能说明**: 使用指定字体找字

**函数签名**:
```go
func (dm *DmSoft) FindStrWithFont(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| font_name | string | 字体名称 |
| font_size | int32 | 字体大小 |
| flag | int32 | 标志 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的文字索引 |

---

### FindStrWithFontE

**功能说明**: 使用指定字体找字并返回字符串

**函数签名**:
```go
func (dm *DmSoft) FindStrWithFontE(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| font_name | string | 字体名称 |
| font_size | int32 | 字体大小 |
| flag | int32 | 标志 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y" |

---

### FindStrWithFontEx

**功能说明**: 使用指定字体扩展找字

**函数签名**:
```go
func (dm *DmSoft) FindStrWithFontEx(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| font_name | string | 字体名称 |
| font_size | int32 | 字体大小 |
| flag | int32 | 标志 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x,y\|str,x,y..." |

---

## 8. OCR文字识别函数

### Ocr

**功能说明**: OCR文字识别

**函数签名**:
```go
func (dm *DmSoft) Ocr(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 识别出的文字 |

**示例**:
```go
text := dm.Ocr(0, 0, 200, 50, "ffffff-000000", 0.9)
fmt.Printf("识别结果: %s\n", text)
```

---

### OcrEx

**功能说明**: 扩展OCR识别，返回坐标

**函数签名**:
```go
func (dm *DmSoft) OcrEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "文字,x,y\|文字,x,y..." |

---

### OcrExOne

**功能说明**: OCR识别单个文字

**函数签名**:
```go
func (dm *DmSoft) OcrExOne(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "文字,x,y" |

---

### FindStrFastWithFont

**功能说明**: 快速找字（指定字体）

**函数签名**:
```go
func (dm *DmSoft) FindStrFastWithFont(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64, font_name string, font_size int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |
| font_name | string | 字体名称 |
| font_size | int32 | 字体大小 |
| flag | int32 | 标志 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的文字索引 |

---

### FetchWord

**功能说明**: 提取文字

**函数签名**:
```go
func (dm *DmSoft) FetchWord(x1 int32, y1 int32, x2 int32, y2 int32, color string, word string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式 |
| word | string | 文字 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 提取的文字 |

---

### FetchWordEx

**功能说明**: 扩展提取文字

**函数签名**:
```go
func (dm *DmSoft) FetchWordEx(x1 int32, y1 int32, x2 int32, y2 int32, color string, word string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式 |
| word | string | 文字 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 提取的文字 |

---

### FindStrExS

**功能说明**: 扩展找字，返回所有找到的坐标（特殊格式）

**函数签名**:
```go
func (dm *DmSoft) FindStrExS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的文字 |
| color | string | 颜色格式 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "str,x0,y0|str,x1,y1|...|str,xn,yn" |

---

### OcrInFile

**功能说明**: 识别图片文件中指定区域内的文字

**函数签名**:
```go
func (dm *DmSoft) OcrInFile(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, color_format string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片文件名 |
| color_format | string | 颜色格式串 |
| sim | float64 | 相似度(0.1-1.0) |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 识别到的文字 |

---

### RegNoMac

**功能说明**: 注册大漠插件（不含MAC）

**函数签名**:
```go
func (dm *DmSoft) RegNoMac(code string, ver string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | string | 注册码 |
| ver | string | 版本号 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### RegExNoMac

**功能说明**: 扩展注册大漠插件（不含MAC）

**函数签名**:
```go
func (dm *DmSoft) RegExNoMac(code string, ver string, ip string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | string | 注册码 |
| ver | string | 版本号 |
| ip | string | IP地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetPath

**功能说明**: 设置资源文件路径

**函数签名**:
```go
func (dm *DmSoft) SetPath(path string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 资源路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
dm.SetPath("C:\\images")
```

---

### GetPath

**功能说明**: 获取当前资源路径

**函数签名**:
```go
func (dm *DmSoft) GetPath() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 路径字符串 |

---

### GetBasePath

**功能说明**: 获取大漠基础路径

**函数签名**:
```go
func (dm *DmSoft) GetBasePath() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 基础路径字符串 |

---

### GetRealPath

**功能说明**: 获取真实文件路径

**函数签名**:
```go
func (dm *DmSoft) GetRealPath(path string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 资源路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 真实路径字符串 |

---

### GetMemoryUsage

**功能说明**: 获取内存使用量

**函数签名**:
```go
func (dm *DmSoft) GetMemoryUsage() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 内存使用量(KB) |

---

### Delay

**功能说明**: 延迟指定时间

**函数签名**:
```go
func (dm *DmSoft) Delay(mis int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| mis | int32 | 延迟时间(毫秒) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### Delays

**功能说明**: 随机延迟

**函数签名**:
```go
func (dm *DmSoft) Delays(min_s int32, max_s int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| min_s | int32 | 最小延迟时间(毫秒) |
| max_s | int32 | 最大延迟时间(毫秒) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FindStrFastExS

**功能说明**: 快速查找文字扩展S

**函数签名**:
```go
func (dm *DmSoft) FindStrFastExS(x1 int32, y1 int32, x2 int32, y2 int32, str string, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| str | string | 要查找的字符串 |
| color | string | 颜色值 |
| sim | float64 | 相似度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 详细格式的搜索结果字符串 |

---

## 2. 窗口操作函数

### FindWindow

**功能说明**: 查找窗口

**函数签名**:
```go
func (dm *DmSoft) FindWindow(class_name string, title_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| class_name | string | 窗口类名，空字符串表示不限制 |
| title_name | string | 窗口标题，空字符串表示不限制 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄，失败返回0 |

**示例**:
```go
hwnd := dm.FindWindow("", "记事本")
if hwnd != 0 {
    fmt.Printf("找到窗口: %d\n", hwnd)
}
```

---

### FindWindowEx

**功能说明**: 查找子窗口

**函数签名**:
```go
func (dm *DmSoft) FindWindowEx(parent int32, class_name string, title_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| parent | int32 | 父窗口句柄 |
| class_name | string | 窗口类名 |
| title_name | string | 窗口标题 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### FindWindowSuper

**功能说明**: 高级查找窗口

**函数签名**:
```go
func (dm *DmSoft) FindWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| spec1 | string | 条件1字符串 |
| flag1 | int32 | 条件1标志 |
| type1 | int32 | 条件1类型 |
| spec2 | string | 条件2字符串 |
| flag2 | int32 | 条件2标志 |
| type2 | int32 | 条件2类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### FindWindowByProcess

**功能说明**: 通过进程名查找窗口

**函数签名**:
```go
func (dm *DmSoft) FindWindowByProcess(process_name string, class_name string, title_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| process_name | string | 进程名称 |
| class_name | string | 窗口类名 |
| title_name | string | 窗口标题 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### FindWindowByProcessId

**功能说明**: 通过进程ID查找窗口

**函数签名**:
```go
func (dm *DmSoft) FindWindowByProcessId(process_id int32, class_name string, title_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| process_id | int32 | 进程ID |
| class_name | string | 窗口类名 |
| title_name | string | 窗口标题 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### EnumWindow

**功能说明**: 枚举窗口

**函数签名**:
```go
func (dm *DmSoft) EnumWindow(parent int32, title string, class_name string, filter int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| parent | int32 | 父窗口句柄 |
| title | string | 窗口标题 |
| class_name | string | 窗口类名 |
| filter | int32 | 过滤条件 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 窗口句柄列表字符串，格式: "hwnd1,hwnd2,hwnd3" |

---

### EnumWindowSuper

**功能说明**: 高级枚举窗口

**函数签名**:
```go
func (dm *DmSoft) EnumWindowSuper(spec1 string, flag1 int32, type1 int32, spec2 string, flag2 int32, type2 int32, sort int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| spec1 | string | 条件1字符串 |
| flag1 | int32 | 条件1标志 |
| type1 | int32 | 条件1类型 |
| spec2 | string | 条件2字符串 |
| flag2 | int32 | 条件2标志 |
| type2 | int32 | 条件2类型 |
| sort | int32 | 排序方式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 窗口句柄列表字符串 |

---

### EnumWindowByProcess

**功能说明**: 通过进程名枚举窗口

**函数签名**:
```go
func (dm *DmSoft) EnumWindowByProcess(process_name string, title string, class_name string, filter int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| process_name | string | 进程名称 |
| title | string | 窗口标题 |
| class_name | string | 窗口类名 |
| filter | int32 | 过滤条件 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 窗口句柄列表字符串 |

---

### EnumWindowByProcessId

**功能说明**: 通过进程ID枚举窗口

**函数签名**:
```go
func (dm *DmSoft) EnumWindowByProcessId(pid int32, title string, class_name string, filter int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pid | int32 | 进程ID |
| title | string | 窗口标题 |
| class_name | string | 窗口类名 |
| filter | int32 | 过滤条件 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 窗口句柄列表字符串 |

---

### GetWindowTitle

**功能说明**: 获取窗口标题

**函数签名**:
```go
func (dm *DmSoft) GetWindowTitle(hwnd int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 标题字符串 |

---

### GetWindowClass

**功能说明**: 获取窗口类名

**函数签名**:
```go
func (dm *DmSoft) GetWindowClass(hwnd int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 类名字符串 |

---

### GetWindowRect

**功能说明**: 获取窗口矩形

**函数签名**:
```go
func (dm *DmSoft) GetWindowRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| x1 | *int32 | 左上角X坐标（输出参数） |
| y1 | *int32 | 左上角Y坐标（输出参数） |
| x2 | *int32 | 右下角X坐标（输出参数） |
| y2 | *int32 | 右下角Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetClientRect

**功能说明**: 获取客户区矩形

**函数签名**:
```go
func (dm *DmSoft) GetClientRect(hwnd int32, x1 *int32, y1 *int32, x2 *int32, y2 *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| x1 | *int32 | 左上角X坐标（输出参数） |
| y1 | *int32 | 左上角Y坐标（输出参数） |
| x2 | *int32 | 右下角X坐标（输出参数） |
| y2 | *int32 | 右下角Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetClientSize

**功能说明**: 获取客户区大小

**函数签名**:
```go
func (dm *DmSoft) GetClientSize(hwnd int32, width *int32, height *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| width | *int32 | 宽度（输出参数） |
| height | *int32 | 高度（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetWindow

**功能说明**: 获取关联窗口

**函数签名**:
```go
func (dm *DmSoft) GetWindow(hwnd int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| flag | int32 | 查找标志: 0父窗口, 1第一个子窗口, 2最后一个子窗口, 3下一个兄弟窗口, 4上一个兄弟窗口, 5所有者窗口 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### GetWindowProcessId

**功能说明**: 获取窗口进程ID

**函数签名**:
```go
func (dm *DmSoft) GetWindowProcessId(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 进程ID |

---

### GetWindowThreadId

**功能说明**: 获取窗口线程ID

**函数签名**:
```go
func (dm *DmSoft) GetWindowThreadId(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 线程ID |

---

### GetWindowProcessPath

**功能说明**: 获取窗口进程路径

**函数签名**:
```go
func (dm *DmSoft) GetWindowProcessPath(hwnd int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 进程路径字符串 |

---

### GetWindowState

**功能说明**: 获取窗口状态

**函数签名**:
```go
func (dm *DmSoft) GetWindowState(hwnd int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| flag | int32 | 状态标志: 0是否存在, 1是否可见, 2是否最小化, 3是否最大化, 4是否置顶, 5是否挂起, 6是否激活, 7是否关闭 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 状态值，1表示是，0表示否 |

---

### SetWindowState

**功能说明**: 设置窗口状态

**函数签名**:
```go
func (dm *DmSoft) SetWindowState(hwnd int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| flag | int32 | 状态标志: 0关闭, 1激活, 2最小化(不激活), 3最小化(释放内存), 4最大化(激活), 5恢复(不激活), 6隐藏, 7显示, 8置顶, 9取消置顶, 10禁止, 11取消禁止, 12恢复并激活, 13强制结束进程, 14闪烁, 15获取输入焦点 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetWindowText

**功能说明**: 设置窗口标题

**函数签名**:
```go
func (dm *DmSoft) SetWindowText(hwnd int32, text string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| text | string | 标题 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetWindowSize

**功能说明**: 设置窗口大小

**函数签名**:
```go
func (dm *DmSoft) SetWindowSize(hwnd int32, width int32, height int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| width | int32 | 宽度 |
| height | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MoveWindow

**功能说明**: 移动窗口

**函数签名**:
```go
func (dm *DmSoft) MoveWindow(hwnd int32, x int32, y int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetWindowTransparent

**功能说明**: 设置窗口透明度

**函数签名**:
```go
func (dm *DmSoft) SetWindowTransparent(hwnd int32, v int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| v | int32 | 透明度（0-255，0完全透明，255完全不透明） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetForegroundWindow

**功能说明**: 获取前台窗口

**函数签名**:
```go
func (dm *DmSoft) GetForegroundWindow() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### GetForegroundFocus

**功能说明**: 获取前台焦点窗口

**函数签名**:
```go
func (dm *DmSoft) GetForegroundFocus() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### GetMousePointWindow

**功能说明**: 获取鼠标指向窗口

**函数签名**:
```go
func (dm *DmSoft) GetMousePointWindow() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### GetPointWindow

**功能说明**: 获取指定坐标窗口

**函数签名**:
```go
func (dm *DmSoft) GetPointWindow(x int32, y int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### GetSpecialWindow

**功能说明**: 获取特殊窗口

**函数签名**:
```go
func (dm *DmSoft) GetSpecialWindow(flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| flag | int32 | 标志: 0桌面窗口, 1任务栏窗口 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### BindWindow

**功能说明**: 绑定窗口

**函数签名**:
```go
func (dm *DmSoft) BindWindow(hwnd int32, display string, mouse string, keypad string, mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| display | string | 显示模式: "normal", "gdi", "gdi2", "dx", "dx2", "dx3" |
| mouse | string | 鼠标模式: "normal", "windows", "windows3", "dx", "dx2" |
| keypad | string | 键盘模式: "normal", "windows", "dx" |
| mode | int32 | 模式: 0-推荐, 1-兼容模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
hwnd := dm.FindWindow("", "窗口标题")
ret := dm.BindWindow(hwnd, "gdi", "windows", "windows", 0)
if ret == 1 {
    fmt.Println("绑定成功")
}
```

---

### BindWindowEx

**功能说明**: 扩展绑定窗口

**函数签名**:
```go
func (dm *DmSoft) BindWindowEx(hwnd int32, display string, mouse string, keypad string, public_desc string, mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| display | string | 显示模式 |
| mouse | string | 鼠标模式 |
| keypad | string | 键盘模式 |
| public_desc | string | 公共描述 |
| mode | int32 | 模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### UnBindWindow

**功能说明**: 解绑窗口

**函数签名**:
```go
func (dm *DmSoft) UnBindWindow() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ForceUnBindWindow

**功能说明**: 强制解绑窗口

**函数签名**:
```go
func (dm *DmSoft) ForceUnBindWindow(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### IsBind

**功能说明**: 判断是否已绑定窗口

**函数签名**:
```go
func (dm *DmSoft) IsBind(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 已绑定返回1，未绑定返回0 |

---

### GetBindWindow

**功能说明**: 获取绑定的窗口句柄

**函数签名**:
```go
func (dm *DmSoft) GetBindWindow() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 窗口句柄 |

---

### SwitchBindWindow

**功能说明**: 切换绑定窗口

**函数签名**:
```go
func (dm *DmSoft) SwitchBindWindow(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ScreenToClient

**功能说明**: 屏幕坐标转客户区坐标

**函数签名**:
```go
func (dm *DmSoft) ScreenToClient(hwnd int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| x | *int32 | X坐标（输入/输出参数） |
| y | *int32 | Y坐标（输入/输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ClientToScreen

**功能说明**: 客户区坐标转屏幕坐标

**函数签名**:
```go
func (dm *DmSoft) ClientToScreen(hwnd int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| x | *int32 | X坐标（输入/输出参数） |
| y | *int32 | Y坐标（输入/输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetEnumWindowDelay

**功能说明**: 设置枚举窗口延迟

**函数签名**:
```go
func (dm *DmSoft) SetEnumWindowDelay(delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| delay | int32 | 延迟时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableBind

**功能说明**: 启用绑定

**函数签名**:
```go
func (dm *DmSoft) EnableBind(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetClientSize

**功能说明**: 设置客户区大小

**函数签名**:
```go
func (dm *DmSoft) SetClientSize(hwnd int32, width int32, height int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| width | int32 | 宽度 |
| height | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ShowTaskBarIcon

**功能说明**: 显示或隐藏任务栏图标

**函数签名**:
```go
func (dm *DmSoft) ShowTaskBarIcon(hwnd int32, is_show int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| is_show | int32 | 1显示，0隐藏 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 2.6 后台设置函数

### DownCpu

**功能说明**: 降低CPU占用率

**函数签名**:
```go
func (dm *DmSoft) DownCpu() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### EnableFakeActive

**功能说明**: 启用虚假激活

**函数签名**:
```go
func (dm *DmSoft) EnableFakeActive(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableIme

**功能说明**: 启用输入法

**函数签名**:
```go
func (dm *DmSoft) EnableIme(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableSpeedDx

**功能说明**: 启用DX加速

**函数签名**:
```go
func (dm *DmSoft) EnableSpeedDx(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### HackSpeed

**功能说明**: 加速/减速

**函数签名**:
```go
func (dm *DmSoft) HackSpeed(rate int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| rate | int32 | 速度倍率 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### LockDisplay

**功能说明**: 锁定显示

**函数签名**:
```go
func (dm *DmSoft) LockDisplay(lock int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| lock | int32 | 1锁定，0解锁 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### LockInput

**功能说明**: 锁定输入

**函数签名**:
```go
func (dm *DmSoft) LockInput(lock int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| lock | int32 | 1锁定，0解锁 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetAero

**功能说明**: 设置Aero效果

**函数签名**:
```go
func (dm *DmSoft) SetAero(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetDisplayDelay

**功能说明**: 设置显示延迟

**函数签名**:
```go
func (dm *DmSoft) SetDisplayDelay(delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| delay | int32 | 延迟时间(毫秒) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetDisplayRefreshDelay

**功能说明**: 设置显示刷新延迟

**函数签名**:
```go
func (dm *DmSoft) SetDisplayRefreshDelay(delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| delay | int32 | 延迟时间(毫秒) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetInputDm

**功能说明**: 设置输入设备

**函数签名**:
```go
func (dm *DmSoft) SetInputDm(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 3. 鼠标操作函数

### MoveTo

**功能说明**: 移动鼠标到指定坐标

**函数签名**:
```go
func (dm *DmSoft) MoveTo(x int32, y int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
dm.MoveTo(100, 200)
```

---

### MoveR

**功能说明**: 相对移动鼠标

**函数签名**:
```go
func (dm *DmSoft) MoveR(rx int32, ry int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| rx | int32 | 相对X偏移 |
| ry | int32 | 相对Y偏移 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MoveToEx

**功能说明**: 扩展移动鼠标，支持随机偏移

**函数签名**:
```go
func (dm *DmSoft) MoveToEx(x int32, y int32, w int32, h int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| w | int32 | 宽度范围 |
| h | int32 | 高度范围 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 实际移动到的坐标字符串 |

---

### MoveDD

**功能说明**: DD驱动移动鼠标

**函数签名**:
```go
func (dm *DmSoft) MoveDD(dx int32, dy int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| dx | int32 | X偏移 |
| dy | int32 | Y偏移 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### LeftClick

**功能说明**: 鼠标左键单击

**函数签名**:
```go
func (dm *DmSoft) LeftClick() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### RightClick

**功能说明**: 鼠标右键单击

**函数签名**:
```go
func (dm *DmSoft) RightClick() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MiddleClick

**功能说明**: 鼠标中键单击

**函数签名**:
```go
func (dm *DmSoft) MiddleClick() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### LeftDoubleClick

**功能说明**: 鼠标左键双击

**函数签名**:
```go
func (dm *DmSoft) LeftDoubleClick() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### LeftDown

**功能说明**: 鼠标左键按下

**函数签名**:
```go
func (dm *DmSoft) LeftDown() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### LeftUp

**功能说明**: 鼠标左键弹起

**函数签名**:
```go
func (dm *DmSoft) LeftUp() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### RightDown

**功能说明**: 鼠标右键按下

**函数签名**:
```go
func (dm *DmSoft) RightDown() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### RightUp

**功能说明**: 鼠标右键弹起

**函数签名**:
```go
func (dm *DmSoft) RightUp() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MiddleDown

**功能说明**: 鼠标中键按下

**函数签名**:
```go
func (dm *DmSoft) MiddleDown() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MiddleUp

**功能说明**: 鼠标中键弹起

**函数签名**:
```go
func (dm *DmSoft) MiddleUp() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WheelUp

**功能说明**: 鼠标滚轮向上

**函数签名**:
```go
func (dm *DmSoft) WheelUp() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WheelDown

**功能说明**: 鼠标滚轮向下

**函数签名**:
```go
func (dm *DmSoft) WheelDown() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetCursorPos

**功能说明**: 获取鼠标当前位置

**函数签名**:
```go
func (dm *DmSoft) GetCursorPos(x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
var x, y int32
ret := dm.GetCursorPos(&x, &y)
if ret == 1 {
    fmt.Printf("鼠标位置: (%d, %d)\n", x, y)
}
```

---

### GetCursorShape

**功能说明**: 获取鼠标形状特征码

**函数签名**:
```go
func (dm *DmSoft) GetCursorShape() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 鼠标形状特征码字符串 |

---

### GetCursorShapeEx

**功能说明**: 获取鼠标形状（扩展）

**函数签名**:
```go
func (dm *DmSoft) GetCursorShapeEx(type_ int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| type_ | int32 | 类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 鼠标形状字符串 |

---

### GetCursorSpot

**功能说明**: 获取鼠标光点位置

**函数签名**:
```go
func (dm *DmSoft) GetCursorSpot() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 光点位置字符串，格式: "x,y" |

---

### SetMouseSpeed

**功能说明**: 设置鼠标移动速度

**函数签名**:
```go
func (dm *DmSoft) SetMouseSpeed(speed int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| speed | int32 | 移动速度（1-100） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### GetMouseSpeed

**功能说明**: 获取鼠标移动速度

**函数签名**:
```go
func (dm *DmSoft) GetMouseSpeed() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 移动速度 |

---

### EnableRealMouse

**功能说明**: 启用真实鼠标模拟

**函数签名**:
```go
func (dm *DmSoft) EnableRealMouse(en int32, mousedelay int32, mousestep int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| en | int32 | 1启用，0禁用 |
| mousedelay | int32 | 鼠标延迟（毫秒） |
| mousestep | int32 | 鼠标步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableMouseSync

**功能说明**: 启用鼠标同步

**函数签名**:
```go
func (dm *DmSoft) EnableMouseSync(enable int32, time_out int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |
| time_out | int32 | 超时时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableMouseMsg

**功能说明**: 启用鼠标消息模拟

**函数签名**:
```go
func (dm *DmSoft) EnableMouseMsg(en int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| en | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableMouseAccuracy

**功能说明**: 启用鼠标精度

**函数签名**:
```go
func (dm *DmSoft) EnableMouseAccuracy(en int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| en | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetMouseDelay

**功能说明**: 设置鼠标操作延迟

**函数签名**:
```go
func (dm *DmSoft) SetMouseDelay(type_ string, delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| type_ | string | 类型 |
| delay | int32 | 延迟时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### LockMouseRect

**功能说明**: 锁定鼠标移动区域

**函数签名**:
```go
func (dm *DmSoft) LockMouseRect(x1 int32, y1 int32, x2 int32, y2 int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 4. 键盘操作函数

### KeyPress

**功能说明**: 按键（虚拟键码）

**函数签名**:
```go
func (dm *DmSoft) KeyPress(vk int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| vk | int32 | 虚拟键码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
dm.KeyPress(0x41) // A键
```

---

### KeyDown

**功能说明**: 按下按键（虚拟键码）

**函数签名**:
```go
func (dm *DmSoft) KeyDown(vk int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| vk | int32 | 虚拟键码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### KeyUp

**功能说明**: 弹起按键（虚拟键码）

**函数签名**:
```go
func (dm *DmSoft) KeyUp(vk int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| vk | int32 | 虚拟键码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### KeyPressChar

**功能说明**: 按键（字符形式）

**函数签名**:
```go
func (dm *DmSoft) KeyPressChar(key_str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| key_str | string | 按键字符串，如 "a", "enter", "space" |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
dm.KeyPressChar("a")
dm.KeyPressChar("enter")
```

---

### KeyDownChar

**功能说明**: 按下按键（字符形式）

**函数签名**:
```go
func (dm *DmSoft) KeyDownChar(key_str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| key_str | string | 按键字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### KeyUpChar

**功能说明**: 弹起按键（字符形式）

**函数签名**:
```go
func (dm *DmSoft) KeyUpChar(key_str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| key_str | string | 按键字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### KeyPressStr

**功能说明**: 按键字符串序列

**函数签名**:
```go
func (dm *DmSoft) KeyPressStr(key_str string, delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| key_str | string | 按键字符串 |
| delay | int32 | 延迟时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WaitKey

**功能说明**: 等待按键

**函数签名**:
```go
func (dm *DmSoft) WaitKey(key_code int32, time_out int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| key_code | int32 | 键码 |
| time_out | int32 | 超时时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetKeyState

**功能说明**: 获取按键状态

**函数签名**:
```go
func (dm *DmSoft) GetKeyState(vk int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| vk | int32 | 虚拟键码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 0弹起，1按下 |

---

### SetKeypadDelay

**功能说明**: 设置键盘按键延迟

**函数签名**:
```go
func (dm *DmSoft) SetKeypadDelay(type_ string, delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| type_ | string | 类型 |
| delay | int32 | 延迟时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableKeypadSync

**功能说明**: 启用键盘同步

**函数签名**:
```go
func (dm *DmSoft) EnableKeypadSync(enable int32, time_out int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |
| time_out | int32 | 超时时间（毫秒） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableKeypadMsg

**功能说明**: 启用键盘消息

**函数签名**:
```go
func (dm *DmSoft) EnableKeypadMsg(en int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| en | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableRealKeypad

**功能说明**: 启用真实键盘

**函数签名**:
```go
func (dm *DmSoft) EnableRealKeypad(en int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| en | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableKeypadPatch

**功能说明**: 启用键盘补丁

**函数签名**:
```go
func (dm *DmSoft) EnableKeypadPatch(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SendString

**功能说明**: 发送字符串

**函数签名**:
```go
func (dm *DmSoft) SendString(hwnd int32, str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| str | string | 要发送的字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SendString2

**功能说明**: 发送字符串（方式2）

**函数签名**:
```go
func (dm *DmSoft) SendString2(hwnd int32, str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| str | string | 要发送的字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SendStringIme

**功能说明**: 通过输入法发送字符串

**函数签名**:
```go
func (dm *DmSoft) SendStringIme(str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | 要发送的字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SendStringIme2

**功能说明**: 通过输入法发送字符串（方式2）

**函数签名**:
```go
func (dm *DmSoft) SendStringIme2(hwnd int32, str string, mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| str | string | 要发送的字符串 |
| mode | int32 | 模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SendPaste

**功能说明**: 发送粘贴

**函数签名**:
```go
func (dm *DmSoft) SendPaste(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetClipboard

**功能说明**: 获取剪贴板内容

**函数签名**:
```go
func (dm *DmSoft) GetClipboard() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 剪贴板内容 |

---

### SetClipboard

**功能说明**: 设置剪贴板内容

**函数签名**:
```go
func (dm *DmSoft) SetClipboard(data string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| data | string | 数据 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 9. 内存操作函数

### ReadInt

**功能说明**: 读取整数

**函数签名**:
```go
func (dm *DmSoft) ReadInt(hwnd int32, addr string, type_ int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| type_ | int32 | 类型: 0-4字节整数, 1-2字节整数, 2-1字节整数 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 读取的整数值 |

---

## 20. 64位(x64)支持说明

### 概述

本库完整支持64位(x64/amd64)架构，通过跨进程通信机制实现与32位大漠DLL的交互。

### 架构原理

由于大漠插件DLL（dm.dll等）是32位的，无法直接在64位进程中加载，因此采用以下方案：

```
┌─────────────────────────────────────┐
│         64位主进程 (Go程序)          │
│  ┌─────────────────────────────┐   │
│  │     DmSoft 实例              │   │
│  │  - comCall*() 封装函数       │   │
│  │  - gob 编解码器             │   │
│  └──────────┬──────────────────┘   │
│             │ TCP连接               │
│             │ gob序列化             │
└─────────────┼───────────────────────┘
              │
              ▼
┌─────────────────────────────────────┐
│    32位helper进程 (dm_com_server)   │
│  ┌─────────────────────────────┐   │
│  │  - 加载dm.dll (32位)        │   │
│  │  - 可选: 加载crack.dll      │   │
│  │  - 创建COM对象              │   │
│  │  - 通过偏移量调用API        │   │
│  │  - 返回结果                 │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
```

### 使用方式

#### 基础用法

```go
package main

import (
    "fmt"
    "log"
    
    dmsoft "github.com/yuan71058/dm72424-go"
)

func main() {
    // 1. 加载大漠DLL（记录路径）
    _, err := dmsoft.LoadDm("xd47243.dll")
    if err != nil {
        log.Fatalf("LoadDm失败: %v", err)
    }

    // 2. 设置破解DLL路径（可选）
    err = dmsoft.CrackDm("Go.dll")
    if err != nil {
        log.Fatalf("CrackDm失败: %v", err)
    }

    // 3. 创建实例并初始化（会启动32位helper进程）
    dm := dmsoft.New()
    if dm == nil {
        log.Fatal("创建实例失败")
    }
    
    err = dm.Init()  // 启动helper进程，建立TCP连接
    if err != nil {
        log.Fatalf("Init失败: %v", err)
    }
    defer dm.Release()

    // 4. 注册并使用
    ret := dm.Reg("", "")
    fmt.Printf("注册结果: %d\n", ret)
    
    // 5. 调用各种API（自动通过TCP转发到helper）
    version := dm.Ver()
    fmt.Printf("版本号: %s\n", version)
}
```

#### 多线程用法

```go
package main

import (
    "sync"
    
    dmsoft "github.com/yuan71058/dm72424-go"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // 每个线程创建独立的DmSoft实例
    // 每个实例会启动独立的helper进程
    dm := dmsoft.New()
    dm.Init()
    defer dm.Release()
    
    dm.Reg("", "")
    
    // 各线程独立操作，互不干扰
    var x, y int32
    dm.GetCursorPos(&x, &y)
}

func main() {
    dmsoft.LoadDm("xd47243.dll")
    dmsoft.CrackDm("Go.dll")
    
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    wg.Wait()
}
```

### 关键文件说明

| 文件 | 说明 |
|------|------|
| `dm_x64_helpers.go` | 64位模式下的COM调用封装函数（comCallInt32、comCallStr等） |
| `dm_x64_pipe.go` | 管道通信实现（结构体定义、偏移量表、gob编解码） |
| `cmd/dm_com_server/main.go` | 32位helper进程源码（必须以GOARCH=386编译） |

### 编译要求

#### 编译helper进程（必须32位）

```bash
# Windows
set GOARCH=386
go build -o dm_com_server.exe ./cmd/dm_com_server/

# Linux/Mac
GOARCH=386 go build -o dm_com_server ./cmd/dm_com_server/
```

#### 编译64位主程序

```bash
# Windows
set GOARCH=amd64
go build -o myapp.exe .

# Linux/Mac
GOARCH=amd64 go build -o myapp .
```

### 工作流程详解

#### 初始化阶段

1. **LoadDm(path)** 
   - 记录大漠DLL路径
   - 检测当前架构（32/64位）
   - 查找helper可执行文件位置
   
2. **CrackDm(crackPath)**
   - 记录破解DLL路径（不立即加载）
   
3. **New() + Init()**
   - 启动32位helper进程：`dm_com_server.exe <dm.dll> [crack.dll]`
   - helper进程：
     a. LoadLibrary加载dm.dll
     b. （可选）LoadLibrary加载crack.dll并调用Go()函数
     c. 通过偏移量0x18000调用创建COM对象的函数
     d. 监听随机TCP端口
     e. 输出"READY <port>"到stdout
   - 主进程读取端口号
   - 建立TCP连接到127.0.0.1:<port>
   - 创建gob Encoder/Decoder用于后续通信

#### 方法调用阶段

```go
// 用户代码
var x, y int32
result := dm.FindPic(0, 0, 800, 600, "test.bmp", "000000", 0.9, 0, &x, &y)

// 内部执行流程：
// 1. comCallWithOutVars("FindPic", [...], &x, &y)
// 2. getMethodOffset("FindPic") -> 返回104032
// 3. pipeCallWithOut(104032, 0, [...], &x, &y)
// 4. 构建callRequest{Offset:104032, RetType:0, NOut:2, Args:[...]}
// 5. gob.Encode(&req) -> 通过TCP发送到helper
// 6. helper接收并解码
// 7. handleCall(req):
//    - 计算函数地址: dmModule + 104032
//    - 构建参数数组（处理类型转换、UTF-8转GBK）
//    - syscall.Syscall12(fnAddr, ...)
//    - 收集返回值和输出参数
// 8. 构建callResponse并通过gob发送回主进程
// 9. 主进程解码响应，提取IRet和OutVals
// 10. 将OutVals[0]写入*x, OutVals[1]写入*y
// 11. 返回IRet给用户
```

### 数据类型映射

| Go类型 | callArg.Type | 传输方式 | 说明 |
|--------|-------------|---------|------|
| int32 | 0 | IVal | 直接传递整数 |
| string | 1 | SVal | UTF-8→GBK转换后传递 |
| float64 | 2 | FVal | IEEE 754双精度 |
| *int32 (输出参数) | 3 | - | 占位符，helper填充实际值 |
| float32 | 4 | FVal | 转换为float64后传递 |
| int64 | 5 | I64Val | 高低32位分别传递 |

### 返回值类型

| RetType | 含义 | callResponse字段 |
|---------|------|-----------------|
| 0 | int32 | IRet |
| 1 | string | SRet（GBK→UTF-8） |
| 2 | float64 | FRet |
| 3 | int64 | IRet64（高低32位组合） |

### 特殊注意事项

#### 1. 字符串编码

- 大漠DLL使用GBK编码
- 64位模式下自动处理编码转换：
  - 发送请求时：UTF-8 → GBK
  - 接收响应时：GBK → UTF-8
- 无需手动处理编码问题

#### 2. 输出参数处理

- 带`*int32`输出参数的方法（如FindPic、FindColor等）
- 输出参数通过gob序列化跨进程传递
- 在helper进程中收集实际值，返回后在主进程写入指针指向的变量
- 完全透明，使用方式与32位模式相同

#### 3. 多线程安全

- 每个`DmSoft`实例有独立的helper进程和TCP连接
- 使用`sync.RWMutex`保证单个实例的线程安全
- 不同实例之间完全独立，互不影响
- 推荐每个goroutine使用独立的`DmSoft`实例

#### 4. 资源管理

- **必须调用`Release()`**释放资源
- Release会：
  - 关闭TCP连接
  - 关闭helper进程的stdin（发送退出信号）
  - 终止helper进程
- 建议使用`defer dm.Release()`确保释放

#### 5. 性能考虑

- 每次API调用都需要：
  - gob序列化/反序列化
  - TCP网络传输
  - 进程间上下文切换
- 相比32位直接调用有一定开销
- 对于性能敏感场景，建议：
  - 批量操作减少调用次数
  - 使用扩展方法（如FindPicEx代替多次FindPic）
  - 合理设置超时时间

#### 6. 错误处理

常见错误及解决方案：

| 错误信息 | 原因 | 解决方案 |
|---------|------|---------|
| 未找到dm_com_server.exe | helper未编译或路径错误 | 以GOARCH=386编译helper并放在正确位置 |
| 连接helper TCP失败 | helper启动失败或崩溃 | 检查stderr输出，确认DLL路径正确 |
| pipe not connected | 未调用Init()或已Release | 确保先Init()再调用API |
| helper进程无输出 | DLL加载失败 | 检查DLL是否存在、位数是否正确 |
| encode/decode error | 网络断开或数据损坏 | 检查helper进程是否存活 |

### 与32位模式的对比

| 特性 | 32位模式 | 64位模式 |
|------|---------|---------|
| DLL加载方式 | 直接LoadLibrary到当前进程 | 由helper进程加载 |
| COM对象创建 | 当前进程内创建 | helper进程内创建 |
| API调用方式 | 直接syscall | TCP+gob转发到helper |
| 性能 | 无额外开销 | 有网络序列化开销 |
| 多线程 | 需要同步保护 | 天然隔离（每实例独立helper） |
| 内存占用 | 较低 | 较高（每个helper约5-10MB） |
| 兼容性 | 仅限32位系统 | 支持64位系统 |

### 示例项目

项目中包含完整的示例代码：

- **example/x64/** - 单线程64位示例
  - 展示基础功能：找图、取色、OCR、鼠标键盘操作等
  - 详细的控制台输出，展示调用流程
  
- **example/x64_mt/** - 多线程64位示例
  - 展示多线程并发操作不同窗口
  - 每个线程独立helper进程，真正并行执行
  - 包含记事本自动化测试场景

### 常见问题FAQ

**Q: 为什么需要helper进程？**
A: 大漠DLL是32位的，无法在64位进程中直接加载。helper作为32位桥梁进程负责实际的DLL调用。

**Q: 可以在Linux/macOS上运行吗？**
A: 可以，但需要Wine或其他Windows兼容层来运行32位Windows DLL。helper进程本身可以交叉编译为Linux/macOS的32位版本。

**Q: 如何调试helper进程？**
A: helper会将错误信息输出到stderr。也可以设置环境变量`DM_HELPER_DEBUG=1`启用详细日志。

**Q: 支持哪些大漠DLL版本？**
A: 支持所有7.24xx系列的大漠DLL。偏移量表定义了所有已知API的位置。

**Q: 如何更新helper？**
A: 当主库更新时，需要重新编译helper：`GOARCH=386 go build -o dm_com_server.exe ./cmd/dm_com_server/`

---

### ReadIntAddr

**功能说明**: 读取整数（地址形式）

**函数签名**:
```go
func (dm *DmSoft) ReadIntAddr(hwnd int32, addr int64, type_ int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| type_ | int32 | 类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 读取的整数值 |

---

### WriteInt

**功能说明**: 写入整数

**函数签名**:
```go
func (dm *DmSoft) WriteInt(hwnd int32, addr string, v int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| v | int32 | 要写入的值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WriteIntAddr

**功能说明**: 写入整数（地址形式）

**函数签名**:
```go
func (dm *DmSoft) WriteIntAddr(hwnd int32, addr int64, v int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| v | int32 | 要写入的值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadFloat

**功能说明**: 读取浮点数

**函数签名**:
```go
func (dm *DmSoft) ReadFloat(hwnd int32, addr string) float32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| float32 | 读取的浮点数值 |

---

### ReadFloatAddr

**功能说明**: 读取浮点数（地址形式）

**函数签名**:
```go
func (dm *DmSoft) ReadFloatAddr(hwnd int32, addr int64) float32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| float32 | 读取的浮点数值 |

---

### WriteFloat

**功能说明**: 写入浮点数

**函数签名**:
```go
func (dm *DmSoft) WriteFloat(hwnd int32, addr string, v float32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| v | float32 | 要写入的值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadDouble

**功能说明**: 读取双精度浮点数

**函数签名**:
```go
func (dm *DmSoft) ReadDouble(hwnd int32, addr string) float64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| float64 | 读取的双精度浮点数值 |

---

### ReadDoubleAddr

**功能说明**: 读取双精度浮点数（地址形式）

**函数签名**:
```go
func (dm *DmSoft) ReadDoubleAddr(hwnd int32, addr int64) float64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| float64 | 读取的双精度浮点数值 |

---

### WriteDouble

**功能说明**: 写入双精度浮点数

**函数签名**:
```go
func (dm *DmSoft) WriteDouble(hwnd int32, addr string, v float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| v | float64 | 要写入的值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadString

**功能说明**: 读取字符串

**函数签名**:
```go
func (dm *DmSoft) ReadString(hwnd int32, addr string, length int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| length | int32 | 长度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取的字符串 |

---

### ReadStringAddr

**功能说明**: 读取字符串（地址形式）

**函数签名**:
```go
func (dm *DmSoft) ReadStringAddr(hwnd int32, addr int64, length int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| length | int32 | 长度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取的字符串 |

---

### WriteString

**功能说明**: 写入字符串

**函数签名**:
```go
func (dm *DmSoft) WriteString(hwnd int32, addr string, v string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| v | string | 要写入的字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WriteStringAddr

**功能说明**: 写入字符串（地址形式）

**函数签名**:
```go
func (dm *DmSoft) WriteStringAddr(hwnd int32, addr int64, v string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| v | string | 要写入的字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadData

**功能说明**: 读取数据

**函数签名**:
```go
func (dm *DmSoft) ReadData(hwnd int32, addr string, length int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| length | int32 | 长度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取的数据（十六进制字符串） |

---

### ReadDataAddr

**功能说明**: 读取数据（地址形式）

**函数签名**:
```go
func (dm *DmSoft) ReadDataAddr(hwnd int32, addr int64, length int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| length | int32 | 长度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取的数据 |

---

### WriteData

**功能说明**: 写入数据

**函数签名**:
```go
func (dm *DmSoft) WriteData(hwnd int32, addr string, data string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 内存地址 |
| data | string | 数据（十六进制字符串） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WriteDataAddr

**功能说明**: 写入数据（地址形式）

**函数签名**:
```go
func (dm *DmSoft) WriteDataAddr(hwnd int32, addr int64, data string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| data | string | 数据 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### VirtualAllocEx

**功能说明**: 在目标进程分配内存

**函数签名**:
```go
func (dm *DmSoft) VirtualAllocEx(hwnd int32, addr int64, size int32, type_ int32, protect int32) int64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 分配地址 |
| size | int32 | 大小 |
| type_ | int32 | 分配类型 |
| protect | int32 | 保护属性 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int64 | 分配的内存地址 |

---

### VirtualFreeEx

**功能说明**: 释放目标进程内存

**函数签名**:
```go
func (dm *DmSoft) VirtualFreeEx(hwnd int32, addr int64, size int32, type_ int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| size | int32 | 大小 |
| type_ | int32 | 释放类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### VirtualProtectEx

**功能说明**: 修改目标进程内存保护属性

**函数签名**:
```go
func (dm *DmSoft) VirtualProtectEx(hwnd int32, addr int64, size int32, protect int32, old_protect *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| size | int32 | 大小 |
| protect | int32 | 新保护属性 |
| old_protect | *int32 | 原保护属性（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### VirtualQueryEx

**功能说明**: 查询目标进程内存信息

**函数签名**:
```go
func (dm *DmSoft) VirtualQueryEx(hwnd int32, addr int64, pmbi int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int64 | 内存地址 |
| pmbi | int32 | 内存信息缓冲区指针 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FindInt

**功能说明**: 查找整数

**函数签名**:
```go
func (dm *DmSoft) FindInt(hwnd int32, addr_range string, int_value_min int32, int_value_max int32, step int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| int_value_min | int32 | 最小值 |
| int_value_max | int32 | 最大值 |
| step | int32 | 步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 找到的地址 |

---

### FindFloat

**功能说明**: 查找浮点数

**函数签名**:
```go
func (dm *DmSoft) FindFloat(hwnd int32, addr_range string, float_value_min float32, float_value_max float32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| float_value_min | float32 | 最小值 |
| float_value_max | float32 | 最大值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 找到的地址 |

---

### FindDouble

**功能说明**: 查找双精度浮点数

**函数签名**:
```go
func (dm *DmSoft) FindDouble(hwnd int32, addr_range string, double_value_min float64, double_value_max float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| double_value_min | float64 | 最小值 |
| double_value_max | float64 | 最大值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 找到的地址 |

---

### FindString

**功能说明**: 查找字符串

**函数签名**:
```go
func (dm *DmSoft) FindString(hwnd int32, addr_range string, str_value string, type_ int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| str_value | string | 字符串值 |
| type_ | int32 | 类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 找到的地址 |

---

### FindData

**功能说明**: 查找数据

**函数签名**:
```go
func (dm *DmSoft) FindData(hwnd int32, addr_range string, data string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| data | string | 数据 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 找到的地址 |

---

### GetRemoteApiAddress

**功能说明**: 获取远程API地址

**函数签名**:
```go
func (dm *DmSoft) GetRemoteApiAddress(hwnd int32, module_name string, proc_name string) int64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| module_name | string | 模块名 |
| proc_name | string | 函数名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int64 | API地址 |

---

### DoubleToData

**功能说明**: 双精度浮点数转二进制数据

**函数签名**:
```go
func (dm *DmSoft) DoubleToData(value float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| value | float64 | 双精度浮点数值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 二进制数据字符串 |

---

### FindDataEx

**功能说明**: 搜索二进制数据扩展

**函数签名**:
```go
func (dm *DmSoft) FindDataEx(hwnd int32, addr_range string, data string, size int32, step int32, multi_depth int32, match_mode int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| data | string | 要搜索的数据 |
| size | int32 | 数据大小 |
| step | int32 | 搜索步长 |
| multi_depth | int32 | 多级深度 |
| match_mode | int32 | 匹配模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 搜索结果 |

---

### FindDoubleEx

**功能说明**: 搜索双精度浮点数扩展

**函数签名**:
```go
func (dm *DmSoft) FindDoubleEx(hwnd int32, addr_range string, double_value float64, multi_depth int32, step int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| double_value | float64 | 双精度浮点数值 |
| multi_depth | int32 | 多级深度 |
| step | int32 | 搜索步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 搜索结果 |

---

### FindFloatEx

**功能说明**: 搜索单精度浮点数扩展

**函数签名**:
```go
func (dm *DmSoft) FindFloatEx(hwnd int32, addr_range string, float_value float64, multi_depth int32, step int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| float_value | float64 | 单精度浮点数值(用float64传递) |
| multi_depth | int32 | 多级深度 |
| step | int32 | 搜索步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 搜索结果 |

---

### FindIntEx

**功能说明**: 搜索整数扩展

**函数签名**:
```go
func (dm *DmSoft) FindIntEx(hwnd int32, addr_range string, int_value int32, multi_depth int32, step int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| int_value | int32 | 整数值 |
| multi_depth | int32 | 多级深度 |
| step | int32 | 搜索步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 搜索结果 |

---

### FindStringEx

**功能说明**: 搜索字符串扩展

**函数签名**:
```go
func (dm *DmSoft) FindStringEx(hwnd int32, addr_range string, string_value string, multi_depth int32, step int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr_range | string | 地址范围 |
| string_value | string | 字符串值 |
| multi_depth | int32 | 多级深度 |
| step | int32 | 搜索步长 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 搜索结果 |

---

### FloatToData

**功能说明**: 单精度浮点数转二进制数据

**函数签名**:
```go
func (dm *DmSoft) FloatToData(value float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| value | float64 | 单精度浮点数值(用float64传递) |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 二进制数据字符串 |

---

### FreeProcessMemory

**功能说明**: 释放进程内存

**函数签名**:
```go
func (dm *DmSoft) FreeProcessMemory(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetCommandLine

**功能说明**: 获取命令行

**函数签名**:
```go
func (dm *DmSoft) GetCommandLine(hwnd int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 命令行字符串 |

---

### GetModuleBaseAddr

**功能说明**: 获取模块基址

**函数签名**:
```go
func (dm *DmSoft) GetModuleBaseAddr(hwnd int32, module string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| module | string | 模块名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 模块基址 |

---

### GetModuleSize

**功能说明**: 获取模块大小

**函数签名**:
```go
func (dm *DmSoft) GetModuleSize(hwnd int32, module string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| module | string | 模块名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 模块大小 |

---

### Int64ToInt32

**功能说明**: 64位整数转32位整数

**函数签名**:
```go
func (dm *DmSoft) Int64ToInt32(value int64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| value | int64 | 64位整数值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 32位整数 |

---

### IntToData

**功能说明**: 整数转二进制数据

**函数签名**:
```go
func (dm *DmSoft) IntToData(value int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| value | int32 | 整数值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 二进制数据字符串 |

---

### OpenProcess

**功能说明**: 打开进程

**函数签名**:
```go
func (dm *DmSoft) OpenProcess(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadDataAddrToBin

**功能说明**: 读取内存数据到二进制(指定地址)

**函数签名**:
```go
func (dm *DmSoft) ReadDataAddrToBin(hwnd int32, addr int32, size int32, data *int32, data_size *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int32 | 内存地址 |
| size | int32 | 读取大小 |
| data | *int32 | 数据指针 |
| data_size | *int32 | 数据大小指针 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadDataToBin

**功能说明**: 读取内存数据到二进制

**函数签名**:
```go
func (dm *DmSoft) ReadDataToBin(hwnd int32, addr string, size int32, data *int32, data_size *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 地址字符串 |
| size | int32 | 读取大小 |
| data | *int32 | 数据指针 |
| data_size | *int32 | 数据大小指针 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetMemoryFindResultToFile

**功能说明**: 设置内存搜索结果到文件

**函数签名**:
```go
func (dm *DmSoft) SetMemoryFindResultToFile(file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | 文件路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetMemoryHwndAsProcessId

**功能说明**: 设置内存句柄为进程ID

**函数签名**:
```go
func (dm *DmSoft) SetMemoryHwndAsProcessId(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetParam64ToPointer

**功能说明**: 设置64位参数为指针

**函数签名**:
```go
func (dm *DmSoft) SetParam64ToPointer(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### StringToData

**功能说明**: 字符串转二进制数据

**函数签名**:
```go
func (dm *DmSoft) StringToData(value string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| value | string | 字符串值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 二进制数据字符串 |

---

### TerminateProcessTree

**功能说明**: 终止进程树

**函数签名**:
```go
func (dm *DmSoft) TerminateProcessTree(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WriteDataAddrFromBin

**功能说明**: 从二进制写入内存数据(指定地址)

**函数签名**:
```go
func (dm *DmSoft) WriteDataAddrFromBin(hwnd int32, addr int32, data int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | int32 | 内存地址 |
| data | int32 | 数据指针 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### WriteDataFromBin

**功能说明**: 从二进制写入内存数据

**函数签名**:
```go
func (dm *DmSoft) WriteDataFromBin(hwnd int32, addr string, data int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| addr | string | 地址字符串 |
| data | int32 | 数据指针 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 10. 文件操作函数

### ReadFile

**功能说明**: 读取文件内容

**函数签名**:
```go
func (dm *DmSoft) ReadFile(file_name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 文件内容 |

---

### WriteFile

**功能说明**: 写入文件

**函数签名**:
```go
func (dm *DmSoft) WriteFile(file_name string, content string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |
| content | string | 内容 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### AppendFile

**功能说明**: 追加文件内容

**函数签名**:
```go
func (dm *DmSoft) AppendFile(file_name string, content string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |
| content | string | 内容 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### CopyFile

**功能说明**: 复制文件

**函数签名**:
```go
func (dm *DmSoft) CopyFile(src_file string, dst_file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| src_file | string | 源文件 |
| dst_file | string | 目标文件 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### MoveFile

**功能说明**: 移动文件

**函数签名**:
```go
func (dm *DmSoft) MoveFile(src_file string, dst_file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| src_file | string | 源文件 |
| dst_file | string | 目标文件 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### DeleteFile

**功能说明**: 删除文件

**函数签名**:
```go
func (dm *DmSoft) DeleteFile(file_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FileExists

**功能说明**: 判断文件是否存在

**函数签名**:
```go
func (dm *DmSoft) FileExists(file_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 存在返回1，不存在返回0 |

---

### GetFileSize

**功能说明**: 获取文件大小

**函数签名**:
```go
func (dm *DmSoft) GetFileSize(file_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 文件大小（字节） |

---

### GetFileLength

**功能说明**: 获取文件长度

**函数签名**:
```go
func (dm *DmSoft) GetFileLength(file_name string) int64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int64 | 文件长度 |

---

### GetFileModifyTime

**功能说明**: 获取文件修改时间

**函数签名**:
```go
func (dm *DmSoft) GetFileModifyTime(file_name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file_name | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 修改时间 |

---

### EnumFiles

**功能说明**: 枚举文件

**函数签名**:
```go
func (dm *DmSoft) EnumFiles(path string, type_ int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 路径 |
| type_ | int32 | 类型 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 文件列表 |

---

### CreateFolder

**功能说明**: 创建文件夹

**函数签名**:
```go
func (dm *DmSoft) CreateFolder(folder_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| folder_name | string | 文件夹名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### DeleteFolder

**功能说明**: 删除文件夹

**函数签名**:
```go
func (dm *DmSoft) DeleteFolder(folder_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| folder_name | string | 文件夹名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FolderExists

**功能说明**: 判断文件夹是否存在

**函数签名**:
```go
func (dm *DmSoft) FolderExists(folder_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| folder_name | string | 文件夹名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 存在返回1，不存在返回0 |

---

### DecodeFile
解码文件
**函数原型**: `func (dm *DmSoft) DecodeFile(file string, pwd string) int32`
**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | 文件路径 |
| pwd | string | 密码 |
**返回值**: int32 - 成功返回1，失败返回0

---

### DownloadFile
下载文件
**函数原型**: `func (dm *DmSoft) DownloadFile(url string, save_file string, timeout int32) int32`
**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| url | string | 下载地址 |
| save_file | string | 保存路径 |
| timeout | int32 | 超时时间(毫秒) |
**返回值**: int32 - 成功返回1，失败返回0

---

### EncodeFile
编码文件
**函数原型**: `func (dm *DmSoft) EncodeFile(file string, pwd string) int32`
**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | 文件路径 |
| pwd | string | 密码 |
**返回值**: int32 - 成功返回1，失败返回0

---

### IsFileExist
判断文件是否存在
**函数原型**: `func (dm *DmSoft) IsFileExist(file string) int32`
**参数**: file - 文件路径
**返回值**: int32 - 1存在，0不存在

---

### IsFolderExist
判断文件夹是否存在
**函数原型**: `func (dm *DmSoft) IsFolderExist(folder string) int32`
**参数**: folder - 文件夹路径
**返回值**: int32 - 1存在，0不存在

---

### SelectDirectory
选择目录对话框
**函数原型**: `func (dm *DmSoft) SelectDirectory() string`
**返回值**: string - 选中的目录路径

---

### SelectFile
选择文件对话框
**函数原型**: `func (dm *DmSoft) SelectFile() string`
**返回值**: string - 选中的文件路径

---

## 11. 进程操作函数

### EnumProcess

**功能说明**: 枚举进程

**函数签名**:
```go
func (dm *DmSoft) EnumProcess(name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| name | string | 进程名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 进程ID列表，格式: "pid1,pid2,pid3" |

---

### RunApp

**功能说明**: 运行程序

**函数签名**:
```go
func (dm *DmSoft) RunApp(path string, mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 程序路径 |
| mode | int32 | 模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### RunAppEx

**功能说明**: 扩展运行程序

**函数签名**:
```go
func (dm *DmSoft) RunAppEx(path string, cmd string, current_dir string, mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| path | string | 程序路径 |
| cmd | string | 命令行参数 |
| current_dir | string | 工作目录 |
| mode | int32 | 模式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### TerminateProcess

**功能说明**: 结束进程

**函数签名**:
```go
func (dm *DmSoft) TerminateProcess(pid int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pid | int32 | 进程ID |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetProcessInfo

**功能说明**: 获取进程信息

**函数签名**:
```go
func (dm *DmSoft) GetProcessInfo(pid int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pid | int32 | 进程ID |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 进程信息 |

---

## 12. 屏幕截图函数

### Capture

**功能说明**: 屏幕截图

**函数签名**:
```go
func (dm *DmSoft) Capture(x1 int32, y1 int32, x2 int32, y2 int32, file_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| file_name | string | 保存文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

**示例**:
```go
ret := dm.Capture(0, 0, 800, 600, "screenshot.bmp")
```

---

### CaptureJpg

**功能说明**: 截图为JPG格式

**函数签名**:
```go
func (dm *DmSoft) CaptureJpg(x1 int32, y1 int32, x2 int32, y2 int32, file_name string, quality int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| file_name | string | 保存文件名 |
| quality | int32 | 质量(1-100) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### CapturePng

**功能说明**: 截图为PNG格式

**函数签名**:
```go
func (dm *DmSoft) CapturePng(x1 int32, y1 int32, x2 int32, y2 int32, file_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| file_name | string | 保存文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### CaptureGif

**功能说明**: 截图为GIF格式

**函数签名**:
```go
func (dm *DmSoft) CaptureGif(x1 int32, y1 int32, x2 int32, y2 int32, file_name string, delay int32, time int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| file_name | string | 保存文件名 |
| delay | int32 | 延迟 |
| time | int32 | 时间 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetScreenData

**功能说明**: 获取屏幕数据

**函数签名**:
```go
func (dm *DmSoft) GetScreenData(x1 int32, y1 int32, x2 int32, y2 int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 数据指针 |

---

### GetScreenDataBmp

**功能说明**: 获取屏幕数据（BMP格式）

**函数签名**:
```go
func (dm *DmSoft) GetScreenDataBmp(x1 int32, y1 int32, x2 int32, y2 int32, data *string, size *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| data | *string | 数据（输出参数） |
| size | *int32 | 大小（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FreeScreenData

**功能说明**: 释放屏幕数据

**函数签名**:
```go
func (dm *DmSoft) FreeScreenData(handle int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| handle | int32 | 数据句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### GetScreenWidth

**功能说明**: 获取屏幕宽度

**函数签名**:
```go
func (dm *DmSoft) GetScreenWidth() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 屏幕宽度 |

---

### GetScreenHeight

**功能说明**: 获取屏幕高度

**函数签名**:
```go
func (dm *DmSoft) GetScreenHeight() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 屏幕高度 |

---

### GetScreenDepth

**功能说明**: 获取屏幕色深

**函数签名**:
```go
func (dm *DmSoft) GetScreenDepth() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 色深 |

---

### SetScreen

**功能说明**: 设置屏幕分辨率

**函数签名**:
```go
func (dm *DmSoft) SetScreen(width int32, height int32, depth int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| width | int32 | 宽度 |
| height | int32 | 高度 |
| depth | int32 | 色深 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetScreen

**功能说明**: 设置屏幕数据

**函数签名**:
```go
func (dm *DmSoft) SetScreen(data int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| data | int32 | 数据指针 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 系统信息函数

### GetDPI

**功能说明**: 获取系统DPI

**函数签名**:
```go
func (dm *DmSoft) GetDPI() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | DPI值 |

---

### GetOsType

**功能说明**: 获取操作系统类型

**函数签名**:
```go
func (dm *DmSoft) GetOsType() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 操作系统类型编码 |

---

### GetOsBuildNumber

**功能说明**: 获取系统版本号

**函数签名**:
```go
func (dm *DmSoft) GetOsBuildNumber() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 系统版本号 |

---

### GetTime

**功能说明**: 获取当前时间戳

**函数签名**:
```go
func (dm *DmSoft) GetTime() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 时间戳 |

---

### GetNetTime

**功能说明**: 获取网络时间

**函数签名**:
```go
func (dm *DmSoft) GetNetTime() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 网络时间戳 |

---

### GetNetTimeSafe

**功能说明**: 安全获取网络时间

**函数签名**:
```go
func (dm *DmSoft) GetNetTimeSafe() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 网络时间戳 |

---

### GetNetTimeByIp

**功能说明**: 通过IP获取网络时间

**函数签名**:
```go
func (dm *DmSoft) GetNetTimeByIp(ip string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| ip | string | IP地址 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 网络时间戳 |

---

### GetCpuType

**功能说明**: 获取CPU类型

**函数签名**:
```go
func (dm *DmSoft) GetCpuType() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | CPU类型字符串 |

---

### GetCpuUsage

**功能说明**: 获取CPU使用率

**函数签名**:
```go
func (dm *DmSoft) GetCpuUsage() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | CPU使用率百分比 |

---

### GetMemoryUsage

**功能说明**: 获取内存使用量

**函数签名**:
```go
func (dm *DmSoft) GetMemoryUsage() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 内存使用量(KB) |

---

### Is64Bit

**功能说明**: 判断是否64位系统

**函数签名**:
```go
func (dm *DmSoft) Is64Bit() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 64位返回1，32位返回0 |

---

### IsSurrpotVt

**功能说明**: 判断是否支持VT虚拟化技术

**函数签名**:
```go
func (dm *DmSoft) IsSurrpotVt() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 支持返回1，不支持返回0 |

---

### GetSystemInfo

**功能说明**: 获取系统信息

**函数签名**:
```go
func (dm *DmSoft) GetSystemInfo() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 系统信息字符串 |

---

### GetDisplayInfo

**功能说明**: 获取显示器信息

**函数签名**:
```go
func (dm *DmSoft) GetDisplayInfo() string
```

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 显示器信息字符串 |

---

### GetFps

**功能说明**: 获取当前帧率

**函数签名**:
```go
func (dm *DmSoft) GetFps() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 帧率值 |

---

### GetDiskReversion

**功能说明**: 获取磁盘版本信息

**函数签名**:
```go
func (dm *DmSoft) GetDiskReversion(index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 磁盘索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 磁盘版本字符串 |

---

### GetDiskSerial

**功能说明**: 获取磁盘序列号

**函数签名**:
```go
func (dm *DmSoft) GetDiskSerial(index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 磁盘索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 磁盘序列号字符串 |

---

### GetDiskModel

**功能说明**: 获取磁盘型号

**函数签名**:
```go
func (dm *DmSoft) GetDiskModel(index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 磁盘索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 磁盘型号字符串 |

---

### CapturePre

**功能说明**: 预处理截图

**函数签名**:
```go
func (dm *DmSoft) CapturePre(x1 int32, y1 int32, x2 int32, y2 int32, pre_type string, quality int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pre_type | string | 预处理类型 |
| quality | int32 | 质量 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### ImageToBmp

**功能说明**: 图片转BMP格式

**函数签名**:
```go
func (dm *DmSoft) ImageToBmp(pic_name string, bmp_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pic_name | string | 源图片文件名 |
| bmp_name | string | 输出BMP文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 14. AI相关函数

### LoadAi

**功能说明**: 加载AI模块

**函数签名**:
```go
func (dm *DmSoft) LoadAi(file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | AI模块文件路径，如"ai.module" |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 1成功，-1打开文件失败，-2内存初始化失败，-3参数错误，-4加载错误，-5AI模块初始化失败，-6内存分配失败 |

---

### LoadAiMemory

**功能说明**: 从内存加载AI模型

**函数签名**:
```go
func (dm *DmSoft) LoadAiMemory(addr int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| addr | int32 | 内存地址 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 1成功，0失败 |

---

### AiFindPic

**功能说明**: 使用AI模块查找指定区域内的图片，比传统FindPic效果更好，不需要训练

**函数签名**:
```go
func (dm *DmSoft) AiFindPic(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名，多个用\|分隔 |
| sim | float64 | 相似度(0.1-1.0) |
| dir | int32 | 查找方向: 0从左到右从上到下, 1从左到右从下到上, 2从右到左从上到下, 3从右到左从下到上 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引(从0开始)，失败返回-1 |

---

### AiFindPicEx

**功能说明**: AI扩展找图，返回所有找到的坐标

**函数签名**:
```go
func (dm *DmSoft) AiFindPicEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_name string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_name | string | 图片名 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 所有找到的坐标，格式: "id,x,y\|id,x,y..." |

---

### AiFindPicMem

**功能说明**: 从内存中AI找图

**函数签名**:
```go
func (dm *DmSoft) AiFindPicMem(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据（Base64编码） |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 找到的图片索引，失败返回-1 |

---

### AiFindPicMemEx

**功能说明**: 从内存中AI找图（扩展）

**函数签名**:
```go
func (dm *DmSoft) AiFindPicMemEx(x1 int32, y1 int32, x2 int32, y2 int32, pic_info string, sim float64, dir int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| pic_info | string | 图片数据 |
| sim | float64 | 相似度 |
| dir | int32 | 查找方向 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 所有找到的坐标 |

---

### AiEnableFindPicWindow

**功能说明**: 启用AI找图窗口

**函数签名**:
```go
func (dm *DmSoft) AiEnableFindPicWindow(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AiYoloSetVersion

**功能说明**: 设置YOLO版本

**函数签名**:
```go
func (dm *DmSoft) AiYoloSetVersion(version string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| version | string | YOLO版本号，如"v5-7.0" |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AiYoloSetModel

**功能说明**: 设置YOLO模型文件

**函数签名**:
```go
func (dm *DmSoft) AiYoloSetModel(model string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| model | string | 模型文件路径(.onnx或.dmx) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AiYoloSetModelMemory

**功能说明**: 从内存设置YOLO模型

**函数签名**:
```go
func (dm *DmSoft) AiYoloSetModelMemory(addr int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| addr | int32 | 内存地址 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AiYoloUseModel

**功能说明**: 切换当前使用的YOLO模型

**函数签名**:
```go
func (dm *DmSoft) AiYoloUseModel(index int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 模型索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### AiYoloDetectObjects

**功能说明**: 在指定范围内检测YOLO目标对象

**函数签名**:
```go
func (dm *DmSoft) AiYoloDetectObjects(x1 int32, y1 int32, x2 int32, y2 int32, prob float64, iou float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| prob | float64 | 置信度阈值，超过此值的对象才会被检测 |
| iou | float64 | IOU合并阈值，建议0.4-0.6 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 格式: "类名,置信度,x,y,w,h\|..."，未检测到返回空字符串 |

---

### AiYoloDetectObjectsToFile

**功能说明**: 检测YOLO目标并保存标注图到文件

**函数签名**:
```go
func (dm *DmSoft) AiYoloDetectObjectsToFile(x1 int32, y1 int32, x2 int32, y2 int32, prob float64, iou float64, file_name string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| prob | float64 | 置信度阈值 |
| iou | float64 | IOU合并阈值 |
| file_name | string | 保存的文件名(.bmp) |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 检测结果字符串 |

---

### AiYoloObjectsToString

**功能说明**: 将YOLO检测结果中的类名提取连接成字符串

**函数签名**:
```go
func (dm *DmSoft) AiYoloObjectsToString(objects string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| objects | string | AiYoloDetectObjects的返回值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 类名连接字符串 |

---

### AiYoloSortsObjects

**功能说明**: 对YOLO检测结果进行排序（从上到下，从左到右）

**函数签名**:
```go
func (dm *DmSoft) AiYoloSortsObjects(objects string, sort int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| objects | string | AiYoloDetectObjects的返回值 |
| sort | int32 | 排序方式 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 排序后的结果字符串 |

---

### AiYoloFreeModel

**功能说明**: 释放YOLO模型

**函数签名**:
```go
func (dm *DmSoft) AiYoloFreeModel(index int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 模型索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

### AiYoloDetectObjectsToDataBmp
YOLO检测对象到数据位图
**函数原型**: `func (dm *DmSoft) AiYoloDetectObjectsToDataBmp(x1 int32, y1 int32, x2 int32, y2 int32, prob float32, iou float32) int32`
**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| prob | float32 | 概率阈值 |
| iou | float32 | IOU阈值 |
**返回值**: int32 - 数据位图句柄

---

## 15. 汇编相关函数

### AsmAdd

**功能说明**: 添加MASM汇编指令到缓冲区

**函数签名**:
```go
func (dm *DmSoft) AsmAdd(asm_ins string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| asm_ins | string | MASM汇编指令，大小写均可，支持emit直接加入字节 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### AsmCall

**功能说明**: 执行AsmAdd添加到缓冲中的汇编指令

**函数签名**:
```go
func (dm *DmSoft) AsmCall(hwnd int32, mode int32) int64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| mode | int32 | 0本进程执行, 1远程线程, 2绑定后远程线程, 3绑定后hwnd所在线程, 4当前线程, 5APC注入, 6hwnd所在线程 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int64 | 执行后EAX/RAX值，-200执行错误，-201模式5未开启memory盾 |

---

### AsmCallEx

**功能说明**: 同AsmCall，但指定已分配的内存地址执行，避免每次分配内存开销

**函数签名**:
```go
func (dm *DmSoft) AsmCallEx(hwnd int32, mode int32, base_addr string) int64
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| mode | int32 | 模式，同AsmCall |
| base_addr | string | 十六进制格式地址，为空效果同AsmCall |

**返回值**:
| 类型 | 说明 |
|------|------|
| int64 | 执行后EAX/RAX值 |

---

### AsmClear

**功能说明**: 清除汇编指令缓冲区

**函数签名**:
```go
func (dm *DmSoft) AsmClear() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### AsmSetTimeout

**功能说明**: 设置AsmCall/AsmCallEx的超时参数

**函数签名**:
```go
func (dm *DmSoft) AsmSetTimeout(time_out int32, param int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| time_out | int32 | 最长等待时间(毫秒)，-1无限等待，默认10000 |
| param | int32 | 仅影响模式6，默认100毫秒 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### Assemble

**功能说明**: 将汇编缓冲区指令转换为机器码十六进制字符串

**函数签名**:
```go
func (dm *DmSoft) Assemble(base_addr int64, is_64bit int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| base_addr | int64 | 第一条指令所在地址 |
| is_64bit | int32 | 0=32位, 1=64位 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 机器码十六进制字符串，如"aa bb cc" |

---

### DisAssemble

**功能说明**: 将机器码转换为汇编语言

**函数签名**:
```go
func (dm *DmSoft) DisAssemble(asm_code string, base_addr int64, is_64bit int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| asm_code | string | 机器码十六进制字符串 |
| base_addr | int64 | 指令所在地址 |
| is_64bit | int32 | 0=32位, 1=64位 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | MASM汇编语言，多条指令以\|连接 |

---

### SetAsmHwndAsProcessId

**功能说明**: 设置AsmCall的hwnd参数当作进程PID，仅对AsmCall模式1起作用

**函数签名**:
```go
func (dm *DmSoft) SetAsmHwndAsProcessId(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 0关闭，1打开 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetShowAsmErrorMsg

**功能说明**: 设置是否弹出汇编功能中的错误提示，默认打开

**函数签名**:
```go
func (dm *DmSoft) SetShowAsmErrorMsg(show int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| show | int32 | 0不打开，1打开 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 16. 网络相关函数

### FaqSend

**功能说明**: 发送FAQ请求

**函数签名**:
```go
func (dm *DmSoft) FaqSend(server string, url string, data string, timeout int32, time_out_ex int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| server | string | 服务器地址 |
| url | string | URL路径 |
| data | string | 发送数据 |
| timeout | int32 | 超时时间 |
| time_out_ex | int32 | 扩展超时 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 返回结果 |

---

### FaqGetSize

**功能说明**: 获取FAQ大小

**函数签名**:
```go
func (dm *DmSoft) FaqGetSize() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | FAQ大小 |

---

### FaqCapture

**功能说明**: FAQ截图

**函数签名**:
```go
func (dm *DmSoft) FaqCapture(x1 int32, y1 int32, x2 int32, y2 int32, quality int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| quality | int32 | 图片质量 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FaqCaptureFromFile

**功能说明**: 从文件FAQ截图

**函数签名**:
```go
func (dm *DmSoft) FaqCaptureFromFile(x1 int32, y1 int32, x2 int32, y2 int32, quality int32, file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| quality | int32 | 图片质量 |
| file | string | 文件路径 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### FaqPost

**功能说明**: FAQ提交

**函数签名**:
```go
func (dm *DmSoft) FaqPost(data string, time_out int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| data | string | 提交数据 |
| time_out | int32 | 超时时间 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 返回结果 |

---

### FaqIsPosted

**功能说明**: FAQ是否已提交

**函数签名**:
```go
func (dm *DmSoft) FaqIsPosted() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 已提交返回1 |

---

### FaqCancel

**功能说明**: 取消FAQ

**函数签名**:
```go
func (dm *DmSoft) FaqCancel() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 17. INI配置文件函数

### ReadIni

**功能说明**: 从INI文件中读取指定信息

**函数签名**:
```go
func (dm *DmSoft) ReadIni(section string, key string, file string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名 |
| file | string | INI文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取到的内容 |

---

### WriteIni

**功能说明**: 向INI文件写入信息

**函数签名**:
```go
func (dm *DmSoft) WriteIni(section string, key string, value string, file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名 |
| value | string | 变量内容 |
| file | string | INI文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ReadIniPwd

**功能说明**: 从加密INI文件中读取指定信息

**函数签名**:
```go
func (dm *DmSoft) ReadIniPwd(section string, key string, file string, pwd string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名 |
| file | string | INI文件名 |
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 读取到的内容 |

---

### WriteIniPwd

**功能说明**: 向加密INI文件写入信息

**函数签名**:
```go
func (dm *DmSoft) WriteIniPwd(section string, key string, value string, file string, pwd string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名 |
| value | string | 变量内容 |
| file | string | INI文件名 |
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### DeleteIni

**功能说明**: 删除INI文件中的指定小节或键

**函数签名**:
```go
func (dm *DmSoft) DeleteIni(section string, key string, file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名，为空串则删除整个section |
| file | string | INI文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### DeleteIniPwd

**功能说明**: 删除加密INI文件中的指定小节或键

**函数签名**:
```go
func (dm *DmSoft) DeleteIniPwd(section string, key string, file string, pwd string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名 |
| key | string | 变量名，为空串则删除整个section |
| file | string | INI文件名 |
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### EnumIniKey

**功能说明**: 枚举INI文件中指定section的所有key名

**函数签名**:
```go
func (dm *DmSoft) EnumIniKey(section string, file string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名，不可为空 |
| file | string | INI文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 每个key用"\|"连接，如"aaa\|bbb\|ccc" |

---

### EnumIniKeyPwd

**功能说明**: 枚举加密INI文件中指定section的所有key名

**函数签名**:
```go
func (dm *DmSoft) EnumIniKeyPwd(section string, file string, pwd string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| section | string | 小节名，不可为空 |
| file | string | INI文件名 |
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 每个key用"\|"连接 |

---

### EnumIniSection

**功能说明**: 枚举INI文件中所有section

**函数签名**:
```go
func (dm *DmSoft) EnumIniSection(file string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | INI文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 每个小节名用"\|"连接 |

---

### EnumIniSectionPwd

**功能说明**: 枚举加密INI文件中所有section

**函数签名**:
```go
func (dm *DmSoft) EnumIniSectionPwd(file string, pwd string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| file | string | INI文件名 |
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 每个小节名用"\|"连接 |

---

## 13. Foobar绘图函数

### CreateFoobarCustom

**功能说明**: 根据指定位图创建自定义形状窗口

**函数签名**:
```go
func (dm *DmSoft) CreateFoobarCustom(hwnd int32, x int32, y int32, pic_name string, trans_color string, sim float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 父窗口句柄，0为桌面 |
| x | int32 | 左上角X坐标 |
| y | int32 | 左上角Y坐标 |
| pic_name | string | 位图文件名 |
| trans_color | string | 透明色(RRGGBB) |
| sim | float64 | 透明色相似值(0.1-1.0) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 创建成功的窗口句柄 |

---

### CreateFoobarEllipse

**功能说明**: 创建椭圆窗口

**函数签名**:
```go
func (dm *DmSoft) CreateFoobarEllipse(hwnd int32, x int32, y int32, w int32, h int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 父窗口句柄，0为桌面 |
| x | int32 | 左上角X坐标 |
| y | int32 | 左上角Y坐标 |
| w | int32 | 宽度 |
| h | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 创建成功的窗口句柄 |

---

### CreateFoobarRect

**功能说明**: 创建矩形窗口

**函数签名**:
```go
func (dm *DmSoft) CreateFoobarRect(hwnd int32, x int32, y int32, w int32, h int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 父窗口句柄，0为桌面 |
| x | int32 | 左上角X坐标 |
| y | int32 | 左上角Y坐标 |
| w | int32 | 宽度 |
| h | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 创建成功的窗口句柄 |

---

### CreateFoobarRoundRect

**功能说明**: 创建圆角矩形窗口

**函数签名**:
```go
func (dm *DmSoft) CreateFoobarRoundRect(hwnd int32, x int32, y int32, w int32, h int32, rw int32, rh int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 父窗口句柄，0为桌面 |
| x | int32 | 左上角X坐标 |
| y | int32 | 左上角Y坐标 |
| w | int32 | 宽度 |
| h | int32 | 高度 |
| rw | int32 | 圆角宽度 |
| rh | int32 | 圆角高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 创建成功的窗口句柄 |

---

### FoobarClearText

**功能说明**: 清除Foobar滚动文本区

**函数签名**:
```go
func (dm *DmSoft) FoobarClearText(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarClose

**功能说明**: 关闭Foobar窗口。必须调用此函数关闭，否则会造成内存泄漏

**函数签名**:
```go
func (dm *DmSoft) FoobarClose(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarDrawLine

**功能说明**: 在Foobar窗口内画线条

**函数签名**:
```go
func (dm *DmSoft) FoobarDrawLine(hwnd int32, x1 int32, y1 int32, x2 int32, y2 int32, color string, style int32, width int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x1 | int32 | 起点X坐标 |
| y1 | int32 | 起点Y坐标 |
| x2 | int32 | 终点X坐标 |
| y2 | int32 | 终点Y坐标 |
| color | string | 颜色值 |
| style | int32 | 0实线，1虚线 |
| width | int32 | 线条宽度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarDrawPic

**功能说明**: 在Foobar窗口内绘制图像

**函数签名**:
```go
func (dm *DmSoft) FoobarDrawPic(hwnd int32, x int32, y int32, pic_name string, trans_color string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| pic_name | string | 图像文件名 |
| trans_color | string | 透明色 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarDrawText

**功能说明**: 在Foobar窗口内绘制文字

**函数签名**:
```go
func (dm *DmSoft) FoobarDrawText(hwnd int32, x int32, y int32, w int32, h int32, text string, color string, align int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x | int32 | 左上角X坐标 |
| y | int32 | 左上角Y坐标 |
| w | int32 | 宽度 |
| h | int32 | 高度 |
| text | string | 文字内容 |
| color | string | 文字颜色 |
| align | int32 | 1左对齐，2中间对齐，4右对齐 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarFillRect

**功能说明**: 在Foobar窗口内填充矩形

**函数签名**:
```go
func (dm *DmSoft) FoobarFillRect(hwnd int32, x1 int32, y1 int32, x2 int32, y2 int32, color string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 填充颜色 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarLock

**功能说明**: 锁定Foobar窗口，不能通过鼠标移动

**函数签名**:
```go
func (dm *DmSoft) FoobarLock(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarPrintText

**功能说明**: 向Foobar窗口输出滚动文字

**函数签名**:
```go
func (dm *DmSoft) FoobarPrintText(hwnd int32, text string, color string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| text | string | 文本内容 |
| color | string | 文本颜色 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarSetFont

**功能说明**: 设置Foobar窗口的字体

**函数签名**:
```go
func (dm *DmSoft) FoobarSetFont(hwnd int32, font_name string, size int32, flag int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| font_name | string | 系统字体名 |
| size | int32 | 字体大小 |
| flag | int32 | 0正常，1粗体，2斜体，4下划线，可组合 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarSetSave

**功能说明**: 设置Foobar滚动文本区保存到文件

**函数签名**:
```go
func (dm *DmSoft) FoobarSetSave(hwnd int32, file string, enable int32, header string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| file | string | 保存文件名 |
| enable | int32 | 0关闭文件输出，1开启 |
| header | string | 附加头信息格式串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarSetTrans

**功能说明**: 设置Foobar窗口是否透明

**函数签名**:
```go
func (dm *DmSoft) FoobarSetTrans(hwnd int32, is_trans int32, color string, sim float64) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| is_trans | int32 | 0不透明，1透明 |
| color | string | 透明色(RRGGBB) |
| sim | float64 | 透明色相似值(0.1-1.0) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarStartGif

**功能说明**: 在Foobar窗口绘制GIF动画

**函数签名**:
```go
func (dm *DmSoft) FoobarStartGif(hwnd int32, x int32, y int32, pic_name string, repeat_limit int32, delay int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| pic_name | string | GIF文件名 |
| repeat_limit | int32 | 重复次数，0一直循环 |
| delay | int32 | 帧间隔，0使用GIF内置时间 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarStopGif

**功能说明**: 停止Foobar窗口的GIF动画

**函数签名**:
```go
func (dm *DmSoft) FoobarStopGif(hwnd int32, x int32, y int32, pic_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| pic_name | string | GIF文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarTextLineGap

**功能说明**: 设置滚动文本区的文字行间距，默认3

**函数签名**:
```go
func (dm *DmSoft) FoobarTextLineGap(hwnd int32, line_gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| line_gap | int32 | 行间距 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarTextPrintDir

**功能说明**: 设置滚动文本区的文字输出方向，默认0

**函数签名**:
```go
func (dm *DmSoft) FoobarTextPrintDir(hwnd int32, dir int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| dir | int32 | 0向下输出，1向上输出 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarTextRect

**功能说明**: 设置Foobar窗口的滚动文本框范围，默认窗口区域

**函数签名**:
```go
func (dm *DmSoft) FoobarTextRect(hwnd int32, x int32, y int32, w int32, h int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |
| x | int32 | X坐标 |
| y | int32 | Y坐标 |
| w | int32 | 宽度 |
| h | int32 | 高度 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarUnlock

**功能说明**: 解锁Foobar窗口，可通过鼠标移动

**函数签名**:
```go
func (dm *DmSoft) FoobarUnlock(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### FoobarUpdate

**功能说明**: 刷新Foobar窗口，所有绘制完成后必须调用此函数刷新窗口

**函数签名**:
```go
func (dm *DmSoft) FoobarUpdate(hwnd int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | Foobar窗口句柄 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 18. 字库相关函数

### SetDict

**功能说明**: 设置字库文件

**函数签名**:
```go
func (dm *DmSoft) SetDict(index int32, file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |
| file | string | 字库文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### UseDict

**功能说明**: 切换使用哪个字库进行识别

**函数签名**:
```go
func (dm *DmSoft) UseDict(index int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### AddDict

**功能说明**: 给指定字库添加一条字库信息（内存中，立即生效）

**函数签名**:
```go
func (dm *DmSoft) AddDict(index int32, dict_info string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |
| dict_info | string | 字库描述串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ClearDict

**功能说明**: 清空指定字库（内存中的字库，不是文件）

**函数签名**:
```go
func (dm *DmSoft) ClearDict(index int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SaveDict

**功能说明**: 保存指定字库到文件

**函数签名**:
```go
func (dm *DmSoft) SaveDict(index int32, file string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |
| file | string | 文件名 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetDict

**功能说明**: 获取指定字库中指定条目的信息

**函数签名**:
```go
func (dm *DmSoft) GetDict(index int32, font_index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |
| font_index | int32 | 字库条目序号(从0开始) |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 字库条目信息，失败返回空串 |

---

### GetDictCount

**功能说明**: 获取指定字库中的字符数量

**函数签名**:
```go
func (dm *DmSoft) GetDictCount(index int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号(0-99) |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 字库字符数量 |

---

### GetNowDict

**功能说明**: 获取当前使用的字库序号

**函数签名**:
```go
func (dm *DmSoft) GetNowDict() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 字库序号(0-99) |

---

### SetExportDict

**功能说明**: 设置导出字库

**函数签名**:
```go
func (dm *DmSoft) SetExportDict(index int32, dict_name string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库序号 |
| dict_name | string | 字库名称 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### EnableShareDict

**功能说明**: 启用共享字库

**函数签名**:
```go
func (dm *DmSoft) EnableShareDict(enable int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| enable | int32 | 1启用，0禁用 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetWordLineHeight

**功能说明**: 设置文字行高，默认10

**函数签名**:
```go
func (dm *DmSoft) SetWordLineHeight(line_height int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| line_height | int32 | 行高 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetWordLineHeightNoDict

**功能说明**: 设置文字行高（无字典模式）

**函数签名**:
```go
func (dm *DmSoft) SetWordLineHeightNoDict(line_height int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| line_height | int32 | 行高 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetMinColGap

**功能说明**: 设置最小列间距，默认0，可提高识别精度

**函数签名**:
```go
func (dm *DmSoft) SetMinColGap(min_col_gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| min_col_gap | int32 | 最小列间距 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetColGapNoDict

**功能说明**: 设置列间距（无字典模式）

**函数签名**:
```go
func (dm *DmSoft) SetColGapNoDict(col_gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| col_gap | int32 | 列间距 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetMinRowGap

**功能说明**: 设置最小行间距，默认1

**函数签名**:
```go
func (dm *DmSoft) SetMinRowGap(min_row_gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| min_row_gap | int32 | 最小行间距 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetRowGapNoDict

**功能说明**: 设置行间距（无字典模式）

**函数签名**:
```go
func (dm *DmSoft) SetRowGapNoDict(row_gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| row_gap | int32 | 行间距 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetExactOcr

**功能说明**: 设置是否开启精准识别

**函数签名**:
```go
func (dm *DmSoft) SetExactOcr(exact_ocr int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| exact_ocr | int32 | 0关闭精准识别，1开启 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetWords

**功能说明**: 识别指定范围内所有满足条件的词组

**函数签名**:
```go
func (dm *DmSoft) GetWords(x1 int32, y1 int32, x2 int32, y2 int32, color string, sim float64) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式串 |
| sim | float64 | 相似度(0.1-1.0) |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 识别到的格式串，需用专用函数解析 |

---

### GetWordsNoDict

**功能说明**: 同GetWords但不使用字库，只识别大概形状位置

**函数签名**:
```go
func (dm *DmSoft) GetWordsNoDict(x1 int32, y1 int32, x2 int32, y2 int32, color string) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| x1 | int32 | 左上角X坐标 |
| y1 | int32 | 左上角Y坐标 |
| x2 | int32 | 右下角X坐标 |
| y2 | int32 | 右下角Y坐标 |
| color | string | 颜色格式串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 识别到的格式串 |

---

### GetWordResultCount

**功能说明**: 获取词组识别结果的数量

**函数签名**:
```go
func (dm *DmSoft) GetWordResultCount(str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | GetWords的返回值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 词组数量 |

---

### GetWordResultPos

**功能说明**: 获取词组识别结果中各个词组的坐标

**函数签名**:
```go
func (dm *DmSoft) GetWordResultPos(str string, index int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | GetWords的返回值 |
| index | int32 | 第几个词组 |
| x | *int32 | X坐标（输出参数） |
| y | *int32 | Y坐标（输出参数） |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### GetWordResultStr

**功能说明**: 获取词组识别结果中各个词组的内容

**函数签名**:
```go
func (dm *DmSoft) GetWordResultStr(str string, index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | GetWords的返回值 |
| index | int32 | 第几个词组 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 词组内容 |

---

### GetDictInfo

**功能说明**: 获取字库信息

**函数签名**:
```go
func (dm *DmSoft) GetDictInfo(index int32, font_index int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库索引 |
| font_index | int32 | 字体索引 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 字库信息字符串 |

---

### GetResultCount

**功能说明**: 获取结果数量

**函数签名**:
```go
func (dm *DmSoft) GetResultCount(str string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | 结果字符串 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 结果数量 |

---

### GetResultPos

**功能说明**: 获取结果位置

**函数签名**:
```go
func (dm *DmSoft) GetResultPos(str string, index int32, x *int32, y *int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| str | string | 结果字符串 |
| index | int32 | 索引 |
| x | *int32 | 返回X坐标 |
| y | *int32 | 返回Y坐标 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetDictMem

**功能说明**: 从内存设置字库

**函数签名**:
```go
func (dm *DmSoft) SetDictMem(index int32, addr int32, size int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| index | int32 | 字库索引 |
| addr | int32 | 内存地址 |
| size | int32 | 数据大小 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetDictPwd

**功能说明**: 设置字库密码

**函数签名**:
```go
func (dm *DmSoft) SetDictPwd(pwd string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| pwd | string | 密码 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetWordGap

**功能说明**: 设置文字间距

**函数签名**:
```go
func (dm *DmSoft) SetWordGap(gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| gap | int32 | 间距值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

### SetWordGapNoDict

**功能说明**: 设置文字间距(无字典)

**函数签名**:
```go
func (dm *DmSoft) SetWordGapNoDict(gap int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| gap | int32 | 间距值 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1 |

---

## 19. 杂项函数

### ActiveInputMethod

**功能说明**: 激活指定窗口所在进程的输入法

**函数签名**:
```go
func (dm *DmSoft) ActiveInputMethod(hwnd int32, input_method string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| input_method | string | 输入法名字 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### CheckInputMethod

**功能说明**: 检测指定窗口所在线程输入法是否开启

**函数签名**:
```go
func (dm *DmSoft) CheckInputMethod(hwnd int32, input_method string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| hwnd | int32 | 窗口句柄 |
| input_method | string | 输入法名字 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 0未开启，1开启 |

---

### EnterCri

**功能说明**: 进入临界区

**函数签名**:
```go
func (dm *DmSoft) EnterCri() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 0不可进入，1已进入临界区 |

---

### InitCri

**功能说明**: 初始化临界区，必须在脚本开头调用一次

**函数签名**:
```go
func (dm *DmSoft) InitCri() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### LeaveCri

**功能说明**: 离开临界区，释放调用对象占用的互斥信号量

**函数签名**:
```go
func (dm *DmSoft) LeaveCri() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### ExecuteCmd

**功能说明**: 执行CMD指令并返回输出结果

**函数签名**:
```go
func (dm *DmSoft) ExecuteCmd(cmd string, current_dir string, time_out int32) string
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| cmd | string | 需要执行的CMD指令 |
| current_dir | string | 执行时所在目录，空表示当前目录 |
| time_out | int32 | 超时毫秒，0表示一直等待 |

**返回值**:
| 类型 | 说明 |
|------|------|
| string | 执行结果，空表示失败 |

---

### FindInputMethod

**功能说明**: 检测系统中是否安装了指定输入法

**函数签名**:
```go
func (dm *DmSoft) FindInputMethod(input_method string) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| input_method | string | 输入法名字 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 0未安装，1已安装 |

---

### ReleaseRef

**功能说明**: 强制降低对象的引用计数（高级接口）

**函数签名**:
```go
func (dm *DmSoft) ReleaseRef() int32
```

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

### SetExitThread

**功能说明**: 设置当前对象的退出线程标记

**函数签名**:
```go
func (dm *DmSoft) SetExitThread(mode int32) int32
```

**参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| mode | int32 | 1开启(会解绑), 2开启(不会解绑), 0关闭 |

**返回值**:
| 类型 | 说明 |
|------|------|
| int32 | 成功返回1，失败返回0 |

---

## 注意事项

1. **必须调用Init()**: 创建DmSoft对象后必须调用Init()初始化COM对象，否则所有函数调用都会失败。

2. **编码自动转换**: 库已内置UTF-8到GBK的自动编码转换，中文参数可直接使用。

3. **32位程序**: 必须编译为32位程序才能调用大漠DLL。

4. **多线程使用**: 每个线程需独立创建DmSoft实例并各自调用Init()。

5. **资源释放**: 使用完毕后调用Release()释放资源。

6. **内存操作**: 内存操作函数需要目标进程有足够的权限。

7. **找图找色**: 图片文件需要放在SetPath设置的路径下，或使用完整路径。
