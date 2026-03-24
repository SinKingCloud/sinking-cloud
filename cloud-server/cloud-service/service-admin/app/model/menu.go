package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Menu 菜单表
type Menu struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	ParentId   int64        `gorm:"column:parent_id" json:"parent_id"`
	Code       string       `gorm:"column:code" json:"code"`
	Name       string       `gorm:"column:name" json:"name"`
	Type       int          `gorm:"column:type" json:"type"`
	Icon       string       `gorm:"column:icon" json:"icon"`
	Sort       int          `gorm:"column:sort" json:"sort"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Menu) TableName() string {
	return "cloud_menus"
}

// BeforeCreate 创建前
func (t *Menu) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Menu) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
