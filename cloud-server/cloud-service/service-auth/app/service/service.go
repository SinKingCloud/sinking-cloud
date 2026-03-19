package service

import (
	"service-auth/app/repository"
	"service-auth/app/service/auth"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var Auth AuthService

func Init() {
	Auth = auth.NewService(repository.NewSession(global.App.Database), global.App.Cache)
}
