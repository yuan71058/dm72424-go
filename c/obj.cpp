
#include "obj.h"
#include <windows.h>

static long g_dm_hmodule = 0;

DWORD LoadDm(PCSTR path)
{
    if (g_dm_hmodule) return TRUE;

    g_dm_hmodule = (long)(ULONG_PTR)LoadLibraryA(path);
    if (g_dm_hmodule) return g_dm_hmodule;
    return FALSE;
}

BOOL FreeDm(void)
{
    if (g_dm_hmodule == 0) return TRUE;

    FreeLibrary((HMODULE)(ULONG_PTR)g_dm_hmodule);
    g_dm_hmodule = 0;
    return TRUE;
}

dmsoft::dmsoft()
{
    typedef long (WINAPI *TypeCreateObj)(void);

    TypeCreateObj CreateObj = (TypeCreateObj)(ULONG_PTR)(g_dm_hmodule + 98304);
    obj = CreateObj();
}

dmsoft::~dmsoft()
{
    typedef long (WINAPI * TypeReleaseObj)(long);

    TypeReleaseObj ReleaseObj = (TypeReleaseObj)(ULONG_PTR)(g_dm_hmodule + 98400);
    ReleaseObj(obj);
}

CString dmsoft::GetDiskReversion(long index)
{
    typedef PCSTR (WINAPI * TypeGetDiskReversion)(long,long);

    TypeGetDiskReversion fun = (TypeGetDiskReversion)(ULONG_PTR)(g_dm_hmodule + 109040);
    return fun(obj,index);
}

long dmsoft::LoadAiMemory(long addr,long size)
{
    typedef long (WINAPI * TypeLoadAiMemory)(long,long,long);

    TypeLoadAiMemory fun = (TypeLoadAiMemory)(ULONG_PTR)(g_dm_hmodule + 108256);
    return fun(obj,addr,size);
}

CString dmsoft::FaqSend(PCSTR server,long handle,long request_type,long time_out)
{
    typedef PCSTR (WINAPI * TypeFaqSend)(long,PCSTR,long,long,long);

    TypeFaqSend fun = (TypeFaqSend)(ULONG_PTR)(g_dm_hmodule + 114016);
    return fun(obj,server,handle,request_type,time_out);
}

long dmsoft::FindPicSimMem(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,long sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindPicSimMem)(long,long,long,long,long,PCSTR,PCSTR,long,long,long*,long*);

    TypeFindPicSimMem fun = (TypeFindPicSimMem)(ULONG_PTR)(g_dm_hmodule + 121744);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir,x,y);
}

CString dmsoft::Ver()
{
    typedef PCSTR (WINAPI * TypeVer)(long);

    TypeVer fun = (TypeVer)(ULONG_PTR)(g_dm_hmodule + 100320);
    return fun(obj);
}

long dmsoft::SetPath(PCSTR path)
{
    typedef long (WINAPI * TypeSetPath)(long,PCSTR);

    TypeSetPath fun = (TypeSetPath)(ULONG_PTR)(g_dm_hmodule + 123808);
    return fun(obj,path);
}

long dmsoft::SetShowAsmErrorMsg(long show)
{
    typedef long (WINAPI * TypeSetShowAsmErrorMsg)(long,long);

    TypeSetShowAsmErrorMsg fun = (TypeSetShowAsmErrorMsg)(ULONG_PTR)(g_dm_hmodule + 101392);
    return fun(obj,show);
}

CString dmsoft::FindStrS(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,long* x,long* y)
{
    typedef PCSTR (WINAPI * TypeFindStrS)(long,long,long,long,long,PCSTR,PCSTR,double,long*,long*);

    TypeFindStrS fun = (TypeFindStrS)(ULONG_PTR)(g_dm_hmodule + 116832);
    return fun(obj,x1,y1,x2,y2,str,color,sim,x,y);
}

CString dmsoft::GetWordsNoDict(long x1,long y1,long x2,long y2,PCSTR color)
{
    typedef PCSTR (WINAPI * TypeGetWordsNoDict)(long,long,long,long,long,PCSTR);

    TypeGetWordsNoDict fun = (TypeGetWordsNoDict)(ULONG_PTR)(g_dm_hmodule + 99024);
    return fun(obj,x1,y1,x2,y2,color);
}

long dmsoft::GetOsBuildNumber()
{
    typedef long (WINAPI * TypeGetOsBuildNumber)(long);

    TypeGetOsBuildNumber fun = (TypeGetOsBuildNumber)(ULONG_PTR)(g_dm_hmodule + 104240);
    return fun(obj);
}

long dmsoft::GetID()
{
    typedef long (WINAPI * TypeGetID)(long);

    TypeGetID fun = (TypeGetID)(ULONG_PTR)(g_dm_hmodule + 105184);
    return fun(obj);
}

long dmsoft::SetMouseSpeed(long speed)
{
    typedef long (WINAPI * TypeSetMouseSpeed)(long,long);

    TypeSetMouseSpeed fun = (TypeSetMouseSpeed)(ULONG_PTR)(g_dm_hmodule + 124800);
    return fun(obj,speed);
}

CString dmsoft::FindData(long hwnd,PCSTR addr_range,PCSTR data)
{
    typedef PCSTR (WINAPI * TypeFindData)(long,long,PCSTR,PCSTR);

    TypeFindData fun = (TypeFindData)(ULONG_PTR)(g_dm_hmodule + 109760);
    return fun(obj,hwnd,addr_range,data);
}

long dmsoft::SendPaste(long hwnd)
{
    typedef long (WINAPI * TypeSendPaste)(long,long);

    TypeSendPaste fun = (TypeSendPaste)(ULONG_PTR)(g_dm_hmodule + 122944);
    return fun(obj,hwnd);
}

CString dmsoft::GetColor(long x,long y)
{
    typedef PCSTR (WINAPI * TypeGetColor)(long,long,long);

    TypeGetColor fun = (TypeGetColor)(ULONG_PTR)(g_dm_hmodule + 117424);
    return fun(obj,x,y);
}

long dmsoft::LoadPicByte(long addr,long size,PCSTR name)
{
    typedef long (WINAPI * TypeLoadPicByte)(long,long,long,PCSTR);

    TypeLoadPicByte fun = (TypeLoadPicByte)(ULONG_PTR)(g_dm_hmodule + 121408);
    return fun(obj,addr,size,name);
}

long dmsoft::WriteFloatAddr(long hwnd,LONGLONG addr,float float_value)
{
    typedef long (WINAPI * TypeWriteFloatAddr)(long,long,LONGLONG,float);

    TypeWriteFloatAddr fun = (TypeWriteFloatAddr)(ULONG_PTR)(g_dm_hmodule + 117312);
    return fun(obj,hwnd,addr,float_value);
}

long dmsoft::SetWordLineHeight(long line_height)
{
    typedef long (WINAPI * TypeSetWordLineHeight)(long,long);

    TypeSetWordLineHeight fun = (TypeSetWordLineHeight)(ULONG_PTR)(g_dm_hmodule + 101296);
    return fun(obj,line_height);
}

LONGLONG dmsoft::AsmCall(long hwnd,long mode)
{
    typedef LONGLONG (WINAPI * TypeAsmCall)(long,long,long);

    TypeAsmCall fun = (TypeAsmCall)(ULONG_PTR)(g_dm_hmodule + 114656);
    return fun(obj,hwnd,mode);
}

long dmsoft::FindColorBlock(long x1,long y1,long x2,long y2,PCSTR color,double sim,long count,long width,long height,long* x,long* y)
{
    typedef long (WINAPI * TypeFindColorBlock)(long,long,long,long,long,PCSTR,double,long,long,long,long*,long*);

    TypeFindColorBlock fun = (TypeFindColorBlock)(ULONG_PTR)(g_dm_hmodule + 113568);
    return fun(obj,x1,y1,x2,y2,color,sim,count,width,height,x,y);
}

CString dmsoft::DisAssemble(PCSTR asm_code,LONGLONG base_addr,long is_64bit)
{
    typedef PCSTR (WINAPI * TypeDisAssemble)(long,PCSTR,LONGLONG,long);

    TypeDisAssemble fun = (TypeDisAssemble)(ULONG_PTR)(g_dm_hmodule + 112656);
    return fun(obj,asm_code,base_addr,is_64bit);
}

long dmsoft::RegEx(PCSTR code,PCSTR ver,PCSTR ip)
{
    typedef long (WINAPI * TypeRegEx)(long,PCSTR,PCSTR,PCSTR);

    TypeRegEx fun = (TypeRegEx)(ULONG_PTR)(g_dm_hmodule + 98864);
    return fun(obj,code,ver,ip);
}

long dmsoft::EncodeFile(PCSTR file,PCSTR pwd)
{
    typedef long (WINAPI * TypeEncodeFile)(long,PCSTR,PCSTR);

    TypeEncodeFile fun = (TypeEncodeFile)(ULONG_PTR)(g_dm_hmodule + 106528);
    return fun(obj,file,pwd);
}

long dmsoft::WriteString(long hwnd,PCSTR addr,long type,PCSTR v)
{
    typedef long (WINAPI * TypeWriteString)(long,long,PCSTR,long,PCSTR);

    TypeWriteString fun = (TypeWriteString)(ULONG_PTR)(g_dm_hmodule + 115936);
    return fun(obj,hwnd,addr,type,v);
}

CString dmsoft::FindStrFastEx(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrFastEx)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrFastEx fun = (TypeFindStrFastEx)(ULONG_PTR)(g_dm_hmodule + 122000);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

LONGLONG dmsoft::AsmCallEx(long hwnd,long mode,PCSTR base_addr)
{
    typedef LONGLONG (WINAPI * TypeAsmCallEx)(long,long,long,PCSTR);

    TypeAsmCallEx fun = (TypeAsmCallEx)(ULONG_PTR)(g_dm_hmodule + 99632);
    return fun(obj,hwnd,mode,base_addr);
}

CString dmsoft::FindDoubleEx(long hwnd,PCSTR addr_range,double double_value_min,double double_value_max,long step,long multi_thread,long mode)
{
    typedef PCSTR (WINAPI * TypeFindDoubleEx)(long,long,PCSTR,double,double,long,long,long);

    TypeFindDoubleEx fun = (TypeFindDoubleEx)(ULONG_PTR)(g_dm_hmodule + 110416);
    return fun(obj,hwnd,addr_range,double_value_min,double_value_max,step,multi_thread,mode);
}

long dmsoft::SetFindPicMultithreadLimit(long limit)
{
    typedef long (WINAPI * TypeSetFindPicMultithreadLimit)(long,long);

    TypeSetFindPicMultithreadLimit fun = (TypeSetFindPicMultithreadLimit)(ULONG_PTR)(g_dm_hmodule + 107616);
    return fun(obj,limit);
}

long dmsoft::SendString2(long hwnd,PCSTR str)
{
    typedef long (WINAPI * TypeSendString2)(long,long,PCSTR);

    TypeSendString2 fun = (TypeSendString2)(ULONG_PTR)(g_dm_hmodule + 99888);
    return fun(obj,hwnd,str);
}

long dmsoft::DownCpu(long type,long rate)
{
    typedef long (WINAPI * TypeDownCpu)(long,long,long);

    TypeDownCpu fun = (TypeDownCpu)(ULONG_PTR)(g_dm_hmodule + 112960);
    return fun(obj,type,rate);
}

long dmsoft::DmGuard(long enable,PCSTR type)
{
    typedef long (WINAPI * TypeDmGuard)(long,long,PCSTR);

    TypeDmGuard fun = (TypeDmGuard)(ULONG_PTR)(g_dm_hmodule + 103552);
    return fun(obj,enable,type);
}

long dmsoft::SpeedNormalGraphic(long en)
{
    typedef long (WINAPI * TypeSpeedNormalGraphic)(long,long);

    TypeSpeedNormalGraphic fun = (TypeSpeedNormalGraphic)(ULONG_PTR)(g_dm_hmodule + 101184);
    return fun(obj,en);
}

long dmsoft::FindPicSim(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,long sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindPicSim)(long,long,long,long,long,PCSTR,PCSTR,long,long,long*,long*);

    TypeFindPicSim fun = (TypeFindPicSim)(ULONG_PTR)(g_dm_hmodule + 98768);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir,x,y);
}

long dmsoft::WriteInt(long hwnd,PCSTR addr,long type,LONGLONG v)
{
    typedef long (WINAPI * TypeWriteInt)(long,long,PCSTR,long,LONGLONG);

    TypeWriteInt fun = (TypeWriteInt)(ULONG_PTR)(g_dm_hmodule + 112416);
    return fun(obj,hwnd,addr,type,v);
}

long dmsoft::SetMemoryHwndAsProcessId(long en)
{
    typedef long (WINAPI * TypeSetMemoryHwndAsProcessId)(long,long);

    TypeSetMemoryHwndAsProcessId fun = (TypeSetMemoryHwndAsProcessId)(ULONG_PTR)(g_dm_hmodule + 107984);
    return fun(obj,en);
}

long dmsoft::WriteDataFromBin(long hwnd,PCSTR addr,long data,long len)
{
    typedef long (WINAPI * TypeWriteDataFromBin)(long,long,PCSTR,long,long);

    TypeWriteDataFromBin fun = (TypeWriteDataFromBin)(ULONG_PTR)(g_dm_hmodule + 118304);
    return fun(obj,hwnd,addr,data,len);
}

long dmsoft::SetMinColGap(long col_gap)
{
    typedef long (WINAPI * TypeSetMinColGap)(long,long);

    TypeSetMinColGap fun = (TypeSetMinColGap)(ULONG_PTR)(g_dm_hmodule + 110560);
    return fun(obj,col_gap);
}

long dmsoft::KeyPressStr(PCSTR key_str,long delay)
{
    typedef long (WINAPI * TypeKeyPressStr)(long,PCSTR,long);

    TypeKeyPressStr fun = (TypeKeyPressStr)(ULONG_PTR)(g_dm_hmodule + 102528);
    return fun(obj,key_str,delay);
}

long dmsoft::LockDisplay(long lock)
{
    typedef long (WINAPI * TypeLockDisplay)(long,long);

    TypeLockDisplay fun = (TypeLockDisplay)(ULONG_PTR)(g_dm_hmodule + 108304);
    return fun(obj,lock);
}

CString dmsoft::FindStrWithFontE(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,PCSTR font_name,long font_size,long flag)
{
    typedef PCSTR (WINAPI * TypeFindStrWithFontE)(long,long,long,long,long,PCSTR,PCSTR,double,PCSTR,long,long);

    TypeFindStrWithFontE fun = (TypeFindStrWithFontE)(ULONG_PTR)(g_dm_hmodule + 112544);
    return fun(obj,x1,y1,x2,y2,str,color,sim,font_name,font_size,flag);
}

CString dmsoft::EnumIniKey(PCSTR section,PCSTR file)
{
    typedef PCSTR (WINAPI * TypeEnumIniKey)(long,PCSTR,PCSTR);

    TypeEnumIniKey fun = (TypeEnumIniKey)(ULONG_PTR)(g_dm_hmodule + 108032);
    return fun(obj,section,file);
}

CString dmsoft::MatchPicName(PCSTR pic_name)
{
    typedef PCSTR (WINAPI * TypeMatchPicName)(long,PCSTR);

    TypeMatchPicName fun = (TypeMatchPicName)(ULONG_PTR)(g_dm_hmodule + 117984);
    return fun(obj,pic_name);
}

long dmsoft::EnableFakeActive(long en)
{
    typedef long (WINAPI * TypeEnableFakeActive)(long,long);

    TypeEnableFakeActive fun = (TypeEnableFakeActive)(ULONG_PTR)(g_dm_hmodule + 107888);
    return fun(obj,en);
}

long dmsoft::FaqGetSize(long handle)
{
    typedef long (WINAPI * TypeFaqGetSize)(long,long);

    TypeFaqGetSize fun = (TypeFaqGetSize)(ULONG_PTR)(g_dm_hmodule + 103456);
    return fun(obj,handle);
}

CString dmsoft::ExecuteCmd(PCSTR cmd,PCSTR current_dir,long time_out)
{
    typedef PCSTR (WINAPI * TypeExecuteCmd)(long,PCSTR,PCSTR,long);

    TypeExecuteCmd fun = (TypeExecuteCmd)(ULONG_PTR)(g_dm_hmodule + 116928);
    return fun(obj,cmd,current_dir,time_out);
}

long dmsoft::EnableRealKeypad(long en)
{
    typedef long (WINAPI * TypeEnableRealKeypad)(long,long);

    TypeEnableRealKeypad fun = (TypeEnableRealKeypad)(ULONG_PTR)(g_dm_hmodule + 105648);
    return fun(obj,en);
}

long dmsoft::SetDisplayRefreshDelay(long t)
{
    typedef long (WINAPI * TypeSetDisplayRefreshDelay)(long,long);

    TypeSetDisplayRefreshDelay fun = (TypeSetDisplayRefreshDelay)(ULONG_PTR)(g_dm_hmodule + 111344);
    return fun(obj,t);
}

long dmsoft::MiddleClick()
{
    typedef long (WINAPI * TypeMiddleClick)(long);

    TypeMiddleClick fun = (TypeMiddleClick)(ULONG_PTR)(g_dm_hmodule + 108560);
    return fun(obj);
}

CString dmsoft::AiYoloSortsObjects(PCSTR objects,long height)
{
    typedef PCSTR (WINAPI * TypeAiYoloSortsObjects)(long,PCSTR,long);

    TypeAiYoloSortsObjects fun = (TypeAiYoloSortsObjects)(ULONG_PTR)(g_dm_hmodule + 120480);
    return fun(obj,objects,height);
}

