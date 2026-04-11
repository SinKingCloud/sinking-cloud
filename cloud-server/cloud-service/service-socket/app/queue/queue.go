package queue

import "service-socket/app/queue/socket"

var (
	Socket *socket.Queue
)

func Init() {
	Socket = socket.Register()
}
