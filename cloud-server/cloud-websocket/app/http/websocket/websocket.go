package websocket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
	"strconv"
	"time"
)

var Connection *sinking_websocket.WebSocketConnections

// Handler websocket处理
func Handler(c *sinking_server.Context) {
	//获取token
	token := c.Request.Header.Get("Sec-WebSocket-Protocol")
	if token == "" {
		c.NotLogin("您还未登陆,请先登陆账户", nil)
		return
	}
	key := jwt.CheckToken(token)
	if key == nil || key.WebUser == nil || time.Now().Unix() > key.ExpiresAt {
		c.TokenError("登陆超时,请重新登陆", nil)
		return
	}
	//获取uid
	id := strconv.Itoa(key.WebUser.Id)
	//ws推送业务
	wsServer := sinking_websocket.WebSocket{
		OnError: func(err error) {
			Connection.Delete(id)
		},
		OnConnect: func(ws *sinking_websocket.Conn) {
			Connection.Set(id, ws)
			//更新用户ws地址
			addr := util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort)
			util.Cache.SetWithExpire(util2.CacheNameWebSocketUser+id, addr, util2.CacheTimeWebSocketUser)
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
