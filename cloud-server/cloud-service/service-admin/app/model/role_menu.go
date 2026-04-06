package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// RoleMenu 角色菜单表
type RoleMenu struct {
	RoleId     int64        `gorm:"column:role_id;primaryKey" json:"role_id"`
	MenuId     int64        `gorm:"column:menu_id;primaryKey" json:"menu_id"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*RoleMenu) TableName() string {
	return "cloud_role_menus"
}

// BeforeCreate 创建前
func (t *RoleMenu) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *RoleMenu) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
