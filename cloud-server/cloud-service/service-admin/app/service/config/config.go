package config

import (
	"service-admin/app/repository/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// service 注入结构
type service struct {
	repositoryConfig config.Interface
	cache            cache.Interface
}

// NewService 实例化service
func NewService(repository config.Interface, cache cache.Interface) *service {
	return &service{
		repositoryConfig: repository,
		cache:            cache,
	}
}
