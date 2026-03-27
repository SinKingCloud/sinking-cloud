package user

import (
	"service-user/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Session服务
type service struct {
	repositoryUser user.Interface
	cache          cache.Interface
}

// NewService 实例化service
func NewService(repositoryUser user.Interface, cache cache.Interface) *service {
	return &service{
		repositoryUser: repositoryUser,
		cache:          cache,
	}
}
