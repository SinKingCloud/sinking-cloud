package app

import (
	"service-cron/app/http/route"
	"service-cron/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
