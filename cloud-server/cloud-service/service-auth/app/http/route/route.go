package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/system"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/web"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadWebRoute(s)
	loadSystemRoute(s)
}

// loadWebRoute 加载网站授权路由
func loadWebRoute(s *sinking_web.Engine) {
	r := s.Group("/web")
	handle := &web.ControllerWeb{}
	r.GET("/login", sinking_server.HandleFunc(handle.Login))
}

// loadWebRoute 加载后台授权路由
func loadSystemRoute(s *sinking_web.Engine) {
	r := s.Group("/system")
	handle := &system.ControllerSystem{}
	r.GET("/login", sinking_server.HandleFunc(handle.Login))
}
