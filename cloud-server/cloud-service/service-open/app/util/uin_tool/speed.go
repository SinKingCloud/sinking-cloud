package uin_tool

type speed struct {
	addr string //服务器地址
}

// NewUinTool 实例化
func NewUinTool(addr string) *speed {
	return &speed{addr: addr}
}

// Create 创建账号
func (s *speed) Create(uin string, pwd string) bool {
	return s.add(uin, pwd)
}

// Login 登陆账号
func (s *speed) Login(uin string, pwd string) (int, error, string) {
	a, b, c := s.login(uin, "")
	if a == NotCreate {
		if s.add(uin, pwd) {
			a, b, c = s.login(uin, "")
		}
	}
	return a, b, c
}

// CaptchaLogin 验证滑动验证码登陆
func (s *speed) CaptchaLogin(uin string, ticket string) (int, error) {
	a, b, _ := s.login(uin, ticket)
	return a, b
}

// SmsLogin 验证设备锁登陆
func (s *speed) SmsLogin(uin string, code string) (int, error) {
	a, b, _ := s.login(uin, code)
	return a, b
}

// Delete 删除QQ
func (s *speed) Delete(uin string) bool {
	return s.add(uin, "123456")
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

// Speed QQ项目加速
func (s *speed) Speed(uin string, pwd string, jobs []int) map[int]bool {
	res := make(map[int]bool)
	mobileLogin := false //是否登录手机
	pcLogin := false     //是否登录pc
	speedA := func(f func() bool) bool {
		if !mobileLogin {
			if !s.add(uin, pwd) {
				return false
			}
			c, err, _ := s.login(uin, "")
			if err != nil {
				return false
			}
			mobileLogin = c == Success || err == nil
		}
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
			if s.pcSpeed(uin) {
				pcLogin = true
			}
			return pcLogin
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
	//加速项目
	for _, v := range jobs {
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
		case Online:
			res[Online] = speedP(func() bool {
				return pcLogin
			})
			break
		case Medal:
			res[Medal] = speedP(func() bool {
				return pcLogin
			})
			break
		case Zone:
			res[Zone] = speedP(func() bool {
				return s.zoneSpeed(uin)
			})
			break
		case Ws:
			res[Ws] = speedP(func() bool {
				return s.wsSpeed(uin)
			})
			break
		case MicroShow:
			res[MicroShow] = speedA(func() bool {
				return s.microShowSpeed(uin)
			})
			break
		case Safe:
			res[Safe] = speedP(func() bool {
				return s.safeSpeed(uin)
			})
			break
		case Game:
			res[Game] = speedA(func() bool {
				return s.gameSpeed(uin)
			})
			break
		case Sport:
			res[Sport] = speedA(func() bool {
				return s.sportSpeed(uin)
			})
			break
		case MobileClock:
			res[MobileClock] = speedA(func() bool {
				return s.mobileClockSpeed(uin)
			})
			break
		case ZoneClock:
			res[ZoneClock] = speedA(func() bool {
				return s.zoneClockSpeed(uin)
			})
			break
		case Login:
			res[Login] = speedA(func() bool {
				return mobileLogin
			})
			break
		case Friend:
			res[Friend] = speedA(func() bool {
				return s.friendSpeed(uin)
			})
			break
		default:
			res[v] = true
			break
		}
	}
	return res
}
