package app

import (
	"service-socket/app/http/route"
	"service-socket/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
