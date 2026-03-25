package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Permission 权限表
type Permission struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Service    string       `gorm:"column:service" json:"service"`
	Name       string       `gorm:"column:name" json:"name"`
	Code       string       `gorm:"column:code" json:"code"`
	Path       string       `gorm:"column:path" json:"path"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Permission) TableName() string {
	return "cloud_permissions"
}

// BeforeCreate 创建前
func (t *Permission) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Permission) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
