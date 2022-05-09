package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc/web_user_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_app"
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
	err := web_app.GetIns().Create(&model.UserApp{
		UserId:   c.GetWebUserInfo().Id,
		Key:      strings.ToUpper(str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5()),
		Name:     form.Name,
		Describe: form.Describe,
		Status:   form.Status,
	})
	if err == nil {
		rpc.WebUserLog.Create(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), web_user_log.EventCreate, "创建应用", "创建应用信息")
		c.Success("创建应用成功")
	} else {
		c.Error("创建应用失败")
	}
}
