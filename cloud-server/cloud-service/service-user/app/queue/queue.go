package queue

import (
	queueMessage "service-user/app/queue/message"

	queueUtil "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Message *queueUtil.UnicastClient[*queueMessage.Task]
)

func Init() {
	Message = queueMessage.Register()
}
