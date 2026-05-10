package controller

import (
	"service-user/app/http/controller/admin"
	"service-user/app/http/controller/api"
	"service-user/app/http/controller/user"
)

var (
	User  = &user.Controller{}
	Admin = &admin.Controller{}
	Api   = &api.Controller{}
)
