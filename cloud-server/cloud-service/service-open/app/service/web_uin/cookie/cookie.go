package cookie

import "encoding/json"

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

type Cookie struct {
	Uin   string
	Skey  string
	Pskey string
}

// CookieToStr 转str
func (s *service) CookieToStr(data *Cookie) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(str)
}

// StrToCookie 转cookie
func (s *service) StrToCookie(data string) (c *Cookie) {
	err := json.Unmarshal([]byte(data), &c)
	if err != nil || c == nil {
		return nil
	}
	return c
}

const (
	Normal = iota //正常
	Error         //失效
)
