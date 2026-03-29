package controller

import (
	"service-user/app/http/controller/person"
	"service-user/app/http/controller/security"
)

var (
	Person   = &person.Controller{}
	Security = &security.Controller{}
)
