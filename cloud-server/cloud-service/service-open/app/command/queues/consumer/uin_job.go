package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/oicq"
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
	l.autoStartThread()
	//go func() {
	//	list, _ := uin2.GetIns().SelectByUin("2903033707") //44422323
	//	for _, v := range list {
	//		if v.PlayerDay > 0 {
	//			l.run(v)
	//			break
	//		}
	//	}
	//}()
	err := uinJobIns.ListenReceiveMsgDo(constant.UinJobTopic, "", func(mqMsg queue.MqMsg) {
		s := &uin2.AppUinJob{}
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
	serverLock = &sync.RWMutex{}                //并发锁

	serverIds    = make(map[int]*oicq.Proxy) //服务器列表
	serverIdLock = &sync.RWMutex{}           //并发锁

	authThreadWait     = &sync.WaitGroup{} //动态协程池并发锁
	authThreadWaitLock = &sync.RWMutex{}
	taskWait           = &sync.WaitGroup{} //任务并发锁
	taskWaitLock       = &sync.RWMutex{}
)

// autoStartThread 动态协程池
func (l *uinJob) autoStartThread() {
	go func(l *uinJob) {
		for true {
			//上锁,禁止新增任务
			authThreadWaitLock.RLock()
			authThreadWait.Add(1)
			authThreadWaitLock.RUnlock()
			taskWaitLock.RLock()
			taskWait.Wait() //等待协程内所有任务执行完毕
			taskWaitLock.RUnlock()
			//查询所有挂机服务器
			list, err := server2.GetIns().SelectAllServer()
			if err == nil {
				for _, server := range list {
					l.restartThread(server)
				}
			}
			authThreadWaitLock.RLock()
			authThreadWait.Done()
			authThreadWaitLock.RUnlock()
			time.Sleep(24 * time.Hour) //每12小时根据账号数量更新协程池
		}
	}(l)
}

// restartThread 重启协程池
func (l *uinJob) restartThread(server *model.AppUinServer) {
	//查询正常挂机账号数量
	num := uin2.GetIns().Count(map[string]string{"server_id": strconv.Itoa(server.Id), "uin_status": strconv.Itoa(int(uin2.StatusUinNormal))})
	threadNum := int(num / 100) //每100个账号开启一个协程
	//最小一个协程
	if threadNum <= 0 {
		threadNum = 1
	}
	//最大100协程
	if threadNum >= 100 {
		threadNum = 100
	}
	//更新缓存信息
	server, err := server2.GetIns().FindByIdWithCache(server.Id, false)
	if err != nil || server == nil || server.Status != server2.Normal {
		return
	}
	log.Println("创建协程池", server.Id, "数量", threadNum)
	var proxy *oicq.Proxy
	_ = json.Unmarshal([]byte(server.Proxy), &proxy)
	serverIdLock.RLock()
	serverIds[server.Id] = proxy
	serverIdLock.RUnlock()
	serverLock.RLock()
	defer serverLock.RUnlock()
	if servers[server.Id] != nil {
		close(servers[server.Id])
	}
	//实例化
	g := &sync.WaitGroup{}
	g.Add(1)
	go func(g *sync.WaitGroup) {
		(&job.Task{Producer: func(channel chan interface{}) {
			servers[server.Id] = channel
			g.Add(-1)
		}, Consumer: func(param interface{}) {
			if param != nil {
				//执行挂机任务
				taskWaitLock.RLock()
				taskWait.Add(1) //任务数量+1
				taskWaitLock.RUnlock()
				sinking_web.Try(func() {
					l.speed(param.(*uin2.AppUinJob))
				}, func(i interface{}) {
					util.Log.Println(param, i, string(debug.Stack()))
				})
				taskWaitLock.RLock()
				taskWait.Done()
				taskWaitLock.RUnlock()
			}
		}, Thread: threadNum}).Run()
	}(g)
	g.Wait()
}

// run 执行任务
func (l *uinJob) run(uin *uin2.AppUinJob) {
	if uin == nil {
		return
	}
	var userJob map[web_uin.JobName]string
	err := json.Unmarshal([]byte(uin.Setting), &userJob)
	if err == nil && userJob != nil {
		var jobs []web_uin.JobName
		for k, v := range web_uin.Setting {
			if (userJob[k] == "" && v == "1") || userJob[k] != "0" {
				jobs = append(jobs, k)
			}
		}
		//更改状态为已提交
		_ = web_uin.SetUinJobStatus(uin.Uin, jobs, 2)
	}
	authThreadWait.Wait() //等待协程池处理完毕
	for i := 0; i < 30; i++ {
		serverLock.RLock()
		if d, ok := servers[uin.ServerId]; ok && d != nil {
			d <- uin
			serverLock.RUnlock()
			break
		}
		serverLock.RUnlock()
		time.Sleep(time.Second)
	}
}

// speed 执行任务
func (l *uinJob) speed(uin *uin2.AppUinJob) {
	log.Println("执行挂机", uin)
	serverIdLock.Lock()
	var proxy *oicq.Proxy
	if d, ok := serverIds[uin.ServerId]; ok && d != nil {
		proxy = d
	}
	serverIdLock.Unlock()
	var jobs []oicq.Job
	var userJob map[web_uin.JobName]string
	err := json.Unmarshal([]byte(uin.Setting), &userJob)
	if err != nil || userJob == nil {
		jobs = []oicq.Job{}
	}
	isCost := false
	for k, v := range web_uin.Setting {
		//类型判断
		num := oicq.None
		if value, ok := web_uin.JobMap[web_uin.Speed][k]; ok {
			if uin.PlayerDay > 0 {
				isCost = true
				//判断是否开启项目
				if (userJob[k] == "" && v == "1") || userJob[k] != "0" {
					num = value
				}
			}
		} else if value, ok = web_uin.JobMap[web_uin.Sign][k]; ok {
			if uin.SignDay > 0 {
				isCost = true
				//判断是否开启项目
				if (userJob[k] == "" && v == "1") || userJob[k] != "0" {
					num = value
				}
			}
		}
		if num != oicq.None {
			jobs = append(jobs, num)
		}
	}
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
	if isCheck && (stamp2.Unix() <= str.NewTimeTool().GetZeroTime(time.Now()).Unix() || time.Now().Unix() < stamp.Unix()) {
		isSys = false
	}
	if len(jobs) > 0 {
		var fail bool
		var token []byte
		if isCheck {
			//第二次挂机
			fail, _, token = l.speedJob(proxy, uin, jobs, false, isSys)
		} else {
			//第一次挂机
			fail, _, token = l.speedJob(proxy, uin, jobs, true, isSys)
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
}

// speedJob 加速项目
func (l *uinJob) speedJob(proxy *oicq.Proxy, uin *uin2.AppUinJob, jobs []oicq.Job, isFirst bool, isSys bool) (bool, map[int]client.SpeedStatus, []byte) {
	log.Println("开始挂机", uin)
	uinInfo, err := uin2.GetIns().FindById(uin.Id)
	if err != nil || uinInfo == nil {
		return false, nil, nil
	}
	cli, err := web_uin.GetOicqInsWithOption(uinInfo, proxy, false)
	if err != nil {
		log.Println("获取挂机实例失败", uin)
		return false, nil, nil
	}
	defer func(cli *oicq.Oicq) {
		cli.Mobile.Release()
		cli.Mobile.Disconnect()
	}(cli)
	//查询cookie项
	cookies := make(map[oicq.Job]string)
	err, wsCookie := web_uin.GetUinCookies(uin.Uin, "h5.weishi.qq.com", true)
	if err == nil {
		cookies[oicq.Ws] = wsCookie
	}
	//登陆手机
	err, loginResponse := cli.AutoLogin()
	log.Println("挂机登陆结果", uin.Uin, err, loginResponse)
	if err != nil {
		l.saveJobStatus(uin, 0, nil)
		return false, nil, nil
	}
	if !loginResponse.Success {
		l.saveJobStatus(uin, 0, nil)
		return true, nil, nil
	}
	//保存相关cookie
	cookie := cli.Mobile.GetAllCookies()
	for k, v := range cookie {
		_ = web_uin.SetUinCookies(uin.Uin, k, v)
	}
	if !isFirst {
		//获取未执行的项目
		e, s := cli.Mobile.Speed.Info()
		if e == nil {
			jobs = l.getUnDoJobs(jobs, s, isSys)
			if len(jobs) <= 0 {
				//保存加速状态
				l.saveJobStatus(uin, 1, s)
			}
		}
	}
	log.Println("获取挂机任务", uin.Uin, jobs)
	//jobs = []oicq.Job{oicq.Word}
	if len(jobs) <= 0 {
		return false, nil, cli.LoginToken
	}
	//执行挂机项
	log.Println("执行挂机项目", uin.Uin, jobs, cookies)
	r := make(map[oicq.Job]client.SpeedCode)
	resCookie := make(map[oicq.Job]string)
	sinking_web.Try(func() {
		r, resCookie = cli.ActiveProject(jobs, cookies)
	}, func(i interface{}) {
		log.Println("执行挂机项目出错", uin.Uin, jobs, cookies, i, string(debug.Stack()))
	})
	log.Println("执行挂机结果", uin.Uin, r, resCookie)
	if resCookie[oicq.Ws] != "" && (wsCookie == "" || resCookie[oicq.Ws] != wsCookie) {
		_ = web_uin.SetUinCookies(uin.Uin, "h5.weishi.qq.com", resCookie[oicq.Ws])
	}
	//保存加速状态
	if e, s := cli.Mobile.Speed.Info(); e == nil {
		l.saveJobStatus(uin, 1, s)
	}
	return false, nil, cli.LoginToken
}

// saveJobStatus 保存加速状态
func (l *uinJob) saveJobStatus(uin *uin2.AppUinJob, status int, speedStatus *client.SpeedStatus) {
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
}

// getUnDoJobs 获取未执行的任务
func (l *uinJob) getUnDoJobs(jobs []oicq.Job, status *client.SpeedStatus, isSys bool) []oicq.Job {
	var newJobs []oicq.Job
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
		case oicq.Dress:
			isCheck = status.IDress != "1"
			break
		default:
			isCheck = !isSys //用户提交则补挂签到
			break
		}
		if isCheck {
			newJobs = append(newJobs, v)
		}
	}
	return newJobs
}
