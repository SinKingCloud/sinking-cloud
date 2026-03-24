package rbac

import (
	"service-admin/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// service 注入结构
type service struct {
	repositoryRole           repository.Role
	repositoryPermission     repository.Permission
	repositoryUserRole       repository.UserRole
	repositoryRolePermission repository.RolePermission
	cache                    cache.Interface
}

// NewService 实例化service
func NewService(
	repositoryRole repository.Role,
	repositoryPermission repository.Permission,
	repositoryUserRole repository.UserRole,
	repositoryRolePerm repository.RolePermission,
	cache cache.Interface,
) *service {
	return &service{
		repositoryRole:           repositoryRole,
		repositoryPermission:     repositoryPermission,
		repositoryUserRole:       repositoryUserRole,
		repositoryRolePermission: repositoryRolePerm,
		cache:                    cache,
	}
}
