package rbac

import (
	"service-admin/app/repository/menu"
	"service-admin/app/repository/menu_permission"
	"service-admin/app/repository/permission"
	"service-admin/app/repository/role"
	"service-admin/app/repository/role_menu"
	"service-admin/app/repository/role_permission"
	"service-admin/app/repository/user_role"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	RegisterPermission(name string, code string, path string) error //注册权限
	CheckUserHasPermission(userId int64, path string) error         //判断用户是否含有权限
	GetUserMenus(userId int64) []*UserMenu                          //获取用户菜单
	GetUserControl(userId int64, pageCode string) []*UserControl    //获取用户页面控件
}

// service 注入结构
type service struct {
	repositoryRole           role.Interface
	repositoryPermission     permission.Interface
	repositoryUserRole       user_role.Interface
	repositoryRolePermission role_permission.Interface
	repositoryRoleMenu       role_menu.Interface
	repositoryMenu           menu.Interface
	repositoryMenuPermission menu_permission.Interface
	cache                    cache.Interface
}

// NewService 实例化service
func NewService(
	repositoryRole role.Interface,
	repositoryPermission permission.Interface,
	repositoryUserRole user_role.Interface,
	repositoryRolePerm role_permission.Interface,
	repositoryRoleMenu role_menu.Interface,
	repositoryMenu menu.Interface,
	repositoryMenuPermission menu_permission.Interface,
	cache cache.Interface,
) *service {
	return &service{
		repositoryRole:           repositoryRole,
		repositoryPermission:     repositoryPermission,
		repositoryUserRole:       repositoryUserRole,
		repositoryRolePermission: repositoryRolePerm,
		repositoryRoleMenu:       repositoryRoleMenu,
		repositoryMenu:           repositoryMenu,
		repositoryMenuPermission: repositoryMenuPermission,
		cache:                    cache,
	}
}
