package user

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 用户仓储接口
type Interface interface {
	ExistsByPhone(phone string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByAccount(account string) (bool, error)
	FindById(id int64) (*model.User, error)
	FindByIdentifier(identifier string) (*model.User, error)
	Select(where *SelectUser, page *page.Query) ([]*User, int64, error)
	Create(user *model.User) error
	UpdateById(id int64, data *UpdateUser) error
	DeleteById(id int64, tx ...*gorm.DB) error
}

// Repository 仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.User{})}
}
