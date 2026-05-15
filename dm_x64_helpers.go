// Package dmsoft 提供大漠插件的 Go 绑定，支持 32 位和 64 位架构
// 本文件实现了64位模式下的COM调用封装函数
// 通过管道(pipe)与32位helper进程通信，实现跨进程调用大漠插件功能
package dmsoft

import (
	"fmt"
)

// comCallInt32 通过管道调用返回int32类型的方法
// 参数:
//   - method: 方法名称（对应methodOffsets表中的键名）
//   - params: 可变参数列表（支持int32, string, float64等类型）
//
// 返回值: int32类型的调用结果
func (dm *DmSoft) comCallInt32(method string, params ...interface{}) int32 {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return 0
		}
		return dm.pipeCallInt32(offset, params...)
	}
	return 0
}

// comCallStr 通过管道调用返回string类型的方法
// 参数:
//   - method: 方法名称
//   - params: 可变参数列表
//
// 返回值: string类型的调用结果（自动处理GBK到UTF-8的编码转换）
func (dm *DmSoft) comCallStr(method string, params ...interface{}) string {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return ""
		}
		return dm.pipeCallStr(offset, params...)
	}
	return ""
}

// comCallInt64 通过管道调用返回int64类型的方法
// 参数:
//   - method: 方法名称
//   - params: 可变参数列表
//
// 返回值: int64类型的调用结果（用于返回64位整数，如内存地址等）
func (dm *DmSoft) comCallInt64(method string, params ...interface{}) int64 {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return 0
		}
		return dm.pipeCallInt64(offset, params...)
	}
	return 0
}

// comCallFloat64 通过管道调用返回float64类型的方法
// 参数:
//   - method: 方法名称
//   - params: 可变参数列表
//
// 返回值: float64类型的调用结果（用于返回双精度浮点数）
func (dm *DmSoft) comCallFloat64(method string, params ...interface{}) float64 {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return 0
		}
		return dm.pipeCallFloat64(offset, params...)
	}
	return 0
}

// comCallBool 通过管道调用返回bool类型的方法
// 参数:
//   - method: 方法名称
//   - params: 可变参数列表
//
// 返回值: bool类型的调用结果（将int32结果转换为bool，非0为true）
func (dm *DmSoft) comCallBool(method string, params ...interface{}) bool {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return false
		}
		return dm.pipeCallInt32(offset, params...) != 0
	}
	return false
}

// comCallWithOutVars 通过管道调用带有输出参数的方法（返回int32类型）
// 参数:
//   - method: 方法名称
//   - inParams: 输入参数列表
//   - outVars: 输出参数指针列表（用于接收返回值，如坐标等）
//
// 返回值: int32类型的调用结果
// 说明: 输出参数通过gob序列化跨进程回传，实现64位到32位的数据传递
func (dm *DmSoft) comCallWithOutVars(method string, inParams []interface{}, outVars ...*int32) int32 {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return 0
		}
		return dm.pipeCallWithOutVars(offset, inParams, outVars...)
	}
	return 0
}

// comCallStrWithOutVars 通过管道调用带有输出参数的方法（返回string类型）
// 参数:
//   - method: 方法名称
//   - inParams: 输入参数列表
//   - outVars: 输出参数指针列表
//
// 返回值: string类型的调用结果（自动处理编码转换）
func (dm *DmSoft) comCallStrWithOutVars(method string, inParams []interface{}, outVars ...*int32) string {
	if dm.pipeConn != nil {
		offset := dm.getMethodOffset(method)
		if offset == 0 {
			return ""
		}
		return dm.pipeCallStrWithOutVars(offset, inParams, outVars...)
	}
	return ""
}

// getMethodOffset 根据方法名获取其在DLL中的偏移量
// 参数:
//   - method: 方法名称（如"FindPic"、"GetColor"等）
//
// 返回值: uint32类型的偏移量，如果方法不存在则返回0并打印警告
// 说明: 偏移量表(methodOffsets)定义在dm_x64_pipe.go中，包含所有大漠插件API函数的偏移地址
func (dm *DmSoft) getMethodOffset(method string) uint32 {
	if off, ok := methodOffsets[method]; ok {
		return off
	}
	fmt.Printf("警告: 方法 %s 未在偏移量表中找到\n", method)
	return 0
}
