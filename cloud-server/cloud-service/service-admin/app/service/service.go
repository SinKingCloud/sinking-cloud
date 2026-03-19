package service

import (
	"service-admin/app/repository"
	"service-admin/app/service/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var Config ConfigService

func Init() {
	Config = config.NewService(repository.NewConfig(global.App.Database), global.App.Cache)
}
