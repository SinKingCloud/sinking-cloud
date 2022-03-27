package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebAppConfig 应用配置
type WebAppConfig struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	IdentId    int          `gorm:"column:ident_id" json:"ident_id"`
	Key        string       `gorm:"column:key" json:"key"`
	Value      string       `gorm:"column:value" json:"value"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (WebAppConfig) TableName() string {
	return "web_app_configs"
}

// BeforeCreate 创建前
func (t *WebAppConfig) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebAppConfig) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
