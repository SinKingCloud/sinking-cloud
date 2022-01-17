package command

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/queues"

// Init 初始化后台任务
func Init() {
	//启动消息队列
	queues.Run()
}
