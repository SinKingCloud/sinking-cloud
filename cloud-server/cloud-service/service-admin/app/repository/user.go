package repository

import (
	"service-admin/app/model"
	"service-admin/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// User 用户仓储接口
type User interface {
	FindById(id int64) (*model.User, error)
	FindByIdentifier(identifier string) (*model.User, error)
	Create(user *model.User) error
}

// NewUser 创建User仓储
func NewUser(db *database.Database) *user.Repository {
	return &user.Repository{Database: db}
}
