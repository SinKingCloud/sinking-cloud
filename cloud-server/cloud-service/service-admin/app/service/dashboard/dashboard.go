package dashboard

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Service service接口
type Service interface {
	GetServiceCount() map[string]int //获取服务数量
}

// service 注入结构
type service struct {
	consul *consul.Client
	cache  cache.Interface
}

// NewService 实例化service
func NewService(consul *consul.Client, cache cache.Interface) *service {
	return &service{
		consul: consul,
		cache:  cache,
	}
}
