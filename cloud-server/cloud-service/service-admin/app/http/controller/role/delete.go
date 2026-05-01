package role

import (
	"errors"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除角色
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.DeleteRoleByIds(form.Ids)
	if err != nil {
		c.Error(errors.Join(errors.New("删除数据失败"), err))
	} else {
		c.Success("删除数据成功")
	}
}
