# 属性解析器设计

属性解析器目的主要将类HTML标签属性结构解析出来

## 属性可能出现情况

这里我们需要设想几种属性结构可能性：

1.基础
```
gokeep att1=1 att2='string' att3=false att4
```
其中gokeep是标签名称，att1为整数，att2为字符串，att3为bool类型false，att4为空等效`att4=""`

2.字符串中嵌套字符串
```
gokeep att1='func test($me, "hello")'
```
或者
```
gokeep att1="func test($me, 'hello')"
```
再或
```
gokeep att1="func test($me, \"hello\")"
```
再或
```
gokeep att1='func test($me, \'hello\')'
```
再或
```
gokeep att1=`func test($me, 'hello', "ok")`
```
这种情况属于字符串中嵌套字符串
其中gokeep为标签名称