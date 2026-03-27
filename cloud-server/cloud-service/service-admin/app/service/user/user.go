package user

import (
	"service-admin/app/model"
	"service-admin/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	GetUserInfo(userId int64, refresh bool) (user *model.User, err error) //获取用户信息
}

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
