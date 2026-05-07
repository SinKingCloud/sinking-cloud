package user

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 用户仓储接口
type Interface interface {
	ExistsByPhone(phone string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByAccount(account string) (bool, error)
	FindById(id int64) (*model.User, error)
	FindByIdentifier(identifier string) (*model.User, error)
	Create(user *model.User) error
	UpdateById(id int64, data *UpdateUser, tx ...*gorm.DB) error
	Select(where *SelectUser, orderByField string, orderByType string, page int, pageSize int) (list []*User, total int64, err error)
}

// Repository 用户仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
