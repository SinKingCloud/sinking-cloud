package app

import (
	"service-socket/app/http/route"
	"service-socket/app/queue"
	"service-socket/app/service"
)

func Run() {
	service.Init()
	queue.Init()
	route.Init()
}
