package repository

import (
	"service-admin/app/repository/menu_permission"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// MenuPermission 菜单权限仓储接口
type MenuPermission interface {
}

// NewMenuPermission 创建MenuPermission仓储
func NewMenuPermission(db *database.Database) *menu_permission.Repository {
	return &menu_permission.Repository{Database: db}
}
