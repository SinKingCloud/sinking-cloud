package user

import (
	"service-user/app/model"
	"service-user/app/repository/user"
	"service-user/app/service/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	GetUserInfo(userId int64, refresh bool) (user *model.User, err error) //获取用户信息
	UpdateById(id int64, data *user.UpdateUser) error                     //更新用户信息
}

// Session服务
type service struct {
	repositoryUser user.Interface
	config         config.Service
	cache          cache.Interface
}

// NewService 实例化service
func NewService(repositoryUser user.Interface, config config.Service, cache cache.Interface) *service {
	return &service{
		repositoryUser: repositoryUser,
		config:         config,
		cache:          cache,
	}
}
