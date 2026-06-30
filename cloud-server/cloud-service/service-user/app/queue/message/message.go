package message

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

const TopicName = constant.ServiceUser + "-message"

// Task 创建站内信任务。
type Task struct {
	UserIds []int64 `json:"user_ids"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
}

// Register 注册站内信队列。
func Register() *queue.UnicastClient[*Task] {
	ins, err := queue.RegisterNatsUnicast[*Task](TopicName, 1, 1, nil, global.App.GetNatsUrl())
	if err != nil {
		global.App.Log.Println(TopicName+" queue初始化失败", err)
		panic(err)
		return nil
	}
	global.App.Log.Println(TopicName + " queue初始化成功")
	return ins
}
