package service

import (
	mindDataRepository "service-space/app/repository/mind_data"
	mindShareRepository "service-space/app/repository/mind_share"
	mindTypeRepository "service-space/app/repository/mind_type"
	mindVersionRepository "service-space/app/repository/mind_version"
	"service-space/app/service/mind"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Mind mind.Service
)

// Init 初始化服务
func Init() {
	app := global.App
	database := app.Database
	cache := app.Cache
	configService := commonService.Config
	mindDataRepo := mindDataRepository.NewRepository(database)
	mindTypeRepo := mindTypeRepository.NewRepository(database)
	mindVersionRepo := mindVersionRepository.NewRepository(database)
	mindShareRepo := mindShareRepository.NewRepository(database)

	Mind = mind.NewService(mindDataRepo, mindTypeRepo, mindVersionRepo, mindShareRepo, database, cache, configService)
}
