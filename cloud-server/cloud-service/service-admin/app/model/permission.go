package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Permission 权限表
type Permission struct {
	Service    string       `gorm:"column:service" json:"service"`
	Name       string       `gorm:"column:name" json:"name"`
	Code       string       `gorm:"column:code" json:"code"`
	Path       string       `gorm:"column:path" json:"path"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*Permission) TableName() string {
	return "cloud_permissions"
}

// BeforeCreate 创建前
func (t *Permission) BeforeCreate(_ interface{}) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
