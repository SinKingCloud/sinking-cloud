package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Channel 支付通道
type Channel struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       string       `gorm:"column:type" json:"type"`
	Plugin     int          `gorm:"column:plugin" json:"plugin"`
	Name       string       `gorm:"column:name" json:"name"`
	Config     string       `gorm:"column:config" json:"config"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Channel) TableName() string {
	return "cloud_channels"
}

// BeforeCreate 创建前
func (t *Channel) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Channel) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
