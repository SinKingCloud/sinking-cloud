package repository

import (
	"service-admin/app/repository/role_menu"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// RoleMenu 角色菜单仓储接口
type RoleMenu interface {
}

// NewRoleMenu 创建RoleMenu仓储
func NewRoleMenu(db *database.Database) *role_menu.Repository {
	return &role_menu.Repository{Database: db}
}
