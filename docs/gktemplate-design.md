# GoKeep模板引擎设计

用于解析模板

## 用法设计

### 直接解析字符串

```go
import (
	gktpl "gktemplate"
)

var tplStr = `
Hello,{gk:field name="info"/}
{gk:range name='items'}
	<li>[field:id/] - [field:name/]</li>
{/gk:range}
`

var items = []gktpl.D {
  gktpl.D {
    "id" : 1,
    "name" : "GoKeep",
  },
  gktpl.D {
    "id" : 2,
    "name" : "llgoer",
  },
}

var data = gktpl.D {
  "info": "Template engine for GoKeep(GK)",
  "items": items,
}

result,err := gktpl.ParseString(tplStr, data)
// 输出结果：
// Hello,Template engine for GoKeep(GK)
// <li>1 - GoKeep</li>
// <li>2 - llgoer</li>

// 下面是自定义标签
nameSpace := "gokeep"
tagStart := "<{"
tagEnd := "}>"

result,err := gktpl.ParseStringWithNameSpace(tplStr, data, nameSpace, tagStart, tagEnd)

```

`ParseString`用于直接将字符串，将Data结构体，解析并显示。

## 直接解析文件

这种方法类似直接解析字符串，只是从文件中加载模板信息。

```go
result,err := gktpl.ParseFile(filename, data);
```

## 通过模板目录解析

这种是指定模板目录，然后从模板目录预加载模板并解析。

在调用的时候直接进行渲染即可。

假定`templates`模板目录中含有`tplA.htm`、`tplB.htm`两个模板文件。

则可以采用以下示例代码解析：

```go
pattern := "templates/*.htm"
gktpl.LoadDir(pattern)
result,err := gktpl.Parse("templates/tplA.htm", data)
```

这里需要记住的是，LoadDir需要在程序初始化时候进行预处理。

预处理的过程是将模板中的标签解析过来。等到数据渲染的时候可以快速呈现。

## 标签解析过程

这里先以测试字符串为例子

```go
var tplStr = `
Hello,{gk:field name="info"}
{gk:range name='items'}
	<li>[field:id] - [field:name]</li>
{/gk}
`
```

1.判断当前字符是否和标签开始标记的第一个字符匹配，这里标签开始标记第一个字符是`"{"`，如果匹配，则尝试取出标签开始完整字符判断是否是标签开始。如果未开始，则忽略跳过继续，如已开始，则开始对标签进行解析；

2.解析标签，先判断是否以标签结束字符，这里是`"/}"`或者`}`，然后将其中字符串按照属性进行解析。这里存在一个可能的错误，即当前标签还没有以`/}`完成结束，又出现了1的开始标记，则需要进行错误提示，告知在具体位置出现了错误；

