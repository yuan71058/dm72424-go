package dmsoft

import (
	"encoding/gob"
	"fmt"
)

type callRequest struct {
	Offset  uint32
	RetType uint8
	NOut    uint8
	Args    []callArg
}

type callArg struct {
	Type   uint8
	IVal   int32
	SVal   string
	FVal   float64
	I64Val int64
}

type callResponse struct {
	RetType uint8
	IRet    int32
	SRet    string
	FRet    float64
	IRet64  int64
	OutVals []int32
	Err     string
}

var methodOffsets = map[string]uint32{
	"GetDiskReversion":              109040,
	"LoadAiMemory":                  108256,
	"FaqSend":                       114016,
	"FindPicSimMem":                 121744,
	"Ver":                           100320,
	"SetPath":                       123808,
	"SetShowAsmErrorMsg":            101392,
	"FindStrS":                      116832,
	"GetWordsNoDict":                99024,
	"GetOsBuildNumber":              104240,
	"GetID":                         105184,
	"SetMouseSpeed":                 124800,
	"FindData":                      109760,
	"SendPaste":                     122944,
	"GetColor":                      117424,
	"LoadPicByte":                   121408,
	"WriteFloatAddr":                117312,
	"SetWordLineHeight":             101296,
	"AsmCall":                       114656,
	"FindColorBlock":                113568,
	"DisAssemble":                   112656,
	"RegEx":                         98864,
	"EncodeFile":                    106528,
	"WriteString":                   115936,
	"FindStrFastEx":                 122000,
	"AsmCallEx":                     99632,
	"FindDoubleEx":                  110416,
	"SetFindPicMultithreadLimit":    107616,
	"SendString2":                   99888,
	"DownCpu":                       112960,
	"DmGuard":                       103552,
	"SpeedNormalGraphic":            101184,
	"FindPicSim":                    98768,
	"WriteInt":                      112416,
	"SetMemoryHwndAsProcessId":      107984,
	"WriteDataFromBin":              118304,
	"SetMinColGap":                  110560,
	"KeyPressStr":                   102528,
	"LockDisplay":                   108304,
	"FindStrWithFontE":              112544,
	"EnumIniKey":                    108032,
	"MatchPicName":                  117984,
	"EnableFakeActive":              107888,
	"FaqGetSize":                    103456,
	"ExecuteCmd":                    116928,
	"EnableRealKeypad":              105648,
	"SetDisplayRefreshDelay":        111344,
	"MiddleClick":                   108560,
	"AiYoloSortsObjects":            120480,
	"WriteDataAddr":                 105744,
	"RGB2BGR":                       115744,
	"DisablePowerSave":              121952,
	"GetClientSize":                 103344,
	"EnableMouseMsg":                101344,
	"EnableKeypadMsg":               120864,
	"GetFileLength":                 111296,
	"GetRemoteApiAddress":           122192,
	"DmGuardParams":                 105472,
	"DownloadFile":                  123648,
	"WriteDoubleAddr":               115232,
	"EnableIme":                     120192,
	"TerminateProcessTree":          114240,
	"FoobarClose":                   102480,
	"FindNearestPos":                112480,
	"CreateFoobarRect":              119072,
	"GetCursorPos":                  121680,
	"FindColorBlockEx":              103840,
	"FindFloat":                     103216,
	"GetProcessInfo":                119024,
	"ReadFile":                      114464,
	"FindShapeEx":                   99792,
	"SetWindowText":                 113008,
	"ForceUnBindWindow":             120144,
	"ReadIntAddr":                   99712,
	"FindShape":                     123856,
	"GetRealPath":                   105008,
	"EnableSpeedDx":                 115472,
	"UnLoadDriver":                  105696,
	"GetMemoryUsage":                106064,
	"MiddleDown":                    109872,
	"EnumIniSection":                117184,
	"CheckUAC":                      123104,
	"OpenProcess":                   124624,
	"IsDisplayDead":                 114896,
	"WriteIniPwd":                   115872,
	"GetNetTime":                    107712,
	"ReadFloat":                     100976,
	"DisableCloseDisplayAndSleep":   114416,
	"GetWindowTitle":                110816,
	"Assemble":                      119584,
	"GetMousePointWindow":           105424,
	"SetExportDict":                 119392,
	"Delay":                         106480,
	"Reg":                           121344,
	"FoobarStopGif":                 108096,
	"ReadFileData":                  115808,
	"FindPicSimEx":                  113728,
	"Capture":                       119456,
	"GetScreenWidth":                113920,
	"FindStrWithFontEx":             118848,
	"SetLocale":                     100928,
	"AsmAdd":                        121232,
	"GetScreenHeight":               117792,
	"CaptureGif":                    120912,
	"ReadDataAddrToBin":             111792,
	"ReadDataToBin":                 104480,
	"FindPicS":                      101952,
	"FindPic":                       104032,
	"FindMultiColor":                109360,
	"HackSpeed":                     104352,
	"FindPicE":                      114144,
	"MiddleUp":                      115072,
	"GetWindow":                     120752,
	"SetUAC":                        108608,
	"FoobarSetSave":                 124736,
	"WheelDown":                     112848,
	"FloatToData":                   100464,
	"EnableFindPicMultithread":      118048,
	"DisableScreenSave":             112800,
	"AiFindPicEx":                   119136,
	"SendString":                    114832,
	"EnterCri":                      116336,
	"FindPicSimMemE":                113296,
	"Delays":                        123328,
	"CreateFoobarCustom":            105872,
	"FindStringEx":                  124384,
	"GetClientRect":                 105808,
	"AiYoloSetModel":                104416,
	"FoobarSetTrans":                117248,
	"GetForegroundFocus":            108512,
	"GetForegroundWindow":           115360,
	"SetExcludeRegion":              104832,
	"SendStringIme2":                119520,
	"ActiveInputMethod":             124320,
	"FoobarDrawPic":                 114288,
	"AiYoloSetVersion":              118496,
	"FindColorE":                    120384,
	"LeftClick":                     118096,
	"IsFileExist":                   113824,
	"Is64Bit":                       110512,
	"FindShapeE":                    120592,
	"GetDisplayInfo":                122992,
	"SetEnumWindowDelay":            114720,
	"RegNoMac":                      118960,
	"KeyUpChar":                     121904,
	"SetDisplayAcceler":             101088,
	"SetRowGapNoDict":               118256,
	"EnableMouseAccuracy":           123760,
	"MoveTo":                        109088,
	"KeyPressChar":                  116464,
	"RightDown":                     124576,
	"AiYoloSetModelMemory":          117600,
	"WriteIni":                      101232,
	"DmGuardLoadCustom":             106896,
	"CreateFolder":                  113120,
	"EnableRealMouse":               105952,
	"GetBasePath":                   107312,
	"GetFps":                        106016,
	"EnableGetColorByCapture":       109216,
	"SetDisplayInput":               110944,
	"Hex64":                         105296,
	"ScreenToClient":                111392,
	"AiEnableFindPicWindow":         100064,
	"ReadIni":                       102912,
	"ImageToBmp":                    109152,
	"SetDisplayDelay":               122784,
	"WheelUp":                       102688,
	"CopyFile":                      100688,
	"FindWindowEx":                  115408,
	"SetFindPicMultithreadCount":    106784,
	"GetScreenDataBmp":              107136,
	"GetWordResultPos":              114352,
	"LeftDoubleClick":               101136,
	"ReadStringAddr":                118608,
	"ReadData":                      111232,
	"AddDict":                       106336,
	"SetInputDm":                    108656,
	"GetWindowProcessId":            124464,
	"WriteDataAddrFromBin":          121120,
	"AiFindPicMemEx":                102976,
	"TerminateProcess":              112032,
	"VirtualQueryEx":                101632,
	"EnableKeypadSync":              109968,
	"AiYoloUseModel":                110032,
	"DeleteFile":                    99408,
	"GetScreenDepth":                102384,
	"FindColor":                     106112,
	"MoveR":                         113504,
	"LockInput":                     124272,
	"IntToData":                     122272,
	"FaqPost":                       107440,
	"GetColorHSV":                   116192,
	"FindWindowSuper":               108432,
	"EnableBind":                    116576,
	"SetAero":                       102640,
	"DecodeFile":                    122496,
	"FindPicExS":                    100368,
	"WriteStringAddr":               122720,
	"GetCommandLine":                100752,
	"SelectFile":                    118144,
	"FindPicSimMemEx":               124912,
	"GetWordResultStr":              104768,
	"EnablePicCache":                99536,
	"FindStrExS":                    100528,
	"LoadPic":                       124128,
	"FindStrFast":                   115584,
	"FindDouble":                    102192,
	"SetParam64ToPointer":           99952,
	"SetMemoryFindResultToFile":     110704,
	"WaitKey":                       114528,
	"CreateFoobarEllipse":           114592,
	"MoveFile":                      102272,
	"Stop":                          100880,
	"ReleaseRef":                    111072,
	"GetColorBGR":                   100000,
	"EnumIniKeyPwd":                 116768,
	"GetMac":                        123536,
	"UseDict":                       104656,
	"FindDataEx":                    123200,
	"Md5":                           117376,
	"BGR2RGB":                       118736,
	"FindColorEx":                   103600,
	"OcrExOne":                      112080,
	"CmpColor":                      109648,
	"OcrInFile":                     110608,
	"CheckInputMethod":              101792,
	"MoveWindow":                    119648,
	"GetClipboard":                  116624,
	"FindStr":                       110320,
	"FoobarClearText":               113072,
	"ClientToScreen":                116512,
	"GetCursorShape":                111984,
	"GetWordResultCount":            103984,
	"SelectDirectory":               116000,
	"CapturePng":                    114080,
	"KeyDownChar":                   105600,
	"CaptureJpg":                    106400,
	"FindStrEx":                     106640,
	"FaqCapture":                    118416,
	"ShowScrMsg":                    112208,
	"SetKeypadDelay":                110256,
	"SetScreen":                     115168,
	"Play":                          105072,
	"FindWindowByProcessId":         104176,
	"WriteDouble":                   116048,
	"GetWindowThreadId":             107504,
	"GetBindWindow":                 109712,
	"FindWindow":                    104288,
	"AiFindPic":                     121536,
	"FindInt":                       106256,
	"IsBind":                        119232,
	"SetSimMode":                    122896,
	"GetNowDict":                    101584,
	"GetNetTimeSafe":                107760,
	"GetMachineCode":                113456,
	"VirtualAllocEx":                99104,
	"GetPath":                       109600,
	"EnumWindowSuper":               107360,
	"GetModuleBaseAddr":             108848,
	"EnumWindowByProcessId":         124672,
	"UnBindWindow":                  101904,
	"GetLastError":                  107936,
	"FoobarDrawText":                119712,
	"SetMinRowGap":                  122144,
	"LeftUp":                        113680,
	"WriteFile":                     105536,
	"SetWindowSize":                 98560,
	"FaqCaptureFromFile":            116256,
	"ReadDataAddr":                  123584,
	"IsSurrpotVt":                   106992,
	"GetWindowProcessPath":          105232,
	"ClearDict":                     123152,
	"SaveDict":                      115520,
	"ShowTaskBarIcon":               119328,
	"GetAveHSV":                     100176,
	"ReadIniPwd":                    102064,
	"FaqIsPosted":                   102864,
	"LeftDown":                      106736,
	"DmGuardExtract":                112160,
	"ExitOs":                        115024,
	"FetchWord":                     117840,
	"GetDiskSerial":                 112352,
	"GetDictCount":                  99584,
	"GetDict":                       99184,
	"SetDict":                       121280,
	"AiYoloObjectsToString":         111456,
	"GetKeyState":                   103296,
	"RightClick":                    101040,
	"EnumWindowByProcess":           110192,
	"GetDiskModel":                  102128,
	"SendStringIme":                 124000,
	"AppendPicAddr":                 106832,
	"DeleteFolder":                  118800,
	"GetDPI":                        107664,
	"GetCpuType":                    102432,
	"WriteIntAddr":                  100240,
	"GetSpecialWindow":              102336,
	"EnumProcess":                   112288,
	"AsmClear":                      119968,
	"GetWindowState":                100112,
	"FindStrFastE":                  120288,
	"SetColGapNoDict":               102592,
	"AiYoloDetectObjects":           116112,
	"RunApp":                        122832,
	"FindString":                    110752,
	"GetOsType":                     121632,
	"Ocr":                           110992,
	"ReadString":                    121472,
	"ReadFloatAddr":                 100816,
	"Beep":                          104544,
	"LoadAi":                        106944,
	"GetCpuUsage":                   121072,
	"EnableShareDict":               108992,
	"AiYoloDetectObjectsToFile":     109504,
	"FoobarUnlock":                  123952,
	"GetSystemInfo":                 115680,
	"GetResultCount":                116720,
	"EnumWindow":                    115296,
	"GetResultPos":                  102800,
	"KeyDown":                       115120,
	"SetWordLineHeightNoDict":       103792,
	"AiFindPicMem":                  111696,
	"FoobarTextRect":                108784,
	"GetPointWindow":                118544,
	"FindMultiColorEx":              122560,
	"FreeProcessMemory":             111120,
	"GetMachineCodeNoMac":           120544,
	"FindWindowByProcess":           122336,
	"GetColorNum":                   124048,
	"SetWindowState":                102736,
	"CheckFontSmooth":               117552,
	"IsFolderExist":                 121184,
	"FaqCancel":                     113968,
	"SetWindowTransparent":          112896,
	"SwitchBindWindow":              109920,
	"EnableFontSmooth":              103936,
	"StringToData":                  114768,
	"GetWindowRect":                 122656,
	"FindPicEx":                     108160,
	"GetWords":                      107808,
	"SetExactOcr":                   123280,
	"EnableMouseSync":               98496,
	"CapturePre":                    109456,
	"BindWindowEx":                  99456,
	"FaqCaptureString":              106208,
	"FoobarTextLineGap":             124848,
	"FoobarDrawLine":                116384,
	"FindInputMethod":               113872,
	"SetPicPwd":                     123712,
	"GetCursorSpot":                 125056,
	"InitCri":                       120240,
	"FindPicMemE":                   109264,
	"FindStrFastS":                  98672,
	"DeleteIniPwd":                  99344,
	"AiYoloDetectObjectsToDataBmp":  98928,
	"AiYoloFreeModel":               106592,
	"DisableFontSmooth":             118368,
	"SetExitThread":                 101536,
	"FindPicMemEx":                  101440,
	"GetDmCount":                    125008,
	"FindMulColor":                  111552,
	"FaqFetch":                      117744,
	"RegExNoMac":                    107552,
	"FoobarUpdate":                  119280,
	"ReadDouble":                    110128,
	"GetCursorShapeEx":              117488,
	"DoubleToData":                  111856,
	"SetWordGapNoDict":              123392,
	"ReadDoubleAddr":                113392,
	"FoobarLock":                    109824,
	"FindStrFastExS":                124176,
	"FindStrWithFont":               119856,
	"VirtualProtectEx":              108912,
	"GetWindowClass":                117056,
	"SetMouseDelay":                 104592,
	"ReadInt":                       112720,
	"GetAveRGB":                     118192,
	"GetScreenData":                 125104,
	"GetMouseSpeed":                 99248,
	"Int64ToInt32":                  110880,
	"FindFloatEx":                   107040,
	"FoobarPrintText":               108720,
	"OcrEx":                         113168,
	"FreePic":                       103408,
	"WriteData":                     123040,
	"MoveDD":                        121840,
	"SetShowErrorMsg":               101856,
	"SetDictMem":                    104704,
	"SetClipboard":                  104960,
	"FindPicMem":                    103696,
	"CreateFoobarRoundRect":         108352,
	"WriteFloat":                    111920,
	"VirtualFreeEx":                 105120,
	"GetDictInfo":                   100624,
	"KeyPress":                      118688,
	"SetClientSize":                 104896,
	"ExcludePos":                    120992,
	"MoveToEx":                      120688,
	"SetDictPwd":                    104128,
	"FoobarSetFont":                 111632,
	"GetNetTimeByIp":                105360,
	"EnableKeypadPatch":             116672,
	"FoobarStartGif":                117664,
	"FindMultiColorE":               101696,
	"SetWordGap":                    98624,
	"GetLocale":                     122096,
	"GetModuleSize":                 120016,
	"FindStrE":                      122400,
	"KeyUp":                         113248,
	"SortPosDistance":                117120,
	"EnableDisplayDebug":            99296,
	"DeleteIni":                     111168,
	"FindIntEx":                     107216,
	"BindWindow":                    120080,
	"GetPicSize":                    114960,
	"AsmSetTimeout":                 117920,
	"LockMouseRect":                 119792,
	"FindPicSimE":                   123440,
	"EnumIniSectionPwd":             116992,
	"RightUp":                       111504,
	"FoobarTextPrintDir":            103072,
	"GetDir":                        124512,
	"Hex32":                         110080,
	"LeaveCri":                      120816,
	"GetTime":                       103504,
	"FoobarFillRect":                103136,
}

