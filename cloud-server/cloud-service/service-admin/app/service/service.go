package service

import (
	"service-admin/app/repository"
	"service-admin/app/service/auth"
	"service-admin/app/service/config"
	"service-admin/app/service/rbac"
	"service-admin/app/service/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Config ConfigService
	Auth   AuthService
	Rbac   RbacService
	User   UserService
)

func Init() {
	configRepository := repository.NewConfig(global.App.Database)
	userRepository := repository.NewUser(global.App.Database)
	roleRepository := repository.NewRole(global.App.Database)
	sessionRepository := repository.NewSession(global.App.Database)
	permissionRepository := repository.NewPermission(global.App.Database)
	userRoleRepository := repository.NewUserRole(global.App.Database)
	rolePermRepository := repository.NewRolePermission(global.App.Database)
	cache := global.App.Cache
	authService := service.Auth

	Config = config.NewService(configRepository, cache)
	Auth = auth.NewService(authService, sessionRepository, userRepository, cache)
	Rbac = rbac.NewService(roleRepository, permissionRepository, userRoleRepository, rolePermRepository, cache)
	User = user.NewService(userRepository, cache)
}
