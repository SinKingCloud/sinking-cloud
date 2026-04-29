package app

import (
	"service-pay/app/http/route"
	"service-pay/app/service"
	"service-pay/app/task"
)

func Run() {
	service.Init()
	task.Init()
	route.Init()
}
