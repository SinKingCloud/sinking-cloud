package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppPlugin 系统插件
type AppPlugin struct {
	Id           int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name         string       `gorm:"column:name" json:"name"`
	Describe     string       `gorm:"column:describe" json:"describe"`
	Icon         string       `gorm:"column:icon" json:"icon"`
	Price        float64      `gorm:"column:price" json:"price"`
	FilterUrl    string       `gorm:"column:filter_url" json:"filter_url"`
	ServiceGroup string       `gorm:"column:service_group" json:"service_group"`
	ServiceName  string       `gorm:"column:service_name" json:"service_name"`
	IsFree       int          `gorm:"column:is_free" json:"is_free"`
	Status       int          `gorm:"column:status" json:"status"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime   str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AppPlugin) TableName() string {
	return "app_plugins"
}

// BeforeCreate 创建前
func (t *AppPlugin) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppPlugin) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
