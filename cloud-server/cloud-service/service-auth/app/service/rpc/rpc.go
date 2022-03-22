package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/rpc/face"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/rpc/log"
)

var (
	Face = face.GetIns()
	Log  = log.GetIns()
)
