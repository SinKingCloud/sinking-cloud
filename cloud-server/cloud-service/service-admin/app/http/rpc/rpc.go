package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/http/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/http/rpc/rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/http/rpc/user"
)

var (
	Rbac = &rbac.ControllerRbac{}
	User = &user.ControllerUser{}
	Log  = &log.ControllerLog{}
)
