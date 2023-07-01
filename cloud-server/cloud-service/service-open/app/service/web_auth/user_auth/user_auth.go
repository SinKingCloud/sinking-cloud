package user_auth

import (
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

// UserAuth 用户授权列表
type UserAuth struct {
	Id             int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId         int          `gorm:"column:user_id" json:"user_id"`
	AuthId         int          `gorm:"column:auth_id" json:"auth_id"`
	Domain         string       `gorm:"column:domain" json:"domain"`
	Code           string       `gorm:"column:code" json:"code"`
	Ip             string       `gorm:"column:ip" json:"ip"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
	Name           string       `gorm:"column:name" json:"name"`
	Describe       string       `gorm:"column:describe" json:"describe"`
	VersionId      int          `gorm:"column:version_id" json:"version_id"`
	ReleaseDownNum int          `gorm:"column:release_down_num" json:"release_down_num"`
	UpdateDownNum  int          `gorm:"column:update_down_num" json:"update_down_num"`
	Price          float64      `gorm:"column:price" json:"price"`
}

const (
	Normal = iota //正常
	Error         //封禁
)
