package user

import (
	"service-space/app/http/controller/user/mind_data"
	"service-space/app/http/controller/user/mind_share"
	"service-space/app/http/controller/user/mind_type"
	"service-space/app/http/controller/user/mind_version"
)

type Controller struct {
	MindData    mind_data.Controller
	MindType    mind_type.Controller
	MindVersion mind_version.Controller
	MindShare   mind_share.Controller
}
