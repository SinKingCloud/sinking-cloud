package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/user"
)

func (ControllerPerson) Password(c *sinking_server.Context) {
	type Form struct {
		OldPassword string `json:"old_password" default:"" validate:"required" label:"旧密码"`
		Password    string `json:"password" default:"" validate:"required" label:"新密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	u := user.GetIns().GetCurrentUserInfo(c.GetSystemUserInfo().Id, true)
	if !str.NetStringTool(u.Password).CheckPassword(form.OldPassword) {
		c.Error("旧密码验证失败")
		return
	}
	if err := user.GetIns().UpdateByIds([]int{u.Id}, map[string]interface{}{
		"password": str.NetStringTool(form.Password).GetPassword(),
	}); err == nil {
		c.Success("修改密码成功")
	} else {
		c.Error("修改密码失败")
	}
}
