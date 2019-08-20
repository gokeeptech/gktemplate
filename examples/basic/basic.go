// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	gktpl "github.com/gokeeptech/gktemplate"
	"net/http"
)

func main() {
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
