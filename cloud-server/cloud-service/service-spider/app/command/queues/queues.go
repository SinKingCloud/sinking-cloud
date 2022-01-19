package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/queues/consumer"

// Run 执行消息队列
func Run() {
	consumer.NewJobLog().Listen() //启动数据写入队列
}
