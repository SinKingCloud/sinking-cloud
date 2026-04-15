package reload

import (
	"fmt"
	"runtime"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceCron + "-reload"

type Action string

const (
	Update Action = "update"
	Create Action = "create"
	Delete Action = "delete"
)

// Message 任务消息
type Message struct {
	Action  Action  `json:"action"`
	TaskIds []int64 `json:"task_ids"`
}

// Register 注册任务队列
func Register() *queue.BroadcastClient[*Message] {
	producerThreadNum := runtime.NumCPU()
	if producerThreadNum <= 0 {
		producerThreadNum = 1
	}
	consumerThreadNum := producerThreadNum * 2
	ins, err := queue.RegisterNatsBroadcast[*Message](TopicName, producerThreadNum, consumerThreadNum, func(param *Message) error {
		fmt.Println("收到消息,重载数据")
		return nil
	}, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
