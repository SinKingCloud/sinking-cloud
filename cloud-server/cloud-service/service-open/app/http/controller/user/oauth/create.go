package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_oauth/oauth"
)

// Create 添加信息
func (ControllerOauth) Create(c *sinking_server.Context) {
	type Form struct {
		AppId  int    `form:"app_id" default:"" validate:"numeric" label:"应用ID"`
		Name   string `json:"name" validate:"required,max=20" label:"名称"`
		Domain string `json:"domain" validate:"required,max=200" label:"域名"`
		Status int    `json:"status" validate:"oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	//查询应用信息
	if info, e := web_app.GetIns().FindById(form.AppId); e != nil || info == nil || info.UserId != c.GetWebUserInfo().Id {
		c.Error("绑定应用不存在")
		return
	}
	err := oauth.GetIns().Create(&model.AppOauth{
		UserId: c.GetWebUserInfo().Id,
		AppId:  form.AppId,
		Name:   form.Name,
		Domain: form.Domain,
		Status: form.Status,
	})
	if err == nil {
		c.Success("添加接入网站成功")
	} else {
		c.Error("添加接入网站失败")
	}
}
