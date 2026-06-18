package controller

import (
	"service-admin/app/http/controller/config"
	"service-admin/app/http/controller/dashboard"
	"service-admin/app/http/controller/master"
	"service-admin/app/http/controller/menu"
	"service-admin/app/http/controller/page"
	"service-admin/app/http/controller/permission"
	"service-admin/app/http/controller/person"
	"service-admin/app/http/controller/role"
	"service-admin/app/http/controller/security"
)

var (
	Page       = &page.Controller{}
	Dashboard  = &dashboard.Controller{}
	Person     = &person.Controller{}
	Security   = &security.Controller{}
	Permission = &permission.Controller{}
	Menu       = &menu.Controller{}
	Role       = &role.Controller{}
	Master     = &master.Controller{}
	Config     = &config.Controller{}
)
