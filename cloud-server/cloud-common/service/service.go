package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/enum"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user"
)

var (
	Enum   = enum.GetIns()
	Config = config.GetIns()
	Auth   = auth.GetIns()
	Rbac   = rbac.GetIns()
	User   = user.GetIns()
)
