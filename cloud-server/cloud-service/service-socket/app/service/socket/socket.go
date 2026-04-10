package socket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Service service接口
type Service interface {
	Server(token string) (*sinking_websocket.Server, error)
}

// Session服务
type service struct {
	cache cache.Interface
}

// NewService 实例化service
func NewService(cache cache.Interface) *service {
	return &service{
		cache: cache,
	}
}
