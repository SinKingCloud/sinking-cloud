package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// GetCurrentUserInfo  获取登录用户最新信息
func (s *service) GetCurrentUserInfo(id int, refresh bool) *model.SystemUser {
	user, _ := service2.GetSystemUserServiceIns().FindByIdWithCache(id, refresh)
	return user
}
