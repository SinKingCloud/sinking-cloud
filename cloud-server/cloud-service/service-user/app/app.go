package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/command"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/route"
)

func Run() {
	//启动异步任务
	command.Init()
	//初始化api
	route.Init()
}
