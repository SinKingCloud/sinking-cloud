package producer

import (
	queue2 "cloud-gateway/app/command/queue"
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
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
	bytes []byte
}

// Build 构建消息
func (l *logRequest) Build(request *model.SystemLogRequest) *logRequest {
	bytes, err := json.Marshal(request)
	if err == nil {
		l.bytes = bytes
	}
	return l
}

// Send 发送消息
func (l *logRequest) Send() {
	if logRequestIns == nil {
		logRequestIns = GetRocketProducerServer(queue2.LogRequestGroup, 1)
	}
	_, err := logRequestIns.SendByteMsg(queue2.LogRequestTopic, "", l.bytes)
	if err != nil {
		util.Log.Println(err)
		return
	}
}
