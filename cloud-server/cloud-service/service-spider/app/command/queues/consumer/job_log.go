package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/model"
	"log"
)

// jobLogIns 消费者单例
var jobLogIns queue.MqConsumer

func NewJobLog() *logRequest {
	return &logRequest{}
}

// logRequest 结构
type logRequest struct {
}

func (l *logRequest) Listen() {
	if jobLogIns == nil {
		jobLogIns = sinking.GetRocketConsumerServer(constant.JobLogGroup, false)
	}
	err := jobLogIns.ListenReceiveMsgDo(constant.JobLogTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.JobLog{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			var temp model.JobLog
			util.Database.Db.Where("company_name=?", s.CompanyName).First(&temp)
			if temp.Id > 0 {
				util.Database.Db.Where("id=?", temp.Id).Updates(&s)
			} else {
				util.Database.Db.Create(&s)
			}
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
}
