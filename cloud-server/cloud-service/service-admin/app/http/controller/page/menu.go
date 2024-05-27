package page

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
)

func (ControllerPage) Menu(c *sinking_server.Context) {
	user := service.SystemUser.GetCurrentUserInfo(c.GetSystemUserInfo().Id, false)
	role, err := service.SystemRole.FindByIdWithCache(user.RoleId, false)
	if err != nil {
		c.Error("获取权限信息失败")
		return
	}
	menus, err2 := service.SystemMenu.SelectByIdsWithCache(role.Id, role.Menus, false)
	if err2 != nil {
		c.Error("获取菜单数据失败")
		return
	}
	c.SuccessWithData("获取菜单数据成功", menus)
}
