package consumer

import (
	"encoding/json"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/http/websocket"
	sinking_websocket "github.com/SinKingCloud/sinking-go/sinking-websocket"
	"log"
	"sync"
)

// webSocketSenderIns 消费者单例
var (
	webSocketSenderIns     queue.MqConsumer
	onceWebSocketSenderIns sync.Once
)

// NewWebSocketSenderIns 新建单例
func NewWebSocketSenderIns() *messageSender {
	return &messageSender{}
}

// messageSender 结构
type messageSender struct {
}

// Listen 监听并发送消息
func (l *messageSender) Listen() {
	if websocket.Connection == nil {
		websocket.Connection = sinking_websocket.NewWebSocketConnections() //初始化ws连接池
	}
	onceWebSocketSenderIns.Do(func() {
		webSocketSenderIns = sinking.GetRocketConsumerServer(sinking.GetServerHash(queue3.WebSocketSenderGroup), false)
	})
	tag := util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort)
	err := webSocketSenderIns.ListenReceiveMsgDo(queue3.WebSocketSenderTopic, tag, func(mqMsg queue.MqMsg) {
		data := &sinking.WebSocketMessage{}
		err := json.Unmarshal(mqMsg.Body, data)
		if err == nil {
			//此处根据消息结构解析的用户id,并给用户推送消息
			conn := websocket.Connection.Get(data.Id)
			if conn != nil {
				err = conn.WriteMessage(1, []byte(data.Message))
				if err != nil {
					log.Println(err)
				}
			}
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}
