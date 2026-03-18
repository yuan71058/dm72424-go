// 72424_C++.cpp : 此文件包含 "main" 函数。程序执行将在此处开始并结束。
//

#include <iostream>
#include <Windows.h>
#include "obj.h"
using namespace std;
#define 大漠插件路径 "xd47243.dll"
#define 破解Dll路径 "Go.dll"
typedef void (WINAPI* GoFunProc)(DWORD g_dm);
int main()
{

    DWORD dm_hmodule = LoadDm(大漠插件路径);
    if (dm_hmodule == NULL) { return 0; }
    HMODULE go_hmodule  = LoadLibraryA(破解Dll路径);
    GoFunProc GoFun = NULL;
    GoFun = (GoFunProc)GetProcAddress(go_hmodule, "Go");
    if (GoFun == NULL) { return 0; }
    GoFun(dm_hmodule);
    dmsoft* dm = new dmsoft;
    long nret =  dm->Reg("", "");
    if (nret == 1)
    {
        cout << "大漠注册成功" << endl;
    }
    else
    {
        cout << "大漠注册成功失败" << endl;
    }

    std::cout << "Hello World!\n";
}

// 运行程序: Ctrl + F5 或调试 >“开始执行(不调试)”菜单
// 调试程序: F5 或调试 >“开始调试”菜单

// 入门使用技巧: 
//   1. 使用解决方案资源管理器窗口添加/管理文件
//   2. 使用团队资源管理器窗口连接到源代码管理
//   3. 使用输出窗口查看生成输出和其他消息
//   4. 使用错误列表窗口查看错误
//   5. 转到“项目”>“添加新项”以创建新的代码文件，或转到“项目”>“添加现有项”以将现有代码文件添加到项目
//   6. 将来，若要再次打开此项目，请转到“文件”>“打开”>“项目”并选择 .sln 文件
