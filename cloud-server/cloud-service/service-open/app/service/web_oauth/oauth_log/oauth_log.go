package oauth_log

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
	Normal  = iota //未授权
	Success        //已授权
	Expire         //已失效
)
