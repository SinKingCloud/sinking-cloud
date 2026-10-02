package config

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type Controller struct {
}

// config 处理配置获取和更新
func (Controller) config(c *context.Context, group string, name string, keys []string) {
	var form struct {
		Action string            `json:"action" default:"get" validate:"required,oneof=get set" label:"操作类型"`
		Data   map[string]string `json:"data" default:"" label:"配置数据"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if form.Action == "set" {
		data := map[string]string{}
		for _, key := range keys {
			if value, ok := form.Data[key]; ok {
				data[key] = value
			}
		}
		if len(data) == 0 {
			c.Error("配置数据不能为空")
			return
		}
		service.Config.SetSysConfigs(data)
		_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改"+name, "修改"+name)
		c.Success("更新数据成功")
		return
	}
	configs := service.Config.GetSysConfigs(group, false)
	for _, key := range keys {
		if _, ok := configs[key]; !ok {
			configs[key] = ""
		}
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看"+name, "查看"+name)
	c.SuccessWithData("获取数据成功", configs)
}
