package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/command/queues/consumer"

// Run 执行消息队列
func Run() {
	consumer.NewWebUserLog().Listen()
}
