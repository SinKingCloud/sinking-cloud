package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/aliyun"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/user"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// AuthInfo 获取实名信息
func (ControllerPerson) AuthInfo(c *sinking_server.Context) {
	user := user2.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	configs := service.GetWebUserConfSrvIns().GetConfigs(user.Id, constant.UserConfigPerson)
	c.SuccessWithData("获取实名信息成功", sinking_web.H{
		"real_auth": user.RealAuth,
		"card_no":   str.NewStringTool(configs[constant.UserConfigPersonCardNo]).Show(4, 2, "*"),
		"card_name": str.NewStringTool(configs[constant.UserConfigPersonCardName]).Show(1, 0, "*"),
	})
}

// AuthInit 初始化实名信息
func (ControllerPerson) AuthInit(c *sinking_server.Context) {
	type Form struct {
		CardType string `json:"card_type" default:"" validate:"required,oneof=IDENTITY_CARD HOME_VISIT_PERMIT_HK_MC HOME_VISIT_PERMIT_TAIWAN RESIDENCE_PERMIT_HK_MC RESIDENCE_PERMIT_TAIWAN" label:"身份证号"`
		CardNo   string `json:"card_no" default:"" validate:"required" label:"身份证号"`
		CardName string `json:"card_name" default:"" validate:"required" label:"姓名"`
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
	id, url, err := aliyun.FaceInit(user.Id, form.CardType, form.CardNo, form.CardName)
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
		CertifyId string `form:"certifyId" validate:"required,len=32" label:"认证ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	if aliyun.FaceCheck(form.CertifyId) {
		c.Success("实名成功")
	} else {
		c.Error("实名失败")
	}
}
