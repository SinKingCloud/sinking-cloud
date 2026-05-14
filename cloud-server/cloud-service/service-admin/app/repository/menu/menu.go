package menu

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 菜单仓储接口
type Interface interface {
	Create(data *model.Menu, tx ...*gorm.DB) error
	UpdateById(id int64, data *UpdateMenu, tx ...*gorm.DB) error
	SelectByIds(ids []int64) ([]*model.Menu, error)
	SelectByIdsWithStatusAndType(ids []int64, status int, types []int) ([]*model.Menu, error)
	SelectAll() ([]*model.Menu, error)
	SelectWithStatus(status int) ([]*model.Menu, error)
	SelectWithType(types []int) ([]*model.Menu, error)
	Select(where *SelectMenu, orderByField string, orderByType string, page int, pageSize int) ([]*model.Menu, int64, error)
	SelectByParentIdWithStatusAndType(parentId int64, status int, types int) ([]*model.Menu, error)
	FindByCode(code string) (*model.Menu, error)
	FindById(id int64) (*model.Menu, error)
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
}

// Repository 用户仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建User仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
