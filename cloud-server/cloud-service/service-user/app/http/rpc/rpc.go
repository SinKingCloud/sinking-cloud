package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/rpc/face"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/rpc/user"
)

var (
	Face = &face.ControllerFace{}
	Log  = &log.ControllerLog{}
	User = &user.ControllerUser{}
)
