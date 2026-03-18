package app

import (
	"service-user/app/http/route"
	"service-user/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
