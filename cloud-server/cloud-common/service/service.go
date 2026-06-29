package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/cron"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/enum"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/message"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user"
)

var (
	Enum    = enum.GetIns()
	Config  = config.GetIns()
	Auth    = auth.GetIns()
	User    = user.GetIns()
	Socket  = socket.GetIns()
	Cron    = cron.GetIns()
	Pay     = pay.GetIns()
	Message = message.GetIns()
)
