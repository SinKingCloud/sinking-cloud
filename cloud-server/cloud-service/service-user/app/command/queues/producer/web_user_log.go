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
var logProducerIns queue.MqProducer

// logSendProducerTask 异步发送
var logSendProducerTask chan *model.WebUserLog

// NewWebUserLogProducer 实例化log
func NewWebUserLogProducer() *log {
	return &log{}
}

// log 结构
type log struct {
	log *model.WebUserLog
}

// Build 构建消息
func (l *log) Build(request *model.WebUserLog) *log {
	l.log = request
	return l
}

// run 异步执行
func (l *log) run() {
	go func(tasks chan *model.WebUserLog) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = logProducerIns.SendByteMsg(queue3.WebLogTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
				return
			}
		}
	}(logSendProducerTask)
}

// Send 发送消息
func (l *log) Send() {
	if logProducerIns == nil {
		logProducerIns = sinking.GetRocketProducerServer(queue3.WebLogGroup, 1)
		logSendProducerTask = make(chan *model.WebUserLog)
		l.run()
	}
	logSendProducerTask <- l.log
}
