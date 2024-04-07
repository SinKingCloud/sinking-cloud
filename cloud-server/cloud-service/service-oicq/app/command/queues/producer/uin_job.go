package producer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
)

// uinJobProducerIns 生产者单例
var uinJobProducerIns queue.MqProducer

// uinJobProducerTask 异步发送
var uinJobProducerTask chan *uin2.AppUinJob

// NewUinJobProducer 实例化
func NewUinJobProducer() *uinJob {
	return &uinJob{}
}

// uinJob 结构
type uinJob struct {
	uin *uin2.AppUinJob
}

// Build 构建消息
func (l *uinJob) Build(d *uin2.AppUinJob) *uinJob {
	l.uin = d
	return l
}

// run 异步执行
func (l *uinJob) run() {
	go func(tasks chan *uin2.AppUinJob) {
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
		uinJobProducerTask = make(chan *uin2.AppUinJob)
		l.run()
	}
	uinJobProducerTask <- l.uin
}
