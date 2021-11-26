package websocket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
	"log"
	"time"
)

var Connection *sinking_websocket.WebSocketConnections

func Handler(c *sinking_server.Context) {
	//生成uid
	id := c.Param("id")
	//处理认证业务

	//ws推送业务
	wsServer := sinking_websocket.WebSocket{
		OnError: func(err error) {
			Connection.Delete(id)
		},
		OnConnect: func(ws *sinking_websocket.Conn) {
			Connection.Set(id, ws)
			//rpc调用用户服务。更新用户ws地址
			topic := util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort)
			log.Println(str.NetStringTool(topic).Md5())
		},
		OnClose: func(err error) {
			Connection.Delete(id)
		},
		OnMessage: func(ws *sinking_websocket.Conn, messageType int, data []byte) {
			//只返回时间信息,用作心跳判断
			dateTime := time.Now().Format("2006-01-02 15:04:05")
			_ = ws.Conn.WriteMessage(1, []byte(dateTime))
		},
	}
	wsServer.Listen(c.Writer, c.Request)
}
