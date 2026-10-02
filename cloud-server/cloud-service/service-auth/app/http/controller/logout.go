package controller

import (
	"service-auth/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Logout 登陆注销
func Logout(c *context.Context) {
	var form struct {
		IsAdmin bool `json:"is_admin" default:"false" label:"用户标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service2.User.GetJwtUserInfo(c)
	if err != nil {
		c.Error(err.Error())
		return
	}
	if err = service.Auth.Logout(form.IsAdmin, u.Token); err != nil {
		c.Error(err.Error())
		return
	}
	if form.IsAdmin {
		_ = service2.User.CreateSystemUserLog(u.Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventLogin, "管理员退出登录", "管理员退出登录")
	}
	c.Success("注销登录成功")
}
