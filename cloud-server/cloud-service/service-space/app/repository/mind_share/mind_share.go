package mind_share

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 导图分享仓储接口
type Interface interface {
	Create(data *model.MindShare, tx ...*gorm.DB) error
	FindByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) (*model.MindShare, error)
	FindByCode(code string, tx ...*gorm.DB) (*model.MindShare, error)
	SelectByDataId(dataId int64, queryPage *page.Query) (*page.Result[*model.MindShare], error)
	SelectValidDataIds(dataIds []int64) ([]int64, error)
	UpdateByIdAndDataId(id int64, dataId int64, data *UpdateMindShare, tx ...*gorm.DB) error
	DeleteByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) error
	DeleteByDataIds(dataIds []int64, tx ...*gorm.DB) error
}

// Repository 导图分享仓储
type Repository struct {
	*repository.Repository[*model.MindShare]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.MindShare](db)}
}
