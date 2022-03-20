package consumer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
)

// logIns 消费者单例
var logIns queue.MqConsumer

func NewWebUserLog() *log {
	return &log{}
}

// log 结构
type log struct {
}

func (l *log) Listen() {
	if logIns == nil {
		logIns = sinking.GetRocketConsumerServer(queue3.WebLogGroup, false)
	}
	err := logIns.ListenReceiveMsgDo(queue3.WebLogTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.WebUserLog{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			util.Database.Db.Create(s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}
