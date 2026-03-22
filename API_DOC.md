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
| flag | int32 | 状态标志: 0关闭, 1激活, 2最小化, 3最大化, 4还原, 5隐藏, 6显示, 7置顶, 8取消置顶 |

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

## 注意事项

1. **必须调用Init()**: 创建DmSoft对象后必须调用Init()初始化COM对象，否则所有函数调用都会失败。

2. **编码自动转换**: 库已内置UTF-8到GBK的自动编码转换，中文参数可直接使用。

3. **32位程序**: 必须编译为32位程序才能调用大漠DLL。

4. **多线程使用**: 每个线程需独立创建DmSoft实例并各自调用Init()。

5. **资源释放**: 使用完毕后调用Release()释放资源。

6. **内存操作**: 内存操作函数需要目标进程有足够的权限。

7. **找图找色**: 图片文件需要放在SetPath设置的路径下，或使用完整路径。
