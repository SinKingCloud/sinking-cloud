package master

import (
	"service-admin/app/service"

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
		c.Success("分配角色成功")
	}
}
