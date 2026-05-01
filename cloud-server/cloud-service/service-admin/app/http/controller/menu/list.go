package menu

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 菜单列表
func (Controller) List(c *context.Context) {
	list := service.Rbac.GetAllMenuTreeByDirAndPage()
	c.SuccessWithData("获取数据成功", list)
}
