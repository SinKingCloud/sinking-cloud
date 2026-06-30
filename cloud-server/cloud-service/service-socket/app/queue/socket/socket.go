package socket

import (
	"encoding/json"
	"errors"
	"runtime"
	"strings"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Message 消息类型
type Message struct {
	Group  string
	Id     string
	IsText bool
	Data   []byte
}

// Queue 消息队列
type Queue struct {
	client        queue.Queue
	name          string
	once          sync.Once
	handler       *Handler
	searchUnicast *queue.UnicastClient[*Message]
	sendUnicast   *queue.UnicastClient[*Message]
	sendBroadcast *queue.BroadcastClient[*Message]
}

type Handler struct {
	GetConnectionIdAddr func(group string, id string) (string, error)
	SendToLocalId       func(group string, id string, messageType int, payload []byte) error
	SendToLocalAll      func(group string, messageType int, payload []byte) (queued int, dropped int, closed int, err error)
}

// Name 消息队列名称
func (s *Queue) Name(isAll bool, name ...string) string {
	str := strings.Builder{}
	str.WriteString(constant.ServiceSocket)
	if isAll {
		str.WriteString("-all")
	} else {
		str.WriteString("-one")
	}
	if len(name) > 0 {
		str.WriteString("-")
		str.WriteString(strings.ReplaceAll(strings.ReplaceAll(name[0], ":", "_"), ".", "_"))
	} else if !isAll {
		str.WriteString("-")
		str.WriteString(strings.ReplaceAll(strings.ReplaceAll(s.name, ":", "_"), ".", "_"))
	}
	return str.String()
}

// SetHandler 设置socket队列业务处理器。
func (s *Queue) SetHandler(handler *Handler) {
	s.handler = handler
	if s.searchUnicast != nil {
		s.searchUnicast.SetConsumer(s.consumeSearchMessage)
	}
	if s.sendUnicast != nil {
		s.sendUnicast.SetConsumer(s.consumeLocalIdMessage)
	}
	if s.sendBroadcast != nil {
		s.sendBroadcast.SetConsumer(s.consumeLocalAllMessage)
	}
}

// Listen 监听消息队列
func (s *Queue) Listen() {
	s.once.Do(func() {
		//查找指定地址在投递
		searchUnicast, err := queue.RegisterUnicast[*Message](s.client, constant.ServiceSocket+"-send", runtime.NumCPU(), runtime.NumCPU()*2, nil)
		if err != nil {
			panic(err)
			return
		}
		s.searchUnicast = searchUnicast
		//单播消息
		sendUnicast, err := queue.RegisterUnicast[*Message](s.client, s.Name(false), runtime.NumCPU(), runtime.NumCPU()*2, nil)
		if err != nil {
			panic(err)
			return
		}
		s.sendUnicast = sendUnicast
		//广播消息
		sendBroadcast, err := queue.RegisterBroadcast[*Message](s.client, s.Name(true), runtime.NumCPU(), runtime.NumCPU()*2, nil)
		if err != nil {
			panic(err)
			return
		}
		s.sendBroadcast = sendBroadcast
	})
}

// consumeSearchMessage 消费全局单播消息，定位连接所在节点后转发。
func (s *Queue) consumeSearchMessage(message *Message) error {
	if message == nil {
		return errors.New("获取消息失败")
	}
	if s.handler == nil || s.handler.GetConnectionIdAddr == nil {
		return errors.New("获取连接地址处理器未初始化")
	}
	addr, err := s.handler.GetConnectionIdAddr(message.Group, message.Id)
	if err != nil {
		return err
	}
	name := s.Name(false, addr)
	marshal, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = s.client.SendUnicast(name, string(marshal))
	return err
}

// consumeLocalIdMessage 消费本机单播消息。
func (s *Queue) consumeLocalIdMessage(message *Message) error {
	if message == nil {
		return errors.New("获取消息失败")
	}
	messageType := sinking_websocket.BinaryMessage
	if message.IsText {
		messageType = sinking_websocket.TextMessage
	}
	if s.handler == nil || s.handler.SendToLocalId == nil {
		return errors.New("单播处理器未初始化")
	}
	return s.handler.SendToLocalId(message.Group, message.Id, messageType, message.Data)
}

// consumeLocalAllMessage 消费本机广播消息。
func (s *Queue) consumeLocalAllMessage(message *Message) error {
	if message == nil {
		return errors.New("获取消息失败")
	}
	messageType := sinking_websocket.BinaryMessage
	if message.IsText {
		messageType = sinking_websocket.TextMessage
	}
	if s.handler == nil || s.handler.SendToLocalAll == nil {
		return errors.New("广播处理器未初始化")
	}
	_, _, _, err := s.handler.SendToLocalAll(message.Group, messageType, message.Data)
	return err
}

// Register 注册消息队列
func Register() *Queue {
	client, err := queue.RegisterNatsQueueClient(global.App.GetNatsUrl())
	if err != nil {
		panic(err)
		return nil
	}
	name := global.App.Consul.GetAddress()
	ins := &Queue{client: client, name: name}
	ins.Listen()
	global.App.Log.Println(ins.Name(false) + " queue初始化成功")
	global.App.Log.Println(ins.Name(true) + " queue初始化成功")
	return ins
}
