package app

import (
	"service-admin/app/http/route"
	"service-admin/app/queue"
	"service-admin/app/service"
)

func Run() {
	service.Init()
	queue.Init()
	route.Init()
}
