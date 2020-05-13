// Copyright 2020 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 模板引擎
package gktemplate

import (
	"crypto/sha1"
	"errors"
	"fmt"
	attr "github.com/gokeeptech/gktemplate/attribute"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const CharToLow = true  // 是否将属性名称统一转换成小写
const TagMaxLen = 64    // 标签最大字符宽度
const Version = "0.0.9" // 版本号

var (
	defaultNameSpace = "gk" // 默认标签名称
	defaultTagStart  = "<{" // 默认标签开始标记
	defaultTagEnd    = "}>" // 默认标签结束标记
)

// 设置标签
func SetNameSpace(ns, start, end string) {
	if ns != "" && start != "" && (start != end) {
		defaultNameSpace = ns
		defaultTagStart = start
		defaultTagEnd = end
	}
}

var (
	reNameSpace = regexp.MustCompile("[a-zA-Z0-9]") // 标签名称规则
	reTback     = regexp.MustCompile("[\\/ \t\r\n]")
)

// Errors
var (
	errNameSpaceInvalid        = errors.New("gktemplate namespace is invalid")
	errTagStartOrTagEndInvalid = errors.New("tagstart or tagend is invalid")
	errSourceStringInvalid     = errors.New("source string is invalid")
	errNoneFileInDir           = errors.New("template file not exists in dir")
)

// Data结构体缩写
type D map[string]interface{}

// GKTag 标记的数据结构描述
type GKTag struct {
	IsReplace  bool            // 是否替换
	TagName    string          // 标签名称
	InnerText  []rune          // 内部文本
	StartPos   int             // 标签开始位置
	EndPos     int             // 标签结束位置
	CAttribute *attr.Attribute // 属性结构
	TagValue   string          // 标签值
	TagID      int             // 标签ID
}

// GetTagName()的简写
func (gktag *GKTag) GetName() string {
	return gktag.GetTagName()
}

// 获取标签值
func (gktag *GKTag) GetValue() string {
	return gktag.TagValue
}

// 获取标签名称
func (gktag *GKTag) GetTagName() string {
	return strings.ToLower(gktag.TagName)
}

// 获取标签值
func (gktag *GKTag) GetTagValue() string {
	return gktag.TagValue
}

// 属性值是否存在
func (gktag *GKTag) IsAttribute(str string) bool {
	return gktag.CAttribute.IsAttribute(str)
}

// 获取属性
func (gktag *GKTag) GetAttribute(str string) string {
	return gktag.CAttribute.GetAtt(str)
}

// GetAttribute的简写
func (gktag *GKTag) GetAtt(str string) string {
	return gktag.CAttribute.GetAtt(str)
}

// 获取内部文本
func (gktag *GKTag) GetInnerText() []rune {
	return gktag.InnerText
}

// 下面定义一个存储结构体，将解析出来的模板保存下来
type templateStorage struct {
	Items map[string]*GKTemplate // 存储结构
	sync.RWMutex
}

func (as *templateStorage) SetTemplate(k string, v *GKTemplate) {
	as.Lock()
	defer as.Unlock()
	as.Items[k] = v
}

func (as *templateStorage) GetTemplate(k string) *GKTemplate {
	as.RLock()
	defer as.RUnlock()
	if os.Getenv("GKENV") == "dev" {
		return nil
	}
	v := as.Items[k]
	return v
}

// 下面定义一个存储结构体，将读取的模板缓存起来
// 这样就避免重复读取对系统的开销
type templateFileStorage struct {
	Items map[string]*string // 存储结构
	sync.RWMutex
}

func (as *templateFileStorage) SetTemplateFile(k string, v *string) {
	as.Lock()
	defer as.Unlock()
	as.Items[k] = v
}

func (as *templateFileStorage) GetTemplateFile(k string) *string {
	as.RLock()
	defer as.RUnlock()
	if os.Getenv("GKENV") == "dev" {
		return nil
	}
	v := as.Items[k]
	return v
}

var tplStorage templateStorage
var tplFileStorage templateFileStorage

// 处理Tag的函数
type TagLib func(tag *GKTag, data *D) string
type TagFunc func(v *string, args ...interface{}) string

var tagLibs map[string]TagLib   // 模板标签
var tagFuncs map[string]TagFunc // 模板函数

// 初始化模板函数
func init() {
	tplStorage.Items = make(map[string]*GKTemplate)
	tplFileStorage.Items = make(map[string]*string)

	tagLibs = make(map[string]TagLib)
	tagLibs["field"] = TagField
	tagLibs["range"] = TagRange

	tagFuncs = make(map[string]TagFunc)
	tagFuncs["ToUpper"] = FuncToUpper
	tagFuncs["ToLower"] = FuncToLower
}

// 支持模板自定义扩展标签
func ExtLibs(libs *map[string]TagLib) {
	for fname, ff := range *libs {
		_, ok := tagLibs[fname]
		if ok {
			panic(fmt.Sprintf("[GKTemplate]tag:%s exists", fname))
		}
		tagLibs[fname] = ff
	}
}

// 支持模板自定义扩展函数
func ExtFuncs(funcs *map[string]TagFunc) {
	for fname, ff := range *funcs {
		_, ok := tagFuncs[fname]
		if ok {
			panic(fmt.Sprintf("[GKTemplate]func:%s exists", fname))
		}
		tagFuncs[fname] = ff
	}
}

// 一个模板结构体
type GKTemplate struct {
	NameSpace    string
	TagStart     string
	TagEnd       string
	CTags        map[int]*GKTag // 所有标签
	Count        int            // 标签总数 -1:未解析 >0:解析
	SourceString []rune         // 模板字符串
}

// 校验名称和标签
func checkNameSpaceAndTag(tpl *GKTemplate) error {
	if reNameSpace.MatchString(tpl.NameSpace) == false {
		return errNameSpaceInvalid
	}
	if (tpl.TagStart == tpl.TagEnd) || tpl.TagStart == "" || tpl.TagEnd == "" {
		return errTagStartOrTagEndInvalid
	}
	return nil
}

// 判断字符串是否以正斜杠`/`结束
func IsEndOfForwardSlash(str *string) (bool, int) {
	l := len(*str)
	for i := l - 1; i >= 0; i-- {
		s := (*str)[i]
		if s == ' ' || s == '	' {
			continue
		}
		if s == '/' {
			return true, i
		} else {
			return false, -1
		}
	}
	return false, -1
}

// 解析模板
func parseTemplate(tplstr *string, nameSpace, tagStart, tagEnd, cachekey string) (*GKTemplate, error) {
	khash := ""
	h := sha1.New()
	if cachekey == "" {
		h.Write([]byte(*tplstr))
		khash = fmt.Sprintf("%x", h.Sum(nil))
	} else {
		khash = cachekey
	}

	v := tplStorage.GetTemplate(khash)
	if v != nil {
		// 存在缓存则直接返回缓存
		return v, nil
	}

	var gktpl = GKTemplate{}
	gktpl.CTags = make(map[int]*GKTag)
	gktpl.Count = 0
	gktpl.SourceString = []rune(*tplstr)

	if nameSpace == "" {
		gktpl.NameSpace = defaultNameSpace
	} else {
		gktpl.NameSpace = nameSpace
	}

	if tagStart == "" {
		gktpl.TagStart = defaultTagStart
	} else {
		gktpl.TagStart = tagStart
	}

	if tagEnd == "" {
		gktpl.TagEnd = defaultTagEnd
	} else {
		gktpl.TagEnd = tagEnd
	}

	// 校验标签
	err := checkNameSpaceAndTag(&gktpl)
	if err != nil {
		return nil, err
	}

	tagStartWord := gktpl.TagStart        // 标签开始标记，例：<{
	rTagStartWord := []rune(tagStartWord) // rune的标签开始标记
	// lenTagStartWord := len(rTagStartWord)   // 标签开始标记宽度
	firstCharOfStartTag := rTagStartWord[0] // 标签开始的第一个字符，例如：<

	tagEndWord := gktpl.TagEnd          // 标签结束标记，例如：}>
	rtagEndWord := []rune(tagEndWord)   // rune的标签结束标记
	lenTagEndWord := len(rtagEndWord)   // 标签结束标记宽度
	firstCharOfEndTag := rtagEndWord[0] // 标签结束的第一个字符，例如：}

	fullTagStartWord := tagStartWord + gktpl.NameSpace + ":" // 标签完整开始标记，例：<{gk:
	rFullTagStartWord := []rune(fullTagStartWord)            // rune的完整标记
	lenFullTagStartWord := len(rFullTagStartWord)            // 完整标记字符串宽度

	sTagEndWord := tagStartWord + "/" + gktpl.NameSpace + ":" // 内嵌内容结束标记：例：<{/gk:
	rSTagEndWord := []rune(sTagEndWord)                       // rune的内嵌内容结束标记
	lenSTagEndWord := len(rSTagEndWord)                       // 内嵌内容结束标记字符串宽度
	// eTagEndWord := "/" + TagEndWord
	tsLen := len(fullTagStartWord)
	sourceLen := len(gktpl.SourceString)

	sPos := 0 // 标签开始位置
	ePos := 0 // 标签结束位置

	processTag := false       // 是否正在处理标记
	processAttr := false      // 是否正在处理属性
	processInnertext := false // 是否正在处理内嵌文本
	innerTextPos := -1        // 内嵌文本位置，如果-1表示没有，>0则开始收集

	attrPos := -1               // 用户收集属性的位置标记，如果是-1则表示没有属性，如果>0则进行判断处理
	tmpAttr := []rune("")       // 用于收集临时存储的属性
	tmpInnertext := []rune("")  // 用于收集临时存储的内嵌文本
	var tmpCAtt *attr.Attribute // 用于存储临时解析的属性

	if sourceLen <= (tsLen + 3) {
		return nil, errSourceStringInvalid
	}

	for pos, r := range gktpl.SourceString {
		// 对开始标记以及内嵌闭合标记进行判断
		if firstCharOfStartTag == r {

			// 内嵌闭合标记的判断
			if processTag == true {
				// 如果正在处理Tag，再次获取到开始标记，判断是否是结束标记
				if pos+lenSTagEndWord > sourceLen {
					continue
				}

				tmpTag := string(gktpl.SourceString[pos : pos+lenSTagEndWord])
				if tmpTag == sTagEndWord {
					// 内联标记结束，校验是否是正常结束`<{/gk:`
					// 确认内联结束标记的标签名是否一致

					// 取出`<{/gk:`后面到`}>`处中间标记的名称
					i := 0
					tmpTagName := []rune("")
					for {
						if pos+lenSTagEndWord+i > sourceLen {
							break
						}
						tmpC := gktpl.SourceString[pos+lenSTagEndWord+i]

						if tmpC == ' ' || tmpC == '	' {
							i++
							continue
						}
						if tmpC == firstCharOfEndTag {
							// 判断是否以}>结束
							if pos+lenSTagEndWord+i+lenTagEndWord > sourceLen {
								// 超出模板范围
								break
							}
							tt := gktpl.SourceString[pos+lenSTagEndWord+i : pos+lenSTagEndWord+i+lenTagEndWord]
							if string(tt) == tagEndWord {
								ePos = pos + lenSTagEndWord + i + lenTagEndWord
								break
							}
						}
						tmpTagName = append(tmpTagName, tmpC)
						i++
						// 超出标签字符宽度
						if i > TagMaxLen {
							break
						}
					}

					if string(tmpTagName) == tmpCAtt.GetTagName() {

						var gktag = GKTag{
							TagName:    tmpCAtt.GetTagName(),
							CAttribute: tmpCAtt,
							StartPos:   sPos,
							EndPos:     ePos,
							TagID:      gktpl.Count,
							InnerText:  tmpInnertext,
						}
						gktpl.CTags[gktpl.Count] = &gktag
						gktpl.Count++

						processInnertext = false
						processTag = false
						tmpInnertext = []rune("")
						innerTextPos = -1

					} else {
						// 如果两个标签的名称不一致，这里将会报错
						panic(fmt.Sprintf("[GKTemplate]tag character postion %d, '%s' error！", pos, string(tmpTagName)))
					}
				}
			}

			// 开始标记的判断
			if pos+lenFullTagStartWord > sourceLen {
				continue
			}

			tmpTag := string(gktpl.SourceString[pos : pos+lenFullTagStartWord])

			if tmpTag == fullTagStartWord {
				// 标签开始，进行分析
				processTag = true
				processAttr = true
				attrPos = pos + lenFullTagStartWord
				sPos = pos
			}
		}

		// 用于对结束标记`}>`进行判断
		if firstCharOfEndTag == r {
			// 对`}>`标记进行判断
			if pos+lenTagEndWord > sourceLen {
				continue
			}
			tmpTag := string(gktpl.SourceString[pos : pos+lenTagEndWord])
			if tmpTag == tagEndWord && processTag == true && processAttr == true {
				// 这里出现了`}>`有两种情况，如果标签前面非空字符是`/`，则表示结束标记
				// 如果没有非空字符，则表示有一个内嵌文本的标记
				attrPos = -1 // 属性收集结束

				attrStr := string(tmpAttr)
				ok, i := IsEndOfForwardSlash(&attrStr)

				if ok {
					// 标签结束
					attrStr = attrStr[:i] // 去掉/的属性，然后进行属性解析
					tmpCAtt, err := attr.Parse(attrStr)
					if err != nil {
						panic(err)
					}

					var gktag = GKTag{
						TagName:    tmpCAtt.GetTagName(),
						CAttribute: tmpCAtt,
						StartPos:   sPos,
						EndPos:     pos + lenTagEndWord,
						TagID:      gktpl.Count,
					}
					gktpl.CTags[gktpl.Count] = &gktag
					gktpl.Count++
					// 整个标签已经结束
					processTag = false
					sPos = 0
					ePos = 0

				} else {
					// 内嵌标签前半部分结束
					// fmt.Println("attrStr=", attrStr)
					var err error
					tmpCAtt, err = attr.Parse(attrStr)
					if err != nil {
						panic(err)
					}
					processInnertext = true // 开始收集内嵌文本
					innerTextPos = pos + lenTagEndWord
				}

				tmpAttr = []rune("") // 清空属性
				processAttr = false  // 属性获取结束
			}

		}

		if attrPos != -1 {
			if pos >= attrPos {
				tmpAttr = append(tmpAttr, r)
			}
		}

		if innerTextPos != -1 && processInnertext == true {
			if pos >= innerTextPos {
				tmpInnertext = append(tmpInnertext, r)
			}
		}

	}

	tplStorage.SetTemplate(khash, &gktpl)

	return &gktpl, nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// 加载目录中的文件到文件缓存，然后使用Parse方法直接渲染
func LoadDir(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errNoneFileInDir
	}

	for _, f := range matches {
		if isDirectory(f) {
			continue
		}
		_, err := ParseFile(f, nil)
		if err != nil {
			panic(err)
		}
		// 如果是Debug模式开启
		// fmt.Println("[GKTemplate]Load file:", f)
	}
	return nil
}

func Parse(filename string, data D) (string, error) {
	return ParseFile(filename, data)
}

// 解析文件
func ParseFile(filename string, data D) (string, error) {
	// 根据文件名获取存储哈希
	khash := ""
	h := sha1.New()
	h.Write([]byte(filename))
	khash = fmt.Sprintf("%x", h.Sum(nil))

	// 尝试从存储中载入模板
	v := tplFileStorage.GetTemplateFile(khash)
	if v != nil {
		return ParseStringWithNameSpace(v, data, "", "", "", khash)
	}

	// 从文件中载入模板
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	log.Println("load file:", filename)

	tplstr := string(d)

	tplFileStorage.SetTemplateFile(khash, &tplstr)

	return ParseStringWithNameSpace(&tplstr, data, "", "", "", khash)
}

// 解析字符串
func ParseString(tplstr string, data D) (string, error) {
	return ParseStringWithNameSpace(&tplstr, data, "", "", "", "")
}

// 指定标记名称解析字符串
func ParseStringWithNameSpace(tplstr *string, data D, nameSpace, tagStart, tagEnd, cachekey string) (string, error) {
	// 解析模板
	gktp, err := parseTemplate(tplstr, nameSpace, tagStart, tagEnd, cachekey)
	if err != nil {
		return "", err
	}

	// 替换模板转换内容
	for i := 0; i < gktp.Count; i++ {
		taglib, ok := tagLibs[gktp.CTags[i].TagName]
		if ok {
			gktp.CTags[i].IsReplace = true
			gktp.CTags[i].TagValue = taglib(gktp.CTags[i], &data)

			// 处理自定义函数
			tplFunc := gktp.CTags[i].GetAttribute("func")
			if tplFunc != "" {
				// 解析模板函数
				funcName, args, err := attr.FuncParser(tplFunc)
				if err == nil {
					tagfunc, ok := tagFuncs[funcName]
					if ok {
						gktp.CTags[i].TagValue = tagfunc(&gktp.CTags[i].TagValue, args)
					}
				}
			}
		}
	}

	// 这里可以采用协程的方式并发解析模板
	// 这里主要适用于模板标签中含有较多SQL查询、HTTP资源请求的情况
	// 开启协程获取会消耗资源

	// tagFuncChans := make(chan tagFuncChan)

	// go func() {
	// 	for i := 0; i < gktp.Count; i++ {
	// 		tagfunc, ok := tagLibs[gktp.CTags[i].TagName]
	// 		if ok {
	// 			tagFuncChans <- tagFuncChan{
	// 				Idx:     i,
	// 				Tagname: gktp.CTags[i].TagName,
	// 				Result:  tagfunc(gktp.CTags[i], &data),
	// 			}
	// 		}
	// 	}
	// 	close(tagFuncChans)
	// }()

	// for elem := range tagFuncChans {
	// 	gktp.CTags[elem.Idx].IsReplace = true
	// 	gktp.CTags[elem.Idx].TagValue = elem.Result
	// }

	ResultString := ""
	nextTagEnd := 0
	for i := 0; i < gktp.Count; i++ {
		if gktp.CTags[i].GetValue() == "#@Delete@#" {
			gktp.CTags[i].TagValue = ""
		}
		ResultString += string(gktp.SourceString[nextTagEnd : nextTagEnd+gktp.CTags[i].StartPos-nextTagEnd])
		ResultString += gktp.CTags[i].GetValue()
		nextTagEnd = gktp.CTags[i].EndPos
	}
	slen := len(gktp.SourceString)
	if slen > nextTagEnd {
		ResultString += string(gktp.SourceString[nextTagEnd:slen])
	}

	return ResultString, nil
}
