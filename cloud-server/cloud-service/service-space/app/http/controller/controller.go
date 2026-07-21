package controller

import (
	"service-space/app/http/controller/api"
	"service-space/app/http/controller/user"
)

var (
	User = &user.Controller{}
	Api  = &api.Controller{}
)
