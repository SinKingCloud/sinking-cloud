package permission

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 权限仓储接口
type Interface interface {
	FindByPath(path string) (*model.Permission, error)
	FindByCode(code string) (*model.Permission, error)
	Create(perm *model.Permission) error
	UpdateByCode(code string, perm *UpdatePermission) error
	UpdateByIds(ids []int64, perm *UpdatePermission) error
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
	SelectCodeNameMap() (map[string]string, error)
	SelectByCodes(codes []string) ([]*model.Permission, error)
	SelectByIds(ids []int64) ([]*model.Permission, error)
	Select(where *SelectPermission, queryPage *page.Query) (*page.Result[*model.Permission], error)
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.Permission]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.Permission](db)}
}
