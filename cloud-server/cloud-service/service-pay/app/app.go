package app

import (
	"service-pay/app/http/route"
	"service-pay/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
