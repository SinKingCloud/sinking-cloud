package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc"
)

// Delete 删除接口信息
func (ControllerApp) Delete(c *sinking_server.Context) {
	type Form struct {
		Id int `json:"id" default:"" validate:"min=1" label:"ID"`
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
	err = app.GetIns().DeleteById(form.Id)
	if err == nil {
		rpc.Log.Create(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), 2, "删除应用", "删除应用"+info.Name)
		c.Success("删除应用成功")
	} else {
		c.Error("删除应用失败")
	}
}
