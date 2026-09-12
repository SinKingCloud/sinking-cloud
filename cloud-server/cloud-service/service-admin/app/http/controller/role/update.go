package role

import (
	"service-admin/app/repository/role"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 修改角色
func (Controller) Update(c *context.Context) {
	var form struct {
		Id      int64    `json:"id" default:"" validate:"required,min=1" label:"ID列表"`
		Name    string   `json:"name" default:"" validate:"omitempty,max=30" label:"名称"`
		Status  string   `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		MenuIds *[]int64 `json:"menu_ids" validate:"omitempty,unique" label:"菜单ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &role.UpdateRole{}
	if form.Name != "" {
		data.Name = &form.Name
	}
	if form.Status != "" {
		value, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &value
	}
	err := service.Rbac.UpdateRole(form.Id, data, form.MenuIds)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}
