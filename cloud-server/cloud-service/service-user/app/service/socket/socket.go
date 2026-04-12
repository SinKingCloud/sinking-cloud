package socket

import (
	"service-user/app/repository/session"
	"service-user/app/service/auth"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

type Service interface {
	GenToken(token string) (string, error)                 //生成连接token
	SendMessage(message *Message, to ...interface{}) error //发送消息
}

// Session服务
type service struct {
	authService       auth.Service
	socketService     *socket.Service
	repositorySession session.Interface
	cache             cache.Interface
}

// NewService 实例化service
func NewService(authService auth.Service, repositorySession session.Interface, socketService *socket.Service, cache cache.Interface) *service {
	return &service{
		authService:       authService,
		repositorySession: repositorySession,
		socketService:     socketService,
		cache:             cache,
	}
}
