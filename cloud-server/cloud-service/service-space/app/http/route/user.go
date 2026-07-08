package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载user控制器
func loadUser(s *sinking_web.Engine) {
	_ = service.Auth.GetUserRouteGroup(s)
	{

	}
}
