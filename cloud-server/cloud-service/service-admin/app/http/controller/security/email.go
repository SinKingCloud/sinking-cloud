package security

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Email 修改邮箱
func (Controller) Email(c *context.Context) {
	var form struct {
		Mail string `json:"mail" default:"" validate:"required,email" label:"邮箱"`
		Code string `json:"code" default:"" validate:"required,numeric,len=6" label:"验证码"`
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
	if !service.User.CheckUserVerifyStatus(u.Id, false) {
		c.Error("请先完成安全验证")
		return
	}
	if !service2.Auth.CheckMail(form.Mail, form.Code) {
		c.Error("验证码错误")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Email: &form.Mail,
	}); err == nil {
		c.Success("修改邮箱成功")
	} else {
		c.Error("修改邮箱失败")
	}
}
