package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/user"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strconv"
)

// AuthInfo 获取实名信息
func (ControllerPerson) AuthInfo(c *sinking_server.Context) {
	user := user2.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	configs := service.GetWebUserConfSrvIns().GetConfigs(user.Id)
	c.SuccessWithData("获取实名信息成功", sinking_web.H{
		"real_auth": user.RealAuth,
		"card_no":   str.NetStringTool(configs["person.card_no"]).Show(4, 2, "*"),
		"card_name": str.NetStringTool(configs["person.card_name"]).Show(1, 0, "*"),
	})
}

// AuthInit 初始化实名信息
func (ControllerPerson) AuthInit(c *sinking_server.Context) {
	type Form struct {
		CardNo   string `json:"card_no" default:"" validate:"required,numeric,len=18" label:"身份证号"`
		CardName string `json:"card_name" default:"" validate:"required,min=1,max=4" label:"姓名"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	user := user2.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if user.RealAuth == 1 {
		c.Error("您已实名成功,无需重复实名")
		return
	}
	//写入信息
	service.GetWebUserConfSrvIns().SetConfigs(user.Id, map[string]string{
		"person.card_no":   form.CardNo,
		"person.card_name": form.CardName,
	})
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	id, url, err := aliyun.NewFace(configs[config.AliyunAccessKey], configs[config.AliyunAccessSecret]).
		Init(strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10),
			configs[config.AliyunAuthSceneId],
			"FACE",
			"IDENTITY_CARD",
			form.CardNo,
			form.CardName,
			configs[config.SystemUrl],
			configs[config.SystemGateway]+"/auth/callback/face",
		)
	if err != nil {
		c.Error("发起认证请求失败，请检查资料")
		return
	}
	//写入用户配置
	service.GetWebUserConfSrvIns().SetConfig(user.Id, "person.card_id", id)
	c.SuccessWithData("初始化实名数据成功", sinking_web.H{
		"url": url,
	})
}
