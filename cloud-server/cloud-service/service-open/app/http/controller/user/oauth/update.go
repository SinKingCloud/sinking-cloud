package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_oauth/oauth"
	"strconv"
)

// Update 更新信息
func (ControllerOauth) Update(c *sinking_server.Context) {
	type Form struct {
		Id     int    `json:"id" default:"" validate:"min=1" label:"ID"`
		AppId  string `json:"app_id" default:"" validate:"omitempty,numeric,min=1" label:"应用ID"`
		Name   string `json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Domain string `json:"domain" default:"" validate:"omitempty,max=100" label:"域名"`
		Status string `json:"status" default:"" validate:"omitempty,numeric,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	info, err := oauth.GetIns().FindById(form.Id)
	if err != nil || info == nil || info.UserId != c.GetWebUserInfo().Id {
		c.Error("接入网站不存在")
		return
	}
	if info.Status == oauth.Error {
		c.Error("该接入网站已被封禁，无法操作")
		return
	}
	data := make(map[string]interface{})
	if form.AppId != "" {
		appId, _ := strconv.Atoi(form.AppId)
		if appId != info.AppId {
			//查询应用信息
			if info, e := web_app.GetIns().FindByIdWithCache(appId, false); e != nil || info == nil || info.UserId != c.GetWebUserInfo().Id || info.Status != web_app.Normal {
				c.Error("绑定应用不存在")
				return
			}
			data["app_id"] = form.AppId
		}
	}
	if form.Name != "" {
		data["name"] = form.Name
	}
	if form.Domain != "" {
		data["domain"] = form.Domain
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err = oauth.GetIns().UpdateById(form.Id, data)
	if err == nil {
		c.Success("修改接入网站信息成功")
	} else {
		c.Error("修改接入网站信息失败")
	}
}
