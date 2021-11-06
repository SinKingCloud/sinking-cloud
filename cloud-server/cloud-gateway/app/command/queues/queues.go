package queues

import (
	"cloud-gateway/app/command/queues/consumer"
)

// Run 执行消息队列
func Run() {
	consumer.NewLogRequest().Listen()
}
