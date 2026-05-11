package service

import (
	config2 "service-admin/app/repository/config"
	"service-admin/app/repository/menu"
	"service-admin/app/repository/menu_permission"
	"service-admin/app/repository/permission"
	"service-admin/app/repository/role"
	"service-admin/app/repository/role_menu"
	"service-admin/app/repository/session"
	user2 "service-admin/app/repository/user"
	"service-admin/app/repository/user_role"
	"service-admin/app/service/auth"
	"service-admin/app/service/config"
	"service-admin/app/service/rbac"
	"service-admin/app/service/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Config config.Service
	Auth   auth.Service
	Rbac   rbac.Service
	User   user.Service
)

func Init() {
	configRepository := config2.NewRepository(global.App.Database)
	userRepository := user2.NewRepository(global.App.Database)
	roleRepository := role.NewRepository(global.App.Database)
	sessionRepository := session.NewRepository(global.App.Database)
	permissionRepository := permission.NewRepository(global.App.Database)
	userRoleRepository := user_role.NewRepository(global.App.Database)
	roleMenuRepository := role_menu.NewRepository(global.App.Database)
	menuRepository := menu.NewRepository(global.App.Database)
	menuPermRepository := menu_permission.NewRepository(global.App.Database)
	cache := global.App.Cache
	authService := service.Auth
	database := global.App.Database

	Config = config.NewService(configRepository, cache)
	Auth = auth.NewService(authService, sessionRepository, userRepository, database, cache)
	Rbac = rbac.NewService(roleRepository, permissionRepository, userRoleRepository, roleMenuRepository, menuRepository, menuPermRepository, database, cache)
	User = user.NewService(userRepository, sessionRepository, userRoleRepository, configRepository, Rbac, database, cache)
}
