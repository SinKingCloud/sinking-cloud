package queue

import (
	"service-cron/app/queue/log"
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Task   *queue.UnicastClient[*task.Message]
	Reload *queue.BroadcastClient[*reload.Message]
	Log    *queue.UnicastClient[*log.Message]
)

func Init(taskHandle func(param *task.Message) error, reloadHandle func(param *reload.Message) error, logHandle func(param *log.Message) error) {
	Task = task.Register(taskHandle)
	Reload = reload.Register(reloadHandle)
	Log = log.Register(logHandle)
}
