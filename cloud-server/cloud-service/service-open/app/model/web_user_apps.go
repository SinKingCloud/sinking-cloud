package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUserApp 用户应用列表
type WebUserApp struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Key        string       `gorm:"column:key" json:"key"`
	Name       string       `gorm:"column:name" json:"name"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (WebUserApp) TableName() string {
	return "web_user_apps"
}

// BeforeCreate 创建前
func (t *WebUserApp) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebUserApp) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
