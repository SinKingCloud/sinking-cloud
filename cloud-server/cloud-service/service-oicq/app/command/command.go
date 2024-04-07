package command

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/command/queues"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/command/task"
)

func Init() {
	queues.Run() //执行消息队列
	task.Run()   //执行异步任务
}
