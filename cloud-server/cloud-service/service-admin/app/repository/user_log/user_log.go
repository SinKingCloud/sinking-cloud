package user_log

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.UserLog, tx ...*gorm.DB) error
	Select(where *SelectUserLog, queryPage *page.Query) (*page.Result[*model.UserLog], error)
	DeleteByUserId(userId int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.UserLog]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.UserLog](db)}
}
