package system

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/captcha"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/rpc"
)

func (ControllerSystem) Login(c *sinking_server.Context) {
	type Form struct {
		Account  string `form:"account" json:"account" default:"" validate:"required" label:"账户"`
		Password string `form:"password" json:"password" default:"" validate:"required" label:"密码"`
		Device   string `form:"device" json:"device" default:"web" validate:"required,oneof=web pc mobile android" label:"登陆设备"`
		Token    string `form:"token" json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
		Captcha  string `form:"captcha" json:"captcha" default:"" validate:"required,numeric,len=4" label:"验证码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	if !captcha.Check(form.Token, form.Captcha) {
		c.Error("验证码验证失败请重试")
		return
	}
	user, err := auth.GetSystemUserInfoByPwd(form.Account, form.Password)
	if err != nil {
		c.Error(err.Error())
		return
	}
	token, err := auth.GetJwtToken(form.Device, user, nil, c.GetRequestIp())
	if err != nil {
		c.Error(err.Error())
		return
	}
	//写入登录日志
	rpc.Log.CreateAdmin(user.Id, c.GetRequestId(), c.GetRequestIp(), 0, "账户登录", "账户登陆成功")
	c.SuccessWithData("登陆成功", token)
}
