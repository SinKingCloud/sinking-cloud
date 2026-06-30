package queue

import socketQueue "service-socket/app/queue/socket"

var (
	Socket *socketQueue.Queue
)

func Init() {
	Socket = socketQueue.Register()
}
