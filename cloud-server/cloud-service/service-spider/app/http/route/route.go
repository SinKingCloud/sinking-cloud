package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	index2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/http/controller/index"
	job2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/http/controller/job"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	index := &index2.ControllerIndex{}
	s.GET("/", sinking_server.HandleFunc(index.Index))

	job := &job2.ControllerJob{}
	s.GET("/job/reset", sinking_server.HandleFunc(job.Reset))
}
