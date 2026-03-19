package config

import (
	"service-user/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// service 注入结构
type service struct {
	repositoryConfig repository.Config
	cache            cache.Interface
}

// NewService 实例化service
func NewService(repository repository.Config, cache cache.Interface) *service {
	return &service{
		repositoryConfig: repository,
		cache:            cache,
	}
}
