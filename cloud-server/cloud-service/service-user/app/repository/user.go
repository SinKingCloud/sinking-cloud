package repository

import (
	"service-user/app/model"
	"service-user/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// User 用户仓储接口
type User interface {
	FindByIdentifier(identifier string) (*model.User, error)
	CountByPhone(phone string) (int64, error)
	Create(user *model.User) error
}

// NewUser 创建User仓储
func NewUser(db *database.Database) *user.Repository {
	return &user.Repository{Database: db}
}
