package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	Id             int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Account        string       `gorm:"column:account" json:"account"`
	Phone          string       `gorm:"column:phone" json:"phone"`
	Email          string       `gorm:"column:email" json:"email"`
	Password       string       `gorm:"column:password" json:"password"`
	Avatar         string       `gorm:"column:avatar" json:"avatar"`
	NickName       string       `gorm:"column:nick_name" json:"nick_name"`
	RealAuthStatus int          `gorm:"column:real_auth_status" json:"real_auth_status"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*User) TableName() string {
	return "cloud_users"
}

// BeforeCreate 创建前
func (t *User) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *User) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
