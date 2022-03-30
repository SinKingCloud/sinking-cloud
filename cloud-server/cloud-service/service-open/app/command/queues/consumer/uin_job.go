package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// uinJobIns 消费者单例
var uinJobIns queue.MqConsumer

func NewUinJob() *uinJob {
	return &uinJob{}
}

// uinJob 结构
type uinJob struct {
}

func (l *uinJob) Listen() {
	if uinJobIns == nil {
		uinJobIns = sinking.GetRocketConsumerServer(constant.UinJobGroup, false)
	}
	err := uinJobIns.ListenReceiveMsgDo(constant.UinJobTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.AppUin{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			l.run(s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

// run 执行任务
func (l *uinJob) run(uin *model.AppUin) {
	//分协程执行
	fmt.Println("接收到消息", uin)
}
