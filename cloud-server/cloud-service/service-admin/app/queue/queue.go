package queue

import (
	queuePermission "service-admin/app/queue/permission"

	queueUtil "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Permission *queueUtil.UnicastClient[*queuePermission.Task]
)

func Init() {
	Permission = queuePermission.Register()
}
