package permission

import (
	"service-admin/app/repository/permission"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 修改权限
func (Controller) Update(c *context.Context) {
	var form struct {
		Ids  []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Name string  `json:"name" default:"" validate:"omitempty,max=100" label:"名称"`
		Path string  `json:"path" default:"" validate:"omitempty,max=200" label:"描述"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &permission.UpdatePermission{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Path != "" {
		data.Path = form.Path
	}
	err := service.Rbac.UpdatePermissionByIds(form.Ids, data)
	if err != nil {
		c.Error("修改数据失败")
	} else {
		c.Success("修改数据成功")
	}
}
