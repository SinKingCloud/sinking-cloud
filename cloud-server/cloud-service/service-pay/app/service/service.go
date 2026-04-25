package service

import (
	balance2 "service-pay/app/repository/balance"
	"service-pay/app/repository/balance_log"
	channel2 "service-pay/app/repository/channel"
	order2 "service-pay/app/repository/order"
	"service-pay/app/repository/order_log"
	"service-pay/app/service/balance"
	"service-pay/app/service/channel"
	"service-pay/app/service/order"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
	Channel channel.Service
	Balance balance.Service
	Order   order.Service
)

func Init() {
	cache := global.App.Cache
	database := global.App.Database

	balanceRepository := balance2.NewRepository(database)
	balanceLogRepository := balance_log.NewRepository(database)
	orderRepository := order2.NewRepository(database)
	orderLogRepository := order_log.NewRepository(database)
	channelRepository := channel2.NewRepository(database)

	Channel = channel.NewService(channelRepository, cache)
	Balance = balance.NewService(balanceRepository, balanceLogRepository, database, cache)
	Order = order.NewService(orderRepository, orderLogRepository, cache)
}
