package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	r := s.Group("/rpc")
	r.Use(middleware.RpcCheck())
	registerService(r)
}

func registerService(g *sinking_web.RouterGroup) {
	r := g.Group("/oauth")
	r.POST("/auth_ticket", sinking_server.HandleFunc(rpc.Oauth.AuthTicket))
}
