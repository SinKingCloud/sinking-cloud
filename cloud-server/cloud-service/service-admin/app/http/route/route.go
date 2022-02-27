package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/http/controller/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/http/controller/person"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadPageRoute(s)
	loadPersonRoute(s)
}

// loadPageRoute 加载网站授权路由
func loadPageRoute(s *sinking_web.Engine) {
	r := s.Group("/page")
	handle := &page.ControllerPage{}
	r.GET("/menu", sinking_server.HandleFunc(handle.Menu))       //获取菜单列表
	r.GET("/control", sinking_server.HandleFunc(handle.Control)) //获取页面控件
}

// loadPersonRoute 加载网站授权路由
func loadPersonRoute(s *sinking_web.Engine) {
	r := s.Group("/person")
	handle := &person.ControllerPerson{}
	r.GET("/info", sinking_server.HandleFunc(handle.Info))
}
