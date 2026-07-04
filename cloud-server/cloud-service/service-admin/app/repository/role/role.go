package role

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 角色仓储接口
type Interface interface {
	FindById(id int64) (*model.Role, error)
	Create(data *model.Role, tx ...*gorm.DB) error
	UpdateById(id int64, data *UpdateRole, tx ...*gorm.DB) error
	SelectByIds(ids []int64) ([]*model.Role, error)
	SelectIdNameMap() (map[int64]string, error)
	Select(where *SelectRole, page *page.Query) (list []*model.Role, total int64, err error)
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.Role{})}
}
