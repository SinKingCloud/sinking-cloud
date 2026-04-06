package role

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 角色仓储接口
type Interface interface {
	FindById(id int64) (*model.Role, error)
	SelectAll() ([]*model.Role, error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
