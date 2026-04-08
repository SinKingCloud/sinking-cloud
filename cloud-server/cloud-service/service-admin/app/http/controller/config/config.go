package config

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type Controller struct {
}

// Get 获取配置
func (Controller) Get(c *context.Context) {
	var form struct {
		Group string `json:"group" default:"" validate:"omitempty,max=100" label:"配置分组"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := service.Config.GetSysConfigs(form.Group, false)
	c.SuccessWithData("获取数据成功", data)
}

// Set 更新配置
func (Controller) Set(c *context.Context) {
	var form struct {
		Data map[string]string `json:"data" default:"" label:"配置数据"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if len(form.Data) == 0 {
		c.Error("配置数据不能为空")
		return
	}
	service.Config.SetSysConfigs(form.Data)
	c.Success("更新数据成功")
}
