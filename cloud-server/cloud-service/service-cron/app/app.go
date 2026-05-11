package app

import (
	"service-cron/app/http/route"
	"service-cron/app/queue"
	"service-cron/app/service"
	"service-cron/app/task"
)

func Run() {
	service.Init()
	queue.Init(service.Task.TaskQueueHandle, service.Task.ReloadQueueHandle, service.Task.LogQueueHandle, service.Task.RegisterQueueHandle)
	service.Task.SetQueue(queue.Task, queue.Reload, queue.Log, queue.Register)
	service.Task.Start()
	task.Init()
	route.Init()
}
