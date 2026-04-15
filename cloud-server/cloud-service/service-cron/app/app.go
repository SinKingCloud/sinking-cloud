package app

import (
	"service-cron/app/http/route"
	"service-cron/app/queue"
	"service-cron/app/service"
)

func Run() {
	service.Init()
	queue.Init(service.Task.TaskQueueHandle, service.Task.ReloadQueueHandle)
	route.Init()
}
