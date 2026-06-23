package security

import (
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/real_auth_type"
	"service-user/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Certificate 实名认证
func (Controller) Certificate(c *context.Context) {
	var form struct {
		Action      string `json:"action" default:"info" validate:"required,oneof=info init verify" label:"步骤"`
		Type        string `json:"type" default:"" validate:"omitempty,numeric" label:"认证类型"`
		CardName    string `json:"card_name" default:"" validate:"omitempty,min=1,max=20" label:"身份证姓名"`
		CardNo      string `json:"card_no" default:"" validate:"omitempty,min=15,max=18" label:"身份证号码"`
		CompanyName string `json:"company_name" default:"" validate:"omitempty,min=1,max=100" label:"企业名称"`
		CompanyNo   string `json:"company_no" default:"" validate:"omitempty,min=1,max=50" label:"企业编码"`
		Env         string `json:"env" default:"" validate:"omitempty" label:"环境信息"`
		Id          string `json:"id" default:"" validate:"omitempty" label:"认证ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	userId := c.GetUserInfo().Id
	switch form.Action {
	case "info":
		info, err := service.User.GetRealAuthInfo(userId, false)
		if err != nil {
			c.Error(err.Error())
			return
		}
		if info == nil {
			c.Error("未获取到实名信息")
			return
		}
		data := sinking_web.H{
			"status": info.Status == real_auth_status.Authed,
		}
		if info.Status == real_auth_status.Authed {
			data["auth_time"] = info.AuthTime.String()
			data["card_name"] = str.NewStringTool(info.CardName).Show(1, 0, "*")
			data["card_no"] = str.NewStringTool(info.CardNo).Show(4, 2, "*")
			if info.Type == real_auth_type.Company {
				data["company_name"] = info.CompanyName
				data["company_no"] = info.CompanyNo
			}
		}
		c.SuccessWithData("获取数据成功", data)
	case "init":
		authType, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("认证类型不合法")
			return
		}
		certifyId, certifyUrl, err := service.User.InitRealAuth(userId, authType, form.CardName, form.CardNo, form.CompanyName, form.CompanyNo, form.Env)
		if err != nil {
			c.Error(err.Error())
			return
		}
		c.SuccessWithData("发起实名成功", sinking_web.H{
			"id":  certifyId,
			"url": certifyUrl,
		})
	case "verify":
		data, result, err := service.User.VerifyRealAuth(form.Id)
		status := false
		breakLoop := true
		message := "实名认证成功"
		if data != nil {
			status = data.Status == real_auth_status.Authed
		}
		if err != nil {
			message = err.Error()
			if result != nil {
				breakLoop = result.Code != "403" && result.Code != "424"
				if result.Reason != "" {
					message = result.Reason
				}
			}
		}
		c.SuccessWithData("获取实名状态成功", sinking_web.H{
			"status":  status,
			"break":   breakLoop,
			"message": message,
		})
	default:
		c.Error("步骤不合法")
	}
}
