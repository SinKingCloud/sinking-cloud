package permission

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 权限仓储接口
type Interface interface {
	FindByPath(path string) (*model.Permission, error)
	FindByCode(code string) (*model.Permission, error)
	Create(perm *model.Permission) error
	UpdateByCode(code string, perm *UpdatePermission) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
