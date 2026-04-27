package service

import (
	"service-socket/app/service/socket"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
	Socket socket.Service
)

func Init() {
	app := global.App
	cache := global.App.Cache

	Socket = socket.NewService(app.Consul.GetAddress(), cache)
}
