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
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

var (
	Channel channel.Service
	Balance balance.Service
	Order   order.Service
)

func Init() {
	app := global.App
	cache := app.Cache
	database := app.Database
	consul := app.Consul

	balanceRepository := balance2.NewRepository(database)
	balanceLogRepository := balance_log.NewRepository(database)
	orderRepository := order2.NewRepository(database)
	orderLogRepository := order_log.NewRepository(database)
	channelRepository := channel2.NewRepository(database)
	configService := service.Config

	Channel = channel.NewService(channelRepository, orderRepository, configService, consul, cache)
	Order = order.NewService(orderRepository, orderLogRepository, balanceRepository, balanceLogRepository, Channel, database, cache)
	Balance = balance.NewService(balanceRepository, balanceLogRepository, Order, Channel, consul, database, cache)
}
