package verify

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service"
)

func (ControllerVerify) Sms(c *sinking_server.Context) {
	type Form struct {
		Phone   string `form:"phone" json:"phone" default:"" validate:"required,numeric,len=11" label:"手机号"`
		Token   string `form:"token" json:"token" default:"" validate:"required" label:"验证码标识"`
		Captcha string `form:"captcha" json:"captcha" default:"" validate:"required" label:"验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if !service.Captcha.Check(form.Token, form.Captcha) {
		c.Error("验证码验证失败请重试")
		return
	}
	if service.Auth.SmsSend(form.Phone) {
		c.Success("验证码已发送")
	} else {
		c.Error("验证码发送失败")
	}
}
