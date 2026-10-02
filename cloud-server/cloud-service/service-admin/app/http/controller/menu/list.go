package menu

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 菜单列表
func (Controller) List(c *context.Context) {
	list := service.Rbac.GetAllMenuTreeByDirAndPage()
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看菜单列表", "查看菜单列表")
	c.SuccessWithData("获取数据成功", list)
}
