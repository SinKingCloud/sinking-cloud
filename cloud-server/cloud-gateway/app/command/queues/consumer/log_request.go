package consumer

import (
	queue2 "cloud-gateway/app/util/queue"
	"fmt"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// logRequestIns 消费者单例
var logRequestIns queue.MqConsumer

func NewLogRequest() *logRequest {
	return &logRequest{}
}

// logRequest 结构
type logRequest struct {
}

func (l *logRequest) Listen() {
	if logRequestIns == nil {
		logRequestIns = queue2.GetRocketConsumerServer(queue3.LogRequestGroup, 1)
	}
	err := logRequestIns.ListenReceiveMsgDo(queue3.LogRequestTopic, "", func(mqMsg queue.MqMsg) {
		fmt.Println("收到消息", mqMsg.MsgId, mqMsg.Body)
	})
	if err != nil {
		return
	}
}
