package mind_type

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 导图分类仓储接口
type Interface interface {
	Create(data *model.MindType, tx ...*gorm.DB) error
	CountByUserId(userId int64, tx ...*gorm.DB) (int64, error)
	SelectByUserId(userId int64, queryPage *page.Query) (*page.Result[*model.MindType], error)
	FindByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (*model.MindType, error)
	FindByNameAndUserId(name string, userId int64, tx ...*gorm.DB) (*model.MindType, error)
	UpdateByIdAndUserId(id int64, userId int64, name string, tx ...*gorm.DB) error
	DeleteByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) error
}

// Repository 导图分类仓储
type Repository struct {
	*repository.Repository[*model.MindType]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.MindType](db)}
}