long dmsoft::WriteDataAddr(long hwnd,LONGLONG addr,PCSTR data)
{
    typedef long (WINAPI * TypeWriteDataAddr)(long,long,LONGLONG,PCSTR);

    TypeWriteDataAddr fun = (TypeWriteDataAddr)(ULONG_PTR)(g_dm_hmodule + 105744);
    return fun(obj,hwnd,addr,data);
}

CString dmsoft::RGB2BGR(PCSTR rgb_color)
{
    typedef PCSTR (WINAPI * TypeRGB2BGR)(long,PCSTR);

    TypeRGB2BGR fun = (TypeRGB2BGR)(ULONG_PTR)(g_dm_hmodule + 115744);
    return fun(obj,rgb_color);
}

long dmsoft::DisablePowerSave()
{
    typedef long (WINAPI * TypeDisablePowerSave)(long);

    TypeDisablePowerSave fun = (TypeDisablePowerSave)(ULONG_PTR)(g_dm_hmodule + 121952);
    return fun(obj);
}

long dmsoft::GetClientSize(long hwnd,long* width,long* height)
{
    typedef long (WINAPI * TypeGetClientSize)(long,long,long*,long*);

    TypeGetClientSize fun = (TypeGetClientSize)(ULONG_PTR)(g_dm_hmodule + 103344);
    return fun(obj,hwnd,width,height);
}

long dmsoft::EnableMouseMsg(long en)
{
    typedef long (WINAPI * TypeEnableMouseMsg)(long,long);

    TypeEnableMouseMsg fun = (TypeEnableMouseMsg)(ULONG_PTR)(g_dm_hmodule + 101344);
    return fun(obj,en);
}

long dmsoft::EnableKeypadMsg(long en)
{
    typedef long (WINAPI * TypeEnableKeypadMsg)(long,long);

    TypeEnableKeypadMsg fun = (TypeEnableKeypadMsg)(ULONG_PTR)(g_dm_hmodule + 120864);
    return fun(obj,en);
}

long dmsoft::GetFileLength(PCSTR file)
{
    typedef long (WINAPI * TypeGetFileLength)(long,PCSTR);

    TypeGetFileLength fun = (TypeGetFileLength)(ULONG_PTR)(g_dm_hmodule + 111296);
    return fun(obj,file);
}

LONGLONG dmsoft::GetRemoteApiAddress(long hwnd,LONGLONG base_addr,PCSTR fun_name)
{
    typedef LONGLONG (WINAPI * TypeGetRemoteApiAddress)(long,long,LONGLONG,PCSTR);

    TypeGetRemoteApiAddress fun = (TypeGetRemoteApiAddress)(ULONG_PTR)(g_dm_hmodule + 122192);
    return fun(obj,hwnd,base_addr,fun_name);
}

CString dmsoft::DmGuardParams(PCSTR cmd,PCSTR sub_cmd,PCSTR param)
{
    typedef PCSTR (WINAPI * TypeDmGuardParams)(long,PCSTR,PCSTR,PCSTR);

    TypeDmGuardParams fun = (TypeDmGuardParams)(ULONG_PTR)(g_dm_hmodule + 105472);
    return fun(obj,cmd,sub_cmd,param);
}

long dmsoft::DownloadFile(PCSTR url,PCSTR save_file,long timeout)
{
    typedef long (WINAPI * TypeDownloadFile)(long,PCSTR,PCSTR,long);

    TypeDownloadFile fun = (TypeDownloadFile)(ULONG_PTR)(g_dm_hmodule + 123648);
    return fun(obj,url,save_file,timeout);
}

long dmsoft::WriteDoubleAddr(long hwnd,LONGLONG addr,double double_value)
{
    typedef long (WINAPI * TypeWriteDoubleAddr)(long,long,LONGLONG,double);

    TypeWriteDoubleAddr fun = (TypeWriteDoubleAddr)(ULONG_PTR)(g_dm_hmodule + 115232);
    return fun(obj,hwnd,addr,double_value);
}

long dmsoft::EnableIme(long en)
{
    typedef long (WINAPI * TypeEnableIme)(long,long);

    TypeEnableIme fun = (TypeEnableIme)(ULONG_PTR)(g_dm_hmodule + 120192);
    return fun(obj,en);
}

long dmsoft::TerminateProcessTree(long pid)
{
    typedef long (WINAPI * TypeTerminateProcessTree)(long,long);

    TypeTerminateProcessTree fun = (TypeTerminateProcessTree)(ULONG_PTR)(g_dm_hmodule + 114240);
    return fun(obj,pid);
}

long dmsoft::FoobarClose(long hwnd)
{
    typedef long (WINAPI * TypeFoobarClose)(long,long);

    TypeFoobarClose fun = (TypeFoobarClose)(ULONG_PTR)(g_dm_hmodule + 102480);
    return fun(obj,hwnd);
}

CString dmsoft::FindNearestPos(PCSTR all_pos,long type,long x,long y)
{
    typedef PCSTR (WINAPI * TypeFindNearestPos)(long,PCSTR,long,long,long);

    TypeFindNearestPos fun = (TypeFindNearestPos)(ULONG_PTR)(g_dm_hmodule + 112480);
    return fun(obj,all_pos,type,x,y);
}

long dmsoft::CreateFoobarRect(long hwnd,long x,long y,long w,long h)
{
    typedef long (WINAPI * TypeCreateFoobarRect)(long,long,long,long,long,long);

    TypeCreateFoobarRect fun = (TypeCreateFoobarRect)(ULONG_PTR)(g_dm_hmodule + 119072);
    return fun(obj,hwnd,x,y,w,h);
}

long dmsoft::GetCursorPos(long* x,long* y)
{
    typedef long (WINAPI * TypeGetCursorPos)(long,long*,long*);

    TypeGetCursorPos fun = (TypeGetCursorPos)(ULONG_PTR)(g_dm_hmodule + 121680);
    return fun(obj,x,y);
}

CString dmsoft::FindColorBlockEx(long x1,long y1,long x2,long y2,PCSTR color,double sim,long count,long width,long height)
{
    typedef PCSTR (WINAPI * TypeFindColorBlockEx)(long,long,long,long,long,PCSTR,double,long,long,long);

    TypeFindColorBlockEx fun = (TypeFindColorBlockEx)(ULONG_PTR)(g_dm_hmodule + 103840);
    return fun(obj,x1,y1,x2,y2,color,sim,count,width,height);
}

CString dmsoft::FindFloat(long hwnd,PCSTR addr_range,float float_value_min,float float_value_max)
{
    typedef PCSTR (WINAPI * TypeFindFloat)(long,long,PCSTR,float,float);

    TypeFindFloat fun = (TypeFindFloat)(ULONG_PTR)(g_dm_hmodule + 103216);
    return fun(obj,hwnd,addr_range,float_value_min,float_value_max);
}

CString dmsoft::GetProcessInfo(long pid)
{
    typedef PCSTR (WINAPI * TypeGetProcessInfo)(long,long);

    TypeGetProcessInfo fun = (TypeGetProcessInfo)(ULONG_PTR)(g_dm_hmodule + 119024);
    return fun(obj,pid);
}

CString dmsoft::ReadFile(PCSTR file)
{
    typedef PCSTR (WINAPI * TypeReadFile)(long,PCSTR);

    TypeReadFile fun = (TypeReadFile)(ULONG_PTR)(g_dm_hmodule + 114464);
    return fun(obj,file);
}

CString dmsoft::FindShapeEx(long x1,long y1,long x2,long y2,PCSTR offset_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindShapeEx)(long,long,long,long,long,PCSTR,double,long);

    TypeFindShapeEx fun = (TypeFindShapeEx)(ULONG_PTR)(g_dm_hmodule + 99792);
    return fun(obj,x1,y1,x2,y2,offset_color,sim,dir);
}

long dmsoft::SetWindowText(long hwnd,PCSTR text)
{
    typedef long (WINAPI * TypeSetWindowText)(long,long,PCSTR);

    TypeSetWindowText fun = (TypeSetWindowText)(ULONG_PTR)(g_dm_hmodule + 113008);
    return fun(obj,hwnd,text);
}

long dmsoft::ForceUnBindWindow(long hwnd)
{
    typedef long (WINAPI * TypeForceUnBindWindow)(long,long);

    TypeForceUnBindWindow fun = (TypeForceUnBindWindow)(ULONG_PTR)(g_dm_hmodule + 120144);
    return fun(obj,hwnd);
}

LONGLONG dmsoft::ReadIntAddr(long hwnd,LONGLONG addr,long type)
{
    typedef LONGLONG (WINAPI * TypeReadIntAddr)(long,long,LONGLONG,long);

    TypeReadIntAddr fun = (TypeReadIntAddr)(ULONG_PTR)(g_dm_hmodule + 99712);
    return fun(obj,hwnd,addr,type);
}

long dmsoft::FindShape(long x1,long y1,long x2,long y2,PCSTR offset_color,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindShape)(long,long,long,long,long,PCSTR,double,long,long*,long*);

    TypeFindShape fun = (TypeFindShape)(ULONG_PTR)(g_dm_hmodule + 123856);
    return fun(obj,x1,y1,x2,y2,offset_color,sim,dir,x,y);
}

CString dmsoft::GetRealPath(PCSTR path)
{
    typedef PCSTR (WINAPI * TypeGetRealPath)(long,PCSTR);

    TypeGetRealPath fun = (TypeGetRealPath)(ULONG_PTR)(g_dm_hmodule + 105008);
    return fun(obj,path);
}

long dmsoft::EnableSpeedDx(long en)
{
    typedef long (WINAPI * TypeEnableSpeedDx)(long,long);

    TypeEnableSpeedDx fun = (TypeEnableSpeedDx)(ULONG_PTR)(g_dm_hmodule + 115472);
    return fun(obj,en);
}

long dmsoft::UnLoadDriver()
{
    typedef long (WINAPI * TypeUnLoadDriver)(long);

    TypeUnLoadDriver fun = (TypeUnLoadDriver)(ULONG_PTR)(g_dm_hmodule + 105696);
    return fun(obj);
}

long dmsoft::GetMemoryUsage()
{
    typedef long (WINAPI * TypeGetMemoryUsage)(long);

    TypeGetMemoryUsage fun = (TypeGetMemoryUsage)(ULONG_PTR)(g_dm_hmodule + 106064);
    return fun(obj);
}

long dmsoft::MiddleDown()
{
    typedef long (WINAPI * TypeMiddleDown)(long);

    TypeMiddleDown fun = (TypeMiddleDown)(ULONG_PTR)(g_dm_hmodule + 109872);
    return fun(obj);
}

CString dmsoft::EnumIniSection(PCSTR file)
{
    typedef PCSTR (WINAPI * TypeEnumIniSection)(long,PCSTR);

    TypeEnumIniSection fun = (TypeEnumIniSection)(ULONG_PTR)(g_dm_hmodule + 117184);
    return fun(obj,file);
}

long dmsoft::CheckUAC()
{
    typedef long (WINAPI * TypeCheckUAC)(long);

    TypeCheckUAC fun = (TypeCheckUAC)(ULONG_PTR)(g_dm_hmodule + 123104);
    return fun(obj);
}

long dmsoft::OpenProcess(long pid)
{
    typedef long (WINAPI * TypeOpenProcess)(long,long);

    TypeOpenProcess fun = (TypeOpenProcess)(ULONG_PTR)(g_dm_hmodule + 124624);
    return fun(obj,pid);
}

long dmsoft::IsDisplayDead(long x1,long y1,long x2,long y2,long t)
{
    typedef long (WINAPI * TypeIsDisplayDead)(long,long,long,long,long,long);

    TypeIsDisplayDead fun = (TypeIsDisplayDead)(ULONG_PTR)(g_dm_hmodule + 114896);
    return fun(obj,x1,y1,x2,y2,t);
}

long dmsoft::WriteIniPwd(PCSTR section,PCSTR key,PCSTR v,PCSTR file,PCSTR pwd)
{
    typedef long (WINAPI * TypeWriteIniPwd)(long,PCSTR,PCSTR,PCSTR,PCSTR,PCSTR);

    TypeWriteIniPwd fun = (TypeWriteIniPwd)(ULONG_PTR)(g_dm_hmodule + 115872);
    return fun(obj,section,key,v,file,pwd);
}

CString dmsoft::GetNetTime()
{
    typedef PCSTR (WINAPI * TypeGetNetTime)(long);

    TypeGetNetTime fun = (TypeGetNetTime)(ULONG_PTR)(g_dm_hmodule + 107712);
    return fun(obj);
}

float dmsoft::ReadFloat(long hwnd,PCSTR addr)
{
    typedef float (WINAPI * TypeReadFloat)(long,long,PCSTR);

    TypeReadFloat fun = (TypeReadFloat)(ULONG_PTR)(g_dm_hmodule + 100976);
    return fun(obj,hwnd,addr);
}

long dmsoft::DisableCloseDisplayAndSleep()
{
    typedef long (WINAPI * TypeDisableCloseDisplayAndSleep)(long);

    TypeDisableCloseDisplayAndSleep fun = (TypeDisableCloseDisplayAndSleep)(ULONG_PTR)(g_dm_hmodule + 114416);
    return fun(obj);
}

CString dmsoft::GetWindowTitle(long hwnd)
{
    typedef PCSTR (WINAPI * TypeGetWindowTitle)(long,long);

    TypeGetWindowTitle fun = (TypeGetWindowTitle)(ULONG_PTR)(g_dm_hmodule + 110816);
    return fun(obj,hwnd);
}

CString dmsoft::Assemble(LONGLONG base_addr,long is_64bit)
{
    typedef PCSTR (WINAPI * TypeAssemble)(long,LONGLONG,long);

    TypeAssemble fun = (TypeAssemble)(ULONG_PTR)(g_dm_hmodule + 119584);
    return fun(obj,base_addr,is_64bit);
}

long dmsoft::GetMousePointWindow()
{
    typedef long (WINAPI * TypeGetMousePointWindow)(long);

    TypeGetMousePointWindow fun = (TypeGetMousePointWindow)(ULONG_PTR)(g_dm_hmodule + 105424);
    return fun(obj);
}

long dmsoft::SetExportDict(long index,PCSTR dict_name)
{
    typedef long (WINAPI * TypeSetExportDict)(long,long,PCSTR);

    TypeSetExportDict fun = (TypeSetExportDict)(ULONG_PTR)(g_dm_hmodule + 119392);
    return fun(obj,index,dict_name);
}

long dmsoft::Delay(long mis)
{
    typedef long (WINAPI * TypeDelay)(long,long);

    TypeDelay fun = (TypeDelay)(ULONG_PTR)(g_dm_hmodule + 106480);
    return fun(obj,mis);
}

long dmsoft::Reg(PCSTR code,PCSTR ver)
{
    typedef long (WINAPI * TypeReg)(long,PCSTR,PCSTR);

    TypeReg fun = (TypeReg)(ULONG_PTR)(g_dm_hmodule + 121344);
    return fun(obj,code,ver);
}

long dmsoft::FoobarStopGif(long hwnd,long x,long y,PCSTR pic_name)
{
    typedef long (WINAPI * TypeFoobarStopGif)(long,long,long,long,PCSTR);

    TypeFoobarStopGif fun = (TypeFoobarStopGif)(ULONG_PTR)(g_dm_hmodule + 108096);
    return fun(obj,hwnd,x,y,pic_name);
}

CString dmsoft::ReadFileData(PCSTR file,long start_pos,long end_pos)
{
    typedef PCSTR (WINAPI * TypeReadFileData)(long,PCSTR,long,long);

    TypeReadFileData fun = (TypeReadFileData)(ULONG_PTR)(g_dm_hmodule + 115808);
    return fun(obj,file,start_pos,end_pos);
}

CString dmsoft::FindPicSimEx(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,long sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicSimEx)(long,long,long,long,long,PCSTR,PCSTR,long,long);

    TypeFindPicSimEx fun = (TypeFindPicSimEx)(ULONG_PTR)(g_dm_hmodule + 113728);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir);
}

long dmsoft::Capture(long x1,long y1,long x2,long y2,PCSTR file)
{
    typedef long (WINAPI * TypeCapture)(long,long,long,long,long,PCSTR);

    TypeCapture fun = (TypeCapture)(ULONG_PTR)(g_dm_hmodule + 119456);
    return fun(obj,x1,y1,x2,y2,file);
}

long dmsoft::GetScreenWidth()
{
    typedef long (WINAPI * TypeGetScreenWidth)(long);

    TypeGetScreenWidth fun = (TypeGetScreenWidth)(ULONG_PTR)(g_dm_hmodule + 113920);
    return fun(obj);
}

CString dmsoft::FindStrWithFontEx(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,PCSTR font_name,long font_size,long flag)
{
    typedef PCSTR (WINAPI * TypeFindStrWithFontEx)(long,long,long,long,long,PCSTR,PCSTR,double,PCSTR,long,long);

    TypeFindStrWithFontEx fun = (TypeFindStrWithFontEx)(ULONG_PTR)(g_dm_hmodule + 118848);
    return fun(obj,x1,y1,x2,y2,str,color,sim,font_name,font_size,flag);
}

long dmsoft::SetLocale()
{
    typedef long (WINAPI * TypeSetLocale)(long);

    TypeSetLocale fun = (TypeSetLocale)(ULONG_PTR)(g_dm_hmodule + 100928);
    return fun(obj);
}

long dmsoft::AsmAdd(PCSTR asm_ins)
{
    typedef long (WINAPI * TypeAsmAdd)(long,PCSTR);

    TypeAsmAdd fun = (TypeAsmAdd)(ULONG_PTR)(g_dm_hmodule + 121232);
    return fun(obj,asm_ins);
}

long dmsoft::GetScreenHeight()
{
    typedef long (WINAPI * TypeGetScreenHeight)(long);

    TypeGetScreenHeight fun = (TypeGetScreenHeight)(ULONG_PTR)(g_dm_hmodule + 117792);
    return fun(obj);
}

