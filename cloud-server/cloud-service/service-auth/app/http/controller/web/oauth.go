package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/rpc"
)

// Oauth 授权登录
func (ControllerWeb) Oauth(c *sinking_server.Context) {
	type Form struct {
		Ticket      string `form:"ticket" json:"ticket" default:"" validate:"required,len=32" label:"授权码"`
		RedirectUrl string `form:"redirect_url" json:"redirect_url" default:"" validate:"required,uri" label:"回调网址"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := rpc.Oauth.AuthTicket(c.GetWebUserInfo().Id, form.Ticket, form.RedirectUrl)
	if err == nil {
		c.Success("授权登录成功")
	} else {
		c.Error(err.Error())
	}
}
