package service

import (
	"service-admin/app/repository"
	"service-admin/app/service/auth"
	"service-admin/app/service/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var Config ConfigService
var Auth AuthService

func Init() {
	configRepository := repository.NewConfig(global.App.Database)
	userRepository := repository.NewUser(global.App.Database)
	sessionRepository := repository.NewSession(global.App.Database)
	cache := global.App.Cache
	authService := service.Auth

	Config = config.NewService(configRepository, cache)
	Auth = auth.NewService(authService, sessionRepository, userRepository, cache)
}
