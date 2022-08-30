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
			//fmt.Println(s)
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
		serverLock.Lock()
		//构建消费协程
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
		}, Thread: 100}).Run()
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
			//只允许9点前执行第三次挂机
			if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*9 && time.Now().Unix()-stamp.Unix() > 2*3600 {
				data["check_time"] = str.DateTime(time.Now().Add(4.2 * 3600 * time.Second))
			}
			//检测挂机
			fail, _ = checkJob(server, uin, speedJobs, signJobs, true, false)
		} else {
			//正常挂机
			fail, _ = speedJob(server, uin, speedJobs, true)
			if !fail {
				fail, _ = checkJob(server, uin, speedJobs, signJobs, false, true) //判断未执行的秒加速项目并补挂一次
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

// checkJob 检测挂机失败的项目
func checkJob(server *model.AppUinServer, uin *model.AppUin, speedCheckJobs []int, signCheckJobs []int, lower bool, isFirst bool) (bool, map[int]bool) {
	//获取加速详情
	status, err, isFail, serverError := web_uin.GetUinSpeedStatus(server.Addr, server.Proxy, uin.Uin, true)
	if isFail && !serverError {
		return true, nil
	}
	if err != nil || status == nil || serverError {
		return false, nil
	}
	s := uin_tool.NewUinTool(server.Addr, server.Proxy)
	//执行签到任务
	if uin.UinStatus == uin2.Normal && len(signCheckJobs) > 0 {
		stamp, e2 := time.ParseInLocation("2006-01-02 15:04:05", uin.RunTime.String(), time.Local)
		//5小时执行一次
		if e2 == nil && stamp.Unix()+5*3600 < time.Now().Unix() {
			producer.NewSignJobProducer().Build(&producer.SignJob{
				Uin:   uin,
				Addr:  server.Addr,
				Proxy: server.Proxy,
				Jobs:  signCheckJobs,
			}).Send()
		}
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
			isCheck = status.IMobileQQOnline != "1"
			if !isCheck && lower {
				//下线手机
				s.MobileOffline(uin.Uin)
			}
			break
		case uin_tool.Pc:
			if isFirst {
				continue
			}
			isCheck = status.IPCQQOnline != "1"
			if !isCheck && lower {
				//下线PC
				s.PcOffline(uin.Uin)
			}
			break
		case uin_tool.Online:
			if isFirst {
				continue
			}
			isCheck = status.INoHideOnline != "1"
			break
		case uin_tool.Medal:
			isCheck = status.IMedal != "1"
			if isCheck {
				v = uin_tool.Pc
			}
			break
		case uin_tool.Zone:
			if isFirst {
				continue
			}
			isCheck = status.QzoneVisitor != "1"
			break
		case uin_tool.Ws:
			isCheck = status.WeishiVideoview != "1"
			break
		case uin_tool.MicroShow:
			isCheck = status.ICmshowActive != "1"
			break
		case uin_tool.Word:
			isCheck = status.ISmallWorld != "1"
			break
		case uin_tool.Safe:
			if isFirst {
				continue
			}
			isCheck = status.IPCSafeOnline != "1"
			break
		case uin_tool.Game:
			isCheck = status.IMobileGameOnline != "1"
			break
		case uin_tool.Sport:
			isCheck = false
			if time.Now().Day() <= 8 {
				isCheck = status.IQQSportStep != "1"
			}
			break
		case uin_tool.MobileClock:
			isCheck = status.IDailySign != "1"
			break
		case uin_tool.ZoneClock:
			isCheck = status.IQzoneState != "1"
			break
		case uin_tool.Friend:
			isCheck = status.IAddFriend != "1"
			break
		}
		if isCheck {
			speedJobs = append(speedJobs, v)
		}
	}
	if len(speedJobs) > 0 {
		return speedJob(server, uin, speedJobs, isFirst)
	}
	return false, nil
}

// speedJob 加速项目
func speedJob(server *model.AppUinServer, uin *model.AppUin, speedJobs []int, isFirst bool) (bool, map[int]bool) {
	fail, b := uin_tool.NewUinTool(server.Addr, server.Proxy).Speed(uin.Uin, uin.Password, uin.Phone, uin.PhoneArea, uin.LoginType, speedJobs, isFirst)
	if !fail {
		//访问加速页面
		e, skey, pskey := web_uin.GetUinCookies(server.Addr, server.Proxy, uin.Uin, "ti.qq.com", false)
		if e == nil && skey != "" && pskey != "" {
			t := uin_tool.NewSignTool(uin.Uin, server.Proxy, skey, pskey, "")
			_ = t.VisitLevelPage()
		}
	}
	return fail, b
}
