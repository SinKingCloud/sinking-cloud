package producer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
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
		logRequestIns = sinking.GetRocketProducerServer(queue3.LogRequestGroup, 1)
	}
	_, err := logRequestIns.SendByteMsg(queue3.LogRequestTopic, "", l.bytes)
	if err != nil {
		util.Log.Println(err)
		return
	}
}
