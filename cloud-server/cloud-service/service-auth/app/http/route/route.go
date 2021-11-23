package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/login"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadWebRoute(s)
}

// loadWebRoute 加载网站授权路由
func loadWebRoute(s *sinking_web.Engine) {
	r := s.Group("/web")
	handle := &login.ControllerLogin{}
	r.GET("/login", sinking_server.HandleFunc(handle.Index))
}
