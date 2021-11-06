package command

import "cloud-gateway/app/command/queues"

// Init 初始化后台任务
func Init() {
	//启动消息队列
	queues.Run()
}
