package app

import (
	"service-admin/app/http/route"
	"service-admin/app/service"
)

func Run() {
	service.Init()
	route.Init()
}
