package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CheckCaptcha 判断滑块验证码
func (*Service) CheckCaptcha(key string, x int, y int) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckCaptcha", &CheckCaptchaParam{
		Key: key,
		X:   x,
		Y:   y,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// CheckSms 判断短信验证码
func (*Service) CheckSms(phone string, code string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckSms", &CheckSmsParam{
		Phone: phone,
		Code:  code,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// CheckMail 判断邮件验证码
func (*Service) CheckMail(mail string, code string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckMail", &CheckMailParam{
		Mail: mail,
		Code: code,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// WebUserLogin 网站用户登陆
func (*Service) WebUserLogin(account string, password string, code string, userAgent string, ip string) (token string, err error) {
	err = rpc.InvokeFunc(constant.ServiceUser, "Login", &WebUserLoginParam{
		Account:   account,
		Password:  password,
		Code:      code,
		UserAgent: userAgent,
		Ip:        ip,
	}, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SystemUserLogin 系统用户登陆
func (*Service) SystemUserLogin(account string, password string, code string, userAgent string, ip string) (token string, err error) {
	err = rpc.InvokeFunc(constant.ServiceAdmin, "Login", &SystemUserLoginParam{
		Account:   account,
		Password:  password,
		Code:      code,
		UserAgent: userAgent,
		Ip:        ip,
	}, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// WebUserLogout 网站用户注销登录
func (*Service) WebUserLogout(token string) error {
	err := rpc.InvokeFunc(constant.ServiceUser, "Logout", token, nil)
	if err != nil {
		return err
	}
	return nil
}

// SystemUserLogout 系统用户注销登录
func (*Service) SystemUserLogout(token string) error {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "Logout", token, nil)
	if err != nil {
		return err
	}
	return nil
}
