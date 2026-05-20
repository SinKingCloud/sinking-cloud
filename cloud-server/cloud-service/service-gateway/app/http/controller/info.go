package controller

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Info 网站信息
func Info(c *context.Context) {
	config := service.Config.GetSystemConfigs([]string{constant.WebGroup, constant.UiGroup}, false)
	web := config[constant.WebGroup]
	ui := config[constant.UiGroup]
	c.SuccessWithData("获取成功", sinking_web.H{
		"title":    c.GetStringWithDefault(web[constant.WebTitle], "默认标题"),
		"name":     c.GetStringWithDefault(web[constant.WebName], "默认名称"),
		"keywords": c.GetStringWithDefault(web[constant.WebKeyWords], "默认名称"),
		"describe": c.GetStringWithDefault(web[constant.WebDescribe], "默认名称"),
		"contact":  c.GetStringWithDefault(web[constant.WebContact], ""),
		"url":      c.GetStringWithDefault(web[constant.WebUrl], ""),
		"api":      c.GetStringWithDefault(web[constant.WebApi], ""),
		"logo":     c.GetStringWithDefault(web[constant.WebLogo], ""),
		"ui": sinking_web.H{
			"layout":    c.GetStringWithDefault(ui[constant.UiLayout], "left"),
			"watermark": c.GetBoolWithDefault(ui[constant.UiWaterMark], false),
			"theme":     c.GetStringWithDefault(ui[constant.UiTheme], "light"),
			"compact":   c.GetBoolWithDefault(ui[constant.UiCompact], false),
			"color":     c.GetStringWithDefault(ui[constant.UiColor], "rgb(0, 83, 251)"),
			"radius":    c.GetIntWithDefault(ui[constant.UiRadius], 0),
		},
	})
}
