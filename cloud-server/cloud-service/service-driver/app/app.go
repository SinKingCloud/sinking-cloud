package app

import (
	"service-driver/app/http/route"
	"service-driver/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
