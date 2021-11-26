package producer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"log"
)

// webSocketProducerIns 生产者单例
var webSocketProducerIns queue.MqProducer

// logSendTask 异步发送
var webSocketProducerTask chan *message

// NewWebSocketProducer 实例化NewWebSocketProducer
func NewWebSocketProducer() *webSocketProducer {
	return &webSocketProducer{}
}

// webSocketProducer 结构
type webSocketProducer struct {
	msg *message
}

// message 消息结构
type message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// Build 构建消息
func (l *webSocketProducer) Build(id string, msg string) *webSocketProducer {
	l.msg = &message{
		Id:      id,
		Message: msg,
	}
	return l
}

// run 异步执行
func (l *webSocketProducer) run() {
	go func(tasks chan *message) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			tag, exists := util.Cache.Get(util2.CacheNameWebSocketUser + task.Id)
			if exists {
				_, err = webSocketProducerIns.SendByteMsg(queue3.WebSocketSenderTopic, tag, bytes)
				if err != nil {
					util.Log.Println(err)
					log.Println(err)
					return
				}
			}
		}
	}(webSocketProducerTask)
}

// Send 发送消息
func (l *webSocketProducer) Send() {
	if webSocketProducerIns == nil {
		webSocketProducerIns = sinking.GetRocketProducerServer(sinking.GetServerHash(queue3.WebSocketSenderGroup), 1)
		webSocketProducerTask = make(chan *message)
		l.run()
	}
	webSocketProducerTask <- l.msg
}
