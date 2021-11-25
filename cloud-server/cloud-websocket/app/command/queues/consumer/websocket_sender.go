package consumer

import (
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"log"
)

// webSocketSenderIns 消费者单例
var webSocketSenderIns queue.MqConsumer

// NewWebSocketSenderIns 新建单例
func NewWebSocketSenderIns() *messageSender {
	return &messageSender{}
}

// messageSender 结构
type messageSender struct {
}

// Listen 监听并发送消息
func (l *messageSender) Listen() {
	if webSocketSenderIns == nil {
		webSocketSenderIns = sinking.GetRocketConsumerServer(queue3.WebSocketSenderGroup, false)
	}
	tag := util.Conf.GetString(config.ConfApplicationAddr) + util.Conf.GetString(config.ConfApplicationPort)
	err := webSocketSenderIns.ListenReceiveMsgDo(queue3.WebSocketSenderTopic, tag, func(mqMsg queue.MqMsg) {
		//log := &model.SystemLogRequest{}
		//err := json.Unmarshal(mqMsg.Body, log)
		//if err == nil {
		//	util.Database.Db.Create(log)
		//}
	})
	if err != nil {
		log.Println(err)
		return
	}
}
