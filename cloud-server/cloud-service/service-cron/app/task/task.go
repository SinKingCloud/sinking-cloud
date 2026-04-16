package task

import (
	"service-cron/app/task/reload"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cron"
)

var task = cron.New()

func Init() {
	defer task.Start()
	_ = task.AddFunc(reload.Register())
}
