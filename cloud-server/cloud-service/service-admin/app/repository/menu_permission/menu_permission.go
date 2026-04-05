package menu_permission

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 菜单权限仓储接口
type Interface interface {
	SelectPermissionCodesByMenuId(menuId int64) ([]string, error)
	SelectByMenuIds(menuIds []int64) ([]*model.MenuPermission, error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
