package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUserConfig 会员用户配置表
type WebUserConfig struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Key        string       `gorm:"column:key" json:"key"`
	Value      string       `gorm:"column:value" json:"value"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int          `gorm:"column:is_delete" json:"is_delete"`
}

// TableName 获取表名
func (WebUserConfig) TableName() string {
	return "web_user_configs"
}

// BeforeCreate 创建前
func (t *WebUserConfig) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// BeforeUpdate 更新前
func (t *WebUserConfig) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
