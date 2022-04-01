package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/util/uin_tool"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"gorm.io/gorm"
	"sync"
	"time"
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

var (
	servers    = make(map[int]chan interface{}) //服务器列表
	serverLock = &sync.Mutex{}                  //并发锁
)

// run 执行任务
func (l *uinJob) run(uin *model.AppUin) {
	if uin == nil {
		return
	}
	if servers[uin.ServerId] == nil {
		//构建消费协程
		serverLock.Lock()
		g := &sync.WaitGroup{}
		g.Add(1)
		go (&job.Task{Producer: func(channel chan interface{}) {
			servers[uin.ServerId] = channel
			g.Add(-1)
		}, Consumer: func(param interface{}) {
			if param != nil {
				//执行挂机任务
				sinking_web.Try(func() {
					speed(param.(*model.AppUin))
				}, func(i interface{}) {
					util.Log.Println(param, i)
				})
			}
		}, Thread: 40}).Run()
		g.Wait()
		serverLock.Unlock()
	}
	//按照服务分发至不同协程执行
	servers[uin.ServerId] <- uin
}

// speed 执行挂机任务
func speed(uin *model.AppUin) {
	//获取挂机服务器信息
	server, err := server2.GetIns().FindByIdWithCache(uin.ServerId, false)
	if err != nil || server == nil {
		return
	}
	//构建加速项目
	var jobs []int
	var userJob map[string]string
	err = json.Unmarshal([]byte(uin.Setting), &userJob)
	if err != nil || userJob == nil {
		jobs = []int{}
	}
	isCost := false
	for k := range web_uin.Setting {
		//类型判断
		num := -1
		if value, ok := web_uin.JobMap["speed"][k]; ok {
			if uin.PlayerDay > 0 {
				isCost = true
				num = value
			}
		} else if value, ok := web_uin.JobMap["sign"][k]; ok {
			if uin.SignDay > 0 {
				isCost = true
				num = value
			}
		}
		if num >= 0 {
			jobs = append(jobs, num)
		}
	}
	data := map[string]interface{}{
		"run_time": str.DateTime(time.Now()),
	}
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uin.RunTime.String(), time.Local)
	if isCost && stamp.Unix() < str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
		if uin.PlayerDay > 0 {
			data["player_day"] = gorm.Expr("`player_day`-1")
		}
		if uin.SignDay > 0 {
			data["sign_day"] = gorm.Expr("`sign_day`-1")
		}
		data["run_day"] = gorm.Expr("`run_day`+1")
	}
	if len(jobs) > 0 {
		s := uin_tool.NewUinTool(server.Addr)
		fail, _ := s.Speed(uin.Uin, uin.Password, jobs)
		if fail {
			data["uin_status"] = 1
		} else {
			data["uin_status"] = 0
		}
	}
	_ = uin2.GetIns().UpdateById(uin.Id, data)
}
