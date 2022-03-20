package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUserPayLog 会员用户消费记录
type WebUserPayLog struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Type       int          `gorm:"column:type" json:"type"`
	LogType    int          `gorm:"column:log_type" json:"log_type"`
	Title      string       `gorm:"column:title" json:"title"`
	Content    string       `gorm:"column:content" json:"content"`
	Money      float64      `gorm:"column:money" json:"money"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (WebUserPayLog) TableName() string {
	return "web_user_pay_logs"
}

// BeforeCreate 创建前
func (t *WebUserPayLog) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebUserPayLog) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
