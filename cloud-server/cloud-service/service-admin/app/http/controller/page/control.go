package page

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
)

func (ControllerPage) Control(c *sinking_server.Context) {
	type Form struct {
		Key string `form:"key" json:"key" default:"" validate:"required" label:"页面ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	user := service.SystemUser.GetCurrentUserInfo(c.GetSystemUserInfo().Id, false)
	role, err := service.SystemRole.FindByIdWithCache(user.RoleId, false)
	if err != nil {
		c.Error("获取权限信息失败")
		return
	}
	page, err := service.SystemMenu.FindByKeyWithCache(form.Key, false)
	if err != nil {
		c.Error("获取页面信息失败")
		return
	}
	controls, err2 := service.SystemMenu.SelectControlWithCache(role.Id, role.Menus, page.Id, false)
	if err2 != nil {
		c.Error("获取控件信息失败")
		return
	}
	c.SuccessWithData("获取页面控件数据成功", controls)
}
