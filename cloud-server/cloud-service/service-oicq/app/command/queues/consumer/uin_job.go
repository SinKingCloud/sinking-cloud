package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/util/oicq"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	"gorm.io/gorm"
	"log"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// uinJobIns 消费者单例
var (
	uinJobIns     queue.MqConsumer
	onceUinJobIns sync.Once
	refreshTask   = utils.NewDelayQueue(10, 100000)
)

func NewUinJob() *uinJob {
	return &uinJob{}
}

// uinJob 结构
type uinJob struct {
}

func (l *uinJob) Listen() {
	onceUinJobIns.Do(func() {
		uinJobIns = sinking.GetRocketConsumerServer(constant.UinJobGroup, false)
	})
	l.autoStartThread()
	//go func() {
	//	list, _ := uin2.GetIns().SelectByUin("209297792")
	//	for _, v := range list {
	//		if v.PlayerDay > 0 && v.UinStatus == 0 {
	//			l.run(v)
	//			break
	//		}
	//	}
	//}()
	err := uinJobIns.ListenReceiveMsgDo(constant.UinJobTopic, "", func(mqMsg queue.MqMsg) {
		sinking_web.Try(func() {
			s := &uin2.AppUinJob{}
			err := json.Unmarshal(mqMsg.Body, s)
			if err == nil && s != nil {
				l.run(s)
			} else {
				log.Println("MqToObject出错", err, s)
			}
		}, func(i interface{}) {
			log.Println("获取MQ出错", i)
		})
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

var (
	servers    = make(map[int]chan interface{}) //服务器列表
	serverLock = &sync.RWMutex{}                //并发锁

	authThreadWait     = &sync.WaitGroup{} //动态协程池并发锁
	authThreadWaitLock = &sync.RWMutex{}
	taskWait           = &sync.WaitGroup{} //任务并发锁
)

// autoStartThread 动态协程池
func (l *uinJob) autoStartThread() {
	go func(l *uinJob) {
		for {
			l.RestartThread()
			time.Sleep(24 * time.Hour) //每24小时根据账号数量更新协程池
		}
	}(l)
}

// RestartThread 重启携程
func (l *uinJob) RestartThread() {
	//上锁,禁止新增任务
	authThreadWait.Add(1)
	defer authThreadWait.Done()
	authThreadWaitLock.Lock()
	defer authThreadWaitLock.Unlock()
	taskWait.Wait() //等待协程内所有任务执行完毕
	//查询所有挂机服务器
	list, err := server2.GetIns().SelectAllServer()
	if err == nil {
		for _, server := range list {
			l.restartThread(server)
		}
	}
}

// restartThread 重启协程池
func (l *uinJob) restartThread(server *model.AppUinServer) {
	//查询正常挂机账号数量
	num := uin2.GetIns().Count(map[string]string{"server_id": strconv.Itoa(server.Id), "uin_status": strconv.Itoa(int(uin2.StatusUinNormal))})
	threadNum := int(num / 120) //每120个账号开启一个协程
	//最小一个协程
	if threadNum <= 0 {
		threadNum = 1
	}
	//最大60协程
	if threadNum >= 60 {
		threadNum = 60
	}
	//更新缓存信息
	server, err := server2.GetIns().FindByIdWithCache(server.Id, true)
	if err != nil || server == nil || server.Status != server2.Normal {
		return
	}
	log.Println("创建协程池", server.Id, "数量", threadNum)
	server2.GetIns().ServerJobReset(strconv.Itoa(server.Id))
	var proxy *oicq.Proxy
	_ = json.Unmarshal([]byte(server.Proxy), &proxy)
	serverLock.RLock()
	defer serverLock.RUnlock()
	if servers[server.Id] != nil {
		close(servers[server.Id])
	}
	//实例化
	g := &sync.WaitGroup{}
	g.Add(1)
	go func(g *sync.WaitGroup, l *uinJob) {
		(&job.Task{Producer: func(channel chan interface{}) {
			servers[server.Id] = channel
			g.Add(-1)
		}, Consumer: func(param interface{}) {
			if param != nil {
				//执行挂机任务
				taskWait.Add(1) //任务数量+1
				sinking_web.Try(func() {
					l.speed(param.(*uin2.AppUinJob))
				}, func(i interface{}) {
					log.Println("执行出错", param, "错误", i)
				})
				taskWait.Done()
			}
		}, Thread: threadNum}).Run()
	}(g, l)
	g.Wait()
}

// run 执行任务
func (l *uinJob) run(uin *uin2.AppUinJob) {
	if uin == nil {
		return
	}
	sinking_web.Try(func() {
		authThreadWait.Wait() //等待协程池处理完毕
		for i := 0; i < 600; i++ {
			serverLock.RLock()
			if d, ok := servers[uin.ServerId]; ok && d != nil {
				d <- uin
				serverLock.RUnlock()
				break
			}
			serverLock.RUnlock()
			time.Sleep(time.Second)
		}
	}, func(i interface{}) {
		serverLock.RUnlock()
		log.Println("执行失败", uin)
	})
}

// speed 执行任务
func (l *uinJob) speed(uin *uin2.AppUinJob) {
	sinking_web.Try(func() {
		//开始限流
		serverId := strconv.Itoa(uin.ServerId)
		for {
			if server2.GetIns().ServerJobNum(serverId) > 100 {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}
		}
		if server2.GetIns().ServerJobNumInc(serverId) > 0 {
			defer func(serverId string) {
				if server2.GetIns().ServerJobNum(serverId) > 0 {
					server2.GetIns().ServerJobNumDec(serverId)
				}
			}(serverId)
		}
		uinInfo, err := uin2.GetIns().FindById(uin.Id)
		if err != nil || uinInfo == nil {
			log.Println("获取uinInfo失败", err, uinInfo)
			return
		}
		//限流结束
		log.Println("执行挂机", uinInfo.Uin)
		var proxy *oicq.Proxy
		server, err := server2.GetIns().FindByIdWithCache(uinInfo.ServerId, false)
		if err == nil && server != nil && server.Status == server2.Normal {
			_ = json.Unmarshal([]byte(server.Proxy), &proxy)
		}
		var jobs []oicq.Job
		var userJob map[web_uin.JobName]string
		var userJobStatus []web_uin.JobName
		err = json.Unmarshal([]byte(uinInfo.Setting), &userJob)
		if err != nil || userJob == nil {
			jobs = []oicq.Job{}
		}
		isCost := false
		for k, v := range web_uin.Setting {
			//类型判断
			num := oicq.None
			status := 0
			if value, ok := web_uin.JobMap[web_uin.Speed][k]; ok {
				if uin.PlayerDay > 0 {
					isCost = true
					//判断是否开启项目
					if d, o := userJob[k]; (o && d == "1") || (!o && v == "1") {
						num = value
						status = 1
					}
				}
			} else if value, ok = web_uin.JobMap[web_uin.Sign][k]; ok {
				if uin.SignDay > 0 {
					isCost = true
					//判断是否开启项目
					if d, o := userJob[k]; (o && d == "1") || (!o && v == "1") {
						num = value
						status = 1
					}
				}
			}
			if num != oicq.None {
				jobs = append(jobs, num)
			}
			if status == 1 {
				userJobStatus = append(userJobStatus, k)
			}
		}
		//更改状态为已提交
		_ = web_uin.SetUinJobStatus(uinInfo.Uin, userJobStatus, 2)
		log.Println("获取挂机开启项目", uinInfo.Uin, jobs)
		data := map[string]interface{}{"run_time": str.DateTime(time.Now())}
		stamp, e := time.ParseInLocation("2006-01-02 15:04:05", uin.RunTime.String(), time.Local)
		if e != nil {
			stamp = str.NewTimeTool().GetZeroTime(time.Now())
		}
		stamp2, e := time.ParseInLocation("2006-01-02 15:04:05", uin.CheckTime.String(), time.Local)
		if e != nil {
			stamp2 = str.NewTimeTool().GetZeroTime(time.Now())
		}
		isCheck := true //是否是检测任务
		if isCost && stamp.Unix() < str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
			data["run_day"] = gorm.Expr("`run_day`+1")
			isCheck = false
		}
		isSys := true //是否系统执行
		if isCheck && (stamp2.Unix() <= str.NewTimeTool().GetZeroTime(time.Now()).Unix() || time.Now().Unix() < stamp.Unix() || time.Now().Unix() > stamp2.Unix()) {
			isSys = false
		}
		if len(jobs) > 0 {
			var fail bool
			var token []byte
			if isCheck {
				//第二次挂机
				fail, _, token = l.speedJob(proxy, uinInfo, jobs, false, isSys)
			} else {
				//第一次挂机
				fail, _, token = l.speedJob(proxy, uinInfo, jobs, true, isSys)
				//6.5小时后补挂一次
				if time.Now().Unix()-str.NewTimeTool().GetZeroTime(time.Now()).Unix() < 3600*17 {
					data["check_time"] = str.DateTime(time.Now().Add(6.5 * 3600 * time.Second))
				}
			}
			if token != nil && len(token) > 0 {
				data["token"] = utils.Bin2HexTo(token)
			}
			if fail {
				data["uin_status"] = uin2.StatusUinError //设置异常状态
			} else {
				data["uin_status"] = uin2.StatusUinNormal //设置正常状态
			}
		}
		_ = uin2.GetIns().UpdateByUin(uin.Uin, data)
	}, func(i interface{}) {
		log.Println("执行任务出错", uin, "错误", i)
	})
}

// speedJob 加速项目
func (l *uinJob) speedJob(proxy *oicq.Proxy, uinInfo *model.AppUin, jobs []oicq.Job, isFirst bool, isSys bool) (bool, map[int]client.SpeedStatus, []byte) {
	log.Println("开始挂机", uinInfo.Uin)
	cli, err := web_uin.GetOicqInsWithOption(uinInfo, proxy, false)
	if err != nil {
		log.Println("获取挂机实例失败", uinInfo.Uin)
		return false, nil, nil
	}
	defer func(cli *oicq.Oicq) {
		cli.Mobile.Release()
		cli.Mobile.Disconnect()
	}(cli)
	//登陆手机
	err, loginResponse := cli.AutoLogin()
	log.Println("挂机登陆结果", uinInfo.Uin, err, loginResponse)
	if err != nil {
		l.saveJobStatus(uinInfo, 0, nil)
		return false, nil, nil
	}
	if !loginResponse.Success {
		l.saveJobStatus(uinInfo, 0, nil)
		return true, nil, nil
	}
	//查询cookie项
	cookies := make(map[oicq.Job]string)
	err, wsCookie := web_uin.GetUinCookies(uinInfo.Uin, "h5.weishi.qq.com", true)
	if err == nil {
		cookies[oicq.Ws] = wsCookie
	}
	err, musicCookie := web_uin.GetUinCookies(uinInfo.Uin, "y.qq.com", true)
	if err == nil {
		cookies[oicq.QQMusic] = musicCookie
	}
	//保存相关cookie
	cookie := cli.Mobile.GetAllCookies()
	for k, v := range cookie {
		_ = web_uin.SetUinCookies(uinInfo.Uin, k, v)
	}
	if !isFirst {
		//获取未执行的项目
		e, s := cli.Mobile.Speed.Info()
		if e == nil {
			jobs = l.getUnDoJobs(jobs, s, isSys)
			if len(jobs) <= 0 {
				//保存加速状态
				l.saveJobStatus(uinInfo, 1, s)
			}
		}
	}
	log.Println("获取挂机任务", uinInfo.Uin, jobs)
	//jobs = []oicq.Job{oicq.Pc}
	if len(jobs) <= 0 {
		return false, nil, cli.LoginToken
	}
	//执行挂机项
	log.Println("执行挂机项目", uinInfo.Uin, jobs, cookies)
	r := make(map[oicq.Job]client.SpeedCode)
	resCookie := make(map[oicq.Job]string)
	sinking_web.Try(func() {
		r, resCookie = cli.ActiveProject(jobs, cookies)
	}, func(i interface{}) {
		log.Println("执行挂机项目出错", uinInfo.Uin, jobs, cookies, i, debug.Stack())
	})
	log.Println("执行挂机结果", uinInfo.Uin, r, resCookie)
	if resCookie[oicq.Ws] != "" && (wsCookie == "" || resCookie[oicq.Ws] != wsCookie) {
		_ = web_uin.SetUinCookies(uinInfo.Uin, "h5.weishi.qq.com", resCookie[oicq.Ws])
	}
	if resCookie[oicq.QQMusic] != "" && (musicCookie == "" || resCookie[oicq.QQMusic] != musicCookie) {
		_ = web_uin.SetUinCookies(uinInfo.Uin, "y.qq.com", resCookie[oicq.QQMusic])
	}
	//保存加速状态
	if e, s := cli.Mobile.Speed.Info(); e == nil {
		l.saveJobStatus(uinInfo, 1, s)
	}
	//延迟刷新状态
	refreshTask.Push(150*time.Minute, func() {
		sinking_web.Try(func() {
			if cli != nil && cli.Mobile != nil {
				cli.Mobile.Speed.RefreshAllTaskStatus()
				if e, s := cli.Mobile.Speed.Info(); e == nil && s != nil {
					_ = web_uin.SetUinSpeedStatus(uinInfo.Uin, s)
				}
			}
		}, func(i interface{}) {
			log.Println("延迟队列出错", i)
		})
	})
	return false, nil, cli.LoginToken
}

// saveJobStatus 保存加速状态
func (l *uinJob) saveJobStatus(uin *model.AppUin, status int, speedStatus *client.SpeedStatus) {
	sinking_web.Try(func() {
		//把所有状态改为已加速
		var userJob map[web_uin.JobName]string
		err2 := json.Unmarshal([]byte(uin.Setting), &userJob)
		var jobs []web_uin.JobName
		for k, v := range web_uin.Setting {
			if err2 != nil || userJob == nil || (userJob[k] == "" && v == "1") || userJob[k] != "0" {
				jobs = append(jobs, k)
			}
		}
		_ = web_uin.SetUinJobStatus(uin.Uin, jobs, status)
		//更改加速状态
		if speedStatus != nil {
			_ = web_uin.SetUinSpeedStatus(uin.Uin, speedStatus)
		}
	}, func(i interface{}) {
		log.Println("保存加速状态错误", i)
	})
}

// getUnDoJobs 获取未执行的任务
func (l *uinJob) getUnDoJobs(jobs []oicq.Job, status *client.SpeedStatus, isSys bool) (newJobs []oicq.Job) {
	sinking_web.Try(func() {
		for _, v := range jobs {
			isCheck := false
			switch v {
			case oicq.Pc:
				isCheck = status.IPCQQOnline != "1"
				break
			case oicq.Ws:
				isCheck = status.WeishiVideoview != "1"
				break
			case oicq.Word:
				isCheck = status.ISmallWorld != "1"
				break
			case oicq.MobileClock:
				isCheck = status.IDailySign != "1"
				break
			case oicq.ZoneClock:
				isCheck = status.IQzoneState != "1"
				break
			case oicq.Friend:
				isCheck = status.IAddFriend != "1"
				break
			case oicq.QQMusic:
				isCheck = status.IQQMusic != "1"
				break
			case oicq.CmsShowPro:
				isCheck = status.ICmsShowPro != "1"
				break
			default:
				isCheck = !isSys //用户提交则补挂签到
				break
			}
			if isCheck {
				newJobs = append(newJobs, v)
			}
		}
	}, func(i interface{}) {
		log.Println("获取未执行的任务错误", i)
	})
	return newJobs
}
