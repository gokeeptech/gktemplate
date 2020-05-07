// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 函数参数解析器
package attribute

import (
	"errors"
	// "fmt"
	"strconv"
	"strings"
)

var (
	errNoFunctionName = errors.New("function name not exists")
	errNoFunctionArgs = errors.New("function args not exists")
)

// 解析出函数名称和参数
func FuncParser(funcStr string) (string, []interface{}, error) {
	funcName := []rune("")
	args := []interface{}{}
	hasFuncName := false
	argStart := false
	// 处理字符串，将属性字符串格式化为单个空格间隔字符
	funcStr = strings.TrimSpace(reSplit.ReplaceAllString(funcStr, " "))

	currentArg := []rune("")

	for _, r := range []rune(funcStr) {

		if hasFuncName == false && r != rune('(') {
			funcName = append(funcName, r)
		}

		if r == rune('(') {
			hasFuncName = true
			argStart = true
			continue
		}

		if argStart == true {
			if r != rune(',') && r != rune(')') {
				currentArg = append(currentArg, r)
			}
		}

		if argStart == true {
			if r == rune(',') || r == rune(')') {

				var thisArg interface{}

				currentStr := strings.TrimSpace(string(currentArg))
				currentArg = []rune(currentStr)

				if strings.ToLower(currentStr) == "true" || strings.ToLower(currentStr) == "false" {
					if s, err := strconv.ParseBool(strings.ToLower(currentStr)); err == nil {
						thisArg = s
					}
				}
				if s, err := strconv.ParseFloat(currentStr, 64); err == nil {
					thisArg = s
				}
				if s, err := strconv.ParseInt(currentStr, 10, 64); err == nil {
					thisArg = s
				}
				if (currentArg[0] == rune('"') && currentArg[len(currentArg)-1] == rune('"')) ||
					(currentArg[0] == rune('\'') && currentArg[len(currentArg)-1] == rune('\'')) ||
					(currentArg[0] == rune('`') && currentArg[len(currentArg)-1] == rune('`')) {
					thisArg = string(currentArg[1 : len(currentArg)-1])
				} else {
					thisArg = currentStr
				}
				args = append(args, thisArg)

				currentArg = []rune("")
			}
		}
	}

	if string(funcName) == "" {
		return "", nil, errNoFunctionName
	}

	if len(args) == 0 {
		return "", nil, errNoFunctionArgs
	}

	return string(funcName), args, nil
}
