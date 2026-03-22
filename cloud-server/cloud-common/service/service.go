package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/enum"
)

var (
	Enum   = enum.GetIns()
	Config = config.GetIns()
	Auth   = auth.GetIns()
)
