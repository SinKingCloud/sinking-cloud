package security

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Account 修改账号
func (Controller) Account(c *context.Context) {
	var form struct {
		Account string `json:"account" default:"" validate:"required,alphanum,min=5,max=20" label:"账号"`
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
	if u.Account != "" {
		c.Error("您的账户已修改过账号,无需修改")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Account: &form.Account,
	}); err == nil {
		_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改账户账号", "修改账户账号为["+form.Account+"]")
		c.Success("修改账号成功")
	} else {
		c.Error("修改账号失败")
	}
}
