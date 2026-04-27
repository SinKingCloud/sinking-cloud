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
	balance := g.Group("/balance")
	{
		balance.ANY("/info", context.HandleFunc(controller.User.Balance.Info))
		balance.ANY("/log", context.HandleFunc(controller.User.Balance.Log))
		balance.ANY("/recharge", context.HandleFunc(controller.User.Balance.Recharge))
	}
	order := g.Group("/order")
	{
		order.ANY("/list", context.HandleFunc(controller.User.Order.List))
		order.ANY("/log", context.HandleFunc(controller.User.Order.Log))
	}
}
