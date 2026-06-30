package app

import (
	"service-cron/app/http/route"
	"service-cron/app/queue"
	"service-cron/app/service"
	"service-cron/app/task"
)

func Run() {
	queue.Init()
	service.Init()
	task.Init()
	route.Init()
}
