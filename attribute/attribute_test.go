// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 属性解析器单元测试
package attribute

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

var testattstr = `gokeep att1="func test($me, \"hello\")" myname='aa'   help="1" test=2   	bbb='世界'`

// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`

// 测试解析属性字符串
func TestParse(t *testing.T) {
	att, err := Parse(testattstr)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println("att=", att)
	fmt.Println("att1=", att.GetAttribute("att1"))
}

// 测试遍历字符串
func TestRangeString(t *testing.T) {
	// 需要保证两个遍历字符串结果一致
	total1 := RangeString1(testattstr)
	total2 := RangeString2(testattstr)

	if total1 != total2 {
		t.Error("Test RangeString not pass")
	}

	fmt.Println("total=", total1)
}

// 这里我们测试两种遍历UTF-8字符串的方式

// 第一种，采用range遍历
func RangeString1(t string) int {
	var total int
	for _, c := range t {
		if c == rune('o') {
			total++
		}
		// fmt.Println(string(c), len(string(c)))
	}
	return total
}

// 第二种，采用utf8.DecodeRuneInString
func RangeString2(t string) int {
	var total int
	for len(t) > 0 {
		r, size := utf8.DecodeRuneInString(t)
		if r == rune('o') {
			total++
		}

		t = t[size:]
	}
	return total
}

// Benchmark测试
// go test -bench=. -benchtime=3s -run=none -count=3
func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		Parse(testattstr)
	}
}

func BenchmarkRangeString1(b *testing.B) {
	b.ResetTimer()
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		RangeString1(testattstr)
	}
}

func BenchmarkRangeString2(b *testing.B) {
	b.ResetTimer()
	// var testattstr = `ddd myname='aa'   help="1" test=2   	bbb='世界'`
	for i := 0; i < b.N; i++ {
		RangeString2(testattstr)
	}
}
