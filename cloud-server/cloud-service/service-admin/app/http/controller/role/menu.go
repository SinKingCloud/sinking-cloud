package role

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Menu 角色菜单树
func (Controller) Menu(c *context.Context) {
	data := service.Rbac.GetAllMenuTreeByStatusEnabled()
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看角色可选菜单", "查看角色可选菜单列表")
	c.SuccessWithData("获取数据成功", data)
}
