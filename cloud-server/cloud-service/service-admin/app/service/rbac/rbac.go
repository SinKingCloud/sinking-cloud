package rbac

import (
	"service-admin/app/model"
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
	RegisterPermission(name string, code string, path string) error                                                                                           //注册权限
	CheckUserHasPermission(userId int64, path string) error                                                                                                   //判断用户是否含有权限
	GetUserMenuTree(userId int64) []*MenuTree                                                                                                                 //获取用户菜单
	GetAllMenuTree() []*MenuTree                                                                                                                              //获取所有菜单
	GetUserControl(userId int64, pageCode string) []*Control                                                                                                  //获取用户控件
	UpdatePermissionByIds(ids []int64, data *permission.UpdatePermission) error                                                                               //根据id集合更新权限
	DeletePermissionByIds(ids []int64) error                                                                                                                  //根据ID集合删除权限
	CreatePermission(data *model.Permission) error                                                                                                            //新增权限数据
	SelectPermission(where *permission.SelectPermission, orderByField string, orderByType string, page int, pageSize int) ([]*model.Permission, int64, error) //查询权限数据
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
