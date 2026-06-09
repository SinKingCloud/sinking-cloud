package balance

import (
	"service-pay/app/repository/balance"
	"service-pay/app/repository/balance_log"
	"service-pay/app/service/channel"
	"service-pay/app/service/order"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Service service接口
type Service interface {
	GetByUserId(userId int64, refresh bool) (*Balance, error)
	OperateMoney(userId int64, money float64, isOperateFrozen bool, logType int, describe string) error                                                   //操作用户余额并写入明细
	OperateMoneyCount(userId int64, money float64) error                                                                                                  //操作累计充值或消费金额
	SelectLog(where *balance_log.SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*Log, total int64, err error) //查询用户余额明细
	CreateRechargeOrder(userId int64, money int, payType string) (tradeNo string, err error)                                                              //创建充值订单
	RechargeOrderNotify(order *pay.NotifyParam) error                                                                                                     //充值订单回调
}

// Session服务
type service struct {
	consul               *client.Client
	balanceRepository    balance.Interface
	balanceLogRepository balance_log.Interface
	orderService         order.Service
	channelService       channel.Service
	app                  *global.Application
	database             *database.Database
	cache                cache.Interface
}

// NewService 实例化service
func NewService(balanceRepository balance.Interface, balanceLogRepository balance_log.Interface, orderService order.Service, channelService channel.Service, consul *client.Client, database *database.Database, cache cache.Interface) *service {
	return &service{
		balanceRepository:    balanceRepository,
		balanceLogRepository: balanceLogRepository,
		orderService:         orderService,
		channelService:       channelService,
		consul:               consul,
		database:             database,
		cache:                cache,
	}
}
