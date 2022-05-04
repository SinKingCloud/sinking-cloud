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
					speed(server, param.(*model.AppUin))
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
func speed(server *model.AppUinServer, uin *model.AppUin) {
	//构建加速项目
	var jobs []int
	var userJob map[string]string
	err := json.Unmarshal([]byte(uin.Setting), &userJob)
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
				//判断是否开启项目
				if userJob[k] == "1" {
					num = value
				}
			}
		} else if value, ok := web_uin.JobMap["sign"][k]; ok {
			if uin.SignDay > 0 {
				isCost = true
				//判断是否开启项目
				if userJob[k] == "1" {
					num = value
				}
			}
		}
		if num >= 0 {
			jobs = append(jobs, num)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(jobs))) //项目降序执行
	data := map[string]interface{}{
		"run_time": str.DateTime(time.Now()),
	}
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uin.RunTime.String(), time.Local)
	isCheck := true //是否是检测任务
	if isCost && stamp.Unix() < str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
		data["run_day"] = gorm.Expr("`run_day`+1")
		isCheck = false
	}
	if len(jobs) > 0 {
		var fail bool
		if isCheck {
			lower := true
			//只允许9点前执行第三次挂机
			if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*9 && time.Now().Unix()-stamp.Unix() > 2*3600 {
				data["check_time"] = str.DateTime(time.Now().Add(4.5 * 3600 * time.Second))
				lower = false
			}
			//检测挂机
			fail, _ = check(server.Addr, server.Key, uin, jobs, lower)
		} else {
			//正常挂机
			fail, _ = uin_tool.NewUinTool(server.Addr, server.Key).Speed(uin.Uin, uin.Password, uin.Phone, uin.LoginType, jobs, true)
			//只允许9点前执行第二次挂机
			if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*9 {
				data["check_time"] = str.DateTime(time.Now().Add(2.5 * 3600 * time.Second))
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
func check(addr string, key string, uin *model.AppUin, jobs []int, lower bool) (bool, map[int]bool) {
	s := uin_tool.NewUinTool(addr, key)
	jobStatus, err, isFail := s.Status(uin.Uin)
	if err != nil || jobStatus == nil {
		return false, nil
	}
	if isFail {
		return true, nil
	}
	var newJobs []int
	//判断项目并加速
	for _, v := range jobs {
		isCheck := false
		switch v {
		case uin_tool.Mobile:
			isCheck = jobStatus.IMobileQQOnline != "1"
			if !isCheck && lower {
				//下线手机
				s.MobileOffline(uin.Uin)
			}
			break
		case uin_tool.Pc:
			isCheck = jobStatus.IPCQQOnline != "1"
			if !isCheck && lower {
				//下线PC
				s.PcOffline(uin.Uin)
			}
			break
		case uin_tool.Online:
			isCheck = jobStatus.INoHideOnline != "1"
			break
		case uin_tool.Medal:
			isCheck = jobStatus.IMedal != "1"
			break
		case uin_tool.Zone:
			isCheck = jobStatus.QzoneVisitor != "1"
			break
		case uin_tool.Ws:
			isCheck = jobStatus.WeishiVideoview != "1"
			break
		case uin_tool.MicroShow:
			isCheck = jobStatus.ICmshowActive != "1"
			break
		case uin_tool.Safe:
			isCheck = jobStatus.IPCSafeOnline != "1"
			break
		case uin_tool.Game:
			isCheck = jobStatus.IMobileGameOnline != "1"
			break
		case uin_tool.Sport:
			isCheck = false //jobStatus.IQQSportStep != "1"
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
			newJobs = append(newJobs, v)
		}
	}
	if len(newJobs) > 0 {
		return s.Speed(uin.Uin, uin.Password, uin.Phone, uin.LoginType, newJobs, false)
	}
	return false, nil
}
