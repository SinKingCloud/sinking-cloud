package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/log"
)

// Update 添加信息
func (ControllerApp) Update(c *sinking_server.Context) {
	type Form struct {
		Id       int    `json:"id" default:"" validate:"min=1" label:"ID"`
		Key      string `json:"key" default:"" validate:"omitempty,len=32" label:"密匙"`
		Name     string `json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric,oneof=0 2" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	info, err := app.GetIns().FindById(form.Id)
	if err != nil || info == nil || info.UserId != c.GetWebUserInfo().Id {
		c.Error("应用不存在")
		return
	}
	if info.Status == 1 {
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
	err = app.GetIns().UpdateById(form.Id, data)
	if err == nil {
		log.GetIns().CreateByQueue(c, 3, "修改应用", "修改应用"+info.Name)
		c.Success("修改应用信息成功")
	} else {
		c.Error("修改应用信息失败")
	}
}
