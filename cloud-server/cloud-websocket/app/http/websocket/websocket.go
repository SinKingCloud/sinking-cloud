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
	token := c.Request.Header.Get("Sec-Websocket-Protocol")
	if token == "" || len(token) != 32 {
		c.NotLogin("token验证失败", nil)
		return
	}
	_, exists := util.Cache.Get(util2.CacheNameWebSocketAuth + token) //获取token信息
	if !exists {
		c.NotLogin("token鉴权失败", nil)
		return
	}
	//ws推送业务
	wsServer := sinking_websocket.WebSocket{
		Id: token,
		OnError: func(id string, err error) {
			Connection.Delete(id)
			util.Cache.Delete(util2.CacheNameWebSocketAuth + id)
		},
		OnConnect: func(id string, ws *sinking_websocket.Conn) {
			Connection.Set(id, ws)
			//更新用户ws地址
			addr := util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort)
			if !util.Cache.SetWithExpire(util2.CacheNameWebSocketAuth+id, addr, util2.CacheTimeWebSocketAuth) {
				Connection.Delete(id)
			}
		},
		OnClose: func(id string, err error) {
			Connection.Delete(id)
			util.Cache.Delete(util2.CacheNameWebSocketAuth + id)
		},
		OnMessage: func(id string, ws *sinking_websocket.Conn, messageType int, data []byte) {
			if string(data) == "ping" {
				dateTime := time.Now().Format("2006-01-02 15:04:05")
				_ = ws.Conn.WriteMessage(1, []byte(dateTime))
			}
		},
	}
	response := c.Writer.Header().Clone()
	response.Set("Sec-Websocket-Protocol", token)
	wsServer.Listen(c.Writer, c.Request, response)
}
