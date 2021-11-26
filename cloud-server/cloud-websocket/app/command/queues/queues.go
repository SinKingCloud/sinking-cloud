package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/command/queues/consumer"

// Run 执行消息队列
func Run() {
	//监听websocket消息
	consumer.NewWebSocketSenderIns().Listen()
}
