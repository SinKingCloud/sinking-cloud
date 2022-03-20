package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	r := s.Group("/rpc")
	r.Use(middleware.RpcCheck())
	registerService(r)
}

func registerService(g *sinking_web.RouterGroup) {
	r := g.Group("/sms")
	r.POST("/check", sinking_server.HandleFunc(rpc.Sms.Check))
	r.POST("/send", sinking_server.HandleFunc(rpc.Sms.Send))

	r2 := g.Group("/captcha")
	r2.POST("/check", sinking_server.HandleFunc(rpc.Captcha.Check))
}
