package real_auth

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 实名认证仓储接口
type Interface interface {
	FindOrCreateByUserId(userId int64, tx ...*gorm.DB) (*model.RealAuth, error)
	UpdateById(id int64, data *UpdateRealAuth, tx ...*gorm.DB) error
}

// Repository 实名认证仓储
type Repository struct {
	*repository.Repository[*model.RealAuth]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.RealAuth](db)}
}
