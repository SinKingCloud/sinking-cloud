package route

import (
	"service-user/app/http/middleware"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadRpc(s)
	loadUser(s)
	loadAdmin(s)
}

// loadMiddleware 加载中间件
func loadMiddleware(s *sinking_web.Engine) {
	checkWhitePath := func(path string) bool {
		return rpc.IsRpcPath(path) || strings.HasPrefix(path, "/"+constant.ServiceAdmin)
	}
	s.Use(context.HandleFunc(func(c *context.Context) {
		if checkWhitePath(c.Request.URL.Path) {
			c.Next()
			return
		}
		middleware.CheckLogin(c)
	}), context.HandleFunc(func(c *context.Context) {
		if checkWhitePath(c.Request.URL.Path) {
			c.Next()
			return
		}
		middleware.Proxy(c)
	}))
}
