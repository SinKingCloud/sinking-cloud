package mind_data

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 导图数据仓储接口
type Interface interface {
	Create(data *model.MindData, tx ...*gorm.DB) error
	CountByUserId(userId int64, tx ...*gorm.DB) (int64, error)
	FindById(id int64) (*model.MindData, error)
	FindByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (*model.MindData, error)
	ExistsByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (bool, error)
	Select(where *SelectMindData, queryPage *page.Query) (*page.Result[*MindData], error)
	UpdateByIdAndUserId(id int64, userId int64, data *UpdateMindData, tx ...*gorm.DB) error
	UpdateByTypeIdAndUserId(typeId int64, userId int64, data *UpdateMindData, tx ...*gorm.DB) error
	DeleteByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) error
}

// Repository 导图数据仓储
type Repository struct {
	*repository.Repository[*MindData]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*MindData](db)}
}
