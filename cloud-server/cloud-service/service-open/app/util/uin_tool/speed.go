package uin_tool

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"sync"
	"time"
)

type speed struct {
	addr  string //服务器地址
	proxy string //使用代理
}

// NewUinTool 实例化
func NewUinTool(addr string, proxy string) *speed {
	return &speed{addr: addr, proxy: proxy}
}

// Create 创建账号
func (s *speed) Create(uin string, pwd string, phone string, phoneArea string, proto int) bool {
	return s.add(uin, pwd, phone, phoneArea, proto)
}

// Login 登陆账号
func (s *speed) Login(uin string, pwd string, phone string, phoneArea string, proto int, code string) (int, error, string, string) {
	a, b, c, d := s.login(uin, code)
	if a == NotCreate {
		s.add(uin, pwd, phone, phoneArea, proto)
		a, b, c, d = s.login(uin, code)
	}
	if a == Success {
		go func(uin string) {
			time.Sleep(1 * time.Second)
			s.mobileOnline(uin)
		}(uin)
	}
	return a, b, c, d
}

// VerifyLogin 验证码登陆
func (s *speed) VerifyLogin(uin string, code string) (int, error, string, string) {
	a, b, c, d := s.login(uin, code)
	if a == Success {
		go func(uin string) {
			time.Sleep(3 * time.Second)
			s.pcForceSpeed(uin)
		}(uin)
	}
	return a, b, c, d
}

// Delete 删除QQ
func (s *speed) Delete(uin string) bool {
	return s.add(uin, "123456", "18888888888", "86", 0)
}

// MobileOffline 下线手机
func (s *speed) MobileOffline(uin string) bool {
	return s.mobileLower(uin)
}

// PcOffline 下线电脑
func (s *speed) PcOffline(uin string) bool {
	return s.pcLower(uin)
}

const (
	Mobile      = iota //0.手机QQ加速
	Pc                 //1.Pc加速
	Online             //2.非隐身加速
	Zone               //3.访客加速
	Medal              //4.勋章墙加速
	Ws                 //5.微视加速
	MicroShow          //6.厘米秀加速
	Safe               //7.管家加速
	Game               //8.游戏加速
	Sport              //9.运动加速
	MobileClock        //10.手Q打卡
	ZoneClock          //11.空间打卡
	Login              //12.连续登陆加速
	Friend             //13.加好友加速
	Word               //14.小世界加速

	Vip        //14.会员签到
	BigVip     //15.大会员签到
	Hz         //16.黄钻签到
	Lz         //17.蓝钻签到
	Ld         //18.大乐斗签到
	GameCenter //19.游戏中心签到
	Music      //20.音乐签到
)

// Status 获取挂机状态 状态,错误,是否异常,服务器是否异常
func (s *speed) Status(uin string) (*SpeedStatus, error, bool, bool) {
	status, e := s.getSpeedStatus(uin)
	if e != nil {
		return nil, e, false, false
	}
	return status, e, false, false
}

var (
	pcServers    = make(map[string]chan interface{}) //PC服务器列表
	pcServerLock = &sync.Mutex{}                     //PC并发锁
)

type speedJob struct {
	Uin       string
	Pwd       string
	Phone     string
	PhoneArea string
	Proto     int
	SpeedJobs []int
	IsFirst   bool
}

func (s *speed) speedByPc(uin *speedJob) {
	pcServerLock.Lock()
	if pcServers[s.addr] == nil {
		//构建其他项目消费协程
		g := &sync.WaitGroup{}
		g.Add(1)
		go (&job.Task{Producer: func(channel chan interface{}) {
			pcServers[s.addr] = channel
			g.Add(-1)
		}, Consumer: func(param interface{}) {
			if param != nil {
				//执行挂机任务
				sinking_web.Try(func() {
					//执行
					info := param.(*speedJob)
					res := make(map[int]bool)
					mobileLogin := false //是否登录手机
					pcLogin := false     //是否登录pc
					isBreak := false     //是否异常
					speedA := func(f func() bool) bool {
						if !mobileLogin {
							c, _, _, _ := s.login(uin.Uin, "")
							if c == NotCreate {
								if s.add(uin.Uin, uin.Pwd, uin.Phone, uin.PhoneArea, uin.Proto) {
									c, _, _, _ = s.login(uin.Uin, "")
								}
							} else if c == EnvError {
								c, _, _, _ = s.login(uin.Uin, "")
							}
							mobileLogin = c == Success
						}
						isBreak = !mobileLogin
						if mobileLogin {
							return f()
						}
						return false
					}
					speedP := func(f func() bool) bool {
						if !pcLogin {
							if !mobileLogin {
								speedA(func() bool {
									return mobileLogin
								})
							}
							if !mobileLogin {
								return false
							}
							if s.pcForceSpeed(uin.Uin) {
								pcLogin = true
							}
						}
						if pcLogin {
							return f()
						}
						return false
					}
					speed := func(jobs []int) {
						//加速项目
						for _, v := range jobs {
							switch v {
							case Pc:
								res[Pc] = speedP(func() bool {
									time.Sleep(3 * time.Second)
									return s.medalSpeed(uin.Uin)
								})
								break
							default:
								res[v] = true
								break
							}
						}
					}
					if !isBreak {
						speed(info.SpeedJobs)
					}
				}, func(i interface{}) {
					util.Log.Println(param, i)
				})
			}
		}, Thread: 5}).Run()
		g.Wait()
	}
	pcServerLock.Unlock()
	//按照服务分发至不同协程执行
	pcServers[s.addr] <- uin
}

