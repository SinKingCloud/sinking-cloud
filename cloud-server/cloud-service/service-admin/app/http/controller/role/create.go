package role

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建角色
func (Controller) Create(c *context.Context) {
	var form struct {
		Name    string  `json:"name" default:"" validate:"required,max=30" label:"名称"`
		Status  int     `json:"status" default:"0" validate:"required,numeric" label:"状态"`
		MenuIds []int64 `json:"menu_ids" default:"" validate:"omitempty,unique" label:"菜单ID列表"`
	}

	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.CreateRole(form.Name, form.Status, form.MenuIds)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventCreate, "创建角色", "创建角色["+form.Name+"]")
	c.Success("创建数据成功")
}
