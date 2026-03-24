package repository

import (
	"service-admin/app/repository/role_permission"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// RolePermission 角色权限仓储接口
type RolePermission interface {
	FindPermissionCodesByRoleId(roleId int64) ([]string, error)
}

// NewRolePermission 创建RolePermission仓储
func NewRolePermission(db *database.Database) *role_permission.Repository {
	return &role_permission.Repository{Database: db}
}
