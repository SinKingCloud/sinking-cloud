package menu_permission

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 菜单权限仓储接口
type Interface interface {
	Creates(data []*model.MenuPermission, tx ...*gorm.DB) error
	SelectPermissionCodesByMenuId(menuId int64) ([]string, error)
	SelectByMenuIds(menuIds []int64) ([]*model.MenuPermission, error)
	DeleteByMenuIds(menuIds []int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
