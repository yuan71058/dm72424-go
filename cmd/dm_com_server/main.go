package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type CallRequest struct {
	Offset  uint32
	RetType uint8
	NOut    uint8
	Args    []CallArg
}

type CallArg struct {
	Type  uint8
	IVal  int32
	SVal  string
	FVal  float64
	I64Val int64
}

type CallResponse struct {
	RetType uint8
	IRet    int32
	SRet    string
	FRet    float64
	IRet64  int64
	OutVals []int32
	Err     string
}

var dmObj uintptr

var keepAlive [][]byte

func utf8ToGBK(s string) ([]byte, error) {
	return simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
}

func gbkToUTF8(b []byte) ([]byte, error) {
	return simplifiedchinese.GBK.NewDecoder().Bytes(b)
}

func readCString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	var n int
	for *(*byte)(unsafe.Pointer(ptr + uintptr(n))) != 0 {
		n++
	}
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
	}
	out, err := gbkToUTF8(buf)
	if err != nil {
		return string(buf)
	}
	return string(out)
}

var dmModule syscall.Handle

func handleCall(req CallRequest) CallResponse {
	resp := CallResponse{
		RetType: req.RetType,
	}

	fnAddr := uintptr(dmModule) + uintptr(req.Offset)

	keepAlive = nil

	outCount := 0
	for _, arg := range req.Args {
		if arg.Type == 3 {
			outCount++
		}
	}
	outVars := make([]int32, outCount)
	outIdx := 0

	const maxSlots = 16
	args := make([]uintptr, maxSlots)
	args[0] = dmObj
	idx := 1

	for _, arg := range req.Args {
		switch arg.Type {
		case 0:
			args[idx] = uintptr(arg.IVal)
			idx++
		case 1:
			gbk, err := utf8ToGBK(arg.SVal)
			if err != nil {
				resp.Err = fmt.Sprintf("utf8ToGBK: %v", err)
				return resp
			}
			cstr := make([]byte, len(gbk)+1)
			copy(cstr, gbk)
			keepAlive = append(keepAlive, cstr)
			args[idx] = uintptr(unsafe.Pointer(&cstr[0]))
			idx++
		case 2:
			bits := math.Float64bits(arg.FVal)
			args[idx] = uintptr(bits)
			args[idx+1] = uintptr(bits >> 32)
			idx += 2
		case 3:
			args[idx] = uintptr(unsafe.Pointer(&outVars[outIdx]))
			outIdx++
			idx++
		case 4:
			bits := math.Float32bits(float32(arg.FVal))
			args[idx] = uintptr(bits)
			idx++
		case 5:
			args[idx] = uintptr(arg.I64Val)
			args[idx+1] = uintptr(arg.I64Val >> 32)
			idx += 2
		default:
			args[idx] = 0
			idx++
		}
	}

	nargs := idx

	var r1, r2 uintptr
	switch {
	case nargs <= 3:
		r1, r2, _ = syscall.Syscall(fnAddr, uintptr(nargs), args[0], args[1], args[2])
	case nargs <= 6:
		r1, r2, _ = syscall.Syscall6(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5])
	case nargs <= 9:
		r1, r2, _ = syscall.Syscall9(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	case nargs <= 12:
		r1, r2, _ = syscall.Syscall12(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11])
	case nargs <= 15:
		r1, r2, _ = syscall.Syscall15(fnAddr, uintptr(nargs), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14])
	default:
		resp.Err = fmt.Sprintf("too many args: %d", nargs)
		return resp
	}

	switch req.RetType {
	case 0:
		resp.IRet = int32(r1)
	case 1:
		resp.SRet = readCString(r1)
	case 2:
		resp.FRet = math.Float64frombits(uint64(r1))
	case 3:
		resp.IRet64 = int64(r2)<<32 | int64(r1)
	}

	resp.OutVals = make([]int32, outCount)
	copy(resp.OutVals, outVars)

	return resp
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: dm_com_server <dm.dll> [crack.dll]\n")
		os.Exit(1)
	}
	dmPath := os.Args[1]
	var crackPath string
	if len(os.Args) >= 3 {
		crackPath = os.Args[2]
	}

	module, err := syscall.LoadLibrary(dmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadLibrary dm.dll failed: %v\n", err)
		os.Exit(1)
	}
	dmModule = module
	defer syscall.FreeLibrary(dmModule)

	if crackPath != "" {
		crackModule, err := syscall.LoadLibrary(crackPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadLibrary crack failed: %v\n", err)
			os.Exit(1)
		}
		defer syscall.FreeLibrary(crackModule)

		goProc, err := syscall.GetProcAddress(crackModule, "Go")
		if err != nil {
			fmt.Fprintf(os.Stderr, "GetProcAddress Go failed: %v\n", err)
			os.Exit(1)
		}
		syscall.Syscall(uintptr(goProc), 1, uintptr(dmModule), 0, 0)
	}

	createObjAddr := uintptr(dmModule) + 98304
	dmObj, _, _ = syscall.Syscall(createObjAddr, 0, 0, 0, 0)
	if dmObj == 0 {
		fmt.Fprintf(os.Stderr, "CreateDmObject(offset 0x18000) failed\n")
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Listen failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("READY %d\n", ln.Addr().(*net.TCPAddr).Port)
	os.Stdout.Sync()

	go func() {
		var buf [1]byte
		for {
			_, err := os.Stdin.Read(buf[:])
			if err != nil {
				os.Exit(0)
			}
		}
	}()

	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	for {
		var req CallRequest
		if err := dec.Decode(&req); err != nil {
			return
		}
		resp := handleCall(req)
		if err := enc.Encode(resp); err != nil {
			return
		}
	}
}
