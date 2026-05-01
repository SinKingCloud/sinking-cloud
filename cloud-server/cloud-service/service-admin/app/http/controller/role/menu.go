package role

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Menu 角色菜单树
func (Controller) Menu(c *context.Context) {
	data := service.Rbac.GetAllMenuTreeByStatusEnabled()
	c.SuccessWithData("获取数据成功", data)
}
