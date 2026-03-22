package route

import (
	"service-gateway/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadApi(s)
}

// loadMiddleware 中间件
func loadMiddleware(s *sinking_web.Engine) {
	s.Use(context.HandleFunc(middleware.Cors))
}

// loadApi api路由
func loadApi(s *sinking_web.Engine) {
	s.ANY("/info", context.HandleFunc(controller.Info))
	s.ANY("/enum", context.HandleFunc(controller.Enum))
	s.ANY("/*", context.HandleFunc(controller.Gateway))
	s.ANY("/", context.HandleFunc(controller.Gateway))
}