long dmsoft::CaptureGif(long x1,long y1,long x2,long y2,PCSTR file,long delay,long time)
{
    typedef long (WINAPI * TypeCaptureGif)(long,long,long,long,long,PCSTR,long,long);

    TypeCaptureGif fun = (TypeCaptureGif)(ULONG_PTR)(g_dm_hmodule + 120912);
    return fun(obj,x1,y1,x2,y2,file,delay,time);
}

long dmsoft::ReadDataAddrToBin(long hwnd,LONGLONG addr,long len)
{
    typedef long (WINAPI * TypeReadDataAddrToBin)(long,long,LONGLONG,long);

    TypeReadDataAddrToBin fun = (TypeReadDataAddrToBin)(ULONG_PTR)(g_dm_hmodule + 111792);
    return fun(obj,hwnd,addr,len);
}

long dmsoft::ReadDataToBin(long hwnd,PCSTR addr,long len)
{
    typedef long (WINAPI * TypeReadDataToBin)(long,long,PCSTR,long);

    TypeReadDataToBin fun = (TypeReadDataToBin)(ULONG_PTR)(g_dm_hmodule + 104480);
    return fun(obj,hwnd,addr,len);
}

CString dmsoft::FindPicS(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,double sim,long dir,long* x,long* y)
{
    typedef PCSTR (WINAPI * TypeFindPicS)(long,long,long,long,long,PCSTR,PCSTR,double,long,long*,long*);

    TypeFindPicS fun = (TypeFindPicS)(ULONG_PTR)(g_dm_hmodule + 101952);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir,x,y);
}

long dmsoft::FindPic(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindPic)(long,long,long,long,long,PCSTR,PCSTR,double,long,long*,long*);

    TypeFindPic fun = (TypeFindPic)(ULONG_PTR)(g_dm_hmodule + 104032);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir,x,y);
}

long dmsoft::FindMultiColor(long x1,long y1,long x2,long y2,PCSTR first_color,PCSTR offset_color,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindMultiColor)(long,long,long,long,long,PCSTR,PCSTR,double,long,long*,long*);

    TypeFindMultiColor fun = (TypeFindMultiColor)(ULONG_PTR)(g_dm_hmodule + 109360);
    return fun(obj,x1,y1,x2,y2,first_color,offset_color,sim,dir,x,y);
}

long dmsoft::HackSpeed(double rate)
{
    typedef long (WINAPI * TypeHackSpeed)(long,double);

    TypeHackSpeed fun = (TypeHackSpeed)(ULONG_PTR)(g_dm_hmodule + 104352);
    return fun(obj,rate);
}

CString dmsoft::FindPicE(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicE)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindPicE fun = (TypeFindPicE)(ULONG_PTR)(g_dm_hmodule + 114144);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir);
}

long dmsoft::MiddleUp()
{
    typedef long (WINAPI * TypeMiddleUp)(long);

    TypeMiddleUp fun = (TypeMiddleUp)(ULONG_PTR)(g_dm_hmodule + 115072);
    return fun(obj);
}

long dmsoft::GetWindow(long hwnd,long flag)
{
    typedef long (WINAPI * TypeGetWindow)(long,long,long);

    TypeGetWindow fun = (TypeGetWindow)(ULONG_PTR)(g_dm_hmodule + 120752);
    return fun(obj,hwnd,flag);
}

long dmsoft::SetUAC(long uac)
{
    typedef long (WINAPI * TypeSetUAC)(long,long);

    TypeSetUAC fun = (TypeSetUAC)(ULONG_PTR)(g_dm_hmodule + 108608);
    return fun(obj,uac);
}

long dmsoft::FoobarSetSave(long hwnd,PCSTR file,long en,PCSTR header)
{
    typedef long (WINAPI * TypeFoobarSetSave)(long,long,PCSTR,long,PCSTR);

    TypeFoobarSetSave fun = (TypeFoobarSetSave)(ULONG_PTR)(g_dm_hmodule + 124736);
    return fun(obj,hwnd,file,en,header);
}

long dmsoft::WheelDown()
{
    typedef long (WINAPI * TypeWheelDown)(long);

    TypeWheelDown fun = (TypeWheelDown)(ULONG_PTR)(g_dm_hmodule + 112848);
    return fun(obj);
}

CString dmsoft::FloatToData(float float_value)
{
    typedef PCSTR (WINAPI * TypeFloatToData)(long,float);

    TypeFloatToData fun = (TypeFloatToData)(ULONG_PTR)(g_dm_hmodule + 100464);
    return fun(obj,float_value);
}

long dmsoft::EnableFindPicMultithread(long en)
{
    typedef long (WINAPI * TypeEnableFindPicMultithread)(long,long);

    TypeEnableFindPicMultithread fun = (TypeEnableFindPicMultithread)(ULONG_PTR)(g_dm_hmodule + 118048);
    return fun(obj,en);
}

long dmsoft::DisableScreenSave()
{
    typedef long (WINAPI * TypeDisableScreenSave)(long);

    TypeDisableScreenSave fun = (TypeDisableScreenSave)(ULONG_PTR)(g_dm_hmodule + 112800);
    return fun(obj);
}

CString dmsoft::AiFindPicEx(long x1,long y1,long x2,long y2,PCSTR pic_name,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeAiFindPicEx)(long,long,long,long,long,PCSTR,double,long);

    TypeAiFindPicEx fun = (TypeAiFindPicEx)(ULONG_PTR)(g_dm_hmodule + 119136);
    return fun(obj,x1,y1,x2,y2,pic_name,sim,dir);
}

long dmsoft::SendString(long hwnd,PCSTR str)
{
    typedef long (WINAPI * TypeSendString)(long,long,PCSTR);

    TypeSendString fun = (TypeSendString)(ULONG_PTR)(g_dm_hmodule + 114832);
    return fun(obj,hwnd,str);
}

long dmsoft::EnterCri()
{
    typedef long (WINAPI * TypeEnterCri)(long);

    TypeEnterCri fun = (TypeEnterCri)(ULONG_PTR)(g_dm_hmodule + 116336);
    return fun(obj);
}

CString dmsoft::FindPicSimMemE(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,long sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicSimMemE)(long,long,long,long,long,PCSTR,PCSTR,long,long);

    TypeFindPicSimMemE fun = (TypeFindPicSimMemE)(ULONG_PTR)(g_dm_hmodule + 113296);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir);
}

long dmsoft::Delays(long min_s,long max_s)
{
    typedef long (WINAPI * TypeDelays)(long,long,long);

    TypeDelays fun = (TypeDelays)(ULONG_PTR)(g_dm_hmodule + 123328);
    return fun(obj,min_s,max_s);
}

long dmsoft::CreateFoobarCustom(long hwnd,long x,long y,PCSTR pic,PCSTR trans_color,double sim)
{
    typedef long (WINAPI * TypeCreateFoobarCustom)(long,long,long,long,PCSTR,PCSTR,double);

    TypeCreateFoobarCustom fun = (TypeCreateFoobarCustom)(ULONG_PTR)(g_dm_hmodule + 105872);
    return fun(obj,hwnd,x,y,pic,trans_color,sim);
}

CString dmsoft::FindStringEx(long hwnd,PCSTR addr_range,PCSTR string_value,long type,long step,long multi_thread,long mode)
{
    typedef PCSTR (WINAPI * TypeFindStringEx)(long,long,PCSTR,PCSTR,long,long,long,long);

    TypeFindStringEx fun = (TypeFindStringEx)(ULONG_PTR)(g_dm_hmodule + 124384);
    return fun(obj,hwnd,addr_range,string_value,type,step,multi_thread,mode);
}

long dmsoft::GetClientRect(long hwnd,long* x1,long* y1,long* x2,long* y2)
{
    typedef long (WINAPI * TypeGetClientRect)(long,long,long*,long*,long*,long*);

    TypeGetClientRect fun = (TypeGetClientRect)(ULONG_PTR)(g_dm_hmodule + 105808);
    return fun(obj,hwnd,x1,y1,x2,y2);
}

long dmsoft::AiYoloSetModel(long index,PCSTR file,PCSTR pwd)
{
    typedef long (WINAPI * TypeAiYoloSetModel)(long,long,PCSTR,PCSTR);

    TypeAiYoloSetModel fun = (TypeAiYoloSetModel)(ULONG_PTR)(g_dm_hmodule + 104416);
    return fun(obj,index,file,pwd);
}

long dmsoft::FoobarSetTrans(long hwnd,long trans,PCSTR color,double sim)
{
    typedef long (WINAPI * TypeFoobarSetTrans)(long,long,long,PCSTR,double);

    TypeFoobarSetTrans fun = (TypeFoobarSetTrans)(ULONG_PTR)(g_dm_hmodule + 117248);
    return fun(obj,hwnd,trans,color,sim);
}

long dmsoft::GetForegroundFocus()
{
    typedef long (WINAPI * TypeGetForegroundFocus)(long);

    TypeGetForegroundFocus fun = (TypeGetForegroundFocus)(ULONG_PTR)(g_dm_hmodule + 108512);
    return fun(obj);
}

long dmsoft::GetForegroundWindow()
{
    typedef long (WINAPI * TypeGetForegroundWindow)(long);

    TypeGetForegroundWindow fun = (TypeGetForegroundWindow)(ULONG_PTR)(g_dm_hmodule + 115360);
    return fun(obj);
}

long dmsoft::SetExcludeRegion(long type,PCSTR info)
{
    typedef long (WINAPI * TypeSetExcludeRegion)(long,long,PCSTR);

    TypeSetExcludeRegion fun = (TypeSetExcludeRegion)(ULONG_PTR)(g_dm_hmodule + 104832);
    return fun(obj,type,info);
}

long dmsoft::SendStringIme2(long hwnd,PCSTR str,long mode)
{
    typedef long (WINAPI * TypeSendStringIme2)(long,long,PCSTR,long);

    TypeSendStringIme2 fun = (TypeSendStringIme2)(ULONG_PTR)(g_dm_hmodule + 119520);
    return fun(obj,hwnd,str,mode);
}

long dmsoft::ActiveInputMethod(long hwnd,PCSTR id)
{
    typedef long (WINAPI * TypeActiveInputMethod)(long,long,PCSTR);

    TypeActiveInputMethod fun = (TypeActiveInputMethod)(ULONG_PTR)(g_dm_hmodule + 124320);
    return fun(obj,hwnd,id);
}

long dmsoft::FoobarDrawPic(long hwnd,long x,long y,PCSTR pic,PCSTR trans_color)
{
    typedef long (WINAPI * TypeFoobarDrawPic)(long,long,long,long,PCSTR,PCSTR);

    TypeFoobarDrawPic fun = (TypeFoobarDrawPic)(ULONG_PTR)(g_dm_hmodule + 114288);
    return fun(obj,hwnd,x,y,pic,trans_color);
}

long dmsoft::AiYoloSetVersion(PCSTR ver)
{
    typedef long (WINAPI * TypeAiYoloSetVersion)(long,PCSTR);

    TypeAiYoloSetVersion fun = (TypeAiYoloSetVersion)(ULONG_PTR)(g_dm_hmodule + 118496);
    return fun(obj,ver);
}

CString dmsoft::FindColorE(long x1,long y1,long x2,long y2,PCSTR color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindColorE)(long,long,long,long,long,PCSTR,double,long);

    TypeFindColorE fun = (TypeFindColorE)(ULONG_PTR)(g_dm_hmodule + 120384);
    return fun(obj,x1,y1,x2,y2,color,sim,dir);
}

long dmsoft::LeftClick()
{
    typedef long (WINAPI * TypeLeftClick)(long);

    TypeLeftClick fun = (TypeLeftClick)(ULONG_PTR)(g_dm_hmodule + 118096);
    return fun(obj);
}

long dmsoft::IsFileExist(PCSTR file)
{
    typedef long (WINAPI * TypeIsFileExist)(long,PCSTR);

    TypeIsFileExist fun = (TypeIsFileExist)(ULONG_PTR)(g_dm_hmodule + 113824);
    return fun(obj,file);
}

long dmsoft::Is64Bit()
{
    typedef long (WINAPI * TypeIs64Bit)(long);

    TypeIs64Bit fun = (TypeIs64Bit)(ULONG_PTR)(g_dm_hmodule + 110512);
    return fun(obj);
}

CString dmsoft::FindShapeE(long x1,long y1,long x2,long y2,PCSTR offset_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindShapeE)(long,long,long,long,long,PCSTR,double,long);

    TypeFindShapeE fun = (TypeFindShapeE)(ULONG_PTR)(g_dm_hmodule + 120592);
    return fun(obj,x1,y1,x2,y2,offset_color,sim,dir);
}

CString dmsoft::GetDisplayInfo()
{
    typedef PCSTR (WINAPI * TypeGetDisplayInfo)(long);

    TypeGetDisplayInfo fun = (TypeGetDisplayInfo)(ULONG_PTR)(g_dm_hmodule + 122992);
    return fun(obj);
}

long dmsoft::SetEnumWindowDelay(long delay)
{
    typedef long (WINAPI * TypeSetEnumWindowDelay)(long,long);

    TypeSetEnumWindowDelay fun = (TypeSetEnumWindowDelay)(ULONG_PTR)(g_dm_hmodule + 114720);
    return fun(obj,delay);
}

long dmsoft::RegNoMac(PCSTR code,PCSTR ver)
{
    typedef long (WINAPI * TypeRegNoMac)(long,PCSTR,PCSTR);

    TypeRegNoMac fun = (TypeRegNoMac)(ULONG_PTR)(g_dm_hmodule + 118960);
    return fun(obj,code,ver);
}

long dmsoft::KeyUpChar(PCSTR key_str)
{
    typedef long (WINAPI * TypeKeyUpChar)(long,PCSTR);

    TypeKeyUpChar fun = (TypeKeyUpChar)(ULONG_PTR)(g_dm_hmodule + 121904);
    return fun(obj,key_str);
}

long dmsoft::SetDisplayAcceler(long level)
{
    typedef long (WINAPI * TypeSetDisplayAcceler)(long,long);

    TypeSetDisplayAcceler fun = (TypeSetDisplayAcceler)(ULONG_PTR)(g_dm_hmodule + 101088);
    return fun(obj,level);
}

long dmsoft::SetRowGapNoDict(long row_gap)
{
    typedef long (WINAPI * TypeSetRowGapNoDict)(long,long);

    TypeSetRowGapNoDict fun = (TypeSetRowGapNoDict)(ULONG_PTR)(g_dm_hmodule + 118256);
    return fun(obj,row_gap);
}

long dmsoft::EnableMouseAccuracy(long en)
{
    typedef long (WINAPI * TypeEnableMouseAccuracy)(long,long);

    TypeEnableMouseAccuracy fun = (TypeEnableMouseAccuracy)(ULONG_PTR)(g_dm_hmodule + 123760);
    return fun(obj,en);
}

long dmsoft::MoveTo(long x,long y)
{
    typedef long (WINAPI * TypeMoveTo)(long,long,long);

    TypeMoveTo fun = (TypeMoveTo)(ULONG_PTR)(g_dm_hmodule + 109088);
    return fun(obj,x,y);
}

long dmsoft::KeyPressChar(PCSTR key_str)
{
    typedef long (WINAPI * TypeKeyPressChar)(long,PCSTR);

    TypeKeyPressChar fun = (TypeKeyPressChar)(ULONG_PTR)(g_dm_hmodule + 116464);
    return fun(obj,key_str);
}

long dmsoft::RightDown()
{
    typedef long (WINAPI * TypeRightDown)(long);

    TypeRightDown fun = (TypeRightDown)(ULONG_PTR)(g_dm_hmodule + 124576);
    return fun(obj);
}

long dmsoft::AiYoloSetModelMemory(long index,long addr,long size,PCSTR pwd)
{
    typedef long (WINAPI * TypeAiYoloSetModelMemory)(long,long,long,long,PCSTR);

    TypeAiYoloSetModelMemory fun = (TypeAiYoloSetModelMemory)(ULONG_PTR)(g_dm_hmodule + 117600);
    return fun(obj,index,addr,size,pwd);
}

long dmsoft::WriteIni(PCSTR section,PCSTR key,PCSTR v,PCSTR file)
{
    typedef long (WINAPI * TypeWriteIni)(long,PCSTR,PCSTR,PCSTR,PCSTR);

    TypeWriteIni fun = (TypeWriteIni)(ULONG_PTR)(g_dm_hmodule + 101232);
    return fun(obj,section,key,v,file);
}

long dmsoft::DmGuardLoadCustom(PCSTR type,PCSTR path)
{
    typedef long (WINAPI * TypeDmGuardLoadCustom)(long,PCSTR,PCSTR);

    TypeDmGuardLoadCustom fun = (TypeDmGuardLoadCustom)(ULONG_PTR)(g_dm_hmodule + 106896);
    return fun(obj,type,path);
}

long dmsoft::CreateFolder(PCSTR folder_name)
{
    typedef long (WINAPI * TypeCreateFolder)(long,PCSTR);

    TypeCreateFolder fun = (TypeCreateFolder)(ULONG_PTR)(g_dm_hmodule + 113120);
    return fun(obj,folder_name);
}

long dmsoft::EnableRealMouse(long en,long mousedelay,long mousestep)
{
    typedef long (WINAPI * TypeEnableRealMouse)(long,long,long,long);

    TypeEnableRealMouse fun = (TypeEnableRealMouse)(ULONG_PTR)(g_dm_hmodule + 105952);
    return fun(obj,en,mousedelay,mousestep);
}

CString dmsoft::GetBasePath()
{
    typedef PCSTR (WINAPI * TypeGetBasePath)(long);

    TypeGetBasePath fun = (TypeGetBasePath)(ULONG_PTR)(g_dm_hmodule + 107312);
    return fun(obj);
}

