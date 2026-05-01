package security

import (
	"service-user/app/repository/user"
	"service-user/app/service"

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
	if err = service.User.CheckAccountIsUse(form.Account); err != nil {
		c.Error(err.Error())
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Account: form.Account,
	}); err == nil {
		c.Success("修改账号成功")
	} else {
		c.Error("修改账号失败")
	}
}
