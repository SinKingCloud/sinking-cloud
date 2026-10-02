package security

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Logout 注销会话
func (Controller) Logout(c *context.Context) {
	var form struct {
		Tokens []string `json:"tokens" default:"" validate:"required,min=1,max=1000,unique" label:"会话token列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Auth.LogoutByUserIdAndTokens(c.GetUserInfo().Id, form.Tokens)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventDelete, "注销账户会话", "注销账户会话")
	c.Success("注销会话成功")
}
