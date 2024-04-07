package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AuthLog 记录
type AuthLog struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	AuthId     int          `gorm:"column:auth_id" json:"auth_id"`
	Content    string       `gorm:"column:content" json:"content"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AuthLog) TableName() string {
	return "auth_logs"
}

// BeforeCreate 创建前
func (t *AuthLog) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AuthLog) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
