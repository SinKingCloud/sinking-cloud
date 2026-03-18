package app

import (
	"service-auth/app/http/route"
	"service-auth/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
