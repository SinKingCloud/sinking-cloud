package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_oauth"
)

type ControllerOauth struct {
}

// GenTicket 生成ticket
func (ControllerOauth) GenTicket(c *sinking_server.Context) {
	type Form struct {
		OauthId int `json:"oauth_id" default:"" validate:"min=1" label:"OauthID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	oauthInfo, err := web_oauth.GetOauthInfo(form.OauthId, c.GetWebUserInfo().Id)
	if err != nil || oauthInfo == nil {
		c.Error(err.Error())
		return
	}
	//生成授权码
	ticket, err := web_oauth.GenTicket(oauthInfo.Id)
	if err != nil {
		c.Error("生成授权码失败")
	} else {
		c.SuccessWithData("生成授权码成功", ticket)
	}
}

// GetUserInfo 获取授权用户信息
func (ControllerOauth) GetUserInfo(c *sinking_server.Context) {
	type Form struct {
		OauthId int    `json:"oauth_id" default:"" validate:"min=1" label:"OauthID"`
		Ticket  string `json:"ticket" validate:"required,len=32" label:"授权码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	oauthInfo, err := web_oauth.GetOauthInfo(form.OauthId, c.GetWebUserInfo().Id)
	if err != nil || oauthInfo == nil {
		c.Error(err.Error())
		return
	}
	data, err := web_oauth.GetUserInfo(oauthInfo.Id, form.Ticket)
	if err != nil {
		c.Error(err.Error())
	} else {
		c.SuccessWithData("获取授权用户信息成功", data)
	}
}
