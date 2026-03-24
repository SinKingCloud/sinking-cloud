package repository

import (
	"service-admin/app/repository/user_role"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// UserRole 用户仓储接口
type UserRole interface {
	FindRoleIdsByUserId(userId int64) ([]int64, error)
}

// NewUserRole 创建UserRole仓储
func NewUserRole(db *database.Database) *user_role.Repository {
	return &user_role.Repository{Database: db}
}
