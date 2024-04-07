package queues

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/command/queues/consumer"

func Run() {
	consumer.NewUinJob().Listen()
}
