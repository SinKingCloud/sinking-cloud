package page

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Menu 系统菜单
func (Controller) Menu(c *context.Context) {
	menus := service.Rbac.GetUserMenuTree(c.GetUserInfo().Id)
	if menus == nil {
		c.Error("获取菜单数据失败")
		return
	}
	c.SuccessWithData("获取菜单数据成功", menus)
}
