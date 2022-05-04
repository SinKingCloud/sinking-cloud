package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// uinJobIns 消费者单例
var proxyJobIns queue.MqConsumer

func NewProxyJob() *proxyJob {
	return &proxyJob{}
}

// proxyJob 结构
type proxyJob struct {
}

func (l *proxyJob) Listen() {
	if proxyJobIns == nil {
		proxyJobIns = sinking.GetRocketConsumerServer(constant.ProxyJobGroup, false)
	}
	err := proxyJobIns.ListenReceiveMsgDo(constant.ProxyJobTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.AppUinProxy{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			//判断代理
			//log.Println("判断代理", s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}
