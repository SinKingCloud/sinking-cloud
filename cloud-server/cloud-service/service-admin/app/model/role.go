package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Role 角色表
type Role struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name       string       `gorm:"column:name" json:"name"`
	Code       string       `gorm:"column:code" json:"code"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Role) TableName() string {
	return "cloud_roles"
}

// BeforeCreate 创建前
func (t *Role) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Role) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
