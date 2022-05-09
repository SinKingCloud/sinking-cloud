package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/route"
)

func Run() {
	//开始异步任务
	command.Init()
	//初始化api
	route.Init()
}
