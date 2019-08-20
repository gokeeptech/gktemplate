// Copyright 2019 The GoKeep Authors. All rights reserved.
// license that can be found in the LICENSE file.

// 属性解析器
package attribute

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"regexp"
	"strings"
	"sync"
	// "unicode/utf8"
)

const SourceMaxSize = 1024 // 解析属性标记最大尺寸
const CharToLow = true     // 是否将属性名称统一转换成小写

var (
	reSplit = regexp.MustCompile("[ \t\r\n]{1,}")
)

// Errors
var (
	errStringEmpty      = errors.New("attribute parse string is empty")
	errStringOutExceeds = errors.New("string exceeds maximum limit")
)

// 定义一个属性结构体
type Attribute struct {
	Count int
	Items map[string]string
}

// 根据名称获取属性值
func (att *Attribute) GetAtt(str string) string {
	if str == "" {
		return ""
	}
	if att.Items[str] != "" {
		return att.Items[str]
	} else {
		return ""
	}
}

// GetAtt的全写
func (att *Attribute) GetAttribute(str string) string {
	return att.GetAtt(str)
}

// 某个属性是否存在
func (att *Attribute) IsAttribute(str string) bool {
	if att.Items[str] != "" {
		return true
	} else {
		return false
	}
}

// 获取标签名称
func (att *Attribute) GetTagName() string {
	return att.GetAtt("tagname")
}

// 获取属性总数
func (att *Attribute) GetCount() int {
	return att.Count + 1
}

// 下面定义一个存储结构体，将解析出来的属性进行缓存
type attributeStorage struct {
	Items map[string]*Attribute // 存储结构
	sync.RWMutex
}

func (as *attributeStorage) SetAttribute(k string, v *Attribute) {
	as.Lock()
	defer as.Unlock()
	as.Items[k] = v
}

func (as *attributeStorage) GetAttribute(k string) *Attribute {
	as.RLock()
	defer as.RUnlock()
	v := as.Items[k]
	return v
}

var attStorage attributeStorage
var h hash.Hash

func init() {
	attStorage.Items = make(map[string]*Attribute)
	h = sha1.New()
}

// 从字符串中解析出属性
func Parse(attStr string) (*Attribute, error) {
	if attStr == "" {
		return nil, errStringEmpty
	}

	if len(attStr) > SourceMaxSize {
		return nil, errStringOutExceeds
	}
	var result Attribute

	h.Reset()
	h.Write([]byte(attStr))
	khash := fmt.Sprintf("%x", h.Sum(nil))

	v := attStorage.GetAttribute(khash)
	if v != nil {
		// 存在缓存则直接返回缓存
		return v, nil
	}

	// 处理字符串，将属性字符串格式化为单个空格间隔字符
	attStr = strings.TrimSpace(reSplit.ReplaceAllString(attStr, " "))

	// 初始化默认结果
	result.Count = 0
	result.Items = make(map[string]string)

	var hasTag = false
	var gkStart = -1
	var gkTag = rune(' ')
	var tmpAtt = []rune("")   // 临时存储的属性名
	var tmpValue = []rune("") // 临时存储的值
	var attName = ""
	var attValue = ""
	var preChar rune

	// 进行标签解析处理
	for _, r := range attStr {

		// 先解析出tagname
		if hasTag == false && r != rune(' ') {
			result.Items["tagname"] = result.Items["tagname"] + string(r)
		}

		if hasTag == false && r == rune(' ') {
			hasTag = true
			if CharToLow == true {
				result.Items["tagname"] = strings.TrimSpace(strings.ToLower(result.Items["tagname"]))
			}
		}

		if hasTag == true {
			// 解析属性
			if gkStart == -1 {
				if r != rune('=') {
					tmpAtt = append(tmpAtt, r)
				} else {
					attName = string(tmpAtt)
					if CharToLow {
						attName = strings.TrimSpace(strings.ToLower(attName))
					}
					gkStart = 0
				}
			} else if gkStart == 0 {
				switch r {
				case rune(' '):
					continue
					break
				case rune('\''):
					gkTag = rune('\'')
					gkStart = 1
					break
				case rune('"'):
					gkTag = rune('"')
					gkStart = 1
					break
				case rune('`'):
					gkTag = rune('`')
					gkStart = 1
					break
				default:
					tmpValue = append(tmpValue, r)
					gkTag = rune(' ')
					gkStart = 1
					break
				}
			} else if gkStart == 1 {
				if r == gkTag && preChar != rune('\\') {
					result.Count++
					attValue = string(tmpValue)
					result.Items[attName] = attValue
					tmpAtt = []rune("")
					tmpValue = []rune("")
					attName = ""
					attValue = ""
					gkStart = -1
				} else {
					tmpValue = append(tmpValue, r)
				}
			}
			preChar = r
		}
	}

	attStorage.SetAttribute(khash, &result)

	return &result, nil
}
