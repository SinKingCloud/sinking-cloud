package controller

import (
	"service-driver/app/http/controller/admin"
	"service-driver/app/http/controller/api"
)

var (
	Admin = &admin.Controller{}
	Api   = &api.Controller{}
)
