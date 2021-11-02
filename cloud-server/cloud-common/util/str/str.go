package str

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// StringTool 字符串工具
type StringTool struct {
	string string
}

// NetStringTool 实例化工具类
func NetStringTool(str string) *StringTool {
	return &StringTool{string: str}
}

func (stringTool *StringTool) GetBetween(start string, end string) string {
	n := strings.Index(stringTool.string, start)
	if n == -1 {
		n = 0
	}
	stringTool.string = string([]byte(stringTool.string)[n:])
	m := strings.Index(stringTool.string, end)
	if m == -1 {
		m = len(stringTool.string)
	}
	stringTool.string = string([]byte(stringTool.string)[len(start):m])
	return stringTool.string
}

// Md5 获取字符串md5值
func (stringTool *StringTool) Md5() string {
	h := md5.New()
	h.Write([]byte(stringTool.string))
	return hex.EncodeToString(h.Sum(nil))
}
