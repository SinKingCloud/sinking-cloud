package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

// GetCurrentUserInfo  获取用户信息
func (s *service) GetCurrentUserInfo(id int, refresh bool) *model.WebUser {
	user, _ := service2.GetWebUserServiceIns().FindByIdWithCache(id, refresh)
	return user
}

// WebUser 网站用户表
type WebUser struct {
	Id         int          `json:"id"`
	Account    string       `json:"account"`
	NickName   string       `json:"nick_name"`
	Avatar     string       `json:"avatar"`
	Money      float64      `json:"money"`
	Email      string       `json:"email"`
	Phone      int64        `json:"phone"`
	LoginIp    string       `json:"login_ip"`
	LoginTime  str.DateTime `json:"login_time"`
	Status     int          `json:"status"`
	RealAuth   int          `json:"real_auth"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
