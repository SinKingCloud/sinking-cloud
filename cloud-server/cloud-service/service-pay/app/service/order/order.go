package order

import (
	"service-pay/app/model"
	balanceRepo "service-pay/app/repository/balance"
	balanceLogRepo "service-pay/app/repository/balance_log"
	"service-pay/app/repository/order"
	"service-pay/app/repository/order_log"
	"service-pay/app/service/channel"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Service service接口
type Service interface {
	FindById(id int64) (*model.Order, error)                                                                                                                            //根据ID查询订单
	Select(where *order.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*order.Order, total int64, err error)                     //查询订单列表
	SelectByUserId(userId int64, where *order.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*UserOrder, total int64, err error) //查询用户订单列表
	SelectLog(where *order_log.SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*order_log.OrderLog, total int64, err error)    //查询订单日志
	FindLogById(id int64) (*model.OrderLog, error)                                                                                                                      //根据ID查询日志
	Create(userId int64, name string, Type int, payType string, amount float64, param string, callbackRpcService string, callbackRpcFunc string, describe string) (int64, error)
	Callback(orderId int64, outTradeNo string, payAmount float64, payTime str.DateTime) error //订单回调
	Deliver(id int64, operatorId int64) error                                                 //订单发货
	Refound(id int64, money float64, operatorId int64) error                                  //订单退款
	CleanExpiredUnPayOrder(day int) (int64, error)                                            //清理过期未支付订单
}

// Session服务
type service struct {
	cache                cache.Interface
	database             *database.Database
	channelService       channel.Service
	balanceRepository    balanceRepo.Interface
	balanceLogRepository balanceLogRepo.Interface
	orderRepository      order.Interface
	orderLogRepository   order_log.Interface
}

// NewService 实例化service
func NewService(orderRepository order.Interface, orderLogRepository order_log.Interface, balanceRepository balanceRepo.Interface, balanceLogRepository balanceLogRepo.Interface, channelService channel.Service, database *database.Database, cache cache.Interface) *service {
	return &service{
		orderRepository:      orderRepository,
		orderLogRepository:   orderLogRepository,
		balanceRepository:    balanceRepository,
		balanceLogRepository: balanceLogRepository,
		channelService:       channelService,
		database:             database,
		cache:                cache,
	}
}