long dmsoft::GetFps()
{
    typedef long (WINAPI * TypeGetFps)(long);

    TypeGetFps fun = (TypeGetFps)(ULONG_PTR)(g_dm_hmodule + 106016);
    return fun(obj);
}

long dmsoft::EnableGetColorByCapture(long enable)
{
    typedef long (WINAPI * TypeEnableGetColorByCapture)(long,long);

    TypeEnableGetColorByCapture fun = (TypeEnableGetColorByCapture)(ULONG_PTR)(g_dm_hmodule + 109216);
    return fun(obj,enable);
}

long dmsoft::SetDisplayInput(PCSTR mode)
{
    typedef long (WINAPI * TypeSetDisplayInput)(long,PCSTR);

    TypeSetDisplayInput fun = (TypeSetDisplayInput)(ULONG_PTR)(g_dm_hmodule + 110944);
    return fun(obj,mode);
}

CString dmsoft::Hex64(LONGLONG v)
{
    typedef PCSTR (WINAPI * TypeHex64)(long,LONGLONG);

    TypeHex64 fun = (TypeHex64)(ULONG_PTR)(g_dm_hmodule + 105296);
    return fun(obj,v);
}

long dmsoft::ScreenToClient(long hwnd,long* x,long* y)
{
    typedef long (WINAPI * TypeScreenToClient)(long,long,long*,long*);

    TypeScreenToClient fun = (TypeScreenToClient)(ULONG_PTR)(g_dm_hmodule + 111392);
    return fun(obj,hwnd,x,y);
}

long dmsoft::AiEnableFindPicWindow(long enable)
{
    typedef long (WINAPI * TypeAiEnableFindPicWindow)(long,long);

    TypeAiEnableFindPicWindow fun = (TypeAiEnableFindPicWindow)(ULONG_PTR)(g_dm_hmodule + 100064);
    return fun(obj,enable);
}

CString dmsoft::ReadIni(PCSTR section,PCSTR key,PCSTR file)
{
    typedef PCSTR (WINAPI * TypeReadIni)(long,PCSTR,PCSTR,PCSTR);

    TypeReadIni fun = (TypeReadIni)(ULONG_PTR)(g_dm_hmodule + 102912);
    return fun(obj,section,key,file);
}

long dmsoft::ImageToBmp(PCSTR pic_name,PCSTR bmp_name)
{
    typedef long (WINAPI * TypeImageToBmp)(long,PCSTR,PCSTR);

    TypeImageToBmp fun = (TypeImageToBmp)(ULONG_PTR)(g_dm_hmodule + 109152);
    return fun(obj,pic_name,bmp_name);
}

long dmsoft::SetDisplayDelay(long t)
{
    typedef long (WINAPI * TypeSetDisplayDelay)(long,long);

    TypeSetDisplayDelay fun = (TypeSetDisplayDelay)(ULONG_PTR)(g_dm_hmodule + 122784);
    return fun(obj,t);
}

long dmsoft::WheelUp()
{
    typedef long (WINAPI * TypeWheelUp)(long);

    TypeWheelUp fun = (TypeWheelUp)(ULONG_PTR)(g_dm_hmodule + 102688);
    return fun(obj);
}

long dmsoft::CopyFile(PCSTR src_file,PCSTR dst_file,long over)
{
    typedef long (WINAPI * TypeCopyFile)(long,PCSTR,PCSTR,long);

    TypeCopyFile fun = (TypeCopyFile)(ULONG_PTR)(g_dm_hmodule + 100688);
    return fun(obj,src_file,dst_file,over);
}

long dmsoft::FindWindowEx(long parent,PCSTR class_name,PCSTR title_name)
{
    typedef long (WINAPI * TypeFindWindowEx)(long,long,PCSTR,PCSTR);

    TypeFindWindowEx fun = (TypeFindWindowEx)(ULONG_PTR)(g_dm_hmodule + 115408);
    return fun(obj,parent,class_name,title_name);
}

long dmsoft::SetFindPicMultithreadCount(long count)
{
    typedef long (WINAPI * TypeSetFindPicMultithreadCount)(long,long);

    TypeSetFindPicMultithreadCount fun = (TypeSetFindPicMultithreadCount)(ULONG_PTR)(g_dm_hmodule + 106784);
    return fun(obj,count);
}

long dmsoft::GetScreenDataBmp(long x1,long y1,long x2,long y2,long* data,long* size)
{
    typedef long (WINAPI * TypeGetScreenDataBmp)(long,long,long,long,long,long*,long*);

    TypeGetScreenDataBmp fun = (TypeGetScreenDataBmp)(ULONG_PTR)(g_dm_hmodule + 107136);
    return fun(obj,x1,y1,x2,y2,data,size);
}

long dmsoft::GetWordResultPos(PCSTR str,long index,long* x,long* y)
{
    typedef long (WINAPI * TypeGetWordResultPos)(long,PCSTR,long,long*,long*);

    TypeGetWordResultPos fun = (TypeGetWordResultPos)(ULONG_PTR)(g_dm_hmodule + 114352);
    return fun(obj,str,index,x,y);
}

long dmsoft::LeftDoubleClick()
{
    typedef long (WINAPI * TypeLeftDoubleClick)(long);

    TypeLeftDoubleClick fun = (TypeLeftDoubleClick)(ULONG_PTR)(g_dm_hmodule + 101136);
    return fun(obj);
}

CString dmsoft::ReadStringAddr(long hwnd,LONGLONG addr,long type,long len)
{
    typedef PCSTR (WINAPI * TypeReadStringAddr)(long,long,LONGLONG,long,long);

    TypeReadStringAddr fun = (TypeReadStringAddr)(ULONG_PTR)(g_dm_hmodule + 118608);
    return fun(obj,hwnd,addr,type,len);
}

CString dmsoft::ReadData(long hwnd,PCSTR addr,long len)
{
    typedef PCSTR (WINAPI * TypeReadData)(long,long,PCSTR,long);

    TypeReadData fun = (TypeReadData)(ULONG_PTR)(g_dm_hmodule + 111232);
    return fun(obj,hwnd,addr,len);
}

long dmsoft::AddDict(long index,PCSTR dict_info)
{
    typedef long (WINAPI * TypeAddDict)(long,long,PCSTR);

    TypeAddDict fun = (TypeAddDict)(ULONG_PTR)(g_dm_hmodule + 106336);
    return fun(obj,index,dict_info);
}

long dmsoft::SetInputDm(long input_dm,long rx,long ry)
{
    typedef long (WINAPI * TypeSetInputDm)(long,long,long,long);

    TypeSetInputDm fun = (TypeSetInputDm)(ULONG_PTR)(g_dm_hmodule + 108656);
    return fun(obj,input_dm,rx,ry);
}

long dmsoft::GetWindowProcessId(long hwnd)
{
    typedef long (WINAPI * TypeGetWindowProcessId)(long,long);

    TypeGetWindowProcessId fun = (TypeGetWindowProcessId)(ULONG_PTR)(g_dm_hmodule + 124464);
    return fun(obj,hwnd);
}

long dmsoft::WriteDataAddrFromBin(long hwnd,LONGLONG addr,long data,long len)
{
    typedef long (WINAPI * TypeWriteDataAddrFromBin)(long,long,LONGLONG,long,long);

    TypeWriteDataAddrFromBin fun = (TypeWriteDataAddrFromBin)(ULONG_PTR)(g_dm_hmodule + 121120);
    return fun(obj,hwnd,addr,data,len);
}

CString dmsoft::AiFindPicMemEx(long x1,long y1,long x2,long y2,PCSTR pic_info,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeAiFindPicMemEx)(long,long,long,long,long,PCSTR,double,long);

    TypeAiFindPicMemEx fun = (TypeAiFindPicMemEx)(ULONG_PTR)(g_dm_hmodule + 102976);
    return fun(obj,x1,y1,x2,y2,pic_info,sim,dir);
}

long dmsoft::TerminateProcess(long pid)
{
    typedef long (WINAPI * TypeTerminateProcess)(long,long);

    TypeTerminateProcess fun = (TypeTerminateProcess)(ULONG_PTR)(g_dm_hmodule + 112032);
    return fun(obj,pid);
}

CString dmsoft::VirtualQueryEx(long hwnd,LONGLONG addr,long pmbi)
{
    typedef PCSTR (WINAPI * TypeVirtualQueryEx)(long,long,LONGLONG,long);

    TypeVirtualQueryEx fun = (TypeVirtualQueryEx)(ULONG_PTR)(g_dm_hmodule + 101632);
    return fun(obj,hwnd,addr,pmbi);
}

long dmsoft::EnableKeypadSync(long enable,long time_out)
{
    typedef long (WINAPI * TypeEnableKeypadSync)(long,long,long);

    TypeEnableKeypadSync fun = (TypeEnableKeypadSync)(ULONG_PTR)(g_dm_hmodule + 109968);
    return fun(obj,enable,time_out);
}

long dmsoft::AiYoloUseModel(long index)
{
    typedef long (WINAPI * TypeAiYoloUseModel)(long,long);

    TypeAiYoloUseModel fun = (TypeAiYoloUseModel)(ULONG_PTR)(g_dm_hmodule + 110032);
    return fun(obj,index);
}

long dmsoft::DeleteFile(PCSTR file)
{
    typedef long (WINAPI * TypeDeleteFile)(long,PCSTR);

    TypeDeleteFile fun = (TypeDeleteFile)(ULONG_PTR)(g_dm_hmodule + 99408);
    return fun(obj,file);
}

long dmsoft::GetScreenDepth()
{
    typedef long (WINAPI * TypeGetScreenDepth)(long);

    TypeGetScreenDepth fun = (TypeGetScreenDepth)(ULONG_PTR)(g_dm_hmodule + 102384);
    return fun(obj);
}

long dmsoft::FindColor(long x1,long y1,long x2,long y2,PCSTR color,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindColor)(long,long,long,long,long,PCSTR,double,long,long*,long*);

    TypeFindColor fun = (TypeFindColor)(ULONG_PTR)(g_dm_hmodule + 106112);
    return fun(obj,x1,y1,x2,y2,color,sim,dir,x,y);
}

long dmsoft::MoveR(long rx,long ry)
{
    typedef long (WINAPI * TypeMoveR)(long,long,long);

    TypeMoveR fun = (TypeMoveR)(ULONG_PTR)(g_dm_hmodule + 113504);
    return fun(obj,rx,ry);
}

long dmsoft::LockInput(long lock)
{
    typedef long (WINAPI * TypeLockInput)(long,long);

    TypeLockInput fun = (TypeLockInput)(ULONG_PTR)(g_dm_hmodule + 124272);
    return fun(obj,lock);
}

CString dmsoft::IntToData(LONGLONG int_value,long type)
{
    typedef PCSTR (WINAPI * TypeIntToData)(long,LONGLONG,long);

    TypeIntToData fun = (TypeIntToData)(ULONG_PTR)(g_dm_hmodule + 122272);
    return fun(obj,int_value,type);
}

long dmsoft::FaqPost(PCSTR server,long handle,long request_type,long time_out)
{
    typedef long (WINAPI * TypeFaqPost)(long,PCSTR,long,long,long);

    TypeFaqPost fun = (TypeFaqPost)(ULONG_PTR)(g_dm_hmodule + 107440);
    return fun(obj,server,handle,request_type,time_out);
}

CString dmsoft::GetColorHSV(long x,long y)
{
    typedef PCSTR (WINAPI * TypeGetColorHSV)(long,long,long);

    TypeGetColorHSV fun = (TypeGetColorHSV)(ULONG_PTR)(g_dm_hmodule + 116192);
    return fun(obj,x,y);
}

long dmsoft::FindWindowSuper(PCSTR spec1,long flag1,long type1,PCSTR spec2,long flag2,long type2)
{
    typedef long (WINAPI * TypeFindWindowSuper)(long,PCSTR,long,long,PCSTR,long,long);

    TypeFindWindowSuper fun = (TypeFindWindowSuper)(ULONG_PTR)(g_dm_hmodule + 108432);
    return fun(obj,spec1,flag1,type1,spec2,flag2,type2);
}

long dmsoft::EnableBind(long en)
{
    typedef long (WINAPI * TypeEnableBind)(long,long);

    TypeEnableBind fun = (TypeEnableBind)(ULONG_PTR)(g_dm_hmodule + 116576);
    return fun(obj,en);
}

long dmsoft::SetAero(long enable)
{
    typedef long (WINAPI * TypeSetAero)(long,long);

    TypeSetAero fun = (TypeSetAero)(ULONG_PTR)(g_dm_hmodule + 102640);
    return fun(obj,enable);
}

long dmsoft::DecodeFile(PCSTR file,PCSTR pwd)
{
    typedef long (WINAPI * TypeDecodeFile)(long,PCSTR,PCSTR);

    TypeDecodeFile fun = (TypeDecodeFile)(ULONG_PTR)(g_dm_hmodule + 122496);
    return fun(obj,file,pwd);
}

CString dmsoft::FindPicExS(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicExS)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindPicExS fun = (TypeFindPicExS)(ULONG_PTR)(g_dm_hmodule + 100368);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir);
}

long dmsoft::WriteStringAddr(long hwnd,LONGLONG addr,long type,PCSTR v)
{
    typedef long (WINAPI * TypeWriteStringAddr)(long,long,LONGLONG,long,PCSTR);

    TypeWriteStringAddr fun = (TypeWriteStringAddr)(ULONG_PTR)(g_dm_hmodule + 122720);
    return fun(obj,hwnd,addr,type,v);
}

CString dmsoft::GetCommandLine(long hwnd)
{
    typedef PCSTR (WINAPI * TypeGetCommandLine)(long,long);

    TypeGetCommandLine fun = (TypeGetCommandLine)(ULONG_PTR)(g_dm_hmodule + 100752);
    return fun(obj,hwnd);
}

CString dmsoft::SelectFile()
{
    typedef PCSTR (WINAPI * TypeSelectFile)(long);

    TypeSelectFile fun = (TypeSelectFile)(ULONG_PTR)(g_dm_hmodule + 118144);
    return fun(obj);
}

CString dmsoft::FindPicSimMemEx(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,long sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicSimMemEx)(long,long,long,long,long,PCSTR,PCSTR,long,long);

    TypeFindPicSimMemEx fun = (TypeFindPicSimMemEx)(ULONG_PTR)(g_dm_hmodule + 124912);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir);
}

CString dmsoft::GetWordResultStr(PCSTR str,long index)
{
    typedef PCSTR (WINAPI * TypeGetWordResultStr)(long,PCSTR,long);

    TypeGetWordResultStr fun = (TypeGetWordResultStr)(ULONG_PTR)(g_dm_hmodule + 104768);
    return fun(obj,str,index);
}

long dmsoft::EnablePicCache(long en)
{
    typedef long (WINAPI * TypeEnablePicCache)(long,long);

    TypeEnablePicCache fun = (TypeEnablePicCache)(ULONG_PTR)(g_dm_hmodule + 99536);
    return fun(obj,en);
}

CString dmsoft::FindStrExS(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrExS)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrExS fun = (TypeFindStrExS)(ULONG_PTR)(g_dm_hmodule + 100528);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

long dmsoft::LoadPic(PCSTR pic_name)
{
    typedef long (WINAPI * TypeLoadPic)(long,PCSTR);

    TypeLoadPic fun = (TypeLoadPic)(ULONG_PTR)(g_dm_hmodule + 124128);
    return fun(obj,pic_name);
}

long dmsoft::FindStrFast(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,long* x,long* y)
{
    typedef long (WINAPI * TypeFindStrFast)(long,long,long,long,long,PCSTR,PCSTR,double,long*,long*);

    TypeFindStrFast fun = (TypeFindStrFast)(ULONG_PTR)(g_dm_hmodule + 115584);
    return fun(obj,x1,y1,x2,y2,str,color,sim,x,y);
}

CString dmsoft::FindDouble(long hwnd,PCSTR addr_range,double double_value_min,double double_value_max)
{
    typedef PCSTR (WINAPI * TypeFindDouble)(long,long,PCSTR,double,double);

    TypeFindDouble fun = (TypeFindDouble)(ULONG_PTR)(g_dm_hmodule + 102192);
    return fun(obj,hwnd,addr_range,double_value_min,double_value_max);
}

long dmsoft::SetParam64ToPointer()
{
    typedef long (WINAPI * TypeSetParam64ToPointer)(long);

    TypeSetParam64ToPointer fun = (TypeSetParam64ToPointer)(ULONG_PTR)(g_dm_hmodule + 99952);
    return fun(obj);
}

long dmsoft::SetMemoryFindResultToFile(PCSTR file)
{
    typedef long (WINAPI * TypeSetMemoryFindResultToFile)(long,PCSTR);

    TypeSetMemoryFindResultToFile fun = (TypeSetMemoryFindResultToFile)(ULONG_PTR)(g_dm_hmodule + 110704);
    return fun(obj,file);
}

long dmsoft::WaitKey(long key_code,long time_out)
{
    typedef long (WINAPI * TypeWaitKey)(long,long,long);

    TypeWaitKey fun = (TypeWaitKey)(ULONG_PTR)(g_dm_hmodule + 114528);
    return fun(obj,key_code,time_out);
}

long dmsoft::CreateFoobarEllipse(long hwnd,long x,long y,long w,long h)
{
    typedef long (WINAPI * TypeCreateFoobarEllipse)(long,long,long,long,long,long);

    TypeCreateFoobarEllipse fun = (TypeCreateFoobarEllipse)(ULONG_PTR)(g_dm_hmodule + 114592);
    return fun(obj,hwnd,x,y,w,h);
}

