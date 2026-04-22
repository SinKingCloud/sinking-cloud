package order

import (
	"service-pay/app/repository/order"
	"service-pay/app/repository/order_log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
}

// Session服务
type service struct {
	cache              cache.Interface
	orderRepository    order.Interface
	orderLogRepository order_log.Interface
}

// NewService 实例化service
func NewService(orderRepository order.Interface, orderLogRepository order_log.Interface, cache cache.Interface) *service {
	return &service{
		orderRepository:    orderRepository,
		orderLogRepository: orderLogRepository,
		cache:              cache,
	}
}
