package user

import (
	"service-admin/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Session服务
type service struct {
	repositoryUser repository.User
	cache          cache.Interface
}

// NewService 实例化service
func NewService(repositoryUser repository.User, cache cache.Interface) *service {
	return &service{
		repositoryUser: repositoryUser,
		cache:          cache,
	}
}
