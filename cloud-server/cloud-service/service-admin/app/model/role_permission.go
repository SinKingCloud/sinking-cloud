package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleId         int64        `gorm:"column:role_id;primaryKey" json:"role_id"`
	PermissionCode string       `gorm:"column:permission_code;primaryKey" json:"permission_code"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*RolePermission) TableName() string {
	return "cloud_role_permissions"
}

// BeforeCreate 创建前
func (t *RolePermission) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
