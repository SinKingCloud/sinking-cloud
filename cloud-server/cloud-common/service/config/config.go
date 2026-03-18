package config

import (
	"sync"
)

// Service 单例对象
type Service struct {
	// systemConfigLock 系统配置缓存锁
	systemConfigLock sync.RWMutex
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
