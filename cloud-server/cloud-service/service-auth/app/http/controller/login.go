package controller

import (
	"service-auth/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Login 账号登录
func Login(c *context.Context) {
	var form struct {
		Account  string `json:"account" default:"" validate:"required" label:"账号"`
		Password string `json:"password" default:"" validate:"omitempty,min=6,max=20" label:"密码"`
		Code     string `json:"code" default:"" validate:"omitempty,numeric,len=6" label:"验证码"`
		IsAdmin  bool   `json:"is_admin" default:"false" validate:"required,oneof=true false" label:"是否后台"`
		Token    string `json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
		CaptchaX int    `json:"captcha_x" default:"" validate:"required,numeric,min=0,max=999" label:"验证码X坐标"`
		CaptchaY int    `json:"captcha_y" default:"" validate:"required,numeric,min=0,max=999" label:"验证码Y坐标"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if !service.Auth.CheckCaptcha(form.Token, form.CaptchaX, form.CaptchaY) {
		c.Error("验证码验证失败请重试")
		return
	}
	token, err := service.Auth.Login(form.IsAdmin,
		form.Account,
		form.Password,
		form.Code,
		c.Request.Header.Get("user-agent"),
		c.ClientIP(true))
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("登录成功", token)
}
