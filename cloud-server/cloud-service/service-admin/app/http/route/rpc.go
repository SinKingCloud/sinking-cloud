package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/http/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	r := s.Group("/rpc")
	r.Use(middleware.RpcCheck())
	registerService(r)
}

func registerService(g *sinking_web.RouterGroup) {
	r := g.Group("/rbac")
	r.POST("/check", sinking_server.HandleFunc(rpc.Rbac.Check))

	r2 := g.Group("/user")
	r2.POST("/findById", sinking_server.HandleFunc(rpc.User.FindById))

	r3 := g.Group("/log")
	r3.POST("/create", sinking_server.HandleFunc(rpc.Log.Create))
}
