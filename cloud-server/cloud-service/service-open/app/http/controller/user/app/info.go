package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	app2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_app"
)

// Info 应用信息
func (ControllerApp) Info(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	data, err := app2.GetIns().FindById(form.Id)
	if err == nil && data != nil {
		c.SuccessWithData("获取数据成功", data)
	} else {
		c.Error("获取数据失败")
	}
}
