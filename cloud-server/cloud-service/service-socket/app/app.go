package app

import (
	"service-socket/app/http/route"
	"service-socket/app/queue"
	"service-socket/app/service"
)

func Run() {
	queue.Init()
	service.Init()
	route.Init()
}
