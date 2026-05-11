package app

import (
	"service-admin/app/http/route"
	"service-admin/app/queue"
	"service-admin/app/service"
	"service-admin/app/task"
)

func Run() {
	service.Init()
	queue.Init()
	task.Init()
	route.Init()
}
