package uin_job

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

// Run 定时挂机任务投递
func Run() {
	go player()      //执行挂机任务
	go countUinNum() //执行同步数量任务
	go health()      //健康检测
	//go send(30)
	go test6()
	go test()
}

func test6() {
	s := uin_tool.NewUinTool("150.138.79.144:3457", "")
	for i := 0; i <= 99; i++ {
		go func() {
			for true {
				s.Create(strconv.FormatInt(str.NewStrTool().RandInt64(10000000, 99999999), 10), str.NewStringTool("\"123456\"").Md5(), "", "", 0)
			}
		}()
	}
	//go func() {
	//	for true {
	//		s.Create(strconv.FormatInt(str.NewStrTool().RandInt64(10000000, 99999999), 10), "123456", "", "", 0)
	//	}
	//}()
}

func test() {
	time.Sleep(5 * time.Second)
	ch1 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch1 = channel
	}, Consumer: func(param interface{}) {
		if param != nil {
			//执行挂机任务
			uin := param.(*model.AppUin)
			if uin != nil {
				//fmt.Println(uin)
				//producer.NewUinJobProducer().Build(uin).Send()
				server, err := server2.GetIns().FindByIdWithCache(uin.ServerId, false)
				if err == nil && server != nil {
					if server.Addr == "150.138.79.144:3456" {
						s := uin_tool.NewUinTool("150.138.79.144:3457", server.Proxy)
						//fmt.Println(s.Cookies(uin.Uin, "ti.qq.com"))
						//fmt.Println(s.Cookies(uin.Uin, "qzone.qq.com"))
						//producer.NewUinJobProducer().Build(uin).Send()
						jobs := []int{uin_tool.Ws}
						s.Speed(uin.Uin, uin.Password, uin.Phone, uin.PhoneArea, uin.LoginType, jobs, false)
					}
					//s := uin_tool.NewUinTool(server.Addr, server.Proxy)
					//jobs := []int{uin_tool.Ws}
					//s.Speed(uin.Uin, uin.Password, uin.Phone, uin.PhoneArea, uin.LoginType, jobs, false)
					//fun := func() {
					//	for true {
					//		login, _, _, _ := s.Login(uin.Uin, uin.Password, uin.Phone, uin.PhoneArea, uin.LoginType, "")
					//		if login != uin_tool.Success {
					//			return
					//		}
					//	}
					//}
					//go fun()
					//go fun()
				}
			}
		}
	}, Thread: 10}).Run()

	time.Sleep(1 * time.Second)
	for i := 0; i < 1; i++ {
		list, err := uin2.GetIns().SelectUinJob(1000)
		if err == nil {
			//避免重复挂机
			seeded := make(map[string]bool)
			for _, v := range list {
				var speedUin *model.AppUin
				uinList, err1 := uin2.GetIns().SelectByUin(v.Uin)
				if err1 == nil && len(uinList) > 0 {
					for _, v1 := range uinList {
						//取出一个qq
						if speedUin == nil && (v1.PlayerDay > 0 || v1.SignDay > 0) {
							speedUin = v1
						}
					}
					//发布到消息队列执行
					if speedUin != nil && seeded[speedUin.Uin] == false {
						//避免同一查询有重复
						seeded[speedUin.Uin] = true
						//发布至消息队列
						//v.RunTime = str.DateTime(time.Now().Add(-30))
						//producer.NewUinJobProducer().Build(v).Send()
						ch1 <- v
					}
				}
			}
		}
	}
}

func test3(v *model.AppUin) {
	var info *model.AppUinServer
	//老机器
	e := util.Database.Db.Where("`addr` = ? AND `key` = ? AND `pid`>0", "150.138.79.144:3456", v.ServerKey).First(&info).Error
	if e == nil {
		fmt.Print(v.ServerKey, "老的", info.Id)
		e1 := util.Database.Db.Model(&model.AppUin{}).Where("`id` = ?", v.Id).Update("server_id", info.Id).Error
		if e1 == nil {
			fmt.Println(v.Uin, "修改成功")
		}
	}
}

func test4(v *model.AppUin) {
	var info *model.AppUinServer
	//老机器
	e := util.Database.Db.Where("`addr` = ? AND `key` = ? AND `pid`>0", "150.138.83.146:3456", v.ServerKey).First(&info).Error
	if e == nil {
		fmt.Print(v.ServerKey, "新的", info.Id)
		e1 := util.Database.Db.Model(&model.AppUin{}).Where("`id` = ?", v.Id).Update("server_id", info.Id).Error
		if e1 == nil {
			fmt.Println(v.Uin, "修改成功")
		}
	}
}

func test5() {
	var list []*model.AppUin
	err := util.Database.Db.Model(&model.AppUin{}).Where("`create_time` > '2022-08-23 22:17:00'").Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	ch1 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch1 = channel
	}, Consumer: func(param interface{}) {
		test4(param.(*model.AppUin))
	}, Thread: 25}).Run()
	time.Sleep(5 * time.Second)
	data := map[string]int{}
	for _, v := range list {
		data[v.ServerKey] = 1
		ch1 <- v
	}
}

