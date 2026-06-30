package log

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceCron + "-log"

// Message 任务消息
type Message struct {
	Log *model.Log
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
