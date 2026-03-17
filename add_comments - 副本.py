#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
add_comments.py - 为dmsoft_impl.go添加函数注释
"""

import re

# 函数名到中文描述的映射
FUNC_COMMENTS = {
    # 磁盘与硬件信息
    'GetDiskReversion': '获取磁盘版本信息',
    'GetDiskSerial': '获取磁盘序列号',
    'GetDiskModel': '获取磁盘型号',
    
    # AI与内存操作
    'LoadAiMemory': '从内存加载AI模型',
    'LoadAi': '从文件加载AI模型',
    
    # FAQ相关
    'FaqSend': '发送FAQ请求',
    'FaqGetSize': '获取FAQ大小',
    'FaqCapture': 'FAQ截图',
    'FaqCaptureFromFile': '从文件FAQ截图',
    'FaqPost': 'FAQ提交',
    'FaqIsPosted': 'FAQ是否已提交',
    'FaqCancel': '取消FAQ',
    
    # 图片查找
    'FindPic': '查找图片',
    'FindPicEx': '高级查找图片',
    'FindPicE': '查找图片(返回坐标字符串)',
    'FindPicS': '查找图片(返回图片索引)',
    'FindPicSim': '查找图片相似度',
    'FindPicSimEx': '高级查找图片相似度',
    'FindPicSimMem': '在内存中查找图片(相似度)',
    'FindPicSimMemE': '在内存中查找图片(返回坐标字符串)',
    'FindPicSimMemEx': '高级在内存中查找图片',
    'FindPicExS': '高级查找图片(返回详细字符串)',
    'LoadPic': '加载图片到内存',
    'LoadPicByte': '从内存加载图片',
    
    # AI图片查找
    'AiFindPic': 'AI查找图片',
    'AiFindPicEx': 'AI高级查找图片',
    'AiFindPicMem': 'AI内存查找图片',
    'AiFindPicMemEx': 'AI高级内存查找图片',
    'AiEnableFindPicWindow': '启用AI找图窗口',
    
    # YOLO相关
    'AiYoloSetVersion': '设置YOLO版本',
    'AiYoloSetModel': '设置YOLO模型',
    'AiYoloSetModelMemory': '从内存设置YOLO模型',
    'AiYoloUseModel': '使用YOLO模型',
    'AiYoloDetectObjects': 'YOLO检测对象',
    'AiYoloDetectObjectsToFile': 'YOLO检测对象并保存到文件',
    'AiYoloObjectsToString': 'YOLO对象转字符串',
    'AiYoloSortsObjects': 'YOLO排序对象',
    
    # 版本与注册
    'Ver': '获取大漠版本号',
    'GetID': '获取ID',
    'GetMachineCode': '获取机器码',
    'GetMachineCodeNoMac': '获取机器码(不含MAC)',
    'GetMac': '获取MAC地址',
    'Reg': '注册大漠',
    'RegEx': '扩展注册',
    'RegNoMac': '注册(不含MAC)',
    'RegExNoMac': '扩展注册(不含MAC)',
    'GetLastError': '获取最后错误码',
    
    # 路径与配置
    'SetPath': '设置资源路径',
    'GetPath': '获取资源路径',
    'GetBasePath': '获取基础路径',
    'GetRealPath': '获取真实路径',
    
    # 文字识别与OCR
    'FindStrS': '查找文字(返回字符串)',
    'FindStr': '查找文字',
    'FindStrEx': '高级查找文字',
    'FindStrExS': '高级查找文字(返回详细字符串)',
    'FindStrFast': '快速查找文字',
    'FindStrFastE': '快速查找文字(返回坐标字符串)',
    'FindStrFastEx': '高级快速查找文字',
    'FindStrWithFontE': '指定字体查找文字',
    'FindStrWithFontEx': '高级指定字体查找文字',
    'Ocr': 'OCR识别',
    'OcrExOne': 'OCR识别单个',
    'OcrInFile': '从文件OCR识别',
    'GetWords': '获取文字',
    'GetWordsNoDict': '无字典获取文字',
    'FetchWord': '提取文字',
    'GetWordResultCount': '获取文字结果数量',
    'GetWordResultPos': '获取文字结果位置',
    'GetWordResultStr': '获取文字结果字符串',
    
    # 字库操作
    'SetDict': '设置字库',
    'UseDict': '使用字库',
    'AddDict': '添加字库',
    'ClearDict': '清除字库',
    'SaveDict': '保存字库',
    'GetDict': '获取字库',
    'GetDictCount': '获取字库数量',
    'GetNowDict': '获取当前字库索引',
    'SetExportDict': '设置导出字库',
    'EnableShareDict': '启用共享字库',
    
    # 字库参数设置
    'SetWordLineHeight': '设置文字行高',
    'SetWordLineHeightNoDict': '设置文字行高(无字典)',
    'SetMinColGap': '设置最小列间距',
    'SetColGapNoDict': '设置列间距(无字典)',
    'SetMinRowGap': '设置最小行间距',
    'SetRowGapNoDict': '设置行间距(无字典)',
    'SetExactOcr': '设置精确OCR',
    
    # 键盘操作
    'KeyPressChar': '按键(字符)',
    'KeyDownChar': '按下按键(字符)',
    'KeyUpChar': '弹起按键(字符)',
    'KeyPress': '按键',
    'KeyDown': '按下按键',
    'KeyUp': '弹起按键',
    'KeyPressStr': '按键字符串',
    'SetKeypadDelay': '设置键盘延迟',
    'SetSimMode': '设置模拟模式',
    'GetKeyState': '获取按键状态',
    'WaitKey': '等待按键',
    
    # 鼠标操作
    'MoveTo': '移动鼠标',
    'MoveR': '相对移动鼠标',
    'MoveToEx': '扩展移动鼠标',
    'MoveDD': 'DD移动鼠标',
    'LeftClick': '左键单击',
    'LeftDown': '左键按下',
    'LeftUp': '左键弹起',
    'LeftDoubleClick': '左键双击',
    'RightClick': '右键单击',
    'RightDown': '右键按下',
    'RightUp': '右键弹起',
    'MiddleClick': '中键单击',
    'MiddleDown': '中键按下',
    'MiddleUp': '中键弹起',
    'WheelDown': '滚轮向下',
    'WheelUp': '滚轮向上',
    'GetCursorPos': '获取鼠标位置',
    'GetCursorShape': '获取鼠标形状',
    'GetCursorShapeEx': '获取鼠标形状(扩展)',
    'GetCursorSpot': '获取鼠标光点',
    'SetMouseSpeed': '设置鼠标速度',
    'GetMouseSpeed': '获取鼠标速度',
    'SetMouseDelay': '设置鼠标延迟',
    'LockMouseRect': '锁定鼠标区域',
    'EnableMouseMsg': '启用鼠标消息',
    'EnableRealMouse': '启用真实鼠标',
    'EnableMouseAccuracy': '启用鼠标精度',
    
    # 颜色操作
    'GetColor': '获取颜色',
    'GetColorBGR': '获取BGR颜色',
    'GetColorHSV': '获取HSV颜色',
    'GetAveHSV': '获取平均HSV',
    'CmpColor': '比较颜色',
    'FindColor': '查找颜色',
    'FindColorE': '查找颜色(返回坐标字符串)',
    'FindColorEx': '高级查找颜色',
    'FindColorBlock': '查找色块',
    'FindColorBlockEx': '高级查找色块',
    'GetColorNum': '获取颜色数量',
    'RGB2BGR': 'RGB转BGR',
    'BGR2RGB': 'BGR转RGB',
    
    # 截图操作
    'Capture': '截图',
    'CaptureGif': '截图GIF',
    'CaptureJpg': '截图JPG',
    'CapturePng': '截图PNG',
    'GetScreenData': '获取屏幕数据',
    'GetScreenDataBmp': '获取屏幕数据BMP',
    'GetScreenWidth': '获取屏幕宽度',
    'GetScreenHeight': '获取屏幕高度',
    'GetScreenDepth': '获取屏幕深度',
    'GetDisplayInfo': '获取显示信息',
    'SetScreen': '设置屏幕',
    'LockDisplay': '锁定显示',
    'SetDisplayInput': '设置显示输入',
    'SetDisplayDelay': '设置显示延迟',
    'SetDisplayRefreshDelay': '设置显示刷新延迟',
    'SetDisplayAcceler': '设置显示加速',
    'EnableGetColorByCapture': '启用截图取色',
    'IsDisplayDead': '检测屏幕是否死机',
    
    # 窗口操作
    'BindWindow': '绑定窗口',
    'BindWindowEx': '扩展绑定窗口',
    'UnBindWindow': '解绑窗口',
    'ForceUnBindWindow': '强制解绑窗口',
    'SwitchBindWindow': '切换绑定窗口',
    'GetBindWindow': '获取绑定窗口',
    'IsBind': '判断是否绑定',
    'EnableBind': '启用绑定',
    'FindWindow': '查找窗口',
    'FindWindowEx': '扩展查找窗口',
    'FindWindowSuper': '超级查找窗口',
    'FindWindowByProcess': '通过进程查找窗口',
    'FindWindowByProcessId': '通过进程ID查找窗口',
    'EnumWindow': '枚举窗口',
    'EnumWindowByProcess': '通过进程枚举窗口',
    'EnumWindowByProcessId': '通过进程ID枚举窗口',
    'EnumWindowSuper': '超级枚举窗口',
    'GetWindow': '获取窗口',
    'GetWindowRect': '获取窗口矩形',
    'GetWindowClass': '获取窗口类名',
    'GetWindowTitle': '获取窗口标题',
    'SetWindowTitle': '设置窗口标题',
    'GetWindowState': '获取窗口状态',
    'SetWindowState': '设置窗口状态',
    'GetWindowProcessId': '获取窗口进程ID',
    'GetWindowThreadId': '获取窗口线程ID',
    'GetWindowProcessPath': '获取窗口进程路径',
    'GetForegroundWindow': '获取前台窗口',
    'GetForegroundFocus': '获取前台焦点窗口',
    'GetMousePointWindow': '获取鼠标指向窗口',
    'GetPointWindow': '获取指定坐标窗口',
    'ScreenToClient': '屏幕坐标转客户区坐标',
    'ClientToScreen': '客户区坐标转屏幕坐标',
    'SetWindowTransparent': '设置窗口透明度',
    'ShowTaskBarIcon': '显示任务栏图标',
    'SetEnumWindowDelay': '设置枚举窗口延迟',
    'MoveWindow': '移动窗口',
    'SetWindowSize': '设置窗口大小',
    'GetClientRect': '获取客户区矩形',
    'GetClientSize': '获取客户区大小',
    
    # 内存操作
    'ReadInt': '读取整数',
    'ReadIntAddr': '读取整数(地址)',
    'WriteInt': '写入整数',
    'WriteIntAddr': '写入整数(地址)',
    'ReadFloat': '读取浮点数',
    'ReadFloatAddr': '读取浮点数(地址)',
    'WriteFloat': '写入浮点数',
    'WriteFloatAddr': '写入浮点数(地址)',
    'ReadDouble': '读取双精度浮点数',
    'ReadDoubleAddr': '读取双精度浮点数(地址)',
    'WriteDouble': '写入双精度浮点数',
    'WriteDoubleAddr': '写入双精度浮点数(地址)',
    'ReadString': '读取字符串',
    'ReadStringAddr': '读取字符串(地址)',
    'WriteString': '写入字符串',
    'WriteStringAddr': '写入字符串(地址)',
    'ReadData': '读取数据',
    'ReadDataAddr': '读取数据(地址)',
    'ReadDataToBin': '读取数据到二进制',
    'ReadDataAddrToBin': '读取数据到二进制(地址)',
    'WriteData': '写入数据',
    'WriteDataAddr': '写入数据(地址)',
    'WriteDataFromBin': '从二进制写入数据',
    'WriteDataAddrFromBin': '从二进制写入数据(地址)',
    'FindData': '查找数据',
    'FindDataEx': '高级查找数据',
    'FindInt': '查找整数',
    'FindIntEx': '高级查找整数',
    'FindFloat': '查找浮点数',
    'FindFloatEx': '高级查找浮点数',
    'FindDouble': '查找双精度浮点数',
    'FindDoubleEx': '高级查找双精度浮点数',
    'FindString': '查找字符串',
    'FindStringEx': '高级查找字符串',
    'VirtualAllocEx': '虚拟内存分配',
    'VirtualProtectEx': '虚拟内存保护',
    'VirtualQueryEx': '虚拟内存查询',
    'GetModuleBaseAddr': '获取模块基址',
    'GetModuleSize': '获取模块大小',
    'GetRemoteApiAddress': '获取远程API地址',
    'OpenProcess': '打开进程',
    'TerminateProcess': '终止进程',
    'TerminateProcessTree': '终止进程树',
    'FreeProcessMemory': '释放进程内存',
    'GetMemoryUsage': '获取内存使用',
    'SetMemoryHwndAsProcessId': '设置内存句柄为进程ID',
    'SetMemoryFindResultToFile': '设置内存查找结果到文件',
    
    # 汇编相关
    'AsmAdd': '添加汇编指令',
    'AsmClear': '清除汇编指令',
    'AsmCall': '调用汇编代码',
    'AsmCallEx': '扩展调用汇编代码',
    'AsmSetTimeout': '设置汇编超时',
    'DisAssemble': '反汇编',
    'SetShowAsmErrorMsg': '设置显示汇编错误信息',
    
    # 文件操作
    'ReadFile': '读取文件',
    'ReadFileData': '读取文件数据',
    'WriteFile': '写入文件',
    'CopyFile': '复制文件',
    'MoveFile': '移动文件',
    'DeleteFile': '删除文件',
    'IsFileExist': '判断文件是否存在',
    'GetFileLength': '获取文件长度',
    'CreateFolder': '创建文件夹',
    'DeleteFolder': '删除文件夹',
    'IsFolderExist': '判断文件夹是否存在',
    'SelectFile': '选择文件',
    'SelectDirectory': '选择目录',
    'GetDir': '获取目录',
    'EncodeFile': '加密文件',
    'DecodeFile': '解密文件',
    'ImageToBmp': '图片转BMP',
    'DownloadFile': '下载文件',
    
    # INI操作
    'ReadIni': '读取INI',
    'ReadIniPwd': '读取INI(密码)',
    'WriteIni': '写入INI',
    'WriteIniPwd': '写入INI(密码)',
    'DeleteIni': '删除INI',
    'DeleteIniPwd': '删除INI(密码)',
    'EnumIniKey': '枚举INI键',
    'EnumIniKeyPwd': '枚举INI键(密码)',
    'EnumIniSection': '枚举INI节',
    'EnumIniSectionPwd': '枚举INI节(密码)',
    
    # 系统信息
    'GetOsType': '获取操作系统类型',
    'GetOsBuildNumber': '获取操作系统版本号',
    'GetTime': '获取时间',
    'GetNetTime': '获取网络时间',
    'GetNetTimeSafe': '安全获取网络时间',
    'GetCpuType': '获取CPU类型',
    'GetCpuUsage': '获取CPU使用率',
    'GetDPI': '获取DPI',
    'Is64Bit': '判断是否64位',
    'IsSurrpotVt': '是否支持VT',
    'CheckUAC': '检查UAC',
    'SetUAC': '设置UAC',
    'GetSystemInfo': '获取系统信息',
    'GetCommandLine': '获取命令行',
    'ExitOs': '退出系统',
    'SetLocale': '设置区域',
    'GetLocale': '获取区域',
    'DisableScreenSave': '禁用屏保',
    'DisablePowerSave': '禁用节能',
    'DisableCloseDisplayAndSleep': '禁用关闭显示器和睡眠',
    'DownCpu': '降低CPU使用率',
    
    # 进程操作
    'EnumProcess': '枚举进程',
    'GetProcessInfo': '获取进程信息',
    
    # Foobar操作
    'CreateFoobarRect': '创建Foobar矩形',
    'CreateFoobarEllipse': '创建Foobar椭圆',
    'CreateFoobarRoundRect': '创建Foobar圆角矩形',
    'CreateFoobarCustom': '创建自定义Foobar',
    'FoobarClose': '关闭Foobar',
    'FoobarFillRect': 'Foobar填充矩形',
    'FoobarDrawText': 'Foobar绘制文字',
    'FoobarDrawPic': 'Foobar绘制图片',
    'FoobarDrawLine': 'Foobar绘制线条',
    'FoobarLock': '锁定Foobar',
    'FoobarUnlock': '解锁Foobar',
    'FoobarUpdate': '更新Foobar',
    'FoobarSetSave': '设置Foobar保存',
    'FoobarSetTrans': '设置Foobar透明度',
    'FoobarSetFont': '设置Foobar字体',
    'FoobarTextRect': 'Foobar文字矩形',
    'FoobarTextLineGap': 'Foobar文字行间距',
    'FoobarTextPrintDir': 'Foobar文字打印方向',
    'FoobarPrintText': 'Foobar打印文字',
    'FoobarClearText': 'Foobar清除文字',
    'FoobarStartGif': 'Foobar开始GIF',
    'FoobarStopGif': 'Foobar停止GIF',
    
    # 多点找色
    'FindMultiColor': '多点找色',
    'FindMultiColorEx': '高级多点找色',
    
    # 形状查找
    'FindShape': '查找形状',
    'FindShapeE': '查找形状(返回坐标字符串)',
    'FindShapeEx': '高级查找形状',
    
    # 其他
    'SetParam64ToPointer': '设置参数64位转指针',
    'SetFindPicMultithreadCount': '设置找图多线程数量',
    'SetFindPicMultithreadLimit': '设置找图多线程限制',
    'EnableFindPicMultithread': '启用找图多线程',
    'SetExcludeRegion': '设置排除区域',
    'EnablePicCache': '启用图片缓存',
    'MatchPicName': '匹配图片名',
    'AppendPicAddr': '追加图片地址',
    'FindNearestPos': '查找最近位置',
    'SortPosDistance': '按距离排序位置',
    'ExcludePos': '排除位置',
    'Delay': '延迟',
    'Delays': '随机延迟',
    'Play': '播放声音',
    'Beep': '蜂鸣',
    'SetClipboard': '设置剪贴板',
    'GetClipboard': '获取剪贴板',
    'SendString': '发送字符串',
    'SendString2': '发送字符串2',
    'SendStringIme': '发送字符串(输入法)',
    'SendStringIme2': '发送字符串(输入法2)',
    'SendPaste': '发送粘贴',
    'ActiveInputMethod': '激活输入法',
    'CheckInputMethod': '检查输入法',
    'FindInputMethod': '查找输入法',
    'EnableIme': '启用输入法',
    'LockInput': '锁定输入',
    'EnableKeypadSync': '启用键盘同步',
    'EnableKeypadMsg': '启用键盘消息',
    'EnableKeypadPatch': '启用键盘补丁',
    'EnableRealKeypad': '启用真实键盘',
    'EnableFakeActive': '启用假激活',
    'EnableSpeedDx': '启用速度DX',
    'SpeedNormalGraphic': '速度正常图形',
    'SetAero': '设置Aero',
    'SetExitThread': '设置退出线程',
    'DmGuard': '大漠守护',
    'DmGuardParams': '大漠守护参数',
    'DmGuardLoadCustom': '大漠守护加载自定义',
    'DmGuardExtract': '大漠守护解压',
    'UnLoadDriver': '卸载驱动',
    'HackSpeed': '加速',
    'RunApp': '运行程序',
    'ExecuteCmd': '执行命令',
    'Md5': '计算MD5',
    'Hex64': '64位转十六进制',
    'Hex32': '32位转十六进制',
    'IntToData': '整数转数据',
    'FloatToData': '浮点数转数据',
    'DoubleToData': '双精度转数据',
    'StringToData': '字符串转数据',
    'Int64ToInt32': 'int64转int32',
    'GetFps': '获取FPS',
    'SetInputDm': '设置输入大漠',
    'GetDmCount': '获取大漠计数',
    'ReleaseRef': '释放引用',
    'Stop': '停止',
    'EnterCri': '进入临界区',
    'LeaveCri': '离开临界区',
    'InitCri': '初始化临界区',
    'CheckFontSmooth': '检查字体平滑',
    'EnableFontSmooth': '启用字体平滑',
    'DisableFontSmooth': '禁用字体平滑',
    'SetShowErrorMsg': '设置显示错误信息',
    'EnableDisplayDebug': '启用显示调试',
    'SetDictMem': '设置字库内存',
    'SetDictPwd': '设置字库密码',
    'GetDictInfo': '获取字库信息',
    'ShowScrMsg': '显示屏幕消息',
    'FaqCaptureString': 'FAQ捕获字符串',
    'FaqFetch': 'FAQ获取',
    'AiYoloFreeModel': '释放YOLO模型',
    'AiYoloDetectObjectsToDataBmp': 'YOLO检测对象到数据BMP',
    'FindPicMem': '内存查找图片',
    'FindPicMemE': '内存查找图片(返回坐标字符串)',
    'FindPicMemEx': '高级内存查找图片',
    'FindStrFastS': '快速查找文字(返回字符串)',
    'FindStrFastExS': '高级快速查找文字(返回字符串)',
    'FindStrWithFont': '指定字体查找文字',
    'FindStrE': '查找文字(返回坐标字符串)',
    'OcrEx': '高级OCR识别',
    'SetWordGap': '设置文字间距',
    'SetWordGapNoDict': '设置文字间距(无字典)',
}

def get_param_comment(param_name, param_type):
    """根据参数名和类型生成参数注释"""
    param_desc = {
        'hwnd': '窗口句柄',
        'x1': '左上角X坐标',
        'y1': '左上角Y坐标',
        'x2': '右下角X坐标',
        'y2': '右下角Y坐标',
        'x': 'X坐标',
        'y': 'Y坐标',
        'width': '宽度',
        'height': '高度',
        'w': '宽度',
        'h': '高度',
        'color': '颜色',
        'sim': '相似度',
        'dir': '方向',
        'type': '类型',
        'type_': '类型',
        'mode': '模式',
        'index': '索引',
        'file': '文件路径',
        'path': '路径',
        'addr': '地址',
        'size': '大小',
        'len': '长度',
        'length': '长度',
        'data': '数据',
        'str': '字符串',
        'string': '字符串',
        'text': '文本',
        'name': '名称',
        'pic_name': '图片名称',
        'pic_info': '图片信息',
        'delta_color': '颜色偏差',
        'offset_color': '偏移颜色',
        'first_color': '第一个颜色',
        'font_name': '字体名称',
        'font_size': '字体大小',
        'delay': '延迟(毫秒)',
        'timeout': '超时时间',
        'time_out': '超时时间',
        'time': '时间',
        'quality': '质量',
        'en': '启用标志(1启用,0禁用)',
        'enable': '启用标志(1启用,0禁用)',
        'flag': '标志',
        'filter': '过滤条件',
        'sort': '排序',
        'step': '步长',
        'count': '数量',
        'multi_thread': '多线程',
        'request_type': '请求类型',
        'server': '服务器地址',
        'handle': '句柄',
        'code': '注册码',
        'ver': '版本',
        'version': '版本',
        'ip': 'IP地址',
        'pwd': '密码',
        'key': '键',
        'value': '值',
        'v': '值',
        'section': '节',
        'cmd': '命令',
        'command': '命令',
        'current_dir': '当前目录',
        'url': 'URL地址',
        'save_file': '保存文件',
        'pic': '图片',
        'trans_color': '透明色',
        'objects': '对象',
        'prob': '概率阈值',
        'iou': 'IOU阈值',
        'base_addr': '基址',
        'fun_name': '函数名',
        'module_name': '模块名',
        'asm_code': '汇编代码',
        'asm_ins': '汇编指令',
        'is_64bit': '是否64位',
        'float_value': '浮点值',
        'double_value': '双精度值',
        'int_value': '整数值',
        'string_value': '字符串值',
        'float_value_min': '最小浮点值',
        'float_value_max': '最大浮点值',
        'double_value_min': '最小双精度值',
        'double_value_max': '最大双精度值',
        'int_value_min': '最小整数值',
        'int_value_max': '最大整数值',
        'addr_range': '地址范围',
        'process_name': '进程名',
        'process_id': '进程ID',
        'pid': '进程ID',
        'class_name': '类名',
        'title_name': '标题名',
        'title': '标题',
        'parent': '父窗口',
        'spec1': '条件1',
        'spec2': '条件2',
        'flag1': '标志1',
        'flag2': '标志2',
        'type1': '类型1',
        'type2': '类型2',
        'display': '显示模式',
        'mouse': '鼠标模式',
        'keypad': '键盘模式',
        'public_desc': '公共描述',
        'all_pos': '所有位置',
        'msg': '消息',
        'info': '信息',
        'param': '参数',
        'sub_cmd': '子命令',
        'rate': '速率',
        'fre': '频率',
        'line_height': '行高',
        'col_gap': '列间距',
        'row_gap': '行间距',
        'word_gap': '字间距',
        'word': '文字',
        'dict_name': '字库名',
        'dict_info': '字库信息',
        'font_index': '字体索引',
        'exact_ocr': '精确OCR',
        'key_str': '按键字符串',
        'vk': '虚拟键码',
        'key_code': '键码',
        'rx': '相对X',
        'ry': '相对Y',
        'speed': '速度',
        'lock': '锁定标志',
        'limit': '限制',
        't': '时间',
        'level': '级别',
        'uac': 'UAC标志',
        'rate': '比率',
        'over': '覆盖标志',
        'src_file': '源文件',
        'dst_file': '目标文件',
        'start_pos': '起始位置',
        'end_pos': '结束位置',
        'pmbi': '内存信息缓冲区',
        'old_protect': '旧保护属性',
        'mousedelay': '鼠标延迟',
        'mousestep': '鼠标步长',
        'id': 'ID',
    }
    
    # 处理指针类型
    if param_type.startswith('*'):
        desc = param_desc.get(param_name, param_name)
        return f'{desc}(输出参数)'
    
    return param_desc.get(param_name, param_name)

def generate_comment(func_name, params, return_type):
    """生成函数注释"""
    # 获取函数描述
    desc = FUNC_COMMENTS.get(func_name, func_name)
    
    # 构建注释
    lines = []
    lines.append(f'// {func_name} {desc}')
    
    # 添加参数说明
    if params:
        for param_type, param_name in params:
            param_desc = get_param_comment(param_name, param_type)
            lines.append(f'// 参数: {param_name} - {param_desc}')
    
    # 添加返回值说明
    if return_type:
        if return_type == 'int32':
            lines.append('// 返回: 成功返回1,失败返回0')
        elif return_type == 'string':
            lines.append('// 返回: 结果字符串')
        elif return_type == 'int64':
            lines.append('// 返回: 64位整数值')
        elif return_type == 'float32':
            lines.append('// 返回: 32位浮点数')
        elif return_type == 'float64':
            lines.append('// 返回: 64位浮点数')
        elif return_type == 'bool':
            lines.append('// 返回: 布尔值')
    
    return '\n'.join(lines)

def add_comments(file_path):
    """为Go文件添加注释"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 匹配函数定义
    pattern = r'(// [^\n]+\n)?(func \(\*?dm \*DmSoftImpl\) (\w+)\(([^)]*)\)(\s*\([^)]*\)|\s*\w+)?)'
    
    def replace_func(match):
        existing_comment = match.group(1)
        func_def = match.group(2)
        func_name = match.group(3)
        params_str = match.group(4)
        return_type = match.group(5).strip() if match.group(5) else ''
        
        # 如果已有注释，保留
        if existing_comment:
            return match.group(0)
        
        # 解析参数
        params = []
        if params_str.strip():
            for param in params_str.split(','):
                param = param.strip()
                if not param:
                    continue
                words = param.split()
                if len(words) >= 2:
                    param_type = words[-1]
                    param_names = words[:-1]
                    for name in param_names:
                        if name and name != '*':
                            params.append((param_type, name))
        
        # 生成注释
        comment = generate_comment(func_name, params, return_type)
        
        return comment + '\n' + func_def
    
    content = re.sub(pattern, replace_func, content)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("注释添加完成")

if __name__ == '__main__':
    file_path = r'E:\SRC\dm\72424_C++\go\dmsoft_impl.go'
    add_comments(file_path)
