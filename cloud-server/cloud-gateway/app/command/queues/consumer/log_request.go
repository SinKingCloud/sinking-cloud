package consumer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"log"
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
		logRequestIns = sinking.GetRocketConsumerServer(sinking.GetServerHash(queue3.LogRequestGroup), false)
	}
	err := logRequestIns.ListenReceiveMsgDo(queue3.LogRequestTopic, "", func(mqMsg queue.MqMsg) {
		log := &model.SystemLogRequest{}
		err := json.Unmarshal(mqMsg.Body, log)
		if err == nil {
			util.Database.Db.Create(log)
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
}
