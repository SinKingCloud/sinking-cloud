package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUserOrder 会员用户订单记录
type WebUserOrder struct {
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
func (WebUserOrder) TableName() string {
	return "web_user_orders"
}

// BeforeCreate 创建前
func (t *WebUserOrder) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebUserOrder) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
