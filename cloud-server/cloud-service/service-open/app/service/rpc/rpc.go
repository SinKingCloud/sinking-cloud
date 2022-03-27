package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc/system_rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc/system_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc/web_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc/web_user_log"
)

var (
	WebUserLog = web_user_log.GetIns()
	WebUser    = web_user.GetIns()
	SystemRbac = system_rbac.GetIns()
	SystemUser = system_user.GetIns()
	SystemLog  = system_log.GetIns()
)
