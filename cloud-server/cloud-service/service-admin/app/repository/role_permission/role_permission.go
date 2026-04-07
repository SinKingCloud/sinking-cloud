package role_permission

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 角色权限仓储接口
type Interface interface {
	Creates(data []*model.RolePermission, tx ...*gorm.DB) error
	SelectPermissionCodesByRoleId(roleId int64) ([]string, error)
	SelectRoleIdsByPermissionCodes(permissionCodes []string) ([]int64, error)
	DeleteByPermissionCodes(permissionCodes []string, tx ...*gorm.DB) error
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
