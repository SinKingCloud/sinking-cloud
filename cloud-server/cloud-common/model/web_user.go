package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUser 会员用户表
type WebUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
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
func (WebUser) TableName() string {
	return "web_users"
}

// BeforeCreate 创建前
func (t *WebUser) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebUser) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
