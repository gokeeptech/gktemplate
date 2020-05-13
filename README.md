# GKTemplate

GKTemplate - 采用Go开发的DedeCMS模板解析器

是否还记得，在PHP流行的年代，有一款开源免费的网站内容管理系统[DedeCMS](https://github.com/dedetech/DedeCMSv5)风靡全国。买域名，买空间主机，下载程序安装，设置好栏目，换上模板，添加采集规则，诺大一个门户瞬间就搭建好了。

SEO成就了DedeCMS的疯狂，但也是因为这样的疯狂让DedeCMS错过了移动互联网。全新的互联网时代，前后端分离，让模版解析渲染变得不再那么重要。

作为DedeCMS的核心开发者之一，也从PHP转到了Go，为了致敬[DedeCMS](https://github.com/dedetech/DedeCMSv5)，决定采用Go开发了一个类DedeCMS模板解析引擎的库。扩展库将骄傲地采用中国首个开源协议[“木兰宽松许可证”](http://license.coscl.org.cn/MulanPSL/)进行发布。

项目名称叫GoKeep，寓意是能够将开源开发继续下去。

## 背景

GKTemplate是一个Go语言开发的模板引擎，设计思想来源于[DedeCMS](https://github.com/dedetech/DedeCMSv5)，由于Go语言内置的模板引擎自由度过高，导致开发使用相对比较困难，对界面模板制作要求会比较高，GKTemplate是一款基于标签、属性机制的模板引擎，在牺牲部分自由度、性能的前提下，优化模板语义机制，使得开发、制作模板变得更为轻松简单。

## 特点

- UTF-8编码支持：模板引擎要求采用UTF-8编码，便于界面能够国际化支持；

- 简单明了属性标记：类似XML结构的属性标记，上手简单，制作模板轻松自如；

- 错误定位：模板标签错误定位，方便模板制作开发调试；

- 标签化语义：类似XHTML标签语义，降低模板制作难度，减少开发制作成本；

- 自由扩展：留有丰富的标签开发接口，方便进行二次扩展；

- 缓存机制：模板解析进行缓存，模板解析性能达到最高；

- 协程并发：采用Go协程机制，标签解析可并发操作，模板渲染性能最高；

- 最小依赖：模板引擎只依赖Go默认库，不依赖任何第三方库；

## 用途

GKTemplate主要用于采用Go编写的HTTP Server中需要自定义呈现数据结构页面，同时也适用于采用模板机制生成例如：静态文件、静态文本等。

## 性能

该模板引擎性能稳定，符合开发者及用户使用要求，详细可参考模板引擎benchmark测试样例。

## 使用方法

执行`go get -u -v github.com/gokeeptech/gktemplate`

使用方法可以参考[examples](./examples)目录中的例子。

如果开启开发模式（模板实时加载），则运行`GKENV=dev go run main.go`

## 资源

- [Github](https://github.com/gokeeptech/gktemplate)

- [Gitee](https://gitee.com/GoKeep/gktemplate)
