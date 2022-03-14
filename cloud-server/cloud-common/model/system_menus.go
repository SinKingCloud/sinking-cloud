package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemMenu 系统菜单表
type SystemMenu struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Pid        int          `gorm:"column:pid" json:"pid"`
	Title      string       `gorm:"column:title" json:"title"`
	Type       int          `gorm:"column:type" json:"type"`
	Accesses   string       `gorm:"column:accesses" json:"accesses"`
	Key        string       `gorm:"column:key" json:"key"`
	Path       string       `gorm:"column:path" json:"path"`
	Icon       string       `gorm:"column:icon" json:"icon"`
	Status     int          `gorm:"column:status" json:"status"`
	Sort       int          `gorm:"column:sort" json:"sort"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (SystemMenu) TableName() string {
	return "system_menus"
}

// BeforeCreate 创建前
func (t *SystemMenu) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *SystemMenu) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
