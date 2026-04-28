package route

import (
	"service-pay/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载user控制器
func loadUser(s *sinking_web.Engine) {
	g := service.Auth.GetUserRouteGroup(s)
	{
		balance := g.Group("/balance")
		{
			balance.ANY("/info", context.HandleFunc(controller.User.Balance.Info))         //余额信息
			balance.ANY("/log", context.HandleFunc(controller.User.Balance.Log))           //余额日志
			balance.ANY("/recharge", context.HandleFunc(controller.User.Balance.Recharge)) //在线充值
		}
		order := g.Group("/order")
		{
			order.ANY("/list", context.HandleFunc(controller.User.Order.List)) //订单列表
			order.ANY("/log", context.HandleFunc(controller.User.Order.Log))   //订单日志
		}
	}
}
