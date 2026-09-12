package security

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Phone 修改手机
func (Controller) Phone(c *context.Context) {
	var form struct {
		Phone string `json:"phone" default:"" validate:"required,numeric,len=11" label:"手机号"`
		Code  string `json:"code" default:"" validate:"required,numeric,len=6" label:"验证码"`
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
	if !service2.Auth.CheckSms(form.Phone, form.Code) {
		c.Error("验证码错误")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Phone: &form.Phone,
	}); err == nil {
		c.Success("修改手机号成功")
	} else {
		c.Error("修改手机号失败")
	}
}
