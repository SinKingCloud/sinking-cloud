package service

import (
	channel2 "service-driver/app/repository/channel"
	file2 "service-driver/app/repository/file"
	scene2 "service-driver/app/repository/scene"
	"service-driver/app/service/channel"
	"service-driver/app/service/file"
	"service-driver/app/service/scene"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
	Channel channel.Service
	Scene   scene.Service
	File    file.Service
)

func Init() {
	cache := global.App.Cache
	database := global.App.Database

	channelRepository := channel2.NewRepository(global.App.Database)
	sceneRepository := scene2.NewRepository(global.App.Database)
	fileRepository := file2.NewRepository(global.App.Database)

	Channel = channel.NewService(channelRepository, sceneRepository, fileRepository, database, cache)
	Scene = scene.NewService(sceneRepository, fileRepository, database, cache)
	File = file.NewService(channelRepository, sceneRepository, fileRepository, database, cache)
}
