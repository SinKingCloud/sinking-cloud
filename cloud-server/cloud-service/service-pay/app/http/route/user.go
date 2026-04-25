package route

import (
	"service-pay/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载user控制器
func loadUser(s *sinking_web.Engine) {
	balance := s.Group("/balance")
	{
		balance.ANY("/info", context.HandleFunc(controller.User.Balance.Info))
		balance.ANY("/log", context.HandleFunc(controller.User.Balance.Log))
	}
}
