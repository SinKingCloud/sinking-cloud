package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_user"
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
	u := system_user.GetIns().GetCurrentUserInfo(c.GetSystemUserInfo().Id, true)
	if !str.NewStringTool(u.Password).CheckPassword(form.OldPassword) {
		c.Error("旧密码验证失败")
		return
	}
	if err := system_user.GetIns().UpdateByIds([]int{u.Id}, map[string]interface{}{
		"password": str.NewStringTool(form.Password).GetPassword(),
	}); err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventEdit, "修改密码", "修改账户密码")
		c.Success("修改密码成功")
	} else {
		c.Error("修改密码失败")
	}
}
