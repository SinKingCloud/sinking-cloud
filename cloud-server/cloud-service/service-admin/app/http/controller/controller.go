package controller

import (
	"service-admin/app/http/controller/person"
	"service-admin/app/http/controller/security"
)

var (
	Person   = &person.Controller{}
	Security = &security.Controller{}
)
