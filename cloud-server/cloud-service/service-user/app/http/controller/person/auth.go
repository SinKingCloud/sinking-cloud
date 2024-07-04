package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/web_user_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// AuthInfo 获取实名信息
func (ControllerPerson) AuthInfo(c *sinking_server.Context) {
	user := service2.User.GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	configs := service.GetWebUserConfSrvIns().GetConfigs(user.Id, web_user_config.Person)
	c.SuccessWithData("获取实名信息成功", sinking_web.H{
		"real_auth": user.RealAuth,
		"card_no":   str.NewStringTool(configs[web_user_config.PersonCardNo]).Show(4, 2, "*"),
		"card_name": str.NewStringTool(configs[web_user_config.PersonCardName]).Show(1, 0, "*"),
	})
}

// AuthInit 初始化实名信息
func (ControllerPerson) AuthInit(c *sinking_server.Context) {
	type Form struct {
		CardType string `form:"card_type" json:"card_type" default:"" validate:"required,oneof=IDENTITY_CARD HOME_VISIT_PERMIT_HK_MC HOME_VISIT_PERMIT_TAIWAN RESIDENCE_PERMIT_HK_MC RESIDENCE_PERMIT_TAIWAN" label:"身份证号"`
		CardNo   string `form:"card_no" json:"card_no" default:"" validate:"required" label:"身份证号"`
		CardName string `form:"card_name" json:"card_name" default:"" validate:"required" label:"姓名"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	user := service2.User.GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if user.RealAuth == 1 {
		c.Error("您已实名成功,无需重复实名")
		return
	}
	id, url, err := service2.AliYun.FaceInit(user.Id, form.CardType, form.CardNo, form.CardName, "", "")
	if err != nil {
		c.Error("发起认证请求失败，请检查资料")
		return
	}
	c.SuccessWithData("初始化实名数据成功", sinking_web.H{
		"id":  id,
		"url": url,
	})
}

// AuthCheck 实人验证结果
func (ControllerPerson) AuthCheck(c *sinking_server.Context) {
	type Form struct {
		CertifyId string `form:"certifyId" json:"certifyId" validate:"required,len=32" label:"认证ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	user := service2.User.GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if user.RealAuth == 1 {
		c.Error("您已实名成功,无需重复实名")
		return
	}
	if res, userId, cardNo, cardName := service2.AliYun.FaceCheck(form.CertifyId); res && userId == c.GetWebUserInfo().Id {
		if service.GetWebUserServiceIns().SetRealAuthStatu(userId, res, cardNo, cardName) == nil {
			c.Success("实名成功")
		} else {
			c.Error("更改实名状态失败")
		}
	} else {
		c.Error("实名失败")
	}
}
