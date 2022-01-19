package command

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/queues"
)

// Init 初始化后台任务
func Init() {
	//启动消息队列
	queues.Run()
	//启动异步任务
	job.Run()
}
