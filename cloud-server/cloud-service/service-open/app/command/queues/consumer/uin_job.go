package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"gorm.io/gorm"
	"sort"
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
	servers    = make(map[string]chan interface{}) //服务器列表
	serverLock = &sync.Mutex{}                     //并发锁

	serverIds    = make(map[int]*model.AppUinServer) //服务器列表
	serverIdLock = &sync.Mutex{}                     //并发锁
)

// run 执行任务
func (l *uinJob) run(uin *model.AppUin) {
	if uin == nil {
		return
	}
	//获取挂机服务器信息
	server, err := server2.GetIns().FindByIdWithCache(uin.ServerId, false)
	if err != nil || server == nil {
		return
	}
	if serverIds[uin.ServerId] == nil || serverIds[uin.ServerId].Status != server.Status || serverIds[uin.ServerId].Proxy != server.Proxy {
		serverIdLock.Lock()
		serverIds[uin.ServerId] = server
		serverIdLock.Unlock()
	}
	if serverIds[uin.ServerId].Status != server2.Normal {
		serverIdLock.Lock()
		serverIds[uin.ServerId].Proxy = ""
		serverIdLock.Unlock()
	}
	if servers[server.Addr] == nil {
		//构建消费协程
		serverLock.Lock()
		g := &sync.WaitGroup{}
		g.Add(1)
		go (&job.Task{Producer: func(channel chan interface{}) {
			servers[server.Addr] = channel
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
		}, Thread: 50}).Run()
		g.Wait()
		serverLock.Unlock()
	}
	//按照服务分发至不同协程执行
	servers[server.Addr] <- uin
}

// speed 执行挂机任务
func speed(uin *model.AppUin) {
	server := serverIds[uin.ServerId]
	if server == nil {
		return
	}
	//构建加速项目
	var speedJobs []int
	var signJobs []int
	var userJob map[string]string
	err := json.Unmarshal([]byte(uin.Setting), &userJob)
	if err != nil || userJob == nil {
		speedJobs = []int{}
	}
	isCost := false
	for k := range web_uin.Setting {
		//类型判断
		num := -1
		num2 := -1
		if value, ok := web_uin.JobMap["speed"][k]; ok {
			if uin.PlayerDay > 0 {
				isCost = true
				//判断是否开启项目
				if userJob[k] != "0" {
					num = value
				}
			}
		} else if value, ok := web_uin.JobMap["sign"][k]; ok {
			if uin.SignDay > 0 {
				isCost = true
				//判断是否开启项目
				if userJob[k] != "0" {
					num2 = value
				}
			}
		}
		if num >= 0 {
			speedJobs = append(speedJobs, num)
		}
		if num2 >= 0 {
			signJobs = append(signJobs, num2)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(speedJobs))) //项目降序执行
	sort.Sort(sort.Reverse(sort.IntSlice(signJobs)))  //项目降序执行
	data := map[string]interface{}{
		"run_time": str.DateTime(time.Now()),
	}
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uin.RunTime.String(), time.Local)
	isCheck := true //是否是检测任务
	if isCost && stamp.Unix() < str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
		data["run_day"] = gorm.Expr("`run_day`+1")
		isCheck = false
	}
	if len(speedJobs) > 0 || len(signJobs) > 0 {
		var fail bool
		if isCheck {
			lower := false
			if uin.PlayerDay < 1 && uin.SignDay < 1 {
				lower = true
			}
			//只允许9点前执行第三次挂机
			if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*9 && time.Now().Unix()-stamp.Unix() > 2*3600 {
				data["check_time"] = str.DateTime(time.Now().Add(4.2 * 3600 * time.Second))
			}
			//检测挂机
			fail, _ = check(server.Addr, server.Proxy, uin, speedJobs, signJobs, lower, false)
		} else {
			//正常挂机
			fail, _ = uin_tool.NewUinTool(server.Addr, server.Proxy).Speed(uin.Uin, uin.Password, uin.Phone, uin.LoginType, speedJobs, true)
			if !fail {
				fail, _ = check(server.Addr, server.Proxy, uin, speedJobs, signJobs, false, true) //判断未执行的秒加速项目并补挂一次
			}
			//只允许9点前执行第二次挂机
			if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*9 {
				data["check_time"] = str.DateTime(time.Now().Add(2.1 * 3600 * time.Second))
			}
		}
		if fail {
			data["uin_status"] = 1 //设置异常状态
		} else {
			data["uin_status"] = 0 //设置正常状态
		}
	}
	_ = uin2.GetIns().UpdateByUin(uin.Uin, data)
}

// check 检测挂机失败的项目
func check(addr string, proxy string, uin *model.AppUin, speedCheckJobs []int, signCheckJobs []int, lower bool, isFirst bool) (bool, map[int]bool) {
	status, err, isFail := web_uin.GetUinSpeedStatus(addr, proxy, uin.Uin)
	if isFail {
		return true, nil
	}
	if err != nil || status == nil {
		return false, nil
	}
	s := uin_tool.NewUinTool(addr, proxy)
	jobStatus := status.CommonInfoGetQQLevelInfo
	//执行签到任务
	if len(signCheckJobs) > 0 {
		producer.NewSignJobProducer().Build(&producer.SignJob{
			Uin:   uin.Uin,
			Addr:  addr,
			Proxy: proxy,
			Skey:  status.Skey,
			PSkey: status.Pskey,
			Jobs:  signCheckJobs,
		}).Send()
	}
	//判断加速项目并加速
	var speedJobs []int
	for _, v := range speedCheckJobs {
		isCheck := false
		switch v {
		case uin_tool.Mobile:
			if isFirst {
				continue
			}
			isCheck = jobStatus.IMobileQQOnline != "1"
			if !isCheck && lower {
				//下线手机
				s.MobileOffline(uin.Uin)
			}
			break
		case uin_tool.Pc:
			if isFirst {
				continue
			}
			isCheck = jobStatus.IPCQQOnline != "1"
			if !isCheck && lower {
				//下线PC
				s.PcOffline(uin.Uin)
			}
			break
		case uin_tool.Online:
			if isFirst {
				continue
			}
			isCheck = jobStatus.INoHideOnline != "1"
			break
		case uin_tool.Medal:
			isCheck = jobStatus.IMedal != "1"
			break
		case uin_tool.Zone:
			if isFirst {
				continue
			}
			isCheck = jobStatus.QzoneVisitor != "1"
			break
		case uin_tool.Ws:
			isCheck = jobStatus.WeishiVideoview != "1"
			break
		case uin_tool.MicroShow:
			isCheck = jobStatus.ICmshowActive != "1"
			break
		case uin_tool.Safe:
			if isFirst {
				continue
			}
			isCheck = jobStatus.IPCSafeOnline != "1"
			break
		case uin_tool.Game:
			isCheck = jobStatus.IMobileGameOnline != "1"
			break
		case uin_tool.Sport:
			isCheck = false
			if time.Now().Day() <= 8 {
				isCheck = jobStatus.IQQSportStep != "1"
			}
			break
		case uin_tool.MobileClock:
			isCheck = jobStatus.IDailySign != "1"
			break
		case uin_tool.ZoneClock:
			isCheck = jobStatus.IQzoneState != "1"
			break
		case uin_tool.Friend:
			isCheck = jobStatus.IAddFriend != "1"
			break
		}
		if isCheck {
			speedJobs = append(speedJobs, v)
		}
	}
	if len(speedJobs) > 0 {
		return s.Speed(uin.Uin, uin.Password, uin.Phone, uin.LoginType, speedJobs, isFirst)
	}
	return false, nil
}
