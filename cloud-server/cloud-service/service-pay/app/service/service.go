package service

import (
	balance2 "service-pay/app/repository/balance"
	"service-pay/app/repository/balance_log"
	order2 "service-pay/app/repository/order"
	"service-pay/app/repository/order_log"
	"service-pay/app/service/balance"
	"service-pay/app/service/order"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
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

	Balance = balance.NewService(balanceRepository, balanceLogRepository, cache)
	Order = order.NewService(orderRepository, orderLogRepository, cache)
}
