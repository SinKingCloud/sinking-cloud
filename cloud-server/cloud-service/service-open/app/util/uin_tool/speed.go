package uin_tool

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"sync"
	"time"
)

type speed struct {
	addr string //服务器地址
	area string //区域
}

// NewUinTool 实例化
func NewUinTool(addr string, area string) *speed {
	return &speed{addr: addr, area: area}
}

// Create 创建账号
func (s *speed) Create(uin string, pwd string, phone string, proto int) bool {
	return s.add(uin, pwd, phone, proto)
}

// Login 登陆账号
func (s *speed) Login(uin string, pwd string, phone string, proto int, code string) (int, error, string, string) {
	a, b, c, d := s.login(uin, code)
	if a == NotCreate {
		s.add(uin, pwd, phone, proto)
		return s.login(uin, code)
	}
	return a, b, c, d
}

// VerifyLogin 验证码登陆
func (s *speed) VerifyLogin(uin string, code string) (int, error, string, string) {
	return s.login(uin, code)
}

// Delete 删除QQ
func (s *speed) Delete(uin string) bool {
	return s.add(uin, "123456", "18888888888", 0)
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
)

// Status 获取挂机状态 状态,错误,是否异常
func (s *speed) Status(uin string) (*SpeedStatus, error, bool) {
	status, e := s.getSpeedStatus(uin)
	if e != nil {
		return nil, e, false
	}
	if status == nil {
		code, _, _, _ := s.login(uin, "")
		if code == Success {
			status, e = s.getSpeedStatus(uin)
			if e != nil || status == nil {
				s.MobileOffline(uin)
				code, _, _, _ = s.login(uin, "")
				if code == Success {
					status, e = s.getSpeedStatus(uin)
				} else {
					return status, e, true
				}
			}
		} else {
			return status, e, true
		}
	}
	return status, e, false
}

var (
	pcServers    = make(map[string]chan interface{}) //PC服务器列表
	pcServerLock = &sync.Mutex{}                     //PC并发锁
)

type speedJob struct {
	Uin     string
	Pwd     string
	Phone   string
	Proto   int
	Jobs    []int
	IsFirst bool
}

func (s *speed) speedByPc(uin *speedJob) {
	if pcServers[s.addr] == nil {
		//构建其他项目消费协程
		pcServerLock.Lock()
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
					pcLogin := false //是否登录pc
					speedP := func(f func() bool) bool {
						if !pcLogin {
							if s.pcSpeed(info.Uin) {
								pcLogin = true
							}
						}
						if pcLogin {
							r := f()
							if !r {
								r = f()
							}
							return r
						}
						return false
					}
					speed := func(jobs []int) {
						//加速项目
						for _, v := range jobs {
							switch v {
							case Pc:
								res[Pc] = speedP(func() bool {
									return pcLogin
								})
								break
							case Medal:
								res[Medal] = speedP(func() bool {
									time.Sleep(5 * time.Second)
									return s.medalSpeed(info.Uin)
								})
								break
							case Zone:
								res[Zone] = s.zoneSpeed(info.Uin)
								break
							default:
								res[v] = true
								break
							}
						}
					}
					speed(info.Jobs)
				}, func(i interface{}) {
					util.Log.Println(param, i)
				})
			}
		}, Thread: 10}).Run()
		g.Wait()
		pcServerLock.Unlock()
	}
	//按照服务分发至不同协程执行
	pcServers[s.addr] <- uin
}

func (s *speed) SpeedWithPcAsync(uin *speedJob) (bool, map[int]bool) {
	res := make(map[int]bool)
	mobileLogin := false //是否登录手机
	isBreak := false     //是否有必要继续执行
	speedA := func(f func() bool) bool {
		if !mobileLogin {
			c, _, _, _ := s.login(uin.Uin, "")
			if c == NotCreate {
				if s.add(uin.Uin, uin.Pwd, uin.Phone, uin.Proto) {
					c, _, _, _ = s.login(uin.Uin, "")
				}
			} else if c == EnvError {
				c, _, _, _ = s.login(uin.Uin, "")
			}
			mobileLogin = c == Success
		}
		isBreak = !mobileLogin
		if mobileLogin {
			r := f()
			if !r {
				r = f()
			}
			return r
		}
		return false
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
					return mobileLogin
				})
				break
			case Pc:
				res[Pc] = true
				pcJobs = append(pcJobs, Pc)
				break
			case Medal:
				res[Medal] = true
				pcJobs = append(pcJobs, Medal)
				break
			case Zone:
				res[Zone] = true
				pcJobs = append(pcJobs, Zone)
				break
			case Online:
				if !uin.IsFirst {
					res[Online] = speedA(func() bool {
						return mobileLogin
					})
				}
				break
			case Login:
				if !uin.IsFirst {
					res[Login] = speedA(func() bool {
						return mobileLogin
					})
				}
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
					return s.sportSpeed(uin.Uin)
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
	pcJobs := speed(uin.Jobs) //执行挂机
	if len(pcJobs) > 0 {
		s.speedByPc(&speedJob{
			Uin:     uin.Uin,
			Pwd:     uin.Pwd,
			Phone:   uin.Phone,
			Proto:   uin.Proto,
			Jobs:    uin.Jobs,
			IsFirst: uin.IsFirst,
		})
	}
	return isBreak, res
}

func (s *speed) SpeedAll(uin *speedJob) (bool, map[int]bool) {
	res := make(map[int]bool)
	mobileLogin := false //是否登录手机
	pcLogin := false     //是否登录pc
	isBreak := false     //是否有必要继续执行
	speedA := func(f func() bool) bool {
		if !mobileLogin {
			c, _, _, _ := s.login(uin.Uin, "")
			if c == NotCreate {
				if s.add(uin.Uin, uin.Pwd, uin.Phone, uin.Proto) {
					c, _, _, _ = s.login(uin.Uin, "")
				}
			} else if c == EnvError {
				c, _, _, _ = s.login(uin.Uin, "")
			}
			mobileLogin = c == Success
		}
		isBreak = !mobileLogin
		if mobileLogin {
			r := f()
			if !r {
				r = f()
			}
			return r
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
			if s.pcSpeed(uin.Uin) {
				pcLogin = true
			}
		}
		if pcLogin {
			r := f()
			if !r {
				r = f()
			}
			return r
		}
		return false
	}
	speed := func(jobs []int) {
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
					return mobileLogin
				})
				break
			case Pc:
				res[Pc] = speedP(func() bool {
					return pcLogin
				})
				break
			case Medal:
				res[Medal] = speedP(func() bool {
					time.Sleep(5 * time.Second)
					return s.medalSpeed(uin.Uin)
				})
				break
			case Zone:
				res[Zone] = s.zoneSpeed(uin.Uin)
				break
			case Online:
				if !uin.IsFirst {
					res[Online] = speedA(func() bool {
						return mobileLogin
					})
				}
				break
			case Login:
				if !uin.IsFirst {
					res[Login] = speedA(func() bool {
						return mobileLogin
					})
				}
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
					return s.sportSpeed(uin.Uin)
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
	}
	speed(uin.Jobs) //执行挂机
	return isBreak, res
}

// Speed QQ项目加速
func (s *speed) Speed(uin string, pwd string, phone string, proto int, jobs []int, isFirst bool) (bool, map[int]bool) {
	return s.SpeedAll(&speedJob{
		Uin:     uin,
		Pwd:     pwd,
		Phone:   phone,
		Proto:   proto,
		Jobs:    jobs,
		IsFirst: isFirst,
	})
}
