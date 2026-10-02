package person

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新信息
func (Controller) Update(c *context.Context) {
	var form struct {
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"用户昵称"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &user.UpdateUser{}
	if form.Avatar != "" {
		data.Avatar = &form.Avatar
	}
	if form.NickName != "" {
		data.NickName = &form.NickName
	}
	if err := service.User.UpdateById(c.GetUserInfo().Id, data); err == nil {
		_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改账户资料", "修改账户资料")
		c.Success("修改成功")
	} else {
		c.Error("修改失败")
	}
}
