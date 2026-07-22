package log

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Log) error
	DeleteByTaskIds(taskIds []int64, tx ...*gorm.DB) error
	Select(where *SelectLog, query *page.Query) (*page.Result[*model.Log], error)
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.Log]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.Log](db)}
}
