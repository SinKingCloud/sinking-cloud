package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/http/controller/user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/http/websocket"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	sinking_websocket "github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	users := &user.ControllerUser{}
	websocket.Connection = sinking_websocket.NewWebSocketConnections() //初始化ws连接池
	s.GET("/test/:id", sinking_server.HandleFunc(users.Index))
}
