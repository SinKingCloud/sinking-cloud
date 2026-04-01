package menu

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 菜单仓储接口
type Interface interface {
	SelectByIdsWithStatusAndType(ids []int64, status int, types []int) ([]*model.Menu, error)
	SelectByParentIdWithStatusAndType(parentId int64, status int, types int) ([]*model.Menu, error)
	FindByCode(code string) (*model.Menu, error)
}

// Repository 用户仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建User仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
