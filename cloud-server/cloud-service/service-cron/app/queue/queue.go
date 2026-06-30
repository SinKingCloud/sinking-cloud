package queue

import (
	queueLog "service-cron/app/queue/log"
	queueRegister "service-cron/app/queue/register"
	queueReload "service-cron/app/queue/reload"
	queueTask "service-cron/app/queue/task"

	queueUtil "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Task     *queueUtil.UnicastClient[*queueTask.Message]
	Reload   *queueUtil.BroadcastClient[*queueReload.Message]
	Log      *queueUtil.UnicastClient[*queueLog.Message]
	Register *queueUtil.UnicastClient[*queueRegister.Message]
)

func Init() {
	Task = queueTask.Register()
	Reload = queueReload.Register()
	Log = queueLog.Register()
	Register = queueRegister.Register()
}
