// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 函数参数解析器
package attribute

import (
	"testing"
)

var testfuncstr = `ToUpper(@me, 1, 2,    "测试数据", true, 1.321,'okkkk'    ,"哦1121😯")`

// 测试解析属性字符串
func TestFuncParser(t *testing.T) {
	funcName, args, err := FuncParser(testfuncstr)
	if err != nil {
		t.Errorf(err.Error())
	}

	if funcName != "ToUpper" {
		t.Errorf("TestFuncParser not passed,function name is wrong")
	}
	if len(args) != 8 {
		t.Errorf("TestFuncParser not passed,args is wrong")
	}

}
