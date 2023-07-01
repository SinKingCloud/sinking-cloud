package oauth_log

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
	Normal  = iota //未授权
	Success        //已授权
	Expire         //已失效
)
