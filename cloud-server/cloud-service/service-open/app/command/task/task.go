package task

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/command/task/proxy_job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/command/task/uin_job"
)

func Run() {
	uin_job.Run()
	proxy_job.Run()
}
