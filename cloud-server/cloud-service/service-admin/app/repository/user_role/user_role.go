package user_role

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 用户仓储接口
type Interface interface {
	SelectRoleIdsByUserId(userId int64) ([]int64, error)
	SelectByUserIds(userIds []int64) ([]*model.UserRole, error)
	CountByRoleIds(roleIds []int64) (int64, error)
	Creates(data []*model.UserRole, tx ...*gorm.DB) error
	DeleteByUserId(userId int64, tx ...*gorm.DB) error
	DeleteByUserIds(userIds []int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.UserRole]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.UserRole](db)}
}
