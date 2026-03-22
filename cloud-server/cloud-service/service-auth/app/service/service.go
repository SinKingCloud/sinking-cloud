package service

import (
	"service-auth/app/service/auth"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var Auth AuthService

func Init() {
	cache := global.App.Cache
	configService := service.Config
	authService := service.Auth

	Auth = auth.NewService(cache, configService, authService)
}