func buildPipeArgs(params []interface{}) []callArg {
	args := make([]callArg, 0, len(params))
	for _, p := range params {
		switch v := p.(type) {
		case int32:
			args = append(args, callArg{Type: 0, IVal: v})
		case string:
			args = append(args, callArg{Type: 1, SVal: v})
		case float64:
			args = append(args, callArg{Type: 2, FVal: v})
		case *int32:
			args = append(args, callArg{Type: 3})
		case float32:
			args = append(args, callArg{Type: 4, FVal: float64(v)})
		case int64:
			args = append(args, callArg{Type: 5, I64Val: v})
		default:
			args = append(args, callArg{Type: 0, IVal: 0})
		}
	}
	return args
}

func (dm *DmSoft) pipeCall(offset uint32, retType uint8, params []interface{}) callResponse {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.pipeConn == nil || dm.pipeEnc == nil || dm.pipeDec == nil {
		return callResponse{RetType: retType, Err: "pipe not connected"}
	}

	args := buildPipeArgs(params)
	req := callRequest{
		Offset:  offset,
		RetType: retType,
		NOut:    0,
		Args:    args,
	}

	if err := dm.pipeEnc.Encode(&req); err != nil {
		return callResponse{RetType: retType, Err: fmt.Sprintf("encode error: %v", err)}
	}

	var resp callResponse
	if err := dm.pipeDec.Decode(&resp); err != nil {
		return callResponse{RetType: retType, Err: fmt.Sprintf("decode error: %v", err)}
	}

	return resp
}

