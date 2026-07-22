package rbac

import (
	"service-admin/app/model"
	queuePermission "service-admin/app/queue/permission"
	"service-admin/app/repository/menu"
	"service-admin/app/repository/menu_permission"
	"service-admin/app/repository/permission"
	"service-admin/app/repository/role"
	"service-admin/app/repository/role_menu"
	"service-admin/app/repository/user_role"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// Service service接口
type Service interface {
	RegisterPermission(name string, code string, path string) error                                                                            //注册权限
	RegisterPermissionAsync(name string, code string, path string)                                                                             //异步注册权限
	CheckUserHasPermission(userId int64, path string) error                                                                                    //判断用户是否含有权限
	GetPermissionEnum(refresh bool) (map[string]string, error)                                                                                 //获取权限枚举
	GetRoleNameEnum(refresh bool) (map[int64]string, error)                                                                                    //获取角色名称枚举
	GetUserMenuTree(userId int64) []*UserMenuTree                                                                                              //获取用户菜单
	GetAllMenuTreeByDirAndPage() []*MenuTree                                                                                                   //获取所有菜单
	GetAllMenuTreeByStatusEnabled() []*MenuTree                                                                                                //获取角色菜单树
	CreateRole(name string, status int, menuIds []int64) error                                                                                 //创建角色
	UpdateRole(id int64, data *role.UpdateRole, menuIds *[]int64) error                                                                        //修改角色
	GetRoleInfoById(id int64) (*Role, error)                                                                                                   //获取角色详情
	GetMenuOrControlInfoById(id int64) (interface{}, error)                                                                                    //获取页面或控件信息
	GetPageControlsByUserId(userId int64, pageCode string) []*UserControl                                                                      //根据用户ID获取页面控件
	GetPageControlsByMenuId(id int64, where *menu.SelectMenu, queryPage *page.Query) (*page.Result[*Control], error)                           //根据菜单ID获取页面控件
	AssignRoleToUser(userId int64, roleIds []int64) error                                                                                      //分配用户角色
	DeleteMenuOrControlByIds(ids []int64) error                                                                                                //根据ID列表删除菜单或控件数据
	CreateMenuOrControl(parentId int64, code string, name string, icon string, Type int, sort int, status int, permissionCodes []string) error //创建菜单或控件数据
	UpdateMenuOrControl(id int64, data *menu.UpdateMenu, permissionCodes *[]string) error                                                      //更新菜单或控件数据
	UpdatePermissionByIds(ids []int64, data *permission.UpdatePermission) error                                                                //根据id集合更新权限
	DeletePermissionByIds(ids []int64) error                                                                                                   //根据ID集合删除权限
	DeleteRoleByIds(ids []int64) error                                                                                                         //根据ID集合删除角色
	CreatePermission(data *model.Permission) error                                                                                             //新增权限数据
	SelectPermission(where *permission.SelectPermission, queryPage *page.Query) (*page.Result[*model.Permission], error)                       //查询权限数据
	SelectRole(where *role.SelectRole, queryPage *page.Query) (*page.Result[*model.Role], error)                                               //查询角色数据
}

// service 注入结构
type service struct {
	repositoryRole           role.Interface
	repositoryPermission     permission.Interface
	repositoryUserRole       user_role.Interface
	repositoryRoleMenu       role_menu.Interface
	repositoryMenu           menu.Interface
	repositoryMenuPermission menu_permission.Interface
	database                 *database.Database
	cache                    cache.Interface
	permissionQueue          *queue.UnicastClient[*queuePermission.Task]
}

// NewService 实例化service
func NewService(
	repositoryRole role.Interface,
	repositoryPermission permission.Interface,
	repositoryUserRole user_role.Interface,
	repositoryRoleMenu role_menu.Interface,
	repositoryMenu menu.Interface,
	repositoryMenuPermission menu_permission.Interface,
	database *database.Database,
	cache cache.Interface,
	permissionQueue *queue.UnicastClient[*queuePermission.Task],
) *service {
	ins := &service{
		repositoryRole:           repositoryRole,
		repositoryPermission:     repositoryPermission,
		repositoryUserRole:       repositoryUserRole,
		repositoryRoleMenu:       repositoryRoleMenu,
		repositoryMenu:           repositoryMenu,
		repositoryMenuPermission: repositoryMenuPermission,
		database:                 database,
		cache:                    cache,
		permissionQueue:          permissionQueue,
	}
	if permissionQueue != nil {
		permissionQueue.SetConsumer(ins.consumePermissionTask)
	}
	return ins
}
