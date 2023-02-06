package uin_job

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// Run 定时挂机任务投递
func Run() {
	go player()      //执行挂机任务
	go countUinNum() //执行同步数量任务
}

// player 轮询挂机
func player() {
	ch1 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch1 = channel
	}, Consumer: func(param interface{}) {
		send(param.(int64))
	}, Thread: 1}).Run()
	ch2 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch2 = channel
	}, Consumer: func(param interface{}) {
		check(param.(time.Time))
	}, Thread: 1}).Run()
	time.Sleep(1 * time.Second)
	oldMinute := (time.Now().Unix() - str.NewTimeTool().GetZeroTime(time.Now()).Unix()) / 60 //1分钟检测一次
	for {
		sinking_web.Try(func() {
			//获取当前执行分钟数
			zeroTime := str.NewTimeTool().GetZeroTime(time.Now())
			nowTime := time.Now()
			minute := (nowTime.Unix() - zeroTime.Unix()) / 60
			if oldMinute != minute {
				oldMinute = minute
				ch1 <- minute                                            //执行挂机任务
				ch2 <- zeroTime.Add(time.Minute * time.Duration(minute)) //执行检测任务
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
		lock := constant.LockUinJob + strconv.FormatInt(minute, 10)
		//查询当前时间需要代挂的QQ
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinJob) {
				list, err := uin2.GetIns().SelectAllUinJob(5000, minute, str.NewTimeTool().GetZeroTime(time.Now()).Format("2006-01-02 15:04:05"))
				if err == nil && len(list) > 0 {
					//需更改状态id
					var ids []int
					//需扣费id
					var idsPlayer []int
					var idsSign []int
					//避免重复挂机
					seeded := make(map[string]bool)
					//避免重复计费
					costed := make(map[string]bool)
					for _, v := range list {
						var speedUin *uin2.AppUinJob
						//查询重复QQ计费
						uinList, err1 := uin2.GetIns().SelectByUin(v.Uin)
						if err1 == nil && costed[v.Uin] == false && len(uinList) > 0 {
							costed[v.Uin] = true
							for _, v1 := range uinList {
								stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v1.RunTime.String(), time.Local)
								if stamp.Unix() < str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
									if v1.PlayerDay > 0 {
										idsPlayer = append(idsPlayer, v1.Id)
									}
									if v1.SignDay > 0 {
										idsSign = append(idsSign, v1.Id)
									}
								}
								ids = append(ids, v1.Id)
								//取出一个正常qq
								if v1.Status == uin2.Normal && v1.UinStatus == uin2.Normal {
									if speedUin == nil && (v1.PlayerDay > 0 || v1.SignDay > 0) {
										speedUin = v1
									}
									//取出一个剩余套餐天数最多的qq
									if speedUin != nil && (speedUin.PlayerDay < v1.PlayerDay || speedUin.SignDay < v1.SignDay) {
										speedUin = v1
									}
								}
							}
						}
						//发布到消息队列执行
						if speedUin != nil && seeded[speedUin.Uin] == false {
							//避免同一查询有重复
							seeded[speedUin.Uin] = true
							//发布至消息队列
							producer.NewUinJobProducer().Build(speedUin).Send()
						}
					}
					//更改等级挂机天数
					if len(idsPlayer) > 0 {
						splitIds(idsPlayer, 500, func(tempIds []int) {
							_ = uin2.GetIns().UpdateByIds(tempIds, map[string]interface{}{
								"player_day": gorm.Expr("`player_day`-1"),
							})
						})
					}
					//更改扩展挂机天数
					if len(idsSign) > 0 {
						splitIds(idsSign, 500, func(tempIds []int) {
							_ = uin2.GetIns().UpdateByIds(tempIds, map[string]interface{}{
								"sign_day": gorm.Expr("`sign_day`-1"),
							})
						})
					}
					//更改代挂状态
					if len(ids) > 0 {
						splitIds(ids, 500, func(tempIds []int) {
							_ = uin2.GetIns().UpdateByIds(tempIds, map[string]interface{}{
								"run_time": str.DateTime(time.Now()),
							})
						})
					}
					util.Cache.UnLock(lock)
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

// splitIds 分片数组
func splitIds(ids []int, slice int, callback func(tempIds []int)) {
	if slice <= 0 || slice > 1000 {
		slice = 500
	}
	index := 0
	num := len(ids)
	var tempIds []int
	for k, v := range ids {
		index++
		tempIds = append(tempIds, v)
		if index >= slice || k+1 == num {
			callback(tempIds)
			index = 0
			tempIds = []int{}
		}
	}
}

// check 补挂
func check(t time.Time) {
	//计算时间区间
	start := t.Format("2006-01-02 15:04:05")
	end := t.Add(time.Second * 59).Format("2006-01-02 15:04:05")
	sinking_web.Try(func() {
		lock := constant.LockUinCheckJob
		//查询当前时间需要补挂的QQ
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinCheckJob) {
				list, err := uin2.GetIns().SelectUinCheckJob(10000, start, end)
				if err == nil && len(list) > 0 {
					var ids []int //需要更改得ids
					//避免重复挂机
					seeded := make(map[string]bool)
					for _, v := range list {
						var speedUin *uin2.AppUinJob
						//查询重复QQ计费
						uinList, err1 := uin2.GetIns().SelectByUin(v.Uin)
						if err1 == nil && len(uinList) > 0 {
							for _, v1 := range uinList {
								ids = append(ids, v1.Id)
								//取出一个正常qq
								if v1.Status == uin2.Normal && v1.UinStatus == uin2.Normal && speedUin == nil && (v1.PlayerDay > 0 || v1.SignDay > 0) {
									speedUin = v1
								}
							}
						}
						//发布到消息队列执行
						if speedUin != nil && seeded[speedUin.Uin] == false {
							//避免同一查询有重复
							seeded[speedUin.Uin] = true
							//发布至消息队列
							producer.NewUinJobProducer().Build(speedUin).Send()
						}
					}
					//更改代挂状态
					if len(ids) > 0 {
						splitIds(ids, 500, func(tempIds []int) {
							_ = uin2.GetIns().UpdateByIds(tempIds, map[string]interface{}{
								"check_time": str.DateTime(str.NewTimeTool().GetZeroTime(time.Now())),
							})
						})
					}
					util.Cache.UnLock(lock)
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

// countUinNum 统计服务器QQ数量
func countUinNum() {
	for true {
		time.Sleep(600 * time.Second)
		var servers []*model.AppUinServer
		err := util.Database.Db.Model(&model.AppUinServer{}).Where("`pid` > 0").Find(&servers).Error
		if err != nil {
			continue
		}
		for _, v := range servers {
			var total int64
			e := util.Database.Db.Model(&model.AppUin{}).Where("`server_id` = ?", v.Id).Count(&total).Error
			if e != nil {
				continue
			}
			if int64(v.UinNum) != total {
				data := map[string]interface{}{
					"uin_num":     total,
					"update_time": str.DateTime(time.Now()),
				}
				e = util.Database.Db.Model(&model.AppUinServer{}).Where("`id` = ? ", v.Id).Updates(data).Error
				if e != nil {
					continue
				}
			}
		}
	}
}
