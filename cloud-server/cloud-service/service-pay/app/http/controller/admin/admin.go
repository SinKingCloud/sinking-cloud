package admin

import (
	"service-pay/app/http/controller/admin/balance"
	"service-pay/app/http/controller/admin/channel"
	"service-pay/app/http/controller/admin/config"
	"service-pay/app/http/controller/admin/log"
	"service-pay/app/http/controller/admin/order"
)

type Controller struct {
	Balance balance.Controller
	Order   order.Controller
	Log     log.Controller
	Channel channel.Controller
	Config  config.Controller
}
