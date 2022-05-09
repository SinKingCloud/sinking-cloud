package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/consumer"

func Run() {
	consumer.NewUinJob().Listen()
	consumer.NewProxyJob().Listen()
}
