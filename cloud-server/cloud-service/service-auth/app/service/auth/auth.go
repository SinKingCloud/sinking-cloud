package auth

import (
	"service-auth/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Session服务
type service struct {
	repositorySession repository.Session
	cache             cache.Interface
}

// NewService 实例化service
func NewService(repository repository.Session, cache cache.Interface) *service {
	return &service{
		repositorySession: repository,
		cache:             cache,
	}
}
