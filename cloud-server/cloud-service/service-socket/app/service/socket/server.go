package socket

import (
	"net/http"

	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Server ws服务器
func (s *service) Server(token string) (*sinking_websocket.Server, error) {
	server := sinking_websocket.NewServer(
		sinking_websocket.WithConnectionIDResolver(func(*http.Request) (string, error) {
			return token, nil
		}),
		sinking_websocket.WithConnectHandler(func(connection *sinking_websocket.Connection) error {

			return nil
		}),
		sinking_websocket.WithDisconnectHandler(func(connection *sinking_websocket.Connection, err error) {

		}),
	)
	return server, nil
}
