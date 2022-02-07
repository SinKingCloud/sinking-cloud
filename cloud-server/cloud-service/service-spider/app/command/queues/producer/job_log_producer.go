package producer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/model"
	"log"
)

// jobLogProducerIns 生产者单例
var jobLogProducerIns queue.MqProducer

// jobLogSendProducerTask 异步发送
var jobLogSendProducerTask chan *model.JobLog

// NewJobLogProducer 实例化logRequest
func NewJobLogProducer() *logLog {
	return &logLog{}
}

// logLog 结构
type logLog struct {
	log *model.JobLog
}

// Build 构建消息
func (l *logLog) Build(request *model.JobLog) *logLog {
	l.log = request
	return l
}

// run 异步执行
func (l *logLog) run() {
	go func(tasks chan *model.JobLog) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = jobLogProducerIns.SendByteMsg(constant.JobLogTopic, "", bytes)
			if err != nil {
				log.Println(err)
				util.Log.Println(err)
				return
			}
		}
	}(jobLogSendProducerTask)
}

// Send 发送消息
func (l *logLog) Send() {
	if jobLogProducerIns == nil {
		jobLogProducerIns = sinking.GetRocketProducerServer(constant.JobLogGroup, 1)
		jobLogSendProducerTask = make(chan *model.JobLog)
		l.run()
	}
	jobLogSendProducerTask <- l.log
}
