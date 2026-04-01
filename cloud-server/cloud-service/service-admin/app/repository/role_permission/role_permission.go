package role_permission

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 角色权限仓储接口
type Interface interface {
	SelectPermissionCodesByRoleId(roleId int64) ([]string, error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
