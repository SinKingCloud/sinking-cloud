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
		ANY("/config", "支付配置", controller.Admin.Config.Pay)

		ANY("/channel/list", "支付通道列表", controller.Admin.Channel.List)
		ANY("/channel/update", "更新支付通道", controller.Admin.Channel.Update)
		ANY("/channel/create", "创建支付通道", controller.Admin.Channel.Create)
		ANY("/channel/delete", "删除支付通道", controller.Admin.Channel.Delete)
		ANY("/channel/info", "支付通道详情", controller.Admin.Channel.Info)

		ANY("/balance/info", "用户余额详情", controller.Admin.Balance.Info)
		ANY("/balance/operate", "用户余额操作", controller.Admin.Balance.Operate)
		ANY("/balance/log", "用户余额明细", controller.Admin.Balance.Log)

		ANY("/order/list", "系统订单列表", controller.Admin.Order.List)
		ANY("/order/info", "系统订单详情", controller.Admin.Order.Info)
		ANY("/order/deliver", "系统订单发货", controller.Admin.Order.Deliver)
		ANY("/order/refound", "系统订单退款", controller.Admin.Order.Refound)
		ANY("/order/clear", "系统订单清理", controller.Admin.Order.Clear)

		ANY("/log/list", "订单日志列表", controller.Admin.Log.List)
		ANY("/log/info", "订单日志详情", controller.Admin.Log.Info)
	})
}
