package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// UserAuth 用户授权列表
type UserAuth struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	AuthId     int          `gorm:"column:auth_id" json:"auth_id"`
	Domain     string       `gorm:"column:domain" json:"domain"`
	Code       string       `gorm:"column:code" json:"code"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*UserAuth) TableName() string {
	return "user_auths"
}

// BeforeCreate 创建前
func (t *UserAuth) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *UserAuth) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
