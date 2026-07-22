package task

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Task) error
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
	FindById(id int64) (*model.Task, error)
	SelectAll() ([]*model.Task, error)
	UpdateById(id int64, data *UpdateTask) error
	Select(where *SelectTask, query *page.Query) (*page.Result[*Task], error)
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*Task]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*Task](db)}
}
