package role

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 角色仓储接口
type Interface interface {
	FindById(id int64) (*model.Role, error)
	SelectByIds(ids []int64) ([]*model.Role, error)
	SelectAll() ([]*model.Role, error)
	Select(where *SelectRole, orderByField string, orderByType string, page int, pageSize int) (list []*model.Role, total int64, err error)
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
