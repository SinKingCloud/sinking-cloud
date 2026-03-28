package user

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 用户仓储接口
type Interface interface {
	CountByPhone(phone string) (int64, error)
	CountByEmail(email string) (int64, error)
	CountByAccount(account string) (int64, error)
	FindById(id int64) (*model.User, error)
	FindByIdentifier(identifier string) (*model.User, error)
	Create(user *model.User) error
	UpdateById(id int64, data *UpdateUser) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