long dmsoft::MoveFile(PCSTR src_file,PCSTR dst_file)
{
    typedef long (WINAPI * TypeMoveFile)(long,PCSTR,PCSTR);

    TypeMoveFile fun = (TypeMoveFile)(ULONG_PTR)(g_dm_hmodule + 102272);
    return fun(obj,src_file,dst_file);
}

long dmsoft::Stop(long id)
{
    typedef long (WINAPI * TypeStop)(long,long);

    TypeStop fun = (TypeStop)(ULONG_PTR)(g_dm_hmodule + 100880);
    return fun(obj,id);
}

long dmsoft::ReleaseRef()
{
    typedef long (WINAPI * TypeReleaseRef)(long);

    TypeReleaseRef fun = (TypeReleaseRef)(ULONG_PTR)(g_dm_hmodule + 111072);
    return fun(obj);
}

CString dmsoft::GetColorBGR(long x,long y)
{
    typedef PCSTR (WINAPI * TypeGetColorBGR)(long,long,long);

    TypeGetColorBGR fun = (TypeGetColorBGR)(ULONG_PTR)(g_dm_hmodule + 100000);
    return fun(obj,x,y);
}

CString dmsoft::EnumIniKeyPwd(PCSTR section,PCSTR file,PCSTR pwd)
{
    typedef PCSTR (WINAPI * TypeEnumIniKeyPwd)(long,PCSTR,PCSTR,PCSTR);

    TypeEnumIniKeyPwd fun = (TypeEnumIniKeyPwd)(ULONG_PTR)(g_dm_hmodule + 116768);
    return fun(obj,section,file,pwd);
}

CString dmsoft::GetMac()
{
    typedef PCSTR (WINAPI * TypeGetMac)(long);

    TypeGetMac fun = (TypeGetMac)(ULONG_PTR)(g_dm_hmodule + 123536);
    return fun(obj);
}

long dmsoft::UseDict(long index)
{
    typedef long (WINAPI * TypeUseDict)(long,long);

    TypeUseDict fun = (TypeUseDict)(ULONG_PTR)(g_dm_hmodule + 104656);
    return fun(obj,index);
}

CString dmsoft::FindDataEx(long hwnd,PCSTR addr_range,PCSTR data,long step,long multi_thread,long mode)
{
    typedef PCSTR (WINAPI * TypeFindDataEx)(long,long,PCSTR,PCSTR,long,long,long);

    TypeFindDataEx fun = (TypeFindDataEx)(ULONG_PTR)(g_dm_hmodule + 123200);
    return fun(obj,hwnd,addr_range,data,step,multi_thread,mode);
}

CString dmsoft::Md5(PCSTR str)
{
    typedef PCSTR (WINAPI * TypeMd5)(long,PCSTR);

    TypeMd5 fun = (TypeMd5)(ULONG_PTR)(g_dm_hmodule + 117376);
    return fun(obj,str);
}

CString dmsoft::BGR2RGB(PCSTR bgr_color)
{
    typedef PCSTR (WINAPI * TypeBGR2RGB)(long,PCSTR);

    TypeBGR2RGB fun = (TypeBGR2RGB)(ULONG_PTR)(g_dm_hmodule + 118736);
    return fun(obj,bgr_color);
}

CString dmsoft::FindColorEx(long x1,long y1,long x2,long y2,PCSTR color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindColorEx)(long,long,long,long,long,PCSTR,double,long);

    TypeFindColorEx fun = (TypeFindColorEx)(ULONG_PTR)(g_dm_hmodule + 103600);
    return fun(obj,x1,y1,x2,y2,color,sim,dir);
}

CString dmsoft::OcrExOne(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeOcrExOne)(long,long,long,long,long,PCSTR,double);

    TypeOcrExOne fun = (TypeOcrExOne)(ULONG_PTR)(g_dm_hmodule + 112080);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

long dmsoft::CmpColor(long x,long y,PCSTR color,double sim)
{
    typedef long (WINAPI * TypeCmpColor)(long,long,long,PCSTR,double);

    TypeCmpColor fun = (TypeCmpColor)(ULONG_PTR)(g_dm_hmodule + 109648);
    return fun(obj,x,y,color,sim);
}

CString dmsoft::OcrInFile(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeOcrInFile)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeOcrInFile fun = (TypeOcrInFile)(ULONG_PTR)(g_dm_hmodule + 110608);
    return fun(obj,x1,y1,x2,y2,pic_name,color,sim);
}

long dmsoft::CheckInputMethod(long hwnd,PCSTR id)
{
    typedef long (WINAPI * TypeCheckInputMethod)(long,long,PCSTR);

    TypeCheckInputMethod fun = (TypeCheckInputMethod)(ULONG_PTR)(g_dm_hmodule + 101792);
    return fun(obj,hwnd,id);
}

long dmsoft::MoveWindow(long hwnd,long x,long y)
{
    typedef long (WINAPI * TypeMoveWindow)(long,long,long,long);

    TypeMoveWindow fun = (TypeMoveWindow)(ULONG_PTR)(g_dm_hmodule + 119648);
    return fun(obj,hwnd,x,y);
}

CString dmsoft::GetClipboard()
{
    typedef PCSTR (WINAPI * TypeGetClipboard)(long);

    TypeGetClipboard fun = (TypeGetClipboard)(ULONG_PTR)(g_dm_hmodule + 116624);
    return fun(obj);
}

long dmsoft::FindStr(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,long* x,long* y)
{
    typedef long (WINAPI * TypeFindStr)(long,long,long,long,long,PCSTR,PCSTR,double,long*,long*);

    TypeFindStr fun = (TypeFindStr)(ULONG_PTR)(g_dm_hmodule + 110320);
    return fun(obj,x1,y1,x2,y2,str,color,sim,x,y);
}

long dmsoft::FoobarClearText(long hwnd)
{
    typedef long (WINAPI * TypeFoobarClearText)(long,long);

    TypeFoobarClearText fun = (TypeFoobarClearText)(ULONG_PTR)(g_dm_hmodule + 113072);
    return fun(obj,hwnd);
}

long dmsoft::ClientToScreen(long hwnd,long* x,long* y)
{
    typedef long (WINAPI * TypeClientToScreen)(long,long,long*,long*);

    TypeClientToScreen fun = (TypeClientToScreen)(ULONG_PTR)(g_dm_hmodule + 116512);
    return fun(obj,hwnd,x,y);
}

CString dmsoft::GetCursorShape()
{
    typedef PCSTR (WINAPI * TypeGetCursorShape)(long);

    TypeGetCursorShape fun = (TypeGetCursorShape)(ULONG_PTR)(g_dm_hmodule + 111984);
    return fun(obj);
}

long dmsoft::GetWordResultCount(PCSTR str)
{
    typedef long (WINAPI * TypeGetWordResultCount)(long,PCSTR);

    TypeGetWordResultCount fun = (TypeGetWordResultCount)(ULONG_PTR)(g_dm_hmodule + 103984);
    return fun(obj,str);
}

CString dmsoft::SelectDirectory()
{
    typedef PCSTR (WINAPI * TypeSelectDirectory)(long);

    TypeSelectDirectory fun = (TypeSelectDirectory)(ULONG_PTR)(g_dm_hmodule + 116000);
    return fun(obj);
}

long dmsoft::CapturePng(long x1,long y1,long x2,long y2,PCSTR file)
{
    typedef long (WINAPI * TypeCapturePng)(long,long,long,long,long,PCSTR);

    TypeCapturePng fun = (TypeCapturePng)(ULONG_PTR)(g_dm_hmodule + 114080);
    return fun(obj,x1,y1,x2,y2,file);
}

long dmsoft::KeyDownChar(PCSTR key_str)
{
    typedef long (WINAPI * TypeKeyDownChar)(long,PCSTR);

    TypeKeyDownChar fun = (TypeKeyDownChar)(ULONG_PTR)(g_dm_hmodule + 105600);
    return fun(obj,key_str);
}

long dmsoft::CaptureJpg(long x1,long y1,long x2,long y2,PCSTR file,long quality)
{
    typedef long (WINAPI * TypeCaptureJpg)(long,long,long,long,long,PCSTR,long);

    TypeCaptureJpg fun = (TypeCaptureJpg)(ULONG_PTR)(g_dm_hmodule + 106400);
    return fun(obj,x1,y1,x2,y2,file,quality);
}

CString dmsoft::FindStrEx(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrEx)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrEx fun = (TypeFindStrEx)(ULONG_PTR)(g_dm_hmodule + 106640);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

long dmsoft::FaqCapture(long x1,long y1,long x2,long y2,long quality,long delay,long time)
{
    typedef long (WINAPI * TypeFaqCapture)(long,long,long,long,long,long,long,long);

    TypeFaqCapture fun = (TypeFaqCapture)(ULONG_PTR)(g_dm_hmodule + 118416);
    return fun(obj,x1,y1,x2,y2,quality,delay,time);
}

long dmsoft::ShowScrMsg(long x1,long y1,long x2,long y2,PCSTR msg,PCSTR color)
{
    typedef long (WINAPI * TypeShowScrMsg)(long,long,long,long,long,PCSTR,PCSTR);

    TypeShowScrMsg fun = (TypeShowScrMsg)(ULONG_PTR)(g_dm_hmodule + 112208);
    return fun(obj,x1,y1,x2,y2,msg,color);
}

long dmsoft::SetKeypadDelay(PCSTR type,long delay)
{
    typedef long (WINAPI * TypeSetKeypadDelay)(long,PCSTR,long);

    TypeSetKeypadDelay fun = (TypeSetKeypadDelay)(ULONG_PTR)(g_dm_hmodule + 110256);
    return fun(obj,type,delay);
}

long dmsoft::SetScreen(long width,long height,long depth)
{
    typedef long (WINAPI * TypeSetScreen)(long,long,long,long);

    TypeSetScreen fun = (TypeSetScreen)(ULONG_PTR)(g_dm_hmodule + 115168);
    return fun(obj,width,height,depth);
}

long dmsoft::Play(PCSTR file)
{
    typedef long (WINAPI * TypePlay)(long,PCSTR);

    TypePlay fun = (TypePlay)(ULONG_PTR)(g_dm_hmodule + 105072);
    return fun(obj,file);
}

long dmsoft::FindWindowByProcessId(long process_id,PCSTR class_name,PCSTR title_name)
{
    typedef long (WINAPI * TypeFindWindowByProcessId)(long,long,PCSTR,PCSTR);

    TypeFindWindowByProcessId fun = (TypeFindWindowByProcessId)(ULONG_PTR)(g_dm_hmodule + 104176);
    return fun(obj,process_id,class_name,title_name);
}

long dmsoft::WriteDouble(long hwnd,PCSTR addr,double double_value)
{
    typedef long (WINAPI * TypeWriteDouble)(long,long,PCSTR,double);

    TypeWriteDouble fun = (TypeWriteDouble)(ULONG_PTR)(g_dm_hmodule + 116048);
    return fun(obj,hwnd,addr,double_value);
}

long dmsoft::GetWindowThreadId(long hwnd)
{
    typedef long (WINAPI * TypeGetWindowThreadId)(long,long);

    TypeGetWindowThreadId fun = (TypeGetWindowThreadId)(ULONG_PTR)(g_dm_hmodule + 107504);
    return fun(obj,hwnd);
}

long dmsoft::GetBindWindow()
{
    typedef long (WINAPI * TypeGetBindWindow)(long);

    TypeGetBindWindow fun = (TypeGetBindWindow)(ULONG_PTR)(g_dm_hmodule + 109712);
    return fun(obj);
}

long dmsoft::FindWindow(PCSTR class_name,PCSTR title_name)
{
    typedef long (WINAPI * TypeFindWindow)(long,PCSTR,PCSTR);

    TypeFindWindow fun = (TypeFindWindow)(ULONG_PTR)(g_dm_hmodule + 104288);
    return fun(obj,class_name,title_name);
}

long dmsoft::AiFindPic(long x1,long y1,long x2,long y2,PCSTR pic_name,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeAiFindPic)(long,long,long,long,long,PCSTR,double,long,long*,long*);

    TypeAiFindPic fun = (TypeAiFindPic)(ULONG_PTR)(g_dm_hmodule + 121536);
    return fun(obj,x1,y1,x2,y2,pic_name,sim,dir,x,y);
}

CString dmsoft::FindInt(long hwnd,PCSTR addr_range,LONGLONG int_value_min,LONGLONG int_value_max,long type)
{
    typedef PCSTR (WINAPI * TypeFindInt)(long,long,PCSTR,LONGLONG,LONGLONG,long);

    TypeFindInt fun = (TypeFindInt)(ULONG_PTR)(g_dm_hmodule + 106256);
    return fun(obj,hwnd,addr_range,int_value_min,int_value_max,type);
}

long dmsoft::IsBind(long hwnd)
{
    typedef long (WINAPI * TypeIsBind)(long,long);

    TypeIsBind fun = (TypeIsBind)(ULONG_PTR)(g_dm_hmodule + 119232);
    return fun(obj,hwnd);
}

long dmsoft::SetSimMode(long mode)
{
    typedef long (WINAPI * TypeSetSimMode)(long,long);

    TypeSetSimMode fun = (TypeSetSimMode)(ULONG_PTR)(g_dm_hmodule + 122896);
    return fun(obj,mode);
}

long dmsoft::GetNowDict()
{
    typedef long (WINAPI * TypeGetNowDict)(long);

    TypeGetNowDict fun = (TypeGetNowDict)(ULONG_PTR)(g_dm_hmodule + 101584);
    return fun(obj);
}

CString dmsoft::GetNetTimeSafe()
{
    typedef PCSTR (WINAPI * TypeGetNetTimeSafe)(long);

    TypeGetNetTimeSafe fun = (TypeGetNetTimeSafe)(ULONG_PTR)(g_dm_hmodule + 107760);
    return fun(obj);
}

CString dmsoft::GetMachineCode()
{
    typedef PCSTR (WINAPI * TypeGetMachineCode)(long);

    TypeGetMachineCode fun = (TypeGetMachineCode)(ULONG_PTR)(g_dm_hmodule + 113456);
    return fun(obj);
}

LONGLONG dmsoft::VirtualAllocEx(long hwnd,LONGLONG addr,long size,long type)
{
    typedef LONGLONG (WINAPI * TypeVirtualAllocEx)(long,long,LONGLONG,long,long);

    TypeVirtualAllocEx fun = (TypeVirtualAllocEx)(ULONG_PTR)(g_dm_hmodule + 99104);
    return fun(obj,hwnd,addr,size,type);
}

CString dmsoft::GetPath()
{
    typedef PCSTR (WINAPI * TypeGetPath)(long);

    TypeGetPath fun = (TypeGetPath)(ULONG_PTR)(g_dm_hmodule + 109600);
    return fun(obj);
}

CString dmsoft::EnumWindowSuper(PCSTR spec1,long flag1,long type1,PCSTR spec2,long flag2,long type2,long sort)
{
    typedef PCSTR (WINAPI * TypeEnumWindowSuper)(long,PCSTR,long,long,PCSTR,long,long,long);

    TypeEnumWindowSuper fun = (TypeEnumWindowSuper)(ULONG_PTR)(g_dm_hmodule + 107360);
    return fun(obj,spec1,flag1,type1,spec2,flag2,type2,sort);
}

LONGLONG dmsoft::GetModuleBaseAddr(long hwnd,PCSTR module_name)
{
    typedef LONGLONG (WINAPI * TypeGetModuleBaseAddr)(long,long,PCSTR);

    TypeGetModuleBaseAddr fun = (TypeGetModuleBaseAddr)(ULONG_PTR)(g_dm_hmodule + 108848);
    return fun(obj,hwnd,module_name);
}

CString dmsoft::EnumWindowByProcessId(long pid,PCSTR title,PCSTR class_name,long filter)
{
    typedef PCSTR (WINAPI * TypeEnumWindowByProcessId)(long,long,PCSTR,PCSTR,long);

    TypeEnumWindowByProcessId fun = (TypeEnumWindowByProcessId)(ULONG_PTR)(g_dm_hmodule + 124672);
    return fun(obj,pid,title,class_name,filter);
}

long dmsoft::UnBindWindow()
{
    typedef long (WINAPI * TypeUnBindWindow)(long);

    TypeUnBindWindow fun = (TypeUnBindWindow)(ULONG_PTR)(g_dm_hmodule + 101904);
    return fun(obj);
}

long dmsoft::GetLastError()
{
    typedef long (WINAPI * TypeGetLastError)(long);

    TypeGetLastError fun = (TypeGetLastError)(ULONG_PTR)(g_dm_hmodule + 107936);
    return fun(obj);
}

long dmsoft::FoobarDrawText(long hwnd,long x,long y,long w,long h,PCSTR text,PCSTR color,long align)
{
    typedef long (WINAPI * TypeFoobarDrawText)(long,long,long,long,long,long,PCSTR,PCSTR,long);

    TypeFoobarDrawText fun = (TypeFoobarDrawText)(ULONG_PTR)(g_dm_hmodule + 119712);
    return fun(obj,hwnd,x,y,w,h,text,color,align);
}

long dmsoft::SetMinRowGap(long row_gap)
{
    typedef long (WINAPI * TypeSetMinRowGap)(long,long);

    TypeSetMinRowGap fun = (TypeSetMinRowGap)(ULONG_PTR)(g_dm_hmodule + 122144);
    return fun(obj,row_gap);
}

long dmsoft::LeftUp()
{
    typedef long (WINAPI * TypeLeftUp)(long);

    TypeLeftUp fun = (TypeLeftUp)(ULONG_PTR)(g_dm_hmodule + 113680);
    return fun(obj);
}

long dmsoft::WriteFile(PCSTR file,PCSTR content)
{
    typedef long (WINAPI * TypeWriteFile)(long,PCSTR,PCSTR);

    TypeWriteFile fun = (TypeWriteFile)(ULONG_PTR)(g_dm_hmodule + 105536);
    return fun(obj,file,content);
}

