package repository

import (
	"service-admin/app/model"
	"service-admin/app/repository/role"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Role 角色仓储接口
type Role interface {
	FindById(id int64) (*model.Role, error)
}

// NewRole 创建Role仓储
func NewRole(db *database.Database) *role.Repository {
	return &role.Repository{Database: db}
}
