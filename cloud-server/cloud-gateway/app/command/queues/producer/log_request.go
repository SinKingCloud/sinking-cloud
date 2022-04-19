package producer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
)

// logRequestProducerIns 生产者单例
var logRequestProducerIns queue.MqProducer

// logSendProducerTask 异步发送
var logSendProducerTask chan *model.SystemLogRequest

// NewLogRequest 实例化logRequest
func NewLogRequest() *logRequest {
	return &logRequest{}
}

// logRequest 结构
type logRequest struct {
	log *model.SystemLogRequest
}

// Build 构建消息
func (l *logRequest) Build(request *model.SystemLogRequest) *logRequest {
	l.log = request
	return l
}

// run 异步执行
func (l *logRequest) run() {
	go func(tasks chan *model.SystemLogRequest) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = logRequestProducerIns.SendByteMsg(queue3.LogRequestTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
			}
		}
	}(logSendProducerTask)
}

// Send 发送消息
func (l *logRequest) Send() {
	if logRequestProducerIns == nil {
		logRequestProducerIns = sinking.GetRocketProducerServer(queue3.LogRequestGroup, 1)
		logSendProducerTask = make(chan *model.SystemLogRequest)
		l.run()
	}
	logSendProducerTask <- l.log
}
