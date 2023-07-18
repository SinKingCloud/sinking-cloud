package sinking

import (
	"encoding/json"
	"errors"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// webSocketProducerIns 生产者单例
var webSocketProducerIns queue.MqProducer

// logSendTask 异步发送
var webSocketProducerTask chan *WebSocketMessage

// NewWebSocketProducer 实例化NewWebSocketProducer
func NewWebSocketProducer() *webSocketProducer {
	return &webSocketProducer{}
}

// webSocketProducer 结构
type webSocketProducer struct {
	msg *WebSocketMessage
}

// WebSocketMessage 消息结构
type WebSocketMessage struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// Build 构建消息
func (l *webSocketProducer) Build(id string, msg string) *webSocketProducer {
	l.msg = &WebSocketMessage{
		Id:      id,
		Message: msg,
	}
	return l
}

// run 异步执行
func (l *webSocketProducer) run() {
	go func(tasks chan *WebSocketMessage) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			tag, exists := util.Cache.Get(util2.CacheNameWebSocketToken + task.Id)
			if exists {
				_, err = webSocketProducerIns.SendByteMsg(queue3.WebSocketSenderTopic, tag, bytes)
				if err != nil {
					util.Log.Println(err)
				}
			}
		}
	}(webSocketProducerTask)
}

// Send 发送消息
func (l *webSocketProducer) Send() {
	if webSocketProducerIns == nil {
		webSocketProducerIns = GetRocketProducerServer(GetServerHash(queue3.WebSocketSenderGroup), 1)
		webSocketProducerTask = make(chan *WebSocketMessage)
		l.run()
	}
	webSocketProducerTask <- l.msg
}

// Auth 授权token
func (l *webSocketProducer) Auth(token string) error {
	//颁发一个token
	if util.Cache.SetWithExpire(util2.CacheNameWebSocketAuth+token, token, util2.CacheTimeWebSocketAuth) {
		return nil
	}
	return errors.New("授权token失败")
}
