package verify

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service"
	"strconv"
)

func (ControllerVerify) Sms(c *sinking_server.Context) {
	type Form struct {
		Phone   string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Token   string `json:"token" default:"" validate:"required" label:"验证码标识"`
		Captcha string `json:"captcha" default:"" validate:"required" label:"验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Phone == "" {
		u := service.User.GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
		form.Phone = strconv.FormatInt(u.Phone, 10)
	}
	if !rpc.Captcha.Check(form.Token, form.Captcha) {
		c.Error("验证码验证失败请重试")
		return
	}
	if rpc.Sms.Send(form.Phone) {
		c.Success("验证码已发送")
	} else {
		c.Error("验证码发送失败")
	}
}