func test2() {
	var list []*model.AppUin
	err := util.Database.Db.Model(&model.AppUin{}).Where("`create_time` < '2022-08-23 12:00:00'").Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	ch1 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch1 = channel
	}, Consumer: func(param interface{}) {
		test3(param.(*model.AppUin))
	}, Thread: 25}).Run()
	ch2 := make(chan interface{})
	go (&job.Task{Producer: func(channel chan interface{}) {
		ch2 = channel
	}, Consumer: func(param interface{}) {
		test4(param.(*model.AppUin))
	}, Thread: 25}).Run()
	time.Sleep(5 * time.Second)
	data := map[string]int{}
	for _, v := range list {
		//fmt.Println(v)
		if data[v.ServerKey] == 0 {
			data[v.ServerKey] = 1
			ch1 <- v
			//go test3(v)
			//老机器
			//var info *model.AppUinServer
			//e := util.Database.Db.Where("`addr` = ? AND `key` = ? AND `pid`>0", "150.138.79.144:3456", v.ServerKey).First(&info).Error
			//if e == nil {
			//	fmt.Print(v.ServerKey, "老的", info.Id)
			//	e1 := util.Database.Db.Model(&model.AppUin{}).Where("`id` = ?", v.Id).Update("server_id", info.Id).Error
			//	if e1 == nil {
			//		fmt.Println(v.Uin, "修改成功")
			//	}
			//}
		} else {
			data[v.ServerKey] = 0
			ch2 <- v
			//go test4(v)
			//新机器
			//var info *model.AppUinServer
			//e := util.Database.Db.Where("`addr` = ? AND `key` = ? AND `pid`>0", "150.138.83.146:3456", v.ServerKey).First(&info).Error
			//if e == nil {
			//	fmt.Print(v.ServerKey, "新的", info.Id)
			//	e1 := util.Database.Db.Model(&model.AppUin{}).Where("`id` = ?", v.Id).Update("server_id", info.Id).Error
			//	if e1 == nil {
			//		fmt.Println(v.Uin, "修改成功")
			//	}
			//}
		}
		//time.Sleep(10 * time.Millisecond)
	}
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
		lock := constant.LockUinJob
		//查询当前时间需要代挂的QQ
		for {
			if util.Cache.Lock(lock, constant.LockTimeUinJob) {
				list, err := uin2.GetIns().SelectAllUinJob(1000, minute, str.NewTimeTool().GetZeroTime(time.Now()).Format("2006-01-02 15:04:05"))
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
						var speedUin *model.AppUin
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
								if v1.Status == uin2.Normal && v1.UinStatus == uin2.Normal && speedUin == nil && (v1.PlayerDay > 0 || v1.SignDay > 0) {
									speedUin = v1
								}
							}
						}
						//发布到消息队列执行
						if speedUin != nil && seeded[speedUin.Uin] == false {
							//避免同一查询有重复
							seeded[speedUin.Uin] = true
							//下线
							//server, e := server2.GetIns().FindByIdWithCache(speedUin.ServerId, false)
							//if e == nil {
							//	s := uin_tool.NewUinTool(server.Addr, server.Proxy)
							//	go s.MobileOffline(speedUin.Uin)
							//}
							//发布至消息队列
							producer.NewUinJobProducer().Build(speedUin).Send()
							log.Println("发送挂机任务", speedUin.Uin)
						}
					}
					//更改等级挂机天数
					if len(idsPlayer) > 0 {
						_ = uin2.GetIns().UpdateByIds(idsPlayer, map[string]interface{}{
							"player_day": gorm.Expr("`player_day`-1"),
						})
					}
					//更改扩展挂机天数
					if len(idsSign) > 0 {
						_ = uin2.GetIns().UpdateByIds(idsSign, map[string]interface{}{
							"sign_day": gorm.Expr("`sign_day`-1"),
						})
					}
					//更改代挂状态
					if len(ids) > 0 {
						_ = uin2.GetIns().UpdateByIds(ids, map[string]interface{}{
							"run_time": str.DateTime(time.Now()),
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
				list, err := uin2.GetIns().SelectUinCheckJob(1000, start, end)
				if err == nil && len(list) > 0 {
					var ids []int //需要更改得ids
					//避免重复挂机
					seeded := make(map[string]bool)
					for _, v := range list {
						var speedUin *model.AppUin
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
							log.Println("发送补挂任务", speedUin.Uin)
						}
					}
					//更改代挂状态
					if len(ids) > 0 {
						_ = uin2.GetIns().UpdateByIds(ids, map[string]interface{}{
							"check_time": str.DateTime(str.NewTimeTool().GetZeroTime(time.Now())),
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
			return
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

// health 健康检测
func health() {
	g := &sync.WaitGroup{}
	fun := func(addr string) {
		g.Add(1)
		defer g.Add(-1)
		s := uin_tool.NewUinTool(addr, "")
		cpu := s.GetCpu()
		if cpu <= 0 || cpu >= 80 || !s.CheckHealth() {
			//阻塞全局请求
			server2.GetIns().ServerLock("server_" + addr)
		} else {
			//解锁全局请求
			server2.GetIns().ServerUnLock("server_" + addr)
		}
	}
	for true {
		var adders []string
		err := util.Database.Db.Model(&model.AppUinServer{}).Where("`pid` > 0").Group("addr").Select("addr").Find(&adders).Error
		if err == nil {
			for _, v := range adders {
				go fun(v)
			}
		}
		g.Wait()
		time.Sleep(2 * time.Second)
	}
}
