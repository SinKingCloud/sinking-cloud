package queue

import (
	"service-admin/app/queue/permission"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

var (
	Permission *queue.UnicastClient[*permission.Task]
)

func Init() {
	Permission = permission.Register()
}
