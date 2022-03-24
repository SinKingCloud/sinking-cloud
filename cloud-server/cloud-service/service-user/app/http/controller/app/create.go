package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/log"
	"strings"
)

// Create 添加接口信息
func (ControllerApp) Create(c *sinking_server.Context) {
	type Form struct {
		Name     string `json:"name" validate:"required,max=20" label:"名称"`
		Describe string `json:"describe" validate:"required,max=200" label:"请求地址"`
		Status   int    `json:"status" validate:"oneof=0 2" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := app.GetIns().Create(&model.WebUserApp{
		UserId:   c.GetWebUserInfo().Id,
		Key:      strings.ToUpper(str.NetStringTool(str.GetSnowWorkIns().GetUuid()).Md5()),
		Name:     form.Name,
		Describe: form.Describe,
		Status:   form.Status,
	})
	if err == nil {
		log.GetIns().CreateByQueue(c, 4, "创建应用", "创建应用信息")
		c.Success("创建应用成功")
	} else {
		c.Error("创建应用失败")
	}
}
