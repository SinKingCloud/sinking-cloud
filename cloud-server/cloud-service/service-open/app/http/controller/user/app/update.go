package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_app"
)

// Update 更新信息
func (ControllerApp) Update(c *sinking_server.Context) {
	type Form struct {
		Id       int    `json:"id" json:"id" default:"" validate:"min=1" label:"ID"`
		Key      string `json:"key" json:"key" default:"" validate:"omitempty,len=32" label:"密匙"`
		Name     string `json:"name" json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `json:"describe" json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Status   string `json:"status" json:"status" default:"" validate:"omitempty,numeric,oneof=0 2" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	info, err := web_app.GetIns().FindById(form.Id)
	if err != nil || info == nil || info.UserId != c.GetWebUserInfo().Id {
		c.Error("应用不存在")
		return
	}
	if info.Status == web_app.Error {
		c.Error("该应用已被封禁，无法操作")
		return
	}
	data := make(map[string]interface{})
	if form.Key != "" {
		data["key"] = form.Key
	}
	if form.Name != "" {
		data["name"] = form.Name
	}
	if form.Describe != "" {
		data["describe"] = form.Describe
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err = web_app.GetIns().UpdateById(form.Id, data)
	if err == nil {
		rpc.Log.CreateUser(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), log.UserEventEdit, "修改应用", "修改应用"+info.Name)
		c.Success("修改应用信息成功")
	} else {
		c.Error("修改应用信息失败")
	}
}
