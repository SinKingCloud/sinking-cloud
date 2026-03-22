package controller

import (
	"service-auth/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Mail 发送邮箱验证码
func Mail(c *context.Context) {
	type Form struct {
		Mail     string `json:"mail" default:"" validate:"required,email" label:"邮箱"`
		Token    string `json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
		CaptchaX int    `json:"captcha_x" default:"" validate:"required,numeric" label:"验证码X坐标"`
		CaptchaY int    `json:"captcha_y" default:"" validate:"required,numeric" label:"验证码Y坐标"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if !service.Auth.CheckCaptcha(form.Token, form.CaptchaX, form.CaptchaY) {
		c.Error("滑块验证失败请重试")
		return
	}
	if service.Auth.SendMail(form.Mail) {
		c.Success("验证码已发送")
	} else {
		c.Error("验证码发送失败")
	}
}
