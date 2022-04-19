package producer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// uinJobProducerIns 生产者单例
var uinJobProducerIns queue.MqProducer

// uinJobProducerTask 异步发送
var uinJobProducerTask chan *model.AppUin

// NewUinJobProducer 实例化
func NewUinJobProducer() *uinJob {
	return &uinJob{}
}

// uinJob 结构
type uinJob struct {
	log *model.AppUin
}

// Build 构建消息
func (l *uinJob) Build(d *model.AppUin) *uinJob {
	l.log = d
	return l
}

// run 异步执行
func (l *uinJob) run() {
	go func(tasks chan *model.AppUin) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = uinJobProducerIns.SendByteMsg(constant.UinJobTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
			}
		}
	}(uinJobProducerTask)
}

// Send 发送消息
func (l *uinJob) Send() {
	if uinJobProducerIns == nil {
		uinJobProducerIns = sinking.GetRocketProducerServer(constant.UinJobGroup, 1)
		uinJobProducerTask = make(chan *model.AppUin)
		l.run()
	}
	uinJobProducerTask <- l.log
}
