package app

import (
	"service-cron/app/http/route"
	"service-cron/app/service"
	"service-cron/app/task"
)

func Run() {
	service.Init()
	task.Init()
	route.Init()
}
