package user

import (
	"service-user/app/http/controller/user/message"
	"service-user/app/http/controller/user/person"
	"service-user/app/http/controller/user/security"
)

type Controller struct {
	Person   person.Controller
	Security security.Controller
	Message  message.Controller
}
