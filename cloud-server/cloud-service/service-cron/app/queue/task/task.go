package task

import (
	"runtime"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceCron + "-task"

// Message 任务消息
type Message struct {
	Id int64 `json:"id"`
}

// Register 注册任务队列
func Register(consumer func(param *Message) error) *queue.UnicastClient[*Message] {
	producerThreadNum := runtime.NumCPU()
	if producerThreadNum <= 0 {
		producerThreadNum = 1
	}
	consumerThreadNum := producerThreadNum * 2
	ins, err := queue.RegisterNatsUnicast[*Message](TopicName, producerThreadNum, consumerThreadNum, consumer, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
