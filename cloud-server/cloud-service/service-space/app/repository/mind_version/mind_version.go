package mind_version

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 导图版本仓储接口
type Interface interface {
	Create(data *model.MindVersion, tx ...*gorm.DB) error
	FindByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) (*model.MindVersion, error)
	FindLatestByDataId(dataId int64, tx ...*gorm.DB) (*model.MindVersion, error)
	FindLatestIdByDataId(dataId int64) (int64, error)
	Select(where *SelectMindVersion, queryPage *page.Query) (list []*MindVersion, total int64, err error)
	UpdateByIdAndDataId(id int64, dataId int64, data *UpdateMindVersion, tx ...*gorm.DB) error
	DeleteByDataIds(dataIds []int64, tx ...*gorm.DB) error
	DeleteOldByDataIdAndType(dataId int64, versionType int, limit int, tx ...*gorm.DB) error
}

// Repository 导图版本仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.MindVersion{})}
}
