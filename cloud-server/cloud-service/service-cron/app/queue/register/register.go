package register

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceCron + "-register"

// Message 注册任务消息
type Message struct {
	Id         int64  `json:"id"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	Script     string `json:"script"`
	Concurrent int    `json:"concurrent"`
	Status     int    `json:"status"`
}

// Register 注册任务队列
func Register() *queue.UnicastClient[*Message] {
	producerThreadNum := 1
	consumerThreadNum := 1
	ins, err := queue.RegisterNatsUnicast[*Message](TopicName, producerThreadNum, consumerThreadNum, nil, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
