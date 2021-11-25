package websocket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
	"log"
)

var Connection *sinking_websocket.WebSocketConnections

func Handler(c *sinking_server.Context) {
	//生成uid
	uid := c.Param("id")
	wsServer := sinking_websocket.WebSocket{
		OnError: func(err error) {
			Connection.Delete(uid)
			log.Println("websocket错误", err)
		},
		OnConnect: func(ws *sinking_websocket.Conn) {
			Connection.Set(uid, ws)
			log.Println("websocket连接", uid, c.GetRequestId())
		},
		OnClose: func(err error) {
			Connection.Delete(uid)
			log.Println("websocket关闭", err)
		},
		OnMessage: func(ws *sinking_websocket.Conn, messageType int, data []byte) {
			log.Println("websocket消息", string(data), messageType)
			conn := Connection.Get(uid)
			if conn != nil {
				_ = conn.WriteMessage(1, data)
			}
		},
	}
	wsServer.Listen(c.Writer, c.Request)
}
