package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinCookie QQCookie
type AppUinCookie struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Uin        string       `gorm:"column:uin" json:"uin"`
	Domain     string       `gorm:"column:domain" json:"domain"`
	Cookie     string       `gorm:"column:cookie" json:"cookie"`
	Status     int          `gorm:"column:status" json:"status"`
	CheckTime  str.DateTime `gorm:"column:check_time" json:"check_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AppUinCookie) TableName() string {
	return "app_uin_cookies"
}

// BeforeCreate 创建前
func (t *AppUinCookie) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinCookie) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
