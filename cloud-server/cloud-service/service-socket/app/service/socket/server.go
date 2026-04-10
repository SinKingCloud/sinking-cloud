package socket

import (
	"net/http"

	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Server ws服务器
func (s *service) Server(token string) (*sinking_websocket.Server, error) {
	registry := sinking_websocket.NewRegistry()
	server := sinking_websocket.NewServer(
		sinking_websocket.WithConnectionIDResolver(func(request *http.Request) (string, error) {
			return request.URL.Query().Get("id"), nil
		}),
		sinking_websocket.WithConnectHandler(func(connection *sinking_websocket.Connection) error {
			registry.Store(connection.ID(), connection)
			return connection.Send(sinking_websocket.TextMessage, []byte("success"))
		}),
		sinking_websocket.WithMessageHandler(func(connection *sinking_websocket.Connection, message sinking_websocket.Message) error {
			return connection.Send(message.Type, message.Payload)
		}),
		sinking_websocket.WithDisconnectHandler(func(connection *sinking_websocket.Connection, err error) {
			registry.DeleteIfMatch(connection.ID(), connection)
		}),
	)
	return server, nil
}
