package queue

import (
	queuePermission "service-admin/app/queue/permission"
	queueUserLog "service-admin/app/queue/user_log"

	queueUtil "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Permission *queueUtil.UnicastClient[*queuePermission.Task]
	UserLog    *queueUtil.UnicastClient[*queueUserLog.Task]
)

func Init() {
	Permission = queuePermission.Register()
	UserLog = queueUserLog.Register()
}
