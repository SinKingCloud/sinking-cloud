package app

import (
	"service-admin/app/http/route"
	"service-admin/app/queue"
	"service-admin/app/service"
	"service-admin/app/task"
)

func Run() {
	queue.Init()
	service.Init()
	task.Init()
	route.Init()
}
