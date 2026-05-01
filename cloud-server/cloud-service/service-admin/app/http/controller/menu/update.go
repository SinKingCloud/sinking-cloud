package menu

import (
	"service-admin/app/repository/menu"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新菜单或控件
func (Controller) Update(c *context.Context) {
	var form struct {
		Id              int64     `json:"id" validate:"required,min=1" label:"记录ID"`
		ParentId        string    `json:"parent_id" validate:"omitempty,numeric,gte=0" label:"父级"`
		Name            string    `json:"name" validate:"omitempty,max=100" label:"名称"`
		Icon            string    `json:"icon" validate:"omitempty,max=30" label:"图标"`
		Sort            string    `json:"sort" validate:"omitempty,numeric" label:"排序"`
		Status          string    `json:"status" validate:"omitempty,numeric" label:"状态"`
		PermissionCodes *[]string `json:"permission_codes" validate:"omitempty,unique" label:"权限标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &menu.UpdateMenu{}
	if form.ParentId != "" {
		data.ParentId, _ = strconv.ParseInt(form.ParentId, 10, 64)
	}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Icon != "" {
		data.Icon = form.Icon
	}
	if form.Sort != "" {
		data.Sort, _ = strconv.Atoi(form.Sort)
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Rbac.UpdateMenuOrControl(form.Id, data, form.PermissionCodes)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}
