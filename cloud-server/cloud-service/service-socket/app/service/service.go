package service

import (
	"errors"
	socketQueue "service-socket/app/queue/socket"
	"service-socket/app/service/socket"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
	Socket socket.Service
)

func Init() {
	app := global.App
	cache := global.App.Cache
	var socketService socket.Service
	queue := socketQueue.Register(app.Consul.GetAddress(), &socketQueue.Handler{
		GetConnectionIdAddr: func(group string, id string) (string, error) {
			if socketService == nil {
				return "", errors.New("socket服务未初始化")
			}
			return socketService.GetConnectionIdAddr(group, id)
		},
		SendToLocalId: func(group string, id string, messageType int, payload []byte) error {
			if socketService == nil {
				return errors.New("socket服务未初始化")
			}
			return socketService.SendToLocalId(group, id, messageType, payload)
		},
		SendToLocalAll: func(group string, messageType int, payload []byte) (queued int, dropped int, closed int, err error) {
			if socketService == nil {
				return 0, 0, 0, errors.New("socket服务未初始化")
			}
			return socketService.SendToLocalAll(group, messageType, payload)
		},
	})

	socketService = socket.NewService(app.Consul.GetAddress(), cache, queue)
	Socket = socketService
}
