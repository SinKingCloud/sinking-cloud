package websocket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
	"time"
)

var Connection *sinking_websocket.WebSocketConnections

// Handler websocket处理
func Handler(c *sinking_server.Context) {
	//获取token
	token := c.Request.Header.Get("Sec-WebSocket-Protocol")
	if token == "" || len(token) != 32 {
		c.NotLogin("token验证失败", nil)
		return
	}
	id, exists := util.Cache.Get(util2.CacheNameWebSocketAuth + token) //获取token信息
	if !exists {
		c.NotLogin("token鉴权失败", nil)
		return
	}
	//ws推送业务
	wsServer := sinking_websocket.WebSocket{
		OnError: func(err error) {
			Connection.Delete(id)
			util.Cache.Delete(util2.CacheNameWebSocketToken + id)
		},
		OnConnect: func(ws *sinking_websocket.Conn) {
			Connection.Set(id, ws)
			//更新用户ws地址
			addr := util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort)
			if !util.Cache.SetWithExpire(util2.CacheNameWebSocketToken+id, addr, util2.CacheTimeWebSocketToken) {
				Connection.Delete(id)
			}
		},
		OnClose: func(err error) {
			Connection.Delete(id)
			util.Cache.Delete(util2.CacheNameWebSocketToken + id)
		},
		OnMessage: func(ws *sinking_websocket.Conn, messageType int, data []byte) {
			//只返回时间信息,用作心跳判断
			dateTime := time.Now().Format("2006-01-02 15:04:05")
			_ = ws.Conn.WriteMessage(1, []byte(dateTime))
		},
	}
	wsServer.Listen(c.Writer, c.Request)
}
