package controller

import (
	"service-auth/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
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
	if err := service.Auth.Logout(form.IsAdmin, c.Request.Header.Get(constant.JwtTokenName)); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("注销登录成功")
}
