package order

import (
	"service-pay/app/model"
	"service-pay/app/repository/order"
	"service-pay/app/repository/order_log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	FindById(id int64) (*model.Order, error)                                                                                                                            //根据ID查询订单
	Select(where *order.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*order.Order, total int64, err error)                     //查询订单列表
	SelectByUserId(userId int64, where *order.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*UserOrder, total int64, err error) //查询用户订单列表
	SelectLog(where *order_log.SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.OrderLog, total int64, err error)        //查询订单日志
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
