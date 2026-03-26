package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadAdmin(s *sinking_web.Engine) {
	g := global.App.GetAdminRouteGroup(s)
	{
		g.ANY("/test", context.HandleFunc(func(c *context.Context) {
			c.Success("test success")
		}))
	}
}
