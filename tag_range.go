// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// field标签函数
package gktemplate

import (
	"fmt"
)

// 解析range标签内容
func TagRange(tag *GKTag, data *D) string {
	innertText := string(tag.GetInnerText())

	gktp, err := parseTemplate(&innertText, "field", "[", "]", "")
	if err != nil {
		return ""
	}
	name := tag.GetAttribute("name")
	items, ok := (*data)[name]

	if !ok {
		return ""
	}

	var resultString string
	var nextTagEnd int
	for _, item := range items.([]D) {
		nextTagEnd = 0
		for i := 0; i < gktp.Count; i++ {
			rs, ok := item[gktp.CTags[i].TagName]
			if ok {
				resultString += string(gktp.SourceString[nextTagEnd : nextTagEnd+gktp.CTags[i].StartPos-nextTagEnd])

				switch v := rs.(type) { //v表示b1 接口转换成Bag对象的值
				case string:
					resultString += v
				case int:
					resultString += fmt.Sprintf("%d", v)
				default:
					resultString += ""
				}

				nextTagEnd = gktp.CTags[i].EndPos
			}
		}
		slen := len(gktp.SourceString)
		if slen > nextTagEnd {
			resultString += string(gktp.SourceString[nextTagEnd:slen])
		}
	}

	return resultString
}
