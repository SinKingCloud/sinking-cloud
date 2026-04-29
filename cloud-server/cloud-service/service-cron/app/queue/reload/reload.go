package reload

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceCron + "-reload"

type Action string

const (
	Reload Action = "reload"
	Check  Action = "check"
	Remove Action = "remove"
)

// Message 任务消息
type Message struct {
	Action  Action  `json:"action"`
	Version int64   `json:"version"`
	TaskIds []int64 `json:"task_ids"`
}

// Register 注册任务队列
func Register(consumer func(param *Message) error) *queue.BroadcastClient[*Message] {
	producerThreadNum := 1
	consumerThreadNum := 1
	ins, err := queue.RegisterNatsBroadcast[*Message](TopicName, producerThreadNum, consumerThreadNum, consumer, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
