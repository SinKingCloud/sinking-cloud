package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/http/rpc/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/http/rpc/sms"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/http/rpc/upload"
)

var (
	Sms     = &sms.ControllerSms{}
	Captcha = &captcha.ControllerCaptcha{}
	Upload  = &upload.ControllerUpload{}
)
