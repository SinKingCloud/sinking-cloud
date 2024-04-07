package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemUser 系统用户表
type SystemUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	RoleId     int          `gorm:"column:role_id" json:"role_id"`
	Account    string       `gorm:"column:account" json:"account"`
	Password   string       `gorm:"column:password" json:"-"`
	NickName   string       `gorm:"column:nick_name" json:"nick_name"`
	Avatar     string       `gorm:"column:avatar" json:"avatar"`
	Email      string       `gorm:"column:email" json:"email"`
	Phone      int64        `gorm:"column:phone" json:"phone"`
	LoginToken string       `gorm:"column:login_token" json:"login_token"`
	LoginIp    string       `gorm:"column:login_ip" json:"login_ip"`
	LoginTime  str.DateTime `gorm:"column:login_time" json:"login_time"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*SystemUser) TableName() string {
	return "system_users"
}

// BeforeCreate 创建前
func (t *SystemUser) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *SystemUser) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
