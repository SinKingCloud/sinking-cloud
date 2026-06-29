package controller

import (
	"net/http"
	"service-socket/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

func Socket(c *context.Context) {
	token := c.Request.Header.Get("Sec-Websocket-Protocol")
	if token == "" {
		c.Error("token验证失败")
		return
	}
	server, err := service.Socket.Server(token)
	if err != nil {
		c.Error(err.Error())
		return
	}
	resp := make(http.Header)
	resp.Set("Sec-Websocket-Protocol", token)
	_ = server.Handle(c.Writer, c.Request, resp)
}
