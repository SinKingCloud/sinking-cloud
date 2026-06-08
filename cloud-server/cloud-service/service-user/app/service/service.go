package service

import (
	config2 "service-user/app/repository/config"
	realAuth2 "service-user/app/repository/real_auth"
	"service-user/app/repository/session"
	user2 "service-user/app/repository/user"
	"service-user/app/service/auth"
	"service-user/app/service/config"
	"service-user/app/service/socket"
	"service-user/app/service/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Config config.Service
	Auth   auth.Service
	User   user.Service
	Socket socket.Service
)

func Init() {
	configRepository := config2.NewRepository(global.App.Database)
	userRepository := user2.NewRepository(global.App.Database)
	realAuthRepository := realAuth2.NewRepository(global.App.Database)
	sessionRepository := session.NewRepository(global.App.Database)
	cache := global.App.Cache
	authService := service.Auth
	payService := service.Pay
	socketService := service.Socket
	systemConfig := service.Config
	database := global.App.Database

	Config = config.NewService(configRepository, cache)
	Auth = auth.NewService(authService, sessionRepository, userRepository, database, cache)
	User = user.NewService(userRepository, realAuthRepository, systemConfig, payService, global.App.Database, cache)
	Socket = socket.NewService(Auth, sessionRepository, socketService, cache)
}
