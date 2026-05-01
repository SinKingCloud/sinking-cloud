package menu

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建菜单或控件
func (Controller) Create(c *context.Context) {
	var form struct {
		ParentId        int64    `json:"parent_id" default:"0" validate:"numeric,gte=0" label:"父级"`
		Code            string   `json:"code" default:"" validate:"required,max=100" label:"标识"`
		Name            string   `json:"name" default:"" validate:"required,max=100" label:"名称"`
		Icon            string   `json:"icon" default:"" validate:"omitempty,max=30" label:"图标"`
		Type            int      `json:"type" default:"0" validate:"numeric" label:"类型"`
		Sort            int      `json:"sort" default:"0" validate:"numeric" label:"排序"`
		Status          int      `json:"status" default:"0" validate:"numeric" label:"状态"`
		PermissionCodes []string `json:"permission_codes" default:"" validate:"omitempty,unique" label:"权限标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.CreateMenuOrControl(form.ParentId, form.Code, form.Name, form.Icon, form.Type, form.Sort, form.Status, form.PermissionCodes)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}
