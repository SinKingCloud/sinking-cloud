package channel

import (
	"service-pay/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建数据
func (Controller) Create(c *context.Context) {
	var form struct {
		Name   string `json:"name" default:"" validate:"required" label:"通道名称"`
		Type   string `json:"type" default:"" validate:"required" label:"通道类型"`
		Plugin int    `json:"plugin" default:"" label:"使用插件"`
		Status int    `json:"status" default:"" label:"通道状态"`
		Config string `json:"config" default:"" label:"通道配置"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.Create(form.Name, form.Type, form.Plugin, form.Status, form.Config)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventCreate, "创建支付通道", "创建支付通道["+form.Name+"]")
	c.Success("创建数据成功")
}
