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

// Listen 监听消息队列
func (s *Queue) Listen(handler *Handler) {
	s.once.Do(func() {
		//查找指定地址在投递
		searchUnicast, err := queue.RegisterUnicast[*Message](s.client, constant.ServiceSocket+"-send", runtime.NumCPU(), runtime.NumCPU()*2, func(message *Message) error {
			if message == nil {
				return errors.New("获取消息失败")
			}
			if handler == nil || handler.GetConnectionIdAddr == nil {
				return errors.New("获取连接地址处理器未初始化")
			}
			addr, err := handler.GetConnectionIdAddr(message.Group, message.Id)
			if err != nil {
				return err
			}
			name := s.Name(false, addr)
			marshal, err := json.Marshal(message)
			if err != nil {
				return err
			}
			_, err2 := s.client.SendUnicast(name, string(marshal))
			return err2
		})
		if err != nil {
			panic(err)
			return
		}
		s.searchUnicast = searchUnicast
		//单播消息
		sendUnicast, err := queue.RegisterUnicast[*Message](s.client, s.Name(false), runtime.NumCPU(), runtime.NumCPU()*2, func(message *Message) error {
			messageType := sinking_websocket.BinaryMessage
			if message.IsText {
				messageType = sinking_websocket.TextMessage
			}
			if handler == nil || handler.SendToLocalId == nil {
				return errors.New("单播处理器未初始化")
			}
			return handler.SendToLocalId(message.Group, message.Id, messageType, message.Data)
		})
		if err != nil {
			panic(err)
			return
		}
		s.sendUnicast = sendUnicast
		//广播消息
		sendBroadcast, err := queue.RegisterBroadcast[*Message](s.client, s.Name(true), runtime.NumCPU(), runtime.NumCPU()*2, func(message *Message) error {
			messageType := sinking_websocket.BinaryMessage
			if message.IsText {
				messageType = sinking_websocket.TextMessage
			}
			if handler == nil || handler.SendToLocalAll == nil {
				return errors.New("广播处理器未初始化")
			}
			_, _, _, err2 := handler.SendToLocalAll(message.Group, messageType, message.Data)
			return err2
		})
		if err != nil {
			panic(err)
			return
		}
		s.sendBroadcast = sendBroadcast
	})
}

// Register 注册消息队列
func Register(name string, handler *Handler) *Queue {
	client, err := queue.RegisterNatsQueueClient(global.App.GetNatsUrl())
	if err != nil {
		panic(err)
		return nil
	}
	ins := &Queue{client: client, name: name}
	ins.Listen(handler)
	global.App.Log.Println(ins.Name(false) + " queue初始化成功")
	global.App.Log.Println(ins.Name(true) + " queue初始化成功")
	return ins
}
