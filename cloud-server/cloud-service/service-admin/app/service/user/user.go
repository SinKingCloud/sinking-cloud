package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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

// SystemUser 系统用户表
type SystemUser struct {
	Id         int          `json:"id"`
	RoleId     int          `json:"role_id"`
	Account    string       `json:"account"`
	NickName   string       `json:"nick_name"`
	Avatar     string       `json:"avatar"`
	Email      string       `json:"email"`
	Phone      int64        `json:"phone"`
	LoginIp    string       `json:"login_ip"`
	LoginTime  str.DateTime `json:"login_time"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
