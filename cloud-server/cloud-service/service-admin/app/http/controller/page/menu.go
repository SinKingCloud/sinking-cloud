package page

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_menu"
	role2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_role"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_user"
)

func (ControllerPage) Menu(c *sinking_server.Context) {
	user := user2.GetIns().GetCurrentUserInfo(c.GetSystemUserInfo().Id, false)
	role, err := role2.GetIns().FindByIdWithCache(user.RoleId, false)
	if err != nil {
		c.Error("获取权限信息失败")
		return
	}
	menus, err2 := system_menu.GetIns().SelectByIdsWithCache(role.Id, role.Menus, false)
	if err2 != nil {
		c.Error("获取菜单数据失败")
		return
	}
	c.SuccessWithData("获取菜单数据成功", menus)
}
