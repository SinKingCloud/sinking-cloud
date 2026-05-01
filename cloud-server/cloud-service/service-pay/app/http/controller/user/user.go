package user

import (
	"service-pay/app/http/controller/user/balance"
	"service-pay/app/http/controller/user/order"
)

type Controller struct {
	Balance balance.Controller
	Order   order.Controller
}
