package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/server"
)

type ControllerServer struct {
}

func (ControllerServer) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		Key string `json:"key" default:"" validate:"omitempty,max=32" label:"地区标识"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["pid"] = "0"
	where["status"] = "0"
	if form.Key != "" {
		where["key"] = form.Key
	}
	data, total, err := server.GetIns().SelectServers(where, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取服务器列表失败")
	} else {
		c.SuccessWithData("获取服务器列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}
