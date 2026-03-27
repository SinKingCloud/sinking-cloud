package service

import (
	config2 "service-user/app/repository/config"
	"service-user/app/repository/session"
	user2 "service-user/app/repository/user"
	"service-user/app/service/auth"
	"service-user/app/service/config"
	"service-user/app/service/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Config config.Service
	Auth   auth.Service
	User   user.Service
)

func Init() {
	configRepository := config2.NewRepository(global.App.Database)
	userRepository := user2.NewRepository(global.App.Database)
	sessionRepository := session.NewRepository(global.App.Database)
	cache := global.App.Cache
	authService := service.Auth

	Config = config.NewService(configRepository, cache)
	Auth = auth.NewService(authService, sessionRepository, userRepository, cache)
	User = user.NewService(userRepository, Config, cache)
}
