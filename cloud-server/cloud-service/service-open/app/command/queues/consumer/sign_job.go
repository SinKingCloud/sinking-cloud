package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin"
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
			//fmt.Println(s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

func (l *signJob) sign(info *producer.SignJob) {
	if info == nil || info.Uin == nil {
		return
	}
	//获取cookie
	err, skey, pSkey := web_uin.GetUinCookies(info.Addr, info.Proxy, info.Uin.Uin, "qzone.qq.com", false)
	if err != nil || skey == "" || pSkey == "" {
		return
	}
	//执行签到任务
	sign := uin_tool.NewSignTool(info.Uin.Uin, info.Proxy, skey, pSkey, "")
	for _, v := range info.Jobs {
		switch v {
		case uin_tool.Vip:
			sign.VipSign() //会员签到
			break
		case uin_tool.Music:
			sign.QQMusicSign() //音乐签到
			break
		case uin_tool.BigVip:
			sign.BigVipSign() //大会员签到
			break
		case uin_tool.Hz:
			sign.YellowDiaSin() //黄钻签到
			break
		case uin_tool.Lz:
			sign.BlueDiaSin() //蓝钻签到
			break
		case uin_tool.Ld:
			sign.DaleDouSign() //乐斗签到
			break
		case uin_tool.GameCenter:
			sign.QQGamSign() //QQ游戏签到
			e, s, p := web_uin.GetUinCookies(info.Addr, info.Proxy, info.Uin.Uin, "gamecenter.qq.com", false)
			if e != nil || s == "" || p == "" {
				return
			}
			t := uin_tool.NewSignTool(info.Uin.Uin, info.Proxy, s, p, "")
			t.GameCenterSign() //游戏中心签到
			break
		}
	}
}
