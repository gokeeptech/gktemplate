// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 模板引擎
package gktemplate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"text/template"
)

var testtpl = `
Hello，<{gk:field name="info" 名字="中国"/}> 

<{gk:range name='items' test=111 年龄=12 姓名="张三"}>
	<li>[field:id/] - [field:name/]</li>
<{/gk:range}>

<{gk:if condition='a>b'}>
11
{else}
22
<{/gk:if}>
sss

`

// 测试解析字符串
func TestParseString(t *testing.T) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
		{
			"id":   3,
			"name": "GKTemplate",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)，GoKeep模板引擎",
		"items": items,
	}

	rs, err := ParseString(testtpl, data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("[TestParseString]result=", rs)
}

// go test -run TestParseFile
func TestParseFile(t *testing.T) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
		{
			"id":   3,
			"name": "GKTemplate",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)，GoKeep模板引擎",
		"items": items,
	}

	rs, err := ParseFile("./testdata/tpl1.htm", data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("[TestParseFile]result=", rs)
}

// 测试加载目录
func TestLoadDir(t *testing.T) {
	// 载入模板
	LoadDir("testdata/*.htm")
	LoadDir("testdata/**/*.htm")

	// 解析模板
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
		{
			"id":   3,
			"name": "GKTemplate",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)，GoKeep模板引擎",
		"items": items,
	}

	rs, err := Parse("./testdata/tpl1.htm", data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("[TestLoadDir]result=", rs)
	rs, err = Parse("testdata/deep/tpl1.htm", data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("[TestLoadDir]result=", rs)
}

func TestIsEndOfForwardSlash(t *testing.T) {
	var str = "xsada/"
	rs, i := IsEndOfForwardSlash(&str)
	if rs == false || i != 5 {
		t.Errorf("test IsEndOfForwardSlash not passed")
	}
	str = "xsada/ss"
	rs, _ = IsEndOfForwardSlash(&str)
	if rs == true {
		t.Errorf("test IsEndOfForwardSlash not passed")
	}
}

func TestExtFuncs(t *testing.T) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
		{
			"id":   3,
			"name": "GKTemplate",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)，GoKeep模板引擎",
		"items": items,
	}

	var funcs = make(map[string]TagFunc)

	funcs["test"] = func(tag *GKTag, data *D) string {
		name := tag.GetAttribute("name")
		return "string from testtag name is:" + name
	}

	ExtFuncs(&funcs)

	rs, err := ParseFile("./testdata/tpl_extfuncs.htm", data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("[TestExtFuncs]result=", rs)
}

// go test -bench=. -benchtime=3s -run=none -count=3 -benchmem
// v1:25885 ns/op
// v2:13993 ns/op
// v3:6983 ns/op
func BenchmarkParseString(b *testing.B) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)",
		"items": items,
	}
	b.ResetTimer()
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		ParseString(testtpl, data)
	}
}

func BenchmarkParseFile(b *testing.B) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)",
		"items": items,
	}
	b.ResetTimer()
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		ParseFile("./testdata/tpl1.htm", data)
	}
}

func BenchmarkParseLoadDir(b *testing.B) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)",
		"items": items,
	}
	b.ResetTimer()
	// 先加载，后渲染
	LoadDir("testdata/*.htm")
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		Parse("./testdata/tpl1.htm", data)
	}
}

var gtplStr string = `
Hello，{{.info}}

{{ range .items }}
	<li>{{.id}} - {{.name}}</li>
{{ end }}
<{gk:if condition='a>b'}>
11
{else}
22
<{/gk:if}>
sss
`

func TestGoTemplate(t *testing.T) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)",
		"items": items,
	}
	tt := template.Must(template.New("gotpl").Parse(gtplStr))
	err := tt.Execute(os.Stdout, data)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// 下面我们测试下Go自带的模板
// 4124 ns/op	     248 B/op	      14 allocs/op
func BenchmarkGoTemplate(b *testing.B) {
	items := []D{
		{
			"id":   1,
			"name": "GoKeep",
		},
		{
			"id":   2,
			"name": "llgoer",
		},
	}

	data := D{
		"info":  "Template engine for GoKeep(GK)",
		"items": items,
	}
	t := template.Must(template.New("gotpl").Parse(gtplStr))
	for i := 0; i < b.N; i++ {
		t.Execute(ioutil.Discard, data)
	}
}
