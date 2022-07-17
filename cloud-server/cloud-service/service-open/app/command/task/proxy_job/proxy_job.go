package proxy_job

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"log"
	"time"
)

const checkTime = 300 //检测周期(秒)

// Run 定时检测任务投递
func Run() {
	//每10分钟检测代理状态
	go func() {
		oldMinute := (time.Now().Unix() - str.NewTimeTool().GetZeroTime(time.Now()).Unix()) / checkTime
		for {
			sinking_web.Try(func() {
				//获取当前执行分钟数
				zeroTime := str.NewTimeTool().GetZeroTime(time.Now())
				nowTime := time.Now()
				minute := (nowTime.Unix() - zeroTime.Unix()) / checkTime
				if oldMinute != minute {
					oldMinute = minute
					go Check(nowTime) //代理检测
				}
				time.Sleep(5 * time.Second)
			}, func(i interface{}) {
				log.Println("定时任务执行错误", i)
			})
		}
	}()
}

// Check 代理检测
func Check(now time.Time) {
	sinking_web.Try(func() {
		lock := constant.LockUinServerCheck
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinServerCheck) {
				list, err := server2.GetIns().SelectCheckJob(200, str.DateTime(now))
				if err == nil && len(list) > 0 {
					var ids []int
					for _, v := range list {
						//发送至消息队列
						if v.Proxy != "" {
							producer.NewProxyJobProducer().Build(v).Send()
						}
						ids = append(ids, v.Id)
					}
					//更改检测时间
					if len(ids) > 0 {
						_ = server2.GetIns().UpdateByIds(ids, map[string]interface{}{
							"check_time": now,
						})
					}
					util.Cache.UnLock(lock)
				} else {
					util.Cache.UnLock(lock)
					break
				}
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}, func(i interface{}) {
		log.Println("定时任务执行错误,详情", i)
	})
}
