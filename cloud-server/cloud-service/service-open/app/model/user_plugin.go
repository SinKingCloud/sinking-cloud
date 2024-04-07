package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// UserPlugin 用户插件
type UserPlugin struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	PluginId   int          `gorm:"column:plugin_id" json:"plugin_id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Status     int          `gorm:"column:status" json:"status"`
	ExpireTime str.DateTime `gorm:"column:expire_time" json:"expire_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*UserPlugin) TableName() string {
	return "user_plugins"
}

// BeforeCreate 创建前
func (t *UserPlugin) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *UserPlugin) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
