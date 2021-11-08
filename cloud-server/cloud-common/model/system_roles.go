package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemRole 系统角色表
type SystemRole struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name       string       `gorm:"column:name" json:"name"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Accesses   string       `gorm:"column:accesses" json:"accesses"`
	Menus      string       `gorm:"column:menus" json:"menus"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int          `gorm:"column:is_delete" json:"is_delete"`
}

// TableName 获取表名
func (SystemRole) TableName() string {
	return "system_roles"
}

// BeforeCreate 创建前
func (t *SystemRole) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// BeforeUpdate 更新前
func (t *SystemRole) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
