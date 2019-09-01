// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	gktpl "github.com/gokeeptech/gktemplate"
	"net/http"
)

func main() {
	// 设置标签
	gktpl.SetNameSpace("llgoer", "{", "}")

	// 使用扩展函数
	var funcs = make(map[string]gktpl.TagFunc)

	// test标签
	funcs["test"] = func(tag *gktpl.GKTag, data *gktpl.D) string {
		name := tag.GetAttribute("name") // 获取属性
		return "string from testtag name is:" + name
	}

	// 扩展自定义标签函数
	gktpl.ExtFuncs(&funcs)

	// 加载模板
	gktpl.LoadDir("./templates/*.htm")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := gktpl.D{
			"info": "Template engine for GoKeep(GK)，GoKeep模板引擎",
		}

		// 渲染模板
		rs, err := gktpl.Parse("templates/simple.htm", data)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, rs)
	})

	http.ListenAndServe(":8088", nil)
}
