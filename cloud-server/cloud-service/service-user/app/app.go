package app

import (
	"service-user/app/http/route"
	"service-user/app/queue"
	"service-user/app/service"
	"service-user/app/task"
)

func Run() {
	queue.Init()
	service.Init()
	task.Init()
	route.Init()
}
