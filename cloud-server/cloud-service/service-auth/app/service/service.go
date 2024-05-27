package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/order"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/pay_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/user"
)

var (
	Auth    = auth.GetIns()
	Captcha = captcha.GetIns()
	Order   = order.GetIns()
	PayLog  = pay_log.GetIns()
	User    = user.GetIns()
)
