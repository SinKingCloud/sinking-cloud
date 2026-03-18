package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Notice 通知表
type Notice struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	Title      string       `gorm:"column:title" json:"title"`
	Content    string       `gorm:"column:content" json:"content"`
	Status     int          `gorm:"column:status" json:"status"`
	ReadTime   str.DateTime `gorm:"column:read_time" json:"read_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Notice) TableName() string {
	return "cloud_notices"
}

// BeforeCreate 创建前
func (t *Notice) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Notice) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
