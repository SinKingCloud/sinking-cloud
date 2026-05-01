package permission

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除权限
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.DeletePermissionByIds(form.Ids)
	if err != nil {
		c.Error("删除数据失败")
	} else {
		c.Success("删除数据成功")
	}
}
