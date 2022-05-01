package producer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// proxyJobProducerIns 生产者单例
var proxyJobProducerIns queue.MqProducer

// proxyJobProducerTask 异步发送
var proxyJobProducerTask chan *model.AppUinProxy

// NewProxyJobProducer 实例化
func NewProxyJobProducer() *proxyJob {
	return &proxyJob{}
}

// proxyJob 结构
type proxyJob struct {
	data *model.AppUinProxy
}

// Build 构建消息
func (l *proxyJob) Build(d *model.AppUinProxy) *proxyJob {
	l.data = d
	return l
}

// run 异步执行
func (l *proxyJob) run() {
	go func(tasks chan *model.AppUinProxy) {
		for task := range tasks {
			bytes, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, err = proxyJobProducerIns.SendByteMsg(constant.ProxyJobTopic, "", bytes)
			if err != nil {
				util.Log.Println(err)
			}
		}
	}(proxyJobProducerTask)
}

// Send 发送消息
func (l *proxyJob) Send() {
	if proxyJobProducerIns == nil {
		proxyJobProducerIns = sinking.GetRocketProducerServer(constant.ProxyJobGroup, 1)
		proxyJobProducerTask = make(chan *model.AppUinProxy)
		l.run()
	}
	proxyJobProducerTask <- l.data
}
