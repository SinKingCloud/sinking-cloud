package repository

import (
	"service-admin/app/model"
	"service-admin/app/repository/permission"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Permission 权限仓储接口
type Permission interface {
	FindByPath(path string) (*model.Permission, error)
	FindByCode(code string) (*model.Permission, error)
	Create(perm *model.Permission) error
	UpdateByCode(code string, perm *model.Permission) error
}

// NewPermission 创建Permission仓储
func NewPermission(db *database.Database) *permission.Repository {
	return &permission.Repository{Database: db}
}
