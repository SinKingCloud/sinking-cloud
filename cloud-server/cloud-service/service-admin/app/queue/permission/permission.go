package permission

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceAdmin + "-permission"

type Task struct {
	Name string
	Code string
	Path string
}

func Register() *queue.UnicastClient[*Task] {
	ins, err := queue.RegisterNatsUnicast[*Task](TopicName, 1, 1, func(task *Task) error {
		_ = service.Rbac.RegisterPermission(task.Name, task.Code, task.Path)
		return nil
	}, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
		return nil
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
