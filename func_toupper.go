// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// field标签函数
package gktemplate

import (
	"strings"
)

// 字符串转为大写
func FuncToUpper(v *string, args ...interface{}) string {
	return strings.ToUpper(*v)
}
