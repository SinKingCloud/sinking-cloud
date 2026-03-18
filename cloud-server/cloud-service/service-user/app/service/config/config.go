package config

import (
	"service-user/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// service 注入结构
type service struct {
	repo  repository.Config
	cache cache.Interface
}

// NewService 实例化service
func NewService(repo repository.Config, cache cache.Interface) *service {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

// Service service接口
type Service interface {
	GetUserConfigs(userId int64, group string, refresh bool) map[string]string
	SetUserConfig(userId int64, group string, key string, value string) bool
	SetUserConfigs(userId int64, group string, data map[string]string)
}