long dmsoft::SetWindowSize(long hwnd,long width,long height)
{
    typedef long (WINAPI * TypeSetWindowSize)(long,long,long,long);

    TypeSetWindowSize fun = (TypeSetWindowSize)(ULONG_PTR)(g_dm_hmodule + 98560);
    return fun(obj,hwnd,width,height);
}

long dmsoft::FaqCaptureFromFile(long x1,long y1,long x2,long y2,PCSTR file,long quality)
{
    typedef long (WINAPI * TypeFaqCaptureFromFile)(long,long,long,long,long,PCSTR,long);

    TypeFaqCaptureFromFile fun = (TypeFaqCaptureFromFile)(ULONG_PTR)(g_dm_hmodule + 116256);
    return fun(obj,x1,y1,x2,y2,file,quality);
}

CString dmsoft::ReadDataAddr(long hwnd,LONGLONG addr,long len)
{
    typedef PCSTR (WINAPI * TypeReadDataAddr)(long,long,LONGLONG,long);

    TypeReadDataAddr fun = (TypeReadDataAddr)(ULONG_PTR)(g_dm_hmodule + 123584);
    return fun(obj,hwnd,addr,len);
}

long dmsoft::IsSurrpotVt()
{
    typedef long (WINAPI * TypeIsSurrpotVt)(long);

    TypeIsSurrpotVt fun = (TypeIsSurrpotVt)(ULONG_PTR)(g_dm_hmodule + 106992);
    return fun(obj);
}

CString dmsoft::GetWindowProcessPath(long hwnd)
{
    typedef PCSTR (WINAPI * TypeGetWindowProcessPath)(long,long);

    TypeGetWindowProcessPath fun = (TypeGetWindowProcessPath)(ULONG_PTR)(g_dm_hmodule + 105232);
    return fun(obj,hwnd);
}

long dmsoft::ClearDict(long index)
{
    typedef long (WINAPI * TypeClearDict)(long,long);

    TypeClearDict fun = (TypeClearDict)(ULONG_PTR)(g_dm_hmodule + 123152);
    return fun(obj,index);
}

long dmsoft::SaveDict(long index,PCSTR file)
{
    typedef long (WINAPI * TypeSaveDict)(long,long,PCSTR);

    TypeSaveDict fun = (TypeSaveDict)(ULONG_PTR)(g_dm_hmodule + 115520);
    return fun(obj,index,file);
}

long dmsoft::ShowTaskBarIcon(long hwnd,long is_show)
{
    typedef long (WINAPI * TypeShowTaskBarIcon)(long,long,long);

    TypeShowTaskBarIcon fun = (TypeShowTaskBarIcon)(ULONG_PTR)(g_dm_hmodule + 119328);
    return fun(obj,hwnd,is_show);
}

CString dmsoft::GetAveHSV(long x1,long y1,long x2,long y2)
{
    typedef PCSTR (WINAPI * TypeGetAveHSV)(long,long,long,long,long);

    TypeGetAveHSV fun = (TypeGetAveHSV)(ULONG_PTR)(g_dm_hmodule + 100176);
    return fun(obj,x1,y1,x2,y2);
}

CString dmsoft::ReadIniPwd(PCSTR section,PCSTR key,PCSTR file,PCSTR pwd)
{
    typedef PCSTR (WINAPI * TypeReadIniPwd)(long,PCSTR,PCSTR,PCSTR,PCSTR);

    TypeReadIniPwd fun = (TypeReadIniPwd)(ULONG_PTR)(g_dm_hmodule + 102064);
    return fun(obj,section,key,file,pwd);
}

long dmsoft::FaqIsPosted()
{
    typedef long (WINAPI * TypeFaqIsPosted)(long);

    TypeFaqIsPosted fun = (TypeFaqIsPosted)(ULONG_PTR)(g_dm_hmodule + 102864);
    return fun(obj);
}

long dmsoft::LeftDown()
{
    typedef long (WINAPI * TypeLeftDown)(long);

    TypeLeftDown fun = (TypeLeftDown)(ULONG_PTR)(g_dm_hmodule + 106736);
    return fun(obj);
}

long dmsoft::DmGuardExtract(PCSTR type,PCSTR path)
{
    typedef long (WINAPI * TypeDmGuardExtract)(long,PCSTR,PCSTR);

    TypeDmGuardExtract fun = (TypeDmGuardExtract)(ULONG_PTR)(g_dm_hmodule + 112160);
    return fun(obj,type,path);
}

long dmsoft::ExitOs(long type)
{
    typedef long (WINAPI * TypeExitOs)(long,long);

    TypeExitOs fun = (TypeExitOs)(ULONG_PTR)(g_dm_hmodule + 115024);
    return fun(obj,type);
}

CString dmsoft::FetchWord(long x1,long y1,long x2,long y2,PCSTR color,PCSTR word)
{
    typedef PCSTR (WINAPI * TypeFetchWord)(long,long,long,long,long,PCSTR,PCSTR);

    TypeFetchWord fun = (TypeFetchWord)(ULONG_PTR)(g_dm_hmodule + 117840);
    return fun(obj,x1,y1,x2,y2,color,word);
}

CString dmsoft::GetDiskSerial(long index)
{
    typedef PCSTR (WINAPI * TypeGetDiskSerial)(long,long);

    TypeGetDiskSerial fun = (TypeGetDiskSerial)(ULONG_PTR)(g_dm_hmodule + 112352);
    return fun(obj,index);
}

long dmsoft::GetDictCount(long index)
{
    typedef long (WINAPI * TypeGetDictCount)(long,long);

    TypeGetDictCount fun = (TypeGetDictCount)(ULONG_PTR)(g_dm_hmodule + 99584);
    return fun(obj,index);
}

CString dmsoft::GetDict(long index,long font_index)
{
    typedef PCSTR (WINAPI * TypeGetDict)(long,long,long);

    TypeGetDict fun = (TypeGetDict)(ULONG_PTR)(g_dm_hmodule + 99184);
    return fun(obj,index,font_index);
}

long dmsoft::SetDict(long index,PCSTR dict_name)
{
    typedef long (WINAPI * TypeSetDict)(long,long,PCSTR);

    TypeSetDict fun = (TypeSetDict)(ULONG_PTR)(g_dm_hmodule + 121280);
    return fun(obj,index,dict_name);
}

CString dmsoft::AiYoloObjectsToString(PCSTR objects)
{
    typedef PCSTR (WINAPI * TypeAiYoloObjectsToString)(long,PCSTR);

    TypeAiYoloObjectsToString fun = (TypeAiYoloObjectsToString)(ULONG_PTR)(g_dm_hmodule + 111456);
    return fun(obj,objects);
}

long dmsoft::GetKeyState(long vk)
{
    typedef long (WINAPI * TypeGetKeyState)(long,long);

    TypeGetKeyState fun = (TypeGetKeyState)(ULONG_PTR)(g_dm_hmodule + 103296);
    return fun(obj,vk);
}

long dmsoft::RightClick()
{
    typedef long (WINAPI * TypeRightClick)(long);

    TypeRightClick fun = (TypeRightClick)(ULONG_PTR)(g_dm_hmodule + 101040);
    return fun(obj);
}

CString dmsoft::EnumWindowByProcess(PCSTR process_name,PCSTR title,PCSTR class_name,long filter)
{
    typedef PCSTR (WINAPI * TypeEnumWindowByProcess)(long,PCSTR,PCSTR,PCSTR,long);

    TypeEnumWindowByProcess fun = (TypeEnumWindowByProcess)(ULONG_PTR)(g_dm_hmodule + 110192);
    return fun(obj,process_name,title,class_name,filter);
}

CString dmsoft::GetDiskModel(long index)
{
    typedef PCSTR (WINAPI * TypeGetDiskModel)(long,long);

    TypeGetDiskModel fun = (TypeGetDiskModel)(ULONG_PTR)(g_dm_hmodule + 102128);
    return fun(obj,index);
}

long dmsoft::SendStringIme(PCSTR str)
{
    typedef long (WINAPI * TypeSendStringIme)(long,PCSTR);

    TypeSendStringIme fun = (TypeSendStringIme)(ULONG_PTR)(g_dm_hmodule + 124000);
    return fun(obj,str);
}

CString dmsoft::AppendPicAddr(PCSTR pic_info,long addr,long size)
{
    typedef PCSTR (WINAPI * TypeAppendPicAddr)(long,PCSTR,long,long);

    TypeAppendPicAddr fun = (TypeAppendPicAddr)(ULONG_PTR)(g_dm_hmodule + 106832);
    return fun(obj,pic_info,addr,size);
}

long dmsoft::DeleteFolder(PCSTR folder_name)
{
    typedef long (WINAPI * TypeDeleteFolder)(long,PCSTR);

    TypeDeleteFolder fun = (TypeDeleteFolder)(ULONG_PTR)(g_dm_hmodule + 118800);
    return fun(obj,folder_name);
}

long dmsoft::GetDPI()
{
    typedef long (WINAPI * TypeGetDPI)(long);

    TypeGetDPI fun = (TypeGetDPI)(ULONG_PTR)(g_dm_hmodule + 107664);
    return fun(obj);
}

long dmsoft::GetCpuType()
{
    typedef long (WINAPI * TypeGetCpuType)(long);

    TypeGetCpuType fun = (TypeGetCpuType)(ULONG_PTR)(g_dm_hmodule + 102432);
    return fun(obj);
}

long dmsoft::WriteIntAddr(long hwnd,LONGLONG addr,long type,LONGLONG v)
{
    typedef long (WINAPI * TypeWriteIntAddr)(long,long,LONGLONG,long,LONGLONG);

    TypeWriteIntAddr fun = (TypeWriteIntAddr)(ULONG_PTR)(g_dm_hmodule + 100240);
    return fun(obj,hwnd,addr,type,v);
}

long dmsoft::GetSpecialWindow(long flag)
{
    typedef long (WINAPI * TypeGetSpecialWindow)(long,long);

    TypeGetSpecialWindow fun = (TypeGetSpecialWindow)(ULONG_PTR)(g_dm_hmodule + 102336);
    return fun(obj,flag);
}

CString dmsoft::EnumProcess(PCSTR name)
{
    typedef PCSTR (WINAPI * TypeEnumProcess)(long,PCSTR);

    TypeEnumProcess fun = (TypeEnumProcess)(ULONG_PTR)(g_dm_hmodule + 112288);
    return fun(obj,name);
}

long dmsoft::AsmClear()
{
    typedef long (WINAPI * TypeAsmClear)(long);

    TypeAsmClear fun = (TypeAsmClear)(ULONG_PTR)(g_dm_hmodule + 119968);
    return fun(obj);
}

long dmsoft::GetWindowState(long hwnd,long flag)
{
    typedef long (WINAPI * TypeGetWindowState)(long,long,long);

    TypeGetWindowState fun = (TypeGetWindowState)(ULONG_PTR)(g_dm_hmodule + 100112);
    return fun(obj,hwnd,flag);
}

CString dmsoft::FindStrFastE(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrFastE)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrFastE fun = (TypeFindStrFastE)(ULONG_PTR)(g_dm_hmodule + 120288);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

long dmsoft::SetColGapNoDict(long col_gap)
{
    typedef long (WINAPI * TypeSetColGapNoDict)(long,long);

    TypeSetColGapNoDict fun = (TypeSetColGapNoDict)(ULONG_PTR)(g_dm_hmodule + 102592);
    return fun(obj,col_gap);
}

CString dmsoft::AiYoloDetectObjects(long x1,long y1,long x2,long y2,float prob,float iou)
{
    typedef PCSTR (WINAPI * TypeAiYoloDetectObjects)(long,long,long,long,long,float,float);

    TypeAiYoloDetectObjects fun = (TypeAiYoloDetectObjects)(ULONG_PTR)(g_dm_hmodule + 116112);
    return fun(obj,x1,y1,x2,y2,prob,iou);
}

long dmsoft::RunApp(PCSTR path,long mode)
{
    typedef long (WINAPI * TypeRunApp)(long,PCSTR,long);

    TypeRunApp fun = (TypeRunApp)(ULONG_PTR)(g_dm_hmodule + 122832);
    return fun(obj,path,mode);
}

CString dmsoft::FindString(long hwnd,PCSTR addr_range,PCSTR string_value,long type)
{
    typedef PCSTR (WINAPI * TypeFindString)(long,long,PCSTR,PCSTR,long);

    TypeFindString fun = (TypeFindString)(ULONG_PTR)(g_dm_hmodule + 110752);
    return fun(obj,hwnd,addr_range,string_value,type);
}

long dmsoft::GetOsType()
{
    typedef long (WINAPI * TypeGetOsType)(long);

    TypeGetOsType fun = (TypeGetOsType)(ULONG_PTR)(g_dm_hmodule + 121632);
    return fun(obj);
}

CString dmsoft::Ocr(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeOcr)(long,long,long,long,long,PCSTR,double);

    TypeOcr fun = (TypeOcr)(ULONG_PTR)(g_dm_hmodule + 110992);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

CString dmsoft::ReadString(long hwnd,PCSTR addr,long type,long len)
{
    typedef PCSTR (WINAPI * TypeReadString)(long,long,PCSTR,long,long);

    TypeReadString fun = (TypeReadString)(ULONG_PTR)(g_dm_hmodule + 121472);
    return fun(obj,hwnd,addr,type,len);
}

float dmsoft::ReadFloatAddr(long hwnd,LONGLONG addr)
{
    typedef float (WINAPI * TypeReadFloatAddr)(long,long,LONGLONG);

    TypeReadFloatAddr fun = (TypeReadFloatAddr)(ULONG_PTR)(g_dm_hmodule + 100816);
    return fun(obj,hwnd,addr);
}

long dmsoft::Beep(long fre,long delay)
{
    typedef long (WINAPI * TypeBeep)(long,long,long);

    TypeBeep fun = (TypeBeep)(ULONG_PTR)(g_dm_hmodule + 104544);
    return fun(obj,fre,delay);
}

long dmsoft::LoadAi(PCSTR file)
{
    typedef long (WINAPI * TypeLoadAi)(long,PCSTR);

    TypeLoadAi fun = (TypeLoadAi)(ULONG_PTR)(g_dm_hmodule + 106944);
    return fun(obj,file);
}

long dmsoft::GetCpuUsage()
{
    typedef long (WINAPI * TypeGetCpuUsage)(long);

    TypeGetCpuUsage fun = (TypeGetCpuUsage)(ULONG_PTR)(g_dm_hmodule + 121072);
    return fun(obj);
}

long dmsoft::EnableShareDict(long en)
{
    typedef long (WINAPI * TypeEnableShareDict)(long,long);

    TypeEnableShareDict fun = (TypeEnableShareDict)(ULONG_PTR)(g_dm_hmodule + 108992);
    return fun(obj,en);
}

long dmsoft::AiYoloDetectObjectsToFile(long x1,long y1,long x2,long y2,float prob,float iou,PCSTR file,long mode)
{
    typedef long (WINAPI * TypeAiYoloDetectObjectsToFile)(long,long,long,long,long,float,float,PCSTR,long);

    TypeAiYoloDetectObjectsToFile fun = (TypeAiYoloDetectObjectsToFile)(ULONG_PTR)(g_dm_hmodule + 109504);
    return fun(obj,x1,y1,x2,y2,prob,iou,file,mode);
}

long dmsoft::FoobarUnlock(long hwnd)
{
    typedef long (WINAPI * TypeFoobarUnlock)(long,long);

    TypeFoobarUnlock fun = (TypeFoobarUnlock)(ULONG_PTR)(g_dm_hmodule + 123952);
    return fun(obj,hwnd);
}

CString dmsoft::GetSystemInfo(PCSTR type,long method)
{
    typedef PCSTR (WINAPI * TypeGetSystemInfo)(long,PCSTR,long);

    TypeGetSystemInfo fun = (TypeGetSystemInfo)(ULONG_PTR)(g_dm_hmodule + 115680);
    return fun(obj,type,method);
}

long dmsoft::GetResultCount(PCSTR str)
{
    typedef long (WINAPI * TypeGetResultCount)(long,PCSTR);

    TypeGetResultCount fun = (TypeGetResultCount)(ULONG_PTR)(g_dm_hmodule + 116720);
    return fun(obj,str);
}

CString dmsoft::EnumWindow(long parent,PCSTR title,PCSTR class_name,long filter)
{
    typedef PCSTR (WINAPI * TypeEnumWindow)(long,long,PCSTR,PCSTR,long);

    TypeEnumWindow fun = (TypeEnumWindow)(ULONG_PTR)(g_dm_hmodule + 115296);
    return fun(obj,parent,title,class_name,filter);
}

long dmsoft::GetResultPos(PCSTR str,long index,long* x,long* y)
{
    typedef long (WINAPI * TypeGetResultPos)(long,PCSTR,long,long*,long*);

    TypeGetResultPos fun = (TypeGetResultPos)(ULONG_PTR)(g_dm_hmodule + 102800);
    return fun(obj,str,index,x,y);
}

long dmsoft::KeyDown(long vk)
{
    typedef long (WINAPI * TypeKeyDown)(long,long);

    TypeKeyDown fun = (TypeKeyDown)(ULONG_PTR)(g_dm_hmodule + 115120);
    return fun(obj,vk);
}

long dmsoft::SetWordLineHeightNoDict(long line_height)
{
    typedef long (WINAPI * TypeSetWordLineHeightNoDict)(long,long);

    TypeSetWordLineHeightNoDict fun = (TypeSetWordLineHeightNoDict)(ULONG_PTR)(g_dm_hmodule + 103792);
    return fun(obj,line_height);
}

