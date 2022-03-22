package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/rpc/face"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/rpc/log"
)

var (
	Face = &face.ControllerFace{}
	Log  = &log.ControllerLog{}
)
