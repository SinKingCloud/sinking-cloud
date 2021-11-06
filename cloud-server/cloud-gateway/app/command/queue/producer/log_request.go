package producer

import (
	queue2 "cloud-gateway/app/command/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// logRequestIns 生产者单例
var logRequestIns queue.MqProducer

// NewLogRequest 实例化logRequest
func NewLogRequest() *logRequest {
	return &logRequest{}
}

// logRequest 结构
type logRequest struct {
}

func (l *logRequest) Send() {
	if logRequestIns == nil {
		logRequestIns = GetRocketProducerServer(queue2.LogRequestGroup, 1)
	}
	_, err := logRequestIns.SendMsg(queue2.LogRequestTopic, "*", "")
	if err != nil {
		util.Log.Println(err)
		return
	}
}
