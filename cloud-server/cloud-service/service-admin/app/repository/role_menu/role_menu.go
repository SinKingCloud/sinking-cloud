package role_menu

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 角色菜单仓储接口
type Interface interface {
	Creates(data []*model.RoleMenu, tx ...*gorm.DB) error
	SelectMenuIdsByRoleId(roleId int64) ([]int64, error)
	SelectRoleIdsByMenuIds(menuIds []int64) ([]int64, error)
	DeleteByMenuIds(menuIds []int64, tx ...*gorm.DB) error
	DeleteByRoleIds(roleIds []int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
