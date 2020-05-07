// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// field标签函数
package gktemplate

import (
	"strings"
)

// 字符串转为小写
func FuncToLower(v *string, args ...interface{}) string {
	return strings.ToLower(*v)
}
