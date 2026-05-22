package security

import (
	"fmt"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Verify 安全验证
func (Controller) Verify(c *context.Context) {
	var form struct {
		Action string `json:"action" default:"get" validate:"required,oneof=get verify" label:"步骤"`
		Type   string `json:"type" default:"phone" validate:"omitempty,oneof=mail phone" label:"验证方式"`
		Code   string `json:"code" default:"" validate:"omitempty,numeric,len=6" label:"验证码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if form.Action == "get" {
		data := sinking_web.H{
			"status": service.User.CheckUserVerifyStatus(u.Id, false),
		}
		if u.Phone != "" {
			data["phone"] = str.NewStringTool(u.Phone).Show(3, 4, "*")
		}
		if u.Email != "" {
			data["mail"] = str.NewStringTool(u.Email).Show(3, 8, "*")
		}
		c.SuccessWithData("获取验证方式成功", data)
		return
	}
	if form.Code == "" {
		c.Error("请输入验证码")
		return
	}
	account := u.Email
	check := service2.Auth.CheckMail
	if form.Type == "phone" {
		check = service2.Auth.CheckSms
		account = u.Phone
	}
	if account == "" {
		c.Error("当前验证方式未绑定")
		return
	}
	if !check(account, form.Code) {
		c.Error("验证码校验失败")
		return
	}
	minute := 10
	service.User.SetUserVerifyStatus(u.Id, true, minute)
	c.Success(fmt.Sprintf("验证成功,此次验证%d分钟内有效", minute))
}
