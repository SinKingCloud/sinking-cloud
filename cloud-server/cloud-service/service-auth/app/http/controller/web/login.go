package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/captcha"
)

func (ControllerWeb) Login(c *sinking_server.Context) {
	type Form struct {
		Phone    string `form:"phone" json:"phone" default:"" validate:"-" label:"手机号"`
		Code     string `form:"code" json:"code" default:"" validate:"-" label:"短信验证码"`
		Account  string `form:"account" json:"account" default:"" validate:"-" label:"账户"`
		Password string `form:"password" json:"password" default:"" validate:"-" label:"密码"`
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
	var user *model.WebUser
	var err error
	if form.Phone != "" {
		//短信登陆
		if form.Code == "" {
			c.Error("请输入短信验证码")
			return
		}
		user, err = auth.GetWebUserInfoByPhone(form.Phone, form.Code)
	} else {
		//密码登陆
		if form.Account == "" || form.Password == "" {
			c.Error("请输入账户密码")
			return
		}
		user, err = auth.GetWebUserInfoByPwd(form.Account, form.Password)
	}
	if err != nil {
		c.Error(err.Error())
		return
	}
	token, err := auth.GetJwtToken(form.Device, nil, user, c.GetRequestIp())
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("登陆成功", token)
}
