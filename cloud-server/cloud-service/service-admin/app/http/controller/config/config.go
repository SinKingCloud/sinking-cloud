package config

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type Controller struct {
}

// Ui 界面配置
func (s Controller) Ui(c *context.Context) {
	s.config(c, constant.UiGroup, []string{
		constant.UiLayout,
		constant.UiWaterMark,
		constant.UiTheme,
		constant.UiCompact,
		constant.UiColor,
		constant.UiRadius,
	})
}

// Web 网站配置
func (s Controller) Web(c *context.Context) {
	s.config(c, constant.WebGroup, []string{
		constant.WebUrl,
		constant.WebName,
		constant.WebTitle,
		constant.WebKeyWords,
		constant.WebDescribe,
		constant.WebContact,
		constant.WebLogo,
	})
}

// Sms 短信配置
func (s Controller) Sms(c *context.Context) {
	s.config(c, constant.AliGroup, []string{
		constant.AliKey,
		constant.AliSecret,
		constant.AliSmsSign,
		constant.AliSmsTemplate,
		constant.AliSmsVar,
	})
}

// config 处理配置获取和更新
func (Controller) config(c *context.Context, group string, keys []string) {
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
		c.Success("更新数据成功")
		return
	}
	configs := service.Config.GetSysConfigs(group, false)
	for _, key := range keys {
		if _, ok := configs[key]; !ok {
			configs[key] = ""
		}
	}
	c.SuccessWithData("获取数据成功", configs)
}
