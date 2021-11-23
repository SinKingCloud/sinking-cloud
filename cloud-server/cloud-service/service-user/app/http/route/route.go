package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	index2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/controller/index"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	index := &index2.ControllerIndex{}
	s.GET("/", sinking_server.HandleFunc(index.Index))
}
