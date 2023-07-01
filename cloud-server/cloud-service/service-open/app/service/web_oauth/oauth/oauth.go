package oauth

import "sync"

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

const (
	Normal = iota //正常
	Stop          //停止
	Error         //封禁
)
