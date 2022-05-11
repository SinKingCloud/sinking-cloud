package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
)

// signJobIns 消费者单例
var signJobIns queue.MqConsumer

func NewSignJob() *signJob {
	return &signJob{}
}

// proxyJob 结构
type signJob struct {
}

func (l *signJob) Listen() {
	if signJobIns == nil {
		signJobIns = sinking.GetRocketConsumerServer(constant.SignJobGroup, false)
	}
	err := signJobIns.ListenReceiveMsgDo(constant.SignJobTopic, "", func(mqMsg queue.MqMsg) {
		s := &producer.SignJob{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			l.sign(s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

func (l *signJob) sign(info *producer.SignJob) {
	//执行签到任务
	sign := uin_tool.NewSignTool(info.Uin, info.Skey, info.PSkey)
	for _, v := range info.Jobs {
		switch v {
		case uin_tool.Vip:
			sign.VipSign() //会员签到
			break
		case uin_tool.BigVip:
			sign.BigVipSign() //大会员签到
			break
		case uin_tool.Hz:
			sign.YellowDiaSin() //黄钻签到
			break
		}
	}
}
