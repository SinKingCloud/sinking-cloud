package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UserRole 用户角色关联表
type UserRole struct {
	UserId     int64        `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleId     int64        `gorm:"column:role_id;primaryKey" json:"role_id"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*UserRole) TableName() string {
	return "cloud_user_roles"
}

// BeforeCreate 创建前
func (t *UserRole) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *UserRole) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
