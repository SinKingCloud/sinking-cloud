package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/command"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app/http/route"
)

func Run() {
	//初始化异步任务
	command.Init()
	//初始化api
	route.Init()
}
