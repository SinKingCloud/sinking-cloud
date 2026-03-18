package config

import (
	"service-admin/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// service 注入结构
type service struct {
	repository repository.Config
	cache      cache.Interface
}

// NewService 实例化service
func NewService(repo repository.Config, cache cache.Interface) *service {
	return &service{
		repository: repo,
		cache:      cache,
	}
}

// Service service接口
type Service interface {
	GetUserConfigs(userId int64, group string, refresh bool) map[string]string
	SetUserConfig(userId int64, group string, key string, value string) bool
	SetUserConfigs(userId int64, group string, data map[string]string)
	GetSysConfigs(group string, refresh bool) map[string]string
	SetSysConfig(key string, value string) bool
	SetSysConfigs(data map[string]string)
}
