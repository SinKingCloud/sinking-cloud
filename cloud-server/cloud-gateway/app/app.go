package app

import (
	"cloud-gateway/app/command"
	"cloud-gateway/app/http/route"
)

func Run() {
	//初始化异步任务
	command.Init()
	//初始化api
	route.Init()
}
