package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/rpc/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/rpc/sms"
)

var (
	Sms     = &sms.ControllerSms{}
	Captcha = &captcha.ControllerCaptcha{}
)
