package queue

import (
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Task   *queue.UnicastClient[*task.Message]
	Reload *queue.BroadcastClient[*reload.Message]
)

func Init() {
	Task = task.Register()
	Reload = reload.Register()
}