long dmsoft::AiFindPicMem(long x1,long y1,long x2,long y2,PCSTR pic_info,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeAiFindPicMem)(long,long,long,long,long,PCSTR,double,long,long*,long*);

    TypeAiFindPicMem fun = (TypeAiFindPicMem)(ULONG_PTR)(g_dm_hmodule + 111696);
    return fun(obj,x1,y1,x2,y2,pic_info,sim,dir,x,y);
}

long dmsoft::FoobarTextRect(long hwnd,long x,long y,long w,long h)
{
    typedef long (WINAPI * TypeFoobarTextRect)(long,long,long,long,long,long);

    TypeFoobarTextRect fun = (TypeFoobarTextRect)(ULONG_PTR)(g_dm_hmodule + 108784);
    return fun(obj,hwnd,x,y,w,h);
}

long dmsoft::GetPointWindow(long x,long y)
{
    typedef long (WINAPI * TypeGetPointWindow)(long,long,long);

    TypeGetPointWindow fun = (TypeGetPointWindow)(ULONG_PTR)(g_dm_hmodule + 118544);
    return fun(obj,x,y);
}

CString dmsoft::FindMultiColorEx(long x1,long y1,long x2,long y2,PCSTR first_color,PCSTR offset_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindMultiColorEx)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindMultiColorEx fun = (TypeFindMultiColorEx)(ULONG_PTR)(g_dm_hmodule + 122560);
    return fun(obj,x1,y1,x2,y2,first_color,offset_color,sim,dir);
}

long dmsoft::FreeProcessMemory(long hwnd)
{
    typedef long (WINAPI * TypeFreeProcessMemory)(long,long);

    TypeFreeProcessMemory fun = (TypeFreeProcessMemory)(ULONG_PTR)(g_dm_hmodule + 111120);
    return fun(obj,hwnd);
}

CString dmsoft::GetMachineCodeNoMac()
{
    typedef PCSTR (WINAPI * TypeGetMachineCodeNoMac)(long);

    TypeGetMachineCodeNoMac fun = (TypeGetMachineCodeNoMac)(ULONG_PTR)(g_dm_hmodule + 120544);
    return fun(obj);
}

long dmsoft::FindWindowByProcess(PCSTR process_name,PCSTR class_name,PCSTR title_name)
{
    typedef long (WINAPI * TypeFindWindowByProcess)(long,PCSTR,PCSTR,PCSTR);

    TypeFindWindowByProcess fun = (TypeFindWindowByProcess)(ULONG_PTR)(g_dm_hmodule + 122336);
    return fun(obj,process_name,class_name,title_name);
}

long dmsoft::SetWindowState(long hwnd,long flag)
{
    typedef long (WINAPI * TypeSetWindowState)(long,long,long);

    TypeSetWindowState fun = (TypeSetWindowState)(ULONG_PTR)(g_dm_hmodule + 102736);
    return fun(obj,hwnd,flag);
}

long dmsoft::CheckFontSmooth()
{
    typedef long (WINAPI * TypeCheckFontSmooth)(long);

    TypeCheckFontSmooth fun = (TypeCheckFontSmooth)(ULONG_PTR)(g_dm_hmodule + 117552);
    return fun(obj);
}

long dmsoft::IsFolderExist(PCSTR folder)
{
    typedef long (WINAPI * TypeIsFolderExist)(long,PCSTR);

    TypeIsFolderExist fun = (TypeIsFolderExist)(ULONG_PTR)(g_dm_hmodule + 121184);
    return fun(obj,folder);
}

long dmsoft::FaqCancel()
{
    typedef long (WINAPI * TypeFaqCancel)(long);

    TypeFaqCancel fun = (TypeFaqCancel)(ULONG_PTR)(g_dm_hmodule + 113968);
    return fun(obj);
}

long dmsoft::SetWindowTransparent(long hwnd,long v)
{
    typedef long (WINAPI * TypeSetWindowTransparent)(long,long,long);

    TypeSetWindowTransparent fun = (TypeSetWindowTransparent)(ULONG_PTR)(g_dm_hmodule + 112896);
    return fun(obj,hwnd,v);
}

long dmsoft::SwitchBindWindow(long hwnd)
{
    typedef long (WINAPI * TypeSwitchBindWindow)(long,long);

    TypeSwitchBindWindow fun = (TypeSwitchBindWindow)(ULONG_PTR)(g_dm_hmodule + 109920);
    return fun(obj,hwnd);
}

long dmsoft::EnableFontSmooth()
{
    typedef long (WINAPI * TypeEnableFontSmooth)(long);

    TypeEnableFontSmooth fun = (TypeEnableFontSmooth)(ULONG_PTR)(g_dm_hmodule + 103936);
    return fun(obj);
}

CString dmsoft::StringToData(PCSTR string_value,long type)
{
    typedef PCSTR (WINAPI * TypeStringToData)(long,PCSTR,long);

    TypeStringToData fun = (TypeStringToData)(ULONG_PTR)(g_dm_hmodule + 114768);
    return fun(obj,string_value,type);
}

long dmsoft::GetWindowRect(long hwnd,long* x1,long* y1,long* x2,long* y2)
{
    typedef long (WINAPI * TypeGetWindowRect)(long,long,long*,long*,long*,long*);

    TypeGetWindowRect fun = (TypeGetWindowRect)(ULONG_PTR)(g_dm_hmodule + 122656);
    return fun(obj,hwnd,x1,y1,x2,y2);
}

CString dmsoft::FindPicEx(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicEx)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindPicEx fun = (TypeFindPicEx)(ULONG_PTR)(g_dm_hmodule + 108160);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir);
}

CString dmsoft::GetWords(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeGetWords)(long,long,long,long,long,PCSTR,double);

    TypeGetWords fun = (TypeGetWords)(ULONG_PTR)(g_dm_hmodule + 107808);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

long dmsoft::SetExactOcr(long exact_ocr)
{
    typedef long (WINAPI * TypeSetExactOcr)(long,long);

    TypeSetExactOcr fun = (TypeSetExactOcr)(ULONG_PTR)(g_dm_hmodule + 123280);
    return fun(obj,exact_ocr);
}

long dmsoft::EnableMouseSync(long enable,long time_out)
{
    typedef long (WINAPI * TypeEnableMouseSync)(long,long,long);

    TypeEnableMouseSync fun = (TypeEnableMouseSync)(ULONG_PTR)(g_dm_hmodule + 98496);
    return fun(obj,enable,time_out);
}

long dmsoft::CapturePre(PCSTR file)
{
    typedef long (WINAPI * TypeCapturePre)(long,PCSTR);

    TypeCapturePre fun = (TypeCapturePre)(ULONG_PTR)(g_dm_hmodule + 109456);
    return fun(obj,file);
}

long dmsoft::BindWindowEx(long hwnd,PCSTR display,PCSTR mouse,PCSTR keypad,PCSTR public_desc,long mode)
{
    typedef long (WINAPI * TypeBindWindowEx)(long,long,PCSTR,PCSTR,PCSTR,PCSTR,long);

    TypeBindWindowEx fun = (TypeBindWindowEx)(ULONG_PTR)(g_dm_hmodule + 99456);
    return fun(obj,hwnd,display,mouse,keypad,public_desc,mode);
}

long dmsoft::FaqCaptureString(PCSTR str)
{
    typedef long (WINAPI * TypeFaqCaptureString)(long,PCSTR);

    TypeFaqCaptureString fun = (TypeFaqCaptureString)(ULONG_PTR)(g_dm_hmodule + 106208);
    return fun(obj,str);
}

long dmsoft::FoobarTextLineGap(long hwnd,long gap)
{
    typedef long (WINAPI * TypeFoobarTextLineGap)(long,long,long);

    TypeFoobarTextLineGap fun = (TypeFoobarTextLineGap)(ULONG_PTR)(g_dm_hmodule + 124848);
    return fun(obj,hwnd,gap);
}

long dmsoft::FoobarDrawLine(long hwnd,long x1,long y1,long x2,long y2,PCSTR color,long style,long width)
{
    typedef long (WINAPI * TypeFoobarDrawLine)(long,long,long,long,long,long,PCSTR,long,long);

    TypeFoobarDrawLine fun = (TypeFoobarDrawLine)(ULONG_PTR)(g_dm_hmodule + 116384);
    return fun(obj,hwnd,x1,y1,x2,y2,color,style,width);
}

long dmsoft::FindInputMethod(PCSTR id)
{
    typedef long (WINAPI * TypeFindInputMethod)(long,PCSTR);

    TypeFindInputMethod fun = (TypeFindInputMethod)(ULONG_PTR)(g_dm_hmodule + 113872);
    return fun(obj,id);
}

long dmsoft::SetPicPwd(PCSTR pwd)
{
    typedef long (WINAPI * TypeSetPicPwd)(long,PCSTR);

    TypeSetPicPwd fun = (TypeSetPicPwd)(ULONG_PTR)(g_dm_hmodule + 123712);
    return fun(obj,pwd);
}

CString dmsoft::GetCursorSpot()
{
    typedef PCSTR (WINAPI * TypeGetCursorSpot)(long);

    TypeGetCursorSpot fun = (TypeGetCursorSpot)(ULONG_PTR)(g_dm_hmodule + 125056);
    return fun(obj);
}

long dmsoft::InitCri()
{
    typedef long (WINAPI * TypeInitCri)(long);

    TypeInitCri fun = (TypeInitCri)(ULONG_PTR)(g_dm_hmodule + 120240);
    return fun(obj);
}

CString dmsoft::FindPicMemE(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicMemE)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindPicMemE fun = (TypeFindPicMemE)(ULONG_PTR)(g_dm_hmodule + 109264);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir);
}

CString dmsoft::FindStrFastS(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,long* x,long* y)
{
    typedef PCSTR (WINAPI * TypeFindStrFastS)(long,long,long,long,long,PCSTR,PCSTR,double,long*,long*);

    TypeFindStrFastS fun = (TypeFindStrFastS)(ULONG_PTR)(g_dm_hmodule + 98672);
    return fun(obj,x1,y1,x2,y2,str,color,sim,x,y);
}

long dmsoft::DeleteIniPwd(PCSTR section,PCSTR key,PCSTR file,PCSTR pwd)
{
    typedef long (WINAPI * TypeDeleteIniPwd)(long,PCSTR,PCSTR,PCSTR,PCSTR);

    TypeDeleteIniPwd fun = (TypeDeleteIniPwd)(ULONG_PTR)(g_dm_hmodule + 99344);
    return fun(obj,section,key,file,pwd);
}

long dmsoft::GetColorNum(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef long (WINAPI * TypeGetColorNum)(long,long,long,long,long,PCSTR,double);

    TypeGetColorNum fun = (TypeGetColorNum)(ULONG_PTR)(g_dm_hmodule + 124048);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

long dmsoft::AiYoloDetectObjectsToDataBmp(long x1,long y1,long x2,long y2,float prob,float iou,long* data,long* size,long mode)
{
    typedef long (WINAPI * TypeAiYoloDetectObjectsToDataBmp)(long,long,long,long,long,float,float,long*,long*,long);

    TypeAiYoloDetectObjectsToDataBmp fun = (TypeAiYoloDetectObjectsToDataBmp)(ULONG_PTR)(g_dm_hmodule + 98928);
    return fun(obj,x1,y1,x2,y2,prob,iou,data,size,mode);
}

long dmsoft::AiYoloFreeModel(long index)
{
    typedef long (WINAPI * TypeAiYoloFreeModel)(long,long);

    TypeAiYoloFreeModel fun = (TypeAiYoloFreeModel)(ULONG_PTR)(g_dm_hmodule + 106592);
    return fun(obj,index);
}

long dmsoft::DisableFontSmooth()
{
    typedef long (WINAPI * TypeDisableFontSmooth)(long);

    TypeDisableFontSmooth fun = (TypeDisableFontSmooth)(ULONG_PTR)(g_dm_hmodule + 118368);
    return fun(obj);
}

long dmsoft::SetExitThread(long mode)
{
    typedef long (WINAPI * TypeSetExitThread)(long,long);

    TypeSetExitThread fun = (TypeSetExitThread)(ULONG_PTR)(g_dm_hmodule + 101536);
    return fun(obj,mode);
}

CString dmsoft::FindPicMemEx(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicMemEx)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindPicMemEx fun = (TypeFindPicMemEx)(ULONG_PTR)(g_dm_hmodule + 101440);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir);
}

long dmsoft::GetDmCount()
{
    typedef long (WINAPI * TypeGetDmCount)(long);

    TypeGetDmCount fun = (TypeGetDmCount)(ULONG_PTR)(g_dm_hmodule + 125008);
    return fun(obj);
}

long dmsoft::FindMulColor(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef long (WINAPI * TypeFindMulColor)(long,long,long,long,long,PCSTR,double);

    TypeFindMulColor fun = (TypeFindMulColor)(ULONG_PTR)(g_dm_hmodule + 111552);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

CString dmsoft::FaqFetch()
{
    typedef PCSTR (WINAPI * TypeFaqFetch)(long);

    TypeFaqFetch fun = (TypeFaqFetch)(ULONG_PTR)(g_dm_hmodule + 117744);
    return fun(obj);
}

long dmsoft::RegExNoMac(PCSTR code,PCSTR ver,PCSTR ip)
{
    typedef long (WINAPI * TypeRegExNoMac)(long,PCSTR,PCSTR,PCSTR);

    TypeRegExNoMac fun = (TypeRegExNoMac)(ULONG_PTR)(g_dm_hmodule + 107552);
    return fun(obj,code,ver,ip);
}

long dmsoft::FoobarUpdate(long hwnd)
{
    typedef long (WINAPI * TypeFoobarUpdate)(long,long);

    TypeFoobarUpdate fun = (TypeFoobarUpdate)(ULONG_PTR)(g_dm_hmodule + 119280);
    return fun(obj,hwnd);
}

double dmsoft::ReadDouble(long hwnd,PCSTR addr)
{
    typedef double (WINAPI * TypeReadDouble)(long,long,PCSTR);

    TypeReadDouble fun = (TypeReadDouble)(ULONG_PTR)(g_dm_hmodule + 110128);
    return fun(obj,hwnd,addr);
}

CString dmsoft::GetCursorShapeEx(long type)
{
    typedef PCSTR (WINAPI * TypeGetCursorShapeEx)(long,long);

    TypeGetCursorShapeEx fun = (TypeGetCursorShapeEx)(ULONG_PTR)(g_dm_hmodule + 117488);
    return fun(obj,type);
}

CString dmsoft::DoubleToData(double double_value)
{
    typedef PCSTR (WINAPI * TypeDoubleToData)(long,double);

    TypeDoubleToData fun = (TypeDoubleToData)(ULONG_PTR)(g_dm_hmodule + 111856);
    return fun(obj,double_value);
}

long dmsoft::SetWordGapNoDict(long word_gap)
{
    typedef long (WINAPI * TypeSetWordGapNoDict)(long,long);

    TypeSetWordGapNoDict fun = (TypeSetWordGapNoDict)(ULONG_PTR)(g_dm_hmodule + 123392);
    return fun(obj,word_gap);
}

double dmsoft::ReadDoubleAddr(long hwnd,LONGLONG addr)
{
    typedef double (WINAPI * TypeReadDoubleAddr)(long,long,LONGLONG);

    TypeReadDoubleAddr fun = (TypeReadDoubleAddr)(ULONG_PTR)(g_dm_hmodule + 113392);
    return fun(obj,hwnd,addr);
}

long dmsoft::FoobarLock(long hwnd)
{
    typedef long (WINAPI * TypeFoobarLock)(long,long);

    TypeFoobarLock fun = (TypeFoobarLock)(ULONG_PTR)(g_dm_hmodule + 109824);
    return fun(obj,hwnd);
}

CString dmsoft::FindStrFastExS(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrFastExS)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrFastExS fun = (TypeFindStrFastExS)(ULONG_PTR)(g_dm_hmodule + 124176);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

long dmsoft::FindStrWithFont(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim,PCSTR font_name,long font_size,long flag,long* x,long* y)
{
    typedef long (WINAPI * TypeFindStrWithFont)(long,long,long,long,long,PCSTR,PCSTR,double,PCSTR,long,long,long*,long*);

    TypeFindStrWithFont fun = (TypeFindStrWithFont)(ULONG_PTR)(g_dm_hmodule + 119856);
    return fun(obj,x1,y1,x2,y2,str,color,sim,font_name,font_size,flag,x,y);
}

long dmsoft::VirtualProtectEx(long hwnd,LONGLONG addr,long size,long type,long old_protect)
{
    typedef long (WINAPI * TypeVirtualProtectEx)(long,long,LONGLONG,long,long,long);

    TypeVirtualProtectEx fun = (TypeVirtualProtectEx)(ULONG_PTR)(g_dm_hmodule + 108912);
    return fun(obj,hwnd,addr,size,type,old_protect);
}

CString dmsoft::GetWindowClass(long hwnd)
{
    typedef PCSTR (WINAPI * TypeGetWindowClass)(long,long);

    TypeGetWindowClass fun = (TypeGetWindowClass)(ULONG_PTR)(g_dm_hmodule + 117056);
    return fun(obj,hwnd);
}

long dmsoft::SetMouseDelay(PCSTR type,long delay)
{
    typedef long (WINAPI * TypeSetMouseDelay)(long,PCSTR,long);

    TypeSetMouseDelay fun = (TypeSetMouseDelay)(ULONG_PTR)(g_dm_hmodule + 104592);
    return fun(obj,type,delay);
}

LONGLONG dmsoft::ReadInt(long hwnd,PCSTR addr,long type)
{
    typedef LONGLONG (WINAPI * TypeReadInt)(long,long,PCSTR,long);

    TypeReadInt fun = (TypeReadInt)(ULONG_PTR)(g_dm_hmodule + 112720);
    return fun(obj,hwnd,addr,type);
}

CString dmsoft::GetAveRGB(long x1,long y1,long x2,long y2)
{
    typedef PCSTR (WINAPI * TypeGetAveRGB)(long,long,long,long,long);

    TypeGetAveRGB fun = (TypeGetAveRGB)(ULONG_PTR)(g_dm_hmodule + 118192);
    return fun(obj,x1,y1,x2,y2);
}

long dmsoft::GetScreenData(long x1,long y1,long x2,long y2)
{
    typedef long (WINAPI * TypeGetScreenData)(long,long,long,long,long);

    TypeGetScreenData fun = (TypeGetScreenData)(ULONG_PTR)(g_dm_hmodule + 125104);
    return fun(obj,x1,y1,x2,y2);
}

long dmsoft::GetMouseSpeed()
{
    typedef long (WINAPI * TypeGetMouseSpeed)(long);

    TypeGetMouseSpeed fun = (TypeGetMouseSpeed)(ULONG_PTR)(g_dm_hmodule + 99248);
    return fun(obj);
}

long dmsoft::Int64ToInt32(LONGLONG v)
{
    typedef long (WINAPI * TypeInt64ToInt32)(long,LONGLONG);

    TypeInt64ToInt32 fun = (TypeInt64ToInt32)(ULONG_PTR)(g_dm_hmodule + 110880);
    return fun(obj,v);
}

CString dmsoft::FindFloatEx(long hwnd,PCSTR addr_range,float float_value_min,float float_value_max,long step,long multi_thread,long mode)
{
    typedef PCSTR (WINAPI * TypeFindFloatEx)(long,long,PCSTR,float,float,long,long,long);

    TypeFindFloatEx fun = (TypeFindFloatEx)(ULONG_PTR)(g_dm_hmodule + 107040);
    return fun(obj,hwnd,addr_range,float_value_min,float_value_max,step,multi_thread,mode);
}

long dmsoft::FoobarPrintText(long hwnd,PCSTR text,PCSTR color)
{
    typedef long (WINAPI * TypeFoobarPrintText)(long,long,PCSTR,PCSTR);

    TypeFoobarPrintText fun = (TypeFoobarPrintText)(ULONG_PTR)(g_dm_hmodule + 108720);
    return fun(obj,hwnd,text,color);
}

CString dmsoft::OcrEx(long x1,long y1,long x2,long y2,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeOcrEx)(long,long,long,long,long,PCSTR,double);

    TypeOcrEx fun = (TypeOcrEx)(ULONG_PTR)(g_dm_hmodule + 113168);
    return fun(obj,x1,y1,x2,y2,color,sim);
}

