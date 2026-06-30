package service

import (
	queueLog "service-cron/app/queue/log"
	queueRegister "service-cron/app/queue/register"
	queueReload "service-cron/app/queue/reload"
	queueTask "service-cron/app/queue/task"
	log2 "service-cron/app/repository/log"
	task2 "service-cron/app/repository/task"
	"service-cron/app/service/log"
	"service-cron/app/service/task"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var (
	Log  log.Service
	Task task.Service
)

func Init() {
	database := global.App.Database
	cache := global.App.Cache

	logRepository := log2.NewRepository(database)
	taskRepository := task2.NewRepository(database)

	Log = log.NewService(logRepository)
	var taskService task.Service
	logQueue := queueLog.Register(func(message *queueLog.Message) error {
		return taskService.LogQueueHandle(message)
	})
	reloadQueue := queueReload.Register(func(message *queueReload.Message) error {
		return taskService.ReloadQueueHandle(message)
	})
	taskQueue := queueTask.Register(func(message *queueTask.Message) error {
		return taskService.TaskQueueHandle(message)
	})
	registerQueue := queueRegister.Register(func(message *queueRegister.Message) error {
		return taskService.RegisterQueueHandle(message)
	})
	taskService = task.NewService(logRepository, taskRepository, cache, database, taskQueue, reloadQueue, logQueue, registerQueue)
	Task = taskService
}
