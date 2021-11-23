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

// logSendTask 异步发送
var logSendTask chan *model.SystemLogRequest

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
			_, err = logRequestIns.SendByteMsg(queue3.LogRequestTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
				return
			}
		}
	}(logSendTask)
}

// Send 发送消息
func (l *logRequest) Send() {
	if logRequestIns == nil {
		logRequestIns = sinking.GetRocketProducerServer(queue3.LogRequestGroup, 1)
		logSendTask = make(chan *model.SystemLogRequest)
		l.run()
	}
	logSendTask <- l.log
}
