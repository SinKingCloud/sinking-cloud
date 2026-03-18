package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Session 会话表
type Session struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	UserType   int          `gorm:"column:user_type" json:"user_type"`
	Token      string       `gorm:"column:token" json:"token"`
	Device     string       `gorm:"column:device" json:"device"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Location   string       `gorm:"column:location" json:"location"`
	ExpireTime str.DateTime `gorm:"column:expire_time" json:"expire_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Session) TableName() string {
	return "cloud_sessions"
}

// BeforeCreate 创建前
func (t *Session) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Session) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
