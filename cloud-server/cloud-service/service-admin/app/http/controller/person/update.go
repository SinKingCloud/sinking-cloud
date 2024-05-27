package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

func (ControllerPerson) Update(c *sinking_server.Context) {
	type Form struct {
		Phone    string `form:"phone" json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email    string `form:"email" json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		NickName string `form:"nick_name" json:"nick_name" default:"" validate:"omitempty,max=50" label:"用户昵称"`
		Avatar   string `form:"avatar" json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Phone != "" {
		data["phone"] = form.Phone
	}
	if form.Email != "" {
		data["email"] = form.Email
	}
	if form.NickName != "" {
		data["nick_name"] = form.NickName
	}
	if form.Avatar != "" {
		data["avatar"] = form.Avatar
	}
	if err := service2.SystemUser.UpdateByIds([]int{c.GetSystemUserInfo().Id}, data); err == nil {
		service.GetSystemUserServiceIns().ClearCache(c.GetSystemUserInfo().Id)
		service2.SystemLog.CreateByQueue(c, system_log.EventEdit, "修改资料", "修改账户资料")
		c.Success("更新资料成功")
	} else {
		c.Error("更新资料失败")
	}
}
