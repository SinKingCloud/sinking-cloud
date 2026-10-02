package user_log

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceAdmin + "-user-log"

// Task 用户操作日志消息。
type Task struct {
	Log *model.UserLog
}

func Register() *queue.UnicastClient[*Task] {
	ins, err := queue.RegisterNatsUnicast[*Task](TopicName, 1, 1, nil, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
