package page

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/menu"
	role2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/role"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/user"
)

func (ControllerPage) Control(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" json:"id" default:"" validate:"required" label:"页面ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	user := user2.GetCurrentUserInfo(c.GetSystemUserInfo().Id, false)
	role := role2.GetRoleServiceIns().FindByIdWithCache(user.RoleId, false)
	c.SuccessWithData("获取页面控件数据成功", menu.SelectRoleControlWithCache(role.Id, role.Menus, form.Id, false))
}
