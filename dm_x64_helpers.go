package dmsoft

import (
	"fmt"
)

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

func (dm *DmSoft) getMethodOffset(method string) uint32 {
	if off, ok := methodOffsets[method]; ok {
		return off
	}
	fmt.Printf("警告: 方法 %s 未在偏移量表中找到\n", method)
	return 0
}
