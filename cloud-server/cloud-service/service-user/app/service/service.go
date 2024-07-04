package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/message"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/notice"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/order"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/pay_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/user"
)

var (
	User    = user.GetIns()
	PayLog  = pay_log.GetIns()
	Order   = order.GetIns()
	Notice  = notice.GetIns()
	Message = message.GetIns()
	Log     = log.GetIns()
	AliYun  = aliyun.GetIns()
)