long dmsoft::FreePic(PCSTR pic_name)
{
    typedef long (WINAPI * TypeFreePic)(long,PCSTR);

    TypeFreePic fun = (TypeFreePic)(ULONG_PTR)(g_dm_hmodule + 103408);
    return fun(obj,pic_name);
}

long dmsoft::WriteData(long hwnd,PCSTR addr,PCSTR data)
{
    typedef long (WINAPI * TypeWriteData)(long,long,PCSTR,PCSTR);

    TypeWriteData fun = (TypeWriteData)(ULONG_PTR)(g_dm_hmodule + 123040);
    return fun(obj,hwnd,addr,data);
}

long dmsoft::MoveDD(long dx,long dy)
{
    typedef long (WINAPI * TypeMoveDD)(long,long,long);

    TypeMoveDD fun = (TypeMoveDD)(ULONG_PTR)(g_dm_hmodule + 121840);
    return fun(obj,dx,dy);
}

long dmsoft::SetShowErrorMsg(long show)
{
    typedef long (WINAPI * TypeSetShowErrorMsg)(long,long);

    TypeSetShowErrorMsg fun = (TypeSetShowErrorMsg)(ULONG_PTR)(g_dm_hmodule + 101856);
    return fun(obj,show);
}

long dmsoft::SetDictMem(long index,long addr,long size)
{
    typedef long (WINAPI * TypeSetDictMem)(long,long,long,long);

    TypeSetDictMem fun = (TypeSetDictMem)(ULONG_PTR)(g_dm_hmodule + 104704);
    return fun(obj,index,addr,size);
}

long dmsoft::SetClipboard(PCSTR data)
{
    typedef long (WINAPI * TypeSetClipboard)(long,PCSTR);

    TypeSetClipboard fun = (TypeSetClipboard)(ULONG_PTR)(g_dm_hmodule + 104960);
    return fun(obj,data);
}

long dmsoft::FindPicMem(long x1,long y1,long x2,long y2,PCSTR pic_info,PCSTR delta_color,double sim,long dir,long* x,long* y)
{
    typedef long (WINAPI * TypeFindPicMem)(long,long,long,long,long,PCSTR,PCSTR,double,long,long*,long*);

    TypeFindPicMem fun = (TypeFindPicMem)(ULONG_PTR)(g_dm_hmodule + 103696);
    return fun(obj,x1,y1,x2,y2,pic_info,delta_color,sim,dir,x,y);
}

long dmsoft::CreateFoobarRoundRect(long hwnd,long x,long y,long w,long h,long rw,long rh)
{
    typedef long (WINAPI * TypeCreateFoobarRoundRect)(long,long,long,long,long,long,long,long);

    TypeCreateFoobarRoundRect fun = (TypeCreateFoobarRoundRect)(ULONG_PTR)(g_dm_hmodule + 108352);
    return fun(obj,hwnd,x,y,w,h,rw,rh);
}

long dmsoft::WriteFloat(long hwnd,PCSTR addr,float float_value)
{
    typedef long (WINAPI * TypeWriteFloat)(long,long,PCSTR,float);

    TypeWriteFloat fun = (TypeWriteFloat)(ULONG_PTR)(g_dm_hmodule + 111920);
    return fun(obj,hwnd,addr,float_value);
}

long dmsoft::VirtualFreeEx(long hwnd,LONGLONG addr)
{
    typedef long (WINAPI * TypeVirtualFreeEx)(long,long,LONGLONG);

    TypeVirtualFreeEx fun = (TypeVirtualFreeEx)(ULONG_PTR)(g_dm_hmodule + 105120);
    return fun(obj,hwnd,addr);
}

CString dmsoft::GetDictInfo(PCSTR str,PCSTR font_name,long font_size,long flag)
{
    typedef PCSTR (WINAPI * TypeGetDictInfo)(long,PCSTR,PCSTR,long,long);

    TypeGetDictInfo fun = (TypeGetDictInfo)(ULONG_PTR)(g_dm_hmodule + 100624);
    return fun(obj,str,font_name,font_size,flag);
}

long dmsoft::KeyPress(long vk)
{
    typedef long (WINAPI * TypeKeyPress)(long,long);

    TypeKeyPress fun = (TypeKeyPress)(ULONG_PTR)(g_dm_hmodule + 118688);
    return fun(obj,vk);
}

long dmsoft::SetClientSize(long hwnd,long width,long height)
{
    typedef long (WINAPI * TypeSetClientSize)(long,long,long,long);

    TypeSetClientSize fun = (TypeSetClientSize)(ULONG_PTR)(g_dm_hmodule + 104896);
    return fun(obj,hwnd,width,height);
}

CString dmsoft::ExcludePos(PCSTR all_pos,long type,long x1,long y1,long x2,long y2)
{
    typedef PCSTR (WINAPI * TypeExcludePos)(long,PCSTR,long,long,long,long,long);

    TypeExcludePos fun = (TypeExcludePos)(ULONG_PTR)(g_dm_hmodule + 120992);
    return fun(obj,all_pos,type,x1,y1,x2,y2);
}

CString dmsoft::MoveToEx(long x,long y,long w,long h)
{
    typedef PCSTR (WINAPI * TypeMoveToEx)(long,long,long,long,long);

    TypeMoveToEx fun = (TypeMoveToEx)(ULONG_PTR)(g_dm_hmodule + 120688);
    return fun(obj,x,y,w,h);
}

long dmsoft::SetDictPwd(PCSTR pwd)
{
    typedef long (WINAPI * TypeSetDictPwd)(long,PCSTR);

    TypeSetDictPwd fun = (TypeSetDictPwd)(ULONG_PTR)(g_dm_hmodule + 104128);
    return fun(obj,pwd);
}

long dmsoft::FoobarSetFont(long hwnd,PCSTR font_name,long size,long flag)
{
    typedef long (WINAPI * TypeFoobarSetFont)(long,long,PCSTR,long,long);

    TypeFoobarSetFont fun = (TypeFoobarSetFont)(ULONG_PTR)(g_dm_hmodule + 111632);
    return fun(obj,hwnd,font_name,size,flag);
}

CString dmsoft::GetNetTimeByIp(PCSTR ip)
{
    typedef PCSTR (WINAPI * TypeGetNetTimeByIp)(long,PCSTR);

    TypeGetNetTimeByIp fun = (TypeGetNetTimeByIp)(ULONG_PTR)(g_dm_hmodule + 105360);
    return fun(obj,ip);
}

long dmsoft::EnableKeypadPatch(long enable)
{
    typedef long (WINAPI * TypeEnableKeypadPatch)(long,long);

    TypeEnableKeypadPatch fun = (TypeEnableKeypadPatch)(ULONG_PTR)(g_dm_hmodule + 116672);
    return fun(obj,enable);
}

long dmsoft::FoobarStartGif(long hwnd,long x,long y,PCSTR pic_name,long repeat_limit,long delay)
{
    typedef long (WINAPI * TypeFoobarStartGif)(long,long,long,long,PCSTR,long,long);

    TypeFoobarStartGif fun = (TypeFoobarStartGif)(ULONG_PTR)(g_dm_hmodule + 117664);
    return fun(obj,hwnd,x,y,pic_name,repeat_limit,delay);
}

CString dmsoft::FindMultiColorE(long x1,long y1,long x2,long y2,PCSTR first_color,PCSTR offset_color,double sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindMultiColorE)(long,long,long,long,long,PCSTR,PCSTR,double,long);

    TypeFindMultiColorE fun = (TypeFindMultiColorE)(ULONG_PTR)(g_dm_hmodule + 101696);
    return fun(obj,x1,y1,x2,y2,first_color,offset_color,sim,dir);
}

long dmsoft::SetWordGap(long word_gap)
{
    typedef long (WINAPI * TypeSetWordGap)(long,long);

    TypeSetWordGap fun = (TypeSetWordGap)(ULONG_PTR)(g_dm_hmodule + 98624);
    return fun(obj,word_gap);
}

long dmsoft::GetLocale()
{
    typedef long (WINAPI * TypeGetLocale)(long);

    TypeGetLocale fun = (TypeGetLocale)(ULONG_PTR)(g_dm_hmodule + 122096);
    return fun(obj);
}

long dmsoft::GetModuleSize(long hwnd,PCSTR module_name)
{
    typedef long (WINAPI * TypeGetModuleSize)(long,long,PCSTR);

    TypeGetModuleSize fun = (TypeGetModuleSize)(ULONG_PTR)(g_dm_hmodule + 120016);
    return fun(obj,hwnd,module_name);
}

CString dmsoft::FindStrE(long x1,long y1,long x2,long y2,PCSTR str,PCSTR color,double sim)
{
    typedef PCSTR (WINAPI * TypeFindStrE)(long,long,long,long,long,PCSTR,PCSTR,double);

    TypeFindStrE fun = (TypeFindStrE)(ULONG_PTR)(g_dm_hmodule + 122400);
    return fun(obj,x1,y1,x2,y2,str,color,sim);
}

long dmsoft::KeyUp(long vk)
{
    typedef long (WINAPI * TypeKeyUp)(long,long);

    TypeKeyUp fun = (TypeKeyUp)(ULONG_PTR)(g_dm_hmodule + 113248);
    return fun(obj,vk);
}

CString dmsoft::SortPosDistance(PCSTR all_pos,long type,long x,long y)
{
    typedef PCSTR (WINAPI * TypeSortPosDistance)(long,PCSTR,long,long,long);

    TypeSortPosDistance fun = (TypeSortPosDistance)(ULONG_PTR)(g_dm_hmodule + 117120);
    return fun(obj,all_pos,type,x,y);
}

long dmsoft::EnableDisplayDebug(long enable_debug)
{
    typedef long (WINAPI * TypeEnableDisplayDebug)(long,long);

    TypeEnableDisplayDebug fun = (TypeEnableDisplayDebug)(ULONG_PTR)(g_dm_hmodule + 99296);
    return fun(obj,enable_debug);
}

long dmsoft::DeleteIni(PCSTR section,PCSTR key,PCSTR file)
{
    typedef long (WINAPI * TypeDeleteIni)(long,PCSTR,PCSTR,PCSTR);

    TypeDeleteIni fun = (TypeDeleteIni)(ULONG_PTR)(g_dm_hmodule + 111168);
    return fun(obj,section,key,file);
}

CString dmsoft::FindIntEx(long hwnd,PCSTR addr_range,LONGLONG int_value_min,LONGLONG int_value_max,long type,long step,long multi_thread,long mode)
{
    typedef PCSTR (WINAPI * TypeFindIntEx)(long,long,PCSTR,LONGLONG,LONGLONG,long,long,long,long);

    TypeFindIntEx fun = (TypeFindIntEx)(ULONG_PTR)(g_dm_hmodule + 107216);
    return fun(obj,hwnd,addr_range,int_value_min,int_value_max,type,step,multi_thread,mode);
}

long dmsoft::BindWindow(long hwnd,PCSTR display,PCSTR mouse,PCSTR keypad,long mode)
{
    typedef long (WINAPI * TypeBindWindow)(long,long,PCSTR,PCSTR,PCSTR,long);

    TypeBindWindow fun = (TypeBindWindow)(ULONG_PTR)(g_dm_hmodule + 120080);
    return fun(obj,hwnd,display,mouse,keypad,mode);
}

CString dmsoft::GetPicSize(PCSTR pic_name)
{
    typedef PCSTR (WINAPI * TypeGetPicSize)(long,PCSTR);

    TypeGetPicSize fun = (TypeGetPicSize)(ULONG_PTR)(g_dm_hmodule + 114960);
    return fun(obj,pic_name);
}

long dmsoft::AsmSetTimeout(long time_out,long param)
{
    typedef long (WINAPI * TypeAsmSetTimeout)(long,long,long);

    TypeAsmSetTimeout fun = (TypeAsmSetTimeout)(ULONG_PTR)(g_dm_hmodule + 117920);
    return fun(obj,time_out,param);
}

long dmsoft::LockMouseRect(long x1,long y1,long x2,long y2)
{
    typedef long (WINAPI * TypeLockMouseRect)(long,long,long,long,long);

    TypeLockMouseRect fun = (TypeLockMouseRect)(ULONG_PTR)(g_dm_hmodule + 119792);
    return fun(obj,x1,y1,x2,y2);
}

CString dmsoft::FindPicSimE(long x1,long y1,long x2,long y2,PCSTR pic_name,PCSTR delta_color,long sim,long dir)
{
    typedef PCSTR (WINAPI * TypeFindPicSimE)(long,long,long,long,long,PCSTR,PCSTR,long,long);

    TypeFindPicSimE fun = (TypeFindPicSimE)(ULONG_PTR)(g_dm_hmodule + 123440);
    return fun(obj,x1,y1,x2,y2,pic_name,delta_color,sim,dir);
}

CString dmsoft::EnumIniSectionPwd(PCSTR file,PCSTR pwd)
{
    typedef PCSTR (WINAPI * TypeEnumIniSectionPwd)(long,PCSTR,PCSTR);

    TypeEnumIniSectionPwd fun = (TypeEnumIniSectionPwd)(ULONG_PTR)(g_dm_hmodule + 116992);
    return fun(obj,file,pwd);
}

long dmsoft::RightUp()
{
    typedef long (WINAPI * TypeRightUp)(long);

    TypeRightUp fun = (TypeRightUp)(ULONG_PTR)(g_dm_hmodule + 111504);
    return fun(obj);
}

long dmsoft::FoobarTextPrintDir(long hwnd,long dir)
{
    typedef long (WINAPI * TypeFoobarTextPrintDir)(long,long,long);

    TypeFoobarTextPrintDir fun = (TypeFoobarTextPrintDir)(ULONG_PTR)(g_dm_hmodule + 103072);
    return fun(obj,hwnd,dir);
}

CString dmsoft::GetDir(long type)
{
    typedef PCSTR (WINAPI * TypeGetDir)(long,long);

    TypeGetDir fun = (TypeGetDir)(ULONG_PTR)(g_dm_hmodule + 124512);
    return fun(obj,type);
}

CString dmsoft::Hex32(long v)
{
    typedef PCSTR (WINAPI * TypeHex32)(long,long);

    TypeHex32 fun = (TypeHex32)(ULONG_PTR)(g_dm_hmodule + 110080);
    return fun(obj,v);
}

long dmsoft::LeaveCri()
{
    typedef long (WINAPI * TypeLeaveCri)(long);

    TypeLeaveCri fun = (TypeLeaveCri)(ULONG_PTR)(g_dm_hmodule + 120816);
    return fun(obj);
}

long dmsoft::GetTime()
{
    typedef long (WINAPI * TypeGetTime)(long);

    TypeGetTime fun = (TypeGetTime)(ULONG_PTR)(g_dm_hmodule + 103504);
    return fun(obj);
}

long dmsoft::FoobarFillRect(long hwnd,long x1,long y1,long x2,long y2,PCSTR color)
{
    typedef long (WINAPI * TypeFoobarFillRect)(long,long,long,long,long,long,PCSTR);

    TypeFoobarFillRect fun = (TypeFoobarFillRect)(ULONG_PTR)(g_dm_hmodule + 103136);
    return fun(obj,hwnd,x1,y1,x2,y2,color);
}

