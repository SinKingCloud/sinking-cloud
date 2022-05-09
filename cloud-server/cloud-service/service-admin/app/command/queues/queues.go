package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/command/queues/consumer"

// Run 执行消息队列
func Run() {
	consumer.NewLog().Listen()
}
