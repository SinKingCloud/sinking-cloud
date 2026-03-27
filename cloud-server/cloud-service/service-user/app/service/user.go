package service

import "service-user/app/model"

// UserService service接口
type UserService interface {
	GetUserInfo(userId int64, refresh bool) (user *model.User, err error) //获取用户信息
}
