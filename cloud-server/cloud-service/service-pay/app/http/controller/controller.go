package controller

import (
	"service-pay/app/http/controller/admin"
	"service-pay/app/http/controller/api"
	"service-pay/app/http/controller/user"
)

var (
	Admin = &admin.Controller{}
	User  = &user.Controller{}
	Api   = &api.Controller{}
)
