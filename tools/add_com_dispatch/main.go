package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	inputFile  = `e:\SRC\go-dm72424\dmsoft_impl.go`
	backupFile = `e:\SRC\go-dm72424\dmsoft_impl.go.bak`
)

type MethodInfo struct {
	Name       string
	ReturnType string
	InParams   []ParamInfo
	OutParams  []ParamInfo
}

type ParamInfo struct {
	Name string
	Type string
}

var funcRegex = regexp.MustCompile(`^func \(dm \*DmSoft\) (\w+)\(([^)]*)\) (\S+) \{$`)

func main() {
	src, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
		os.Exit(1)
	}

	content := strings.ReplaceAll(string(src), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	err = copyFile(inputFile, backupFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "备份文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("已备份到", backupFile)

	lines := strings.Split(content, "\n")
	var result []string
	i := 0
	modified := 0
	skipped := 0

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimRight(line, " \t\r")
		m := funcRegex.FindStringSubmatch(trimmed)
		if m == nil {
			result = append(result, line)
			i++
			continue
		}

		methodName := m[1]
		paramsStr := m[2]
		returnType := m[3]

		if methodName == "Init" || methodName == "Release" {
			result = append(result, line)
			i++
			continue
		}

		if strings.HasPrefix(methodName, "comCall") || strings.HasPrefix(methodName, "getDispid") {
			result = append(result, line)
			i++
			continue
		}

		method := parseMethod(methodName, paramsStr, returnType)

		dispatchCode := generateDispatch(method)

		result = append(result, line)
		i++

		if i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "if dm.disp != nil") {
			skipped++
			for i < len(lines) {
				result = append(result, lines[i])
				i++
				trimmedNext := strings.TrimSpace(lines[i-1])
				if trimmedNext == "}" {
					break
				}
			}
			continue
		}

		result = append(result, dispatchCode)
		modified++
	}

	out := strings.Join(result, "\n")
	err = os.WriteFile(inputFile, []byte(out), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("完成! 修改了 %d 个方法, 跳过 %d 个已有COM分支的方法\n", modified, skipped)
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

func parseMethod(name, paramsStr, returnType string) MethodInfo {
	method := MethodInfo{
		Name:       name,
		ReturnType: returnType,
	}

	if strings.TrimSpace(paramsStr) == "" {
		return method
	}

	scanner := bufio.NewScanner(strings.NewReader(paramsStr))
	scanner.Split(scanCommaSeparated)

	for scanner.Scan() {
		p := strings.TrimSpace(scanner.Text())
		if p == "" {
			continue
		}
		parts := strings.Fields(p)
		if len(parts) < 2 {
			continue
		}
		pType := parts[len(parts)-1]
		pName := parts[len(parts)-2]
		if strings.HasPrefix(pType, "*") {
			method.OutParams = append(method.OutParams, ParamInfo{Name: pName, Type: pType})
		} else {
			method.InParams = append(method.InParams, ParamInfo{Name: pName, Type: pType})
		}
	}

	return method
}

func scanCommaSeparated(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	depth := 0
	start := 0
	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				return i + 1, data[start:i], nil
			}
		}
	}

	if atEOF {
		return len(data), data[start:], nil
	}
	return 0, nil, nil
}

func generateDispatch(method MethodInfo) string {
	hasOutParams := len(method.OutParams) > 0

	if hasOutParams {
		return generateWithOutVars(method)
	}

	return generateSimple(method)
}

func generateSimple(method MethodInfo) string {
	var callFunc string
	switch method.ReturnType {
	case "int32":
		callFunc = "comCallInt32"
	case "string":
		callFunc = "comCallStr"
	case "int64":
		callFunc = "comCallInt64"
	case "float64":
		callFunc = "comCallFloat64"
	case "float32":
		callFunc = "float32(dm.comCallFloat64"
	default:
		callFunc = "comCallInt32"
	}

	args := make([]string, len(method.InParams))
	for i, p := range method.InParams {
		args[i] = convertParam(p)
	}

	argsStr := strings.Join(args, ", ")
	methodNameQuoted := fmt.Sprintf("%q", method.Name)

	var callArgs string
	if len(args) > 0 {
		callArgs = methodNameQuoted + ", " + argsStr
	} else {
		callArgs = methodNameQuoted
	}

	if method.ReturnType == "float32" {
		return fmt.Sprintf("\tif dm.disp != nil {\n\t\treturn %s(%s))\n\t}", callFunc, callArgs)
	}

	return fmt.Sprintf("\tif dm.disp != nil {\n\t\treturn dm.%s(%s)\n\t}", callFunc, callArgs)
}

func generateWithOutVars(method MethodInfo) string {
	inArgs := make([]string, len(method.InParams))
	for i, p := range method.InParams {
		inArgs[i] = convertParam(p)
	}

	outVarNames := make([]string, len(method.OutParams))
	for i, p := range method.OutParams {
		outVarNames[i] = p.Name
	}

	inSliceElements := make([]string, len(method.InParams))
	for i, p := range method.InParams {
		inSliceElements[i] = convertParamForSlice(p)
	}

	inSliceStr := "[]interface{}{" + strings.Join(inSliceElements, ", ") + "}"
	outVarsStr := strings.Join(outVarNames, ", ")
	methodNameQuoted := fmt.Sprintf("%q", method.Name)

	if method.ReturnType == "string" {
		return fmt.Sprintf(
			"\tif dm.disp != nil {\n\t\treturn dm.comCallStrWithOutVars(%s, %s, %s)\n\t}",
			methodNameQuoted, inSliceStr, outVarsStr,
		)
	}

	return fmt.Sprintf(
		"\tif dm.disp != nil {\n\t\treturn dm.comCallWithOutVars(%s, %s, %s)\n\t}",
		methodNameQuoted, inSliceStr, outVarsStr,
	)
}

func convertParam(p ParamInfo) string {
	switch p.Type {
	case "float32":
		return fmt.Sprintf("float64(%s)", p.Name)
	default:
		return p.Name
	}
}

func convertParamForSlice(p ParamInfo) string {
	switch p.Type {
	case "float32":
		return fmt.Sprintf("float64(%s)", p.Name)
	default:
		return p.Name
	}
}
