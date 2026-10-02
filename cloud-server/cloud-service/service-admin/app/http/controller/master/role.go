package master

import (
	"service-admin/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Role 分配角色
func (Controller) Role(c *context.Context) {
	var form struct {
		Id      int64   `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
		RoleIds []int64 `json:"role_ids" default:"" validate:"required,unique" label:"角色ID列表"`
	}

	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.User.AssignRole(form.Id, form.RoleIds)
	if err != nil {
		c.Error(err.Error())
	} else {
		_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "分配管理员角色", "分配管理员["+strconv.FormatInt(form.Id, 10)+"]角色")
		c.Success("分配角色成功")
	}
}
