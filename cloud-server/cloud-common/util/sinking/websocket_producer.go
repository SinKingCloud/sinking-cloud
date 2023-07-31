package sinking

import (
	"encoding/json"
	"errors"
	queue3 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/sys_user_config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/web_user_config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"strconv"
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

// WebSocketMessageBody 消息内容结构
type WebSocketMessageBody struct {
	Id   int64       `json:"id"`
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

// Build 构建消息 id:user_id msg:消息
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
			tag, exists := util.Cache.Get(util2.CacheNameWebSocketAuth + task.Id)
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

// Auth 授权token id:user_id
func (l *webSocketProducer) Auth(id string) error {
	//颁发一个token
	if util.Cache.SetWithExpire(util2.CacheNameWebSocketAuth+id, id, util2.CacheTimeWebSocketAuth) {
		return nil
	}
	return errors.New("授权token失败")
}

// Push 根据32位token推送 id:user_id msg:消息内容
func (l *webSocketProducer) Push(id string, msg string) {
	l.msg = &WebSocketMessage{
		Id:      id,
		Message: msg,
	}
	l.Send()
}

// GenToken 生成token
func (l *webSocketProducer) GenToken() (string, error) {
	id := str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5()
	if l.Auth(id) == nil {
		return id, nil
	}
	return "", errors.New("授权token失败")
}

// PushByUserId 根据用户ID推送 id:user_id msg:消息内容
func (l *webSocketProducer) PushByUserId(userId int, msg string, isSys bool) error {
	key := "web_"
	if isSys {
		key = "sys_"
	}
	//1.查询用户绑定得token
	token, exists := util.Cache.Get(util2.CacheNameWebSocketUser + key + strconv.Itoa(userId))
	if !exists {
		var d string
		var ok bool
		if isSys {
			configs := service.GetSysUserConfSrvIns().GetConfigs(userId, sys_user_config.WebSocket)
			d, ok = configs[sys_user_config.WebSocketToken]
		} else {
			configs := service.GetWebUserConfSrvIns().GetConfigs(userId, web_user_config.WebSocket)
			d, ok = configs[web_user_config.WebSocketToken]
		}
		if ok {
			if util.Cache.SetWithExpire(util2.CacheNameWebSocketUser+key+strconv.Itoa(userId), d, util2.CacheTimeWebSocketUser) {
				exists = true
				token = d
			}
		}
	}
	if exists && token != "" {
		l.Push(token, msg)
		return nil
	}
	return errors.New("获取绑定关系失败")
}

// GenTokenByUserId 根据用户ID生成token id:user_id msg:消息内容
func (l *webSocketProducer) GenTokenByUserId(userId int, isSys bool) (string, error) {
	key := "web_"
	if isSys {
		key = "sys_"
	}
	//1.生成用户绑定token
	token, err := l.GenToken()
	if err != nil || token == "" {
		return "", errors.New("获取授权token失败")
	}
	//2.写入用户绑定关系
	r := false
	if isSys {
		r = service.GetSysUserConfSrvIns().SetConfig(userId, sys_user_config.WebSocketToken, token)
	} else {
		r = service.GetWebUserConfSrvIns().SetConfig(userId, web_user_config.WebSocketToken, token)
	}
	if r {
		if util.Cache.SetWithExpire(util2.CacheNameWebSocketUser+key+strconv.Itoa(userId), token, util2.CacheTimeWebSocketUser) {
			return token, nil
		}
	}
	return "", errors.New("写入授权token失败")
}
