package cookie

import (
	"encoding/json"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

type Cookie struct {
	Uin   string
	Skey  string
	Pskey string
}

// CookieToStr 转str
func (s *Service) CookieToStr(data *Cookie) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(str)
}

// StrToCookie 转cookie
func (s *Service) StrToCookie(data string) (c *Cookie) {
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
