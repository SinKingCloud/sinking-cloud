package route

import (
	"service-pay/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadAdmin 加载admin控制器
func loadAdmin(s *sinking_web.Engine) {
	g := service.Auth.GetAdminRouteGroup(s)
	service.Auth.LoadAdminRouteWithRegister(g, func(ANY func(string, string, func(*context.Context))) {
		ANY("/channel/list", "通道列表", controller.Admin.Channel.List)
		ANY("/channel/update", "更新通道", controller.Admin.Channel.Update)
		ANY("/channel/create", "创建通道", controller.Admin.Channel.Create)
		ANY("/channel/delete", "删除通道", controller.Admin.Channel.Delete)
		ANY("/channel/info", "通道详情", controller.Admin.Channel.Info)

		ANY("/balance/info", "余额详情", controller.Admin.Balance.Info)
		ANY("/balance/operate", "余额操作", controller.Admin.Balance.Operate)
		ANY("/balance/log", "余额明细", controller.Admin.Balance.Log)

		ANY("/order/list", "订单列表", controller.Admin.Order.List)
		ANY("/order/log", "订单日志", controller.Admin.Order.Log)
	})
}
