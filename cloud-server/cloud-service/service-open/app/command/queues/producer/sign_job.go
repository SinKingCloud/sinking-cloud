package producer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// signJobProducerIns 生产者单例
var signJobProducerIns queue.MqProducer

// signJobProducerTask 异步发送
var signJobProducerTask chan *SignJob

// NewSignJobProducer 实例化
func NewSignJobProducer() *signJob {
	return &signJob{}
}

// signJob 结构
type signJob struct {
	d *SignJob
}

// SignJob 数据结构
type SignJob struct {
	Uin   *model.AppUin
	Addr  string
	Proxy string
	Jobs  []int
}

// Build 构建消息
func (l *signJob) Build(d *SignJob) *signJob {
	l.d = d
	return l
}

// run 异步执行
func (l *signJob) run() {
	go func(tasks chan *SignJob) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = signJobProducerIns.SendByteMsg(constant.SignJobTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
			}
		}
	}(signJobProducerTask)
}

// Send 发送消息
func (l *signJob) Send() {
	if signJobProducerIns == nil {
		signJobProducerIns = sinking.GetRocketProducerServer(constant.SignJobGroup, 1)
		signJobProducerTask = make(chan *SignJob)
		l.run()
	}
	signJobProducerTask <- l.d
}
