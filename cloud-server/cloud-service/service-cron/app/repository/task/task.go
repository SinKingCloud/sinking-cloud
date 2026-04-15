package task

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Task) error
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
	FindById(id int64) (*model.Task, error)
	SelectAll() ([]*Task, error)
	UpdateById(id int64, data *UpdateTask) error
	Select(where *SelectTask, orderByField string, orderByType string, page int, pageSize int) (list []*Task, total int64, err error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
