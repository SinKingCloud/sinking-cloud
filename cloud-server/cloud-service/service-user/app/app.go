package app

import (
	"service-user/app/http/route"
	"service-user/app/service"
	"service-user/app/task"
)

func Run() {
	service.Init()
	task.Init()
	route.Init()
}