func (s *speed) SpeedAll(uin *speedJob) (bool, map[int]bool) {
	res := make(map[int]bool)
	mobileLogin := false //是否登录手机
	pcLogin := false     //是否登录pc
	isBreak := false     //是否异常
	serverError := false //系统异常
	speedA := func(f func() bool) bool {
		if !mobileLogin {
			c, _, _, _ := s.login(uin.Uin, "")
			if c == EnvError {
				c, _, _, _ = s.login(uin.Uin, "")
			} else if c == Error {
				serverError = true
			}
			mobileLogin = c == Success
		}
		isBreak = !mobileLogin
		if mobileLogin {
			return f()
		}
		return false
	}
	speedP := func(f func() bool) bool {
		if !pcLogin {
			if !mobileLogin {
				speedA(func() bool {
					return mobileLogin
				})
			}
			if !mobileLogin {
				return false
			}
			if s.pcForceSpeed(uin.Uin) {
				pcLogin = true
			}
		}
		if pcLogin {
			return f()
		}
		return false
	}
	getTryJobs := func(jobs []int) (bool, []int) {
		status, err, isFail, _ := s.Status(uin.Uin)
		if isFail {
			return true, []int{}
		}
		if err != nil || status == nil {
			return false, []int{}
		}
		var speedJobs []int
		for _, v := range jobs {
			isCheck := false
			switch v {
			case MicroShow:
				isCheck = status.ICmshowActive != "1"
				break
			case Sport:
				isCheck = false
				if time.Now().Day() <= 8 {
					isCheck = status.IQQSportStep != "1"
				}
				break
			}
			if isCheck {
				speedJobs = append(speedJobs, v)
			}
		}
		return false, speedJobs
	}
	speed := func(jobs []int) (pcJobs []int) {
		isBreak = false //初始化
		//加速项目
		for _, v := range jobs {
			if isBreak {
				res[v] = false
				continue
			}
			switch v {
			case Mobile:
				res[Mobile] = speedA(func() bool {
					return s.mobileOnline(uin.Uin)
				})
				break
			case Pc:
				res[Pc] = speedP(func() bool {
					return pcLogin
				})
				break
			case Online:
				if !uin.IsFirst {
					res[Online] = speedP(func() bool {
						return pcLogin
					})
				}
				break
			case Zone:
				res[Zone] = s.zoneSpeed(uin.Uin)
				break
			case Ws:
				res[Ws] = speedA(func() bool {
					return s.wsSpeed(uin.Uin)
				})
				break
			case MicroShow:
				res[MicroShow] = speedA(func() bool {
					return s.microShowSpeed(uin.Uin)
				})
				break
			case Word:
				res[Word] = speedA(func() bool {
					return s.wordSpeed(uin.Uin)
				})
				break
			case Safe:
				res[Safe] = speedA(func() bool {
					return s.safeSpeed(uin.Uin)
				})
				break
			case Game:
				res[Game] = speedA(func() bool {
					return s.gameSpeed(uin.Uin)
				})
				break
			case Sport:
				res[Sport] = speedA(func() bool {
					r := s.sportSpeed(uin.Uin)
					if !r {
						r = s.sportSpeed(uin.Uin)
					}
					return r
				})
				break
			case MobileClock:
				res[MobileClock] = speedA(func() bool {
					return s.mobileClockSpeed(uin.Uin)
				})
				break
			case ZoneClock:
				res[ZoneClock] = speedA(func() bool {
					return s.zoneClockSpeed(uin.Uin)
				})
				break
			case Friend:
				res[Friend] = speedA(func() bool {
					return s.friendSpeed(uin.Uin)
				})
				break
			default:
				res[v] = true
				break
			}
		}
		return
	}
	pcJobs := speed(uin.SpeedJobs) //执行加速任务
	if uin.IsFirst {
		isFail, retryJobs := getTryJobs(uin.SpeedJobs)
		if !isFail && len(retryJobs) > 0 {
			pcJobs = speed(retryJobs) //重新执行秒加速项目
			isFail, retryJobs = getTryJobs(uin.SpeedJobs)
			if !isFail && len(retryJobs) > 0 {
				pcJobs = speed(retryJobs) //重新执行秒加速项目
			}
		}
	}
	if len(pcJobs) > 0 {
		s.speedByPc(&speedJob{
			Uin:       uin.Uin,
			Pwd:       uin.Pwd,
			Phone:     uin.Phone,
			Proto:     uin.Proto,
			SpeedJobs: uin.SpeedJobs,
			IsFirst:   uin.IsFirst,
		})
	}
	if serverError {
		return false, res
	}
	return isBreak, res
}

// Speed QQ项目加速
func (s *speed) Speed(uin string, pwd string, phone string, phoneArea string, proto int, speedJobs []int, isFirst bool) (bool, map[int]bool) {
	if server2.GetIns().ServerLock("speed_" + uin) {
		defer server2.GetIns().ServerUnLock("speed_" + uin)
	} else {
		return false, map[int]bool{}
	}
	return s.SpeedAll(&speedJob{
		Uin:       uin,
		Pwd:       pwd,
		Phone:     phone,
		PhoneArea: phoneArea,
		Proto:     proto,
		SpeedJobs: speedJobs,
		IsFirst:   isFirst,
	})
}
