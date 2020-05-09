// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gin-gonic/gin"
	gkt "github.com/gokeeptech/gktemplate"
	"log"
	"os"
)

const CTYPE = "text/html; charset=utf-8"

func init() {
	// 开发模式下启用日志调试
	if os.Getenv("GKENV") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	r := gin.Default()
	// 加载模板
	gkt.LoadDir("./templates/*.htm")
	r.GET("/", func(c *gin.Context) {
		data := gkt.D{
			"info": "Template engine for GoKeep(GK)，GoKeep模板引擎",
		}
		// 渲染模板
		rs, err := gkt.Parse("templates/simple.htm", data)
		if err != nil {
			log.Println(err)
		}
		c.Data(200, CTYPE, []byte(rs))
	})
	r.Run()
}
