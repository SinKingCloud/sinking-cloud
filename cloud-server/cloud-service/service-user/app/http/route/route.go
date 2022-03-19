package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/controller/person"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadPersonRoute(s)
}

// loadPersonRoute 加载个人操作路由
func loadPersonRoute(s *sinking_web.Engine) {
	r := s.Group("/person")
	handle := &person.ControllerPerson{}
	r.GET("/info", sinking_server.HandleFunc(handle.Info))
	r.POST("/update/info", sinking_server.HandleFunc(handle.Update))
	r.POST("/update/password", sinking_server.HandleFunc(handle.Password))
	r.POST("/update/phone", sinking_server.HandleFunc(handle.Phone))
}
