package permission

import (
	"service-admin/app/model"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建权限
func (Controller) Create(c *context.Context) {
	var form struct {
		Code string `json:"code" default:"" validate:"required,max=50" label:"标识"`
		Name string `json:"name" default:"" validate:"required,max=20" label:"名称"`
		Path string `json:"path" default:"" validate:"required,max=100" label:"路径"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.CreatePermission(&model.Permission{
		Name: form.Name,
		Code: form.Code,
		Path: form.Path,
	})
	if err == nil {
		c.Success("创建数据成功")
	} else {
		c.Error(err.Error())
	}
}
