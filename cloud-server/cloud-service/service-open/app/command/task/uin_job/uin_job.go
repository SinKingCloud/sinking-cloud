package uin_job

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	sinking_web "github.com/SinKingCloud/sinking-go/sinking-web"
	"log"
	"strconv"
	"time"
)

// Run 定时挂机任务投递
func Run() {
	go player() //轮询执行
}

// player 轮询挂机
func player() {
	oldMinute := (time.Now().Unix() - str.NewTimeTool().GetZeroTime(time.Now()).Unix()) / 60
	for {
		sinking_web.Try(func() {
			//获取当前执行分钟数
			zeroTime := str.NewTimeTool().GetZeroTime(time.Now())
			nowTime := time.Now()
			minute := (nowTime.Unix() - zeroTime.Unix()) / 60
			if oldMinute != minute {
				oldMinute = minute
				go send(minute)                                             //执行挂机任务
				go check(zeroTime.Add(time.Minute * time.Duration(minute))) //执行检测任务
			}
			time.Sleep(5 * time.Second)
		}, func(i interface{}) {
			log.Println("定时任务执行错误", i)
		})
	}
}

// send 发送任务至队列
func send(minute int64) {
	sinking_web.Try(func() {
		lock := constant.LockUinJob + strconv.Itoa(int(minute))
		//查询当前时间需要代挂的QQ
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinJob) {
				list, err := uin2.GetIns().SelectUinJob(100, minute, str.NewTimeTool().GetZeroTime(time.Now()).Format("2006-01-02 15:04:05"))
				if err == nil && len(list) > 0 {
					var ids []int
					for _, v := range list {
						//发布到消息队列执行
						producer.NewUinJobProducer().Build(v).Send()
						ids = append(ids, v.Id)
					}
					//更改代挂状态
					_ = uin2.GetIns().UpdateByIds(ids, map[string]interface{}{
						"run_time": str.DateTime(time.Now()),
					})
				} else {
					util.Cache.UnLock(lock)
					break
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}, func(i interface{}) {
		log.Println("定时任务执行错误,位置", minute, "详情", i)
	})
}

// check 补挂
func check(t time.Time) {
	//计算时间区间
	start := t.Format("2006-01-02 15:04:05")
	end := t.Add(time.Second * 59).Format("2006-01-02 15:04:05")
	sinking_web.Try(func() {
		lock := constant.LockUinCheckJob + start
		//查询当前时间需要代挂的QQ
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinCheckJob) {
				list, err := uin2.GetIns().SelectUinCheckJob(100, start, end)
				if err == nil && len(list) > 0 {
					var ids []int
					for _, v := range list {
						//发布到消息队列执行
						producer.NewUinJobProducer().Build(v).Send()
						ids = append(ids, v.Id)
					}
					//更改代挂状态
					_ = uin2.GetIns().UpdateByIds(ids, map[string]interface{}{
						"run_time": str.DateTime(time.Now()),
					})
				} else {
					util.Cache.UnLock(lock)
					break
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}, func(i interface{}) {
		log.Println("定时任务执行错误,位置", start, "详情", i)
	})
}
