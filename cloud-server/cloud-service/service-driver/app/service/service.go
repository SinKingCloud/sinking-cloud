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

	channelRepository := channel2.NewRepository(database)
	sceneRepository := scene2.NewRepository(database)
	fileRepository := file2.NewRepository(database)

	Channel = channel.NewService(channelRepository, sceneRepository, fileRepository, cache)
	Scene = scene.NewService(channelRepository, sceneRepository, fileRepository, cache)
	File = file.NewService(channelRepository, sceneRepository, fileRepository, Channel, Scene, cache)
}
