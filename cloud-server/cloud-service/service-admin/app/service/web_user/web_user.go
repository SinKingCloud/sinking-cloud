package web_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// GetCurrentUserInfo  获取登录用户最新信息
func (s *Service) GetCurrentUserInfo(id int, refresh bool) *model.WebUser {
	user, _ := service2.GetWebUserServiceIns().FindByIdWithCache(id, refresh)
	return user
}

// WebUser 会员用户表
type WebUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Account    string       `gorm:"column:account" json:"account"`
	NickName   string       `gorm:"column:nick_name" json:"nick_name"`
	Avatar     string       `gorm:"column:avatar" json:"avatar"`
	Email      string       `gorm:"column:email" json:"email"`
	Phone      int64        `gorm:"column:phone" json:"phone"`
	Money      float64      `gorm:"column:money" json:"money"`
	LoginIp    string       `gorm:"column:login_ip" json:"login_ip"`
	LoginTime  str.DateTime `gorm:"column:login_time" json:"login_time"`
	Status     int          `gorm:"column:status" json:"status"`
	RealAuth   int          `gorm:"column:real_auth" json:"real_auth"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
