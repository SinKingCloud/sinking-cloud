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
		logRequestIns = sinking.GetRocketConsumerServer(queue3.LogRequestGroup, false)
	}
	err := logRequestIns.ListenReceiveMsgDo(queue3.LogRequestTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.SystemLogRequest{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			util.Database.Db.Create(s)
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
}
