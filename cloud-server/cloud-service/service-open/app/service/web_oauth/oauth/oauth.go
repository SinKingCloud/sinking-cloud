package oauth

import "sync"

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

const (
	Normal = iota //正常
	Stop          //停止
	Error         //封禁
)
