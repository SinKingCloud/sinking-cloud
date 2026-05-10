package route

import (
	constant2 "service-user/app/constant"
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
	loadApi(s)
	loadUser(s)
	loadAdmin(s)
}

// loadMiddleware 加载中间件
func loadMiddleware(s *sinking_web.Engine) {
	checkApiPath := func(path string) bool {
		return strings.HasPrefix(path, "/api/")
	}
	checkWhitePath := func(path string) bool {
		return rpc.IsRpcPath(path) || strings.HasPrefix(path, "/"+constant.ServiceAdmin) || checkApiPath(path)
	}
	s.Use(context.HandleFunc(func(c *context.Context) {
		if !checkWhitePath(c.Request.URL.Path) && !c.IsAborted() {
			middleware.CheckLogin(c)
		}
	}), context.HandleFunc(func(c *context.Context) {
		c.Next()
		if !c.GetBool(constant2.NeedProxy) {
			return
		}
		if checkApiPath(c.Request.URL.Path) && !c.IsAborted() {
			middleware.CheckLogin(c)
		}
		if !c.IsAborted() {
			middleware.Proxy(c)
		}
	}))
}
