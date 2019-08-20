// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// field标签函数
package gktemplate

// import (
// 	"fmt"
// )

// 解析field标签内容
func TagField(tag *GKTag, data *D) string {
	name := tag.GetAttribute("name")
	v, ok := (*data)[name]
	if ok {
		return v.(string)
	} else {
		return ""
	}
}
