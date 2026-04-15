package socket

import (
	"encoding/json"
	"errors"
	"runtime"
	"service-socket/app/service"
	"strings"
	"sync"

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
	once          sync.Once
	searchUnicast *queue.UnicastClient[*Message]
	sendUnicast   *queue.UnicastClient[*Message]
	sendBroadcast *queue.BroadcastClient[*Message]
}

// Name 消息队列名称
func (s *Queue) Name(isAll bool, name ...string) string {
	str := strings.Builder{}
	str.WriteString("socket-")
	if isAll {
		str.WriteString("all")
	} else {
		str.WriteString("one")
	}
	if len(name) > 0 {
		str.WriteString("-")
		str.WriteString(strings.ReplaceAll(strings.ReplaceAll(name[0], ":", "_"), ".", "_"))
	} else if !isAll {
		str.WriteString("-")
		str.WriteString(strings.ReplaceAll(strings.ReplaceAll(service.Socket.Name(), ":", "_"), ".", "_"))
	}
	return str.String()
}

// Listen 监听消息队列
func (s *Queue) Listen() {
	s.once.Do(func() {
		//查找指定地址在投递
		searchUnicast, err := queue.RegisterUnicast[*Message](s.client, "socket-send", runtime.NumCPU(), runtime.NumCPU()*2, func(message *Message) error {
			if message == nil {
				return errors.New("获取消息失败")
			}
			addr, err := service.Socket.GetConnectionIdAddr(message.Group, message.Id)
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
			return service.Socket.SendToLocalId(message.Group, message.Id, messageType, message.Data)
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
			_, _, _, err2 := service.Socket.SendToLocalAll(message.Group, messageType, message.Data)
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
func Register() *Queue {
	client, err := queue.RegisterNatsQueueClient(global.App.GetNatsUrl())
	if err != nil {
		panic(err)
		return nil
	}
	ins := &Queue{client: client}
	ins.Listen()
	global.App.Log.Println(ins.Name(false) + " queue初始化成功")
	global.App.Log.Println(ins.Name(true) + " queue初始化成功")
	return ins
}
