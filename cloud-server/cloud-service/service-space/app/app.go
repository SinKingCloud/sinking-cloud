package app

import (
	"service-space/app/http/route"
	"service-space/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
