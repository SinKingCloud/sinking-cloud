package service

import (
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
	Task = task.NewService(logRepository, taskRepository, cache)
}
