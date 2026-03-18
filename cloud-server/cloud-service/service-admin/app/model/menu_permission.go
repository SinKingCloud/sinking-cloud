package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// MenuPermission 菜单权限表
type MenuPermission struct {
	Id             int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	MenuId         int64        `gorm:"column:menu_id" json:"menu_id"`
	PermissionCode string       `gorm:"column:permission_code" json:"permission_code"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*MenuPermission) TableName() string {
	return "cloud_menu_permissions"
}

// BeforeCreate 创建前
func (t *MenuPermission) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
