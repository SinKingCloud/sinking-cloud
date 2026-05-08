package controller

import (
	"service-auth/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Sms 发送短信验证码
func Sms(c *context.Context) {
	var form struct {
		Phone    string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		IsAdmin  bool   `json:"is_admin" default:"false" label:"用户标识"`
		Token    string `json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
		CaptchaX int    `json:"captcha_x" default:"" validate:"required,numeric,min=0,max=999" label:"验证码X坐标"`
		CaptchaY int    `json:"captcha_y" default:"" validate:"required,numeric,min=0,max=999" label:"验证码Y坐标"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if !service.Auth.CheckCaptcha(form.Token, form.CaptchaX, form.CaptchaY) {
		c.Error("滑块验证失败请重试")
		return
	}
	if form.Phone == "" {
		u, err := service2.User.GetJwtUserInfo(c)
		if err != nil {
			c.Error(err.Error())
			return
		}
		if form.IsAdmin {
			if user, err := service2.User.GetWebUserInfo(u.Id, false); err == nil && user != nil {
				form.Phone = user.Phone
			}
		} else {
			if user, err := service2.User.GetSystemUserInfo(u.Id, false); err == nil && user != nil {
				form.Phone = user.Phone
			}
		}
		if form.Phone == "" {
			c.Error("获取绑定手机号失败")
			return
		}
	}
	if service.Auth.SendSms(form.Phone) {
		c.Success("验证码已发送")
	} else {
		c.Error("验证码发送失败")
	}
}
