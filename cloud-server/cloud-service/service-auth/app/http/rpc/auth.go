package rpc

import (
	"encoding/json"
	"service-auth/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
)

type authService struct {
}

// CheckCaptcha 判断滑块验证码
func (*authService) CheckCaptcha(params json.RawMessage) (interface{}, error) {
	var param *auth.CheckCaptchaParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Auth.CheckCaptcha(param.Key, param.X, param.Y), nil
}

// CheckSms 判断短信验证码
func (*authService) CheckSms(params json.RawMessage) (interface{}, error) {
	var param *auth.CheckSmsParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Auth.CheckSms(param.Phone, param.Code), nil
}

// CheckMail 判断邮件验证码
func (*authService) CheckMail(params json.RawMessage) (interface{}, error) {
	var param *auth.CheckMailParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Auth.CheckMail(param.Mail, param.Code), nil
}
