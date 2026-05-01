package menu

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Controls 控件列表
func (Controller) Controls(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"页面ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	list := service.Rbac.GetPageControlsByMenuId(form.Id)
	c.SuccessWithData("获取数据成功", list)
}
