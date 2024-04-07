package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebNotice 网站通知
type WebNotice struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Place      string       `gorm:"column:place" json:"place"`
	Title      string       `gorm:"column:title" json:"title"`
	Content    string       `gorm:"column:content" json:"content"`
	Status     int          `gorm:"column:status" json:"status"`
	Sort       int          `gorm:"column:sort" json:"sort"`
	LookNum    int          `gorm:"column:look_num" json:"look_num"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*WebNotice) TableName() string {
	return "web_notices"
}

// BeforeCreate 创建前
func (t *WebNotice) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebNotice) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
