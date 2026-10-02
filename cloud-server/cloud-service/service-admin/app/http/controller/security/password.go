package security

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Password 修改密码
func (Controller) Password(c *context.Context) {
	var form struct {
		Password string `json:"password" default:"" validate:"required,min=6,max=20" label:"新密码"`
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
	password := str.NewStringTool(form.Password).GetPassword()
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Password: &password,
	}); err == nil {
		_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改账户密码", "修改账户密码")
		c.Success("修改密码成功")
	} else {
		c.Error("修改密码失败")
	}
}