func (dm *DmSoft) pipeCallWithOut(offset uint32, retType uint8, inParams []interface{}, outVars ...*int32) callResponse {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.pipeConn == nil || dm.pipeEnc == nil || dm.pipeDec == nil {
		return callResponse{RetType: retType, Err: "pipe not connected"}
	}

	allParams := make([]interface{}, 0, len(inParams)+len(outVars))
	allParams = append(allParams, inParams...)
	for range outVars {
		allParams = append(allParams, (*int32)(nil))
	}

	args := buildPipeArgs(allParams)
	req := callRequest{
		Offset:  offset,
		RetType: retType,
		NOut:    uint8(len(outVars)),
		Args:    args,
	}

	if err := dm.pipeEnc.Encode(&req); err != nil {
		return callResponse{RetType: retType, Err: fmt.Sprintf("encode error: %v", err)}
	}

	var resp callResponse
	if err := dm.pipeDec.Decode(&resp); err != nil {
		return callResponse{RetType: retType, Err: fmt.Sprintf("decode error: %v", err)}
	}

	for i, ov := range outVars {
		if i < len(resp.OutVals) {
			*ov = resp.OutVals[i]
		}
	}

	return resp
}

func (dm *DmSoft) pipeCallInt32(offset uint32, params ...interface{}) int32 {
	resp := dm.pipeCall(offset, 0, params)
	return resp.IRet
}

func (dm *DmSoft) pipeCallStr(offset uint32, params ...interface{}) string {
	resp := dm.pipeCall(offset, 1, params)
	return resp.SRet
}

func (dm *DmSoft) pipeCallFloat64(offset uint32, params ...interface{}) float64 {
	resp := dm.pipeCall(offset, 2, params)
	return resp.FRet
}

func (dm *DmSoft) pipeCallInt64(offset uint32, params ...interface{}) int64 {
	resp := dm.pipeCall(offset, 3, params)
	return resp.IRet64
}

func (dm *DmSoft) pipeCallWithOutVars(offset uint32, inParams []interface{}, outVars ...*int32) int32 {
	resp := dm.pipeCallWithOut(offset, 0, inParams, outVars...)
	return resp.IRet
}

func (dm *DmSoft) pipeCallStrWithOutVars(offset uint32, inParams []interface{}, outVars ...*int32) string {
	resp := dm.pipeCallWithOut(offset, 1, inParams, outVars...)
	return resp.SRet
}

func init() {
	gob.Register(callRequest{})
	gob.Register(callArg{})
	gob.Register(callResponse{})
	gob.Register([]callArg{})
	gob.Register([]int32{})
}
