package speed

import "errors"

type speed struct {
	addr string //服务器地址
}

// NewSpeed 实例化
func NewSpeed(addr string) *speed {
	return &speed{addr: addr}
}

// Login 登陆账号
func (s *speed) Login(uin string, pwd string) (int, error) {
	if s.add(uin, pwd) {
		return s.login(uin, "")
	}
	return Error, errors.New("初始化登录请求失败")
}

// CaptchaLogin 验证滑动验证码登陆
func (s *speed) CaptchaLogin(uin string, ticket string) (int, error) {
	return s.login(uin, ticket)
}

// SmsLogin 验证设备锁登陆
func (s *speed) SmsLogin(uin string, code string) (int, error) {
	return s.login(uin, code)
}

// Delete 删除QQ
func (s *speed) Delete(uin string) bool {
	return s.add(uin, "123456")
}

const (
	Mobile      = iota //手机QQ加速
	Pc                 //Pc加速
	Online             //非隐身加速
	Zone               //访客加速
	Medal              //勋章墙加速
	Ws                 //微视加速
	MicroShow          //厘米秀加速
	Safe               //管家加速
	Game               //游戏加速
	Sport              //运动加速
	MobileClock        //手Q打卡
	ZoneClock          //空间打卡
	Login              //连续登陆加速
	Add                //加好友加速
)

// Speed QQ项目加速
func (s *speed) Speed(uin string, pwd string, jobs []int) map[int]bool {
	res := make(map[int]bool)
	login, err := s.Login(uin, pwd)
	if err != nil || login != Success {
		return res
	}
	mobileLogin := false //是否登录手机
	pcLogin := false     //是否登录pc
	loginA := func() bool {
		c, err := s.login(uin, "")
		if err != nil {
			return false
		}
		mobileLogin = c == Success
		return mobileLogin
	}
	loginP := func() bool {
		if !mobileLogin {
			if !loginA() {
				return false
			}
		}
		if s.pcSpeed(uin) {
			pcLogin = true
		}
		return pcLogin
	}
	speedA := func(f func() bool) bool {
		if !mobileLogin {
			loginA()
		}
		if mobileLogin {
			return f()
		}
		return false
	}
	speedP := func(f func() bool) bool {
		if !pcLogin {
			loginP()
		}
		if pcLogin {
			return f()
		}
		return false
	}
	//加速项目
	for _, v := range jobs {
		switch v {
		case Mobile:
			res[Mobile] = loginA()
			break
		case Pc:
			res[Pc] = loginP()
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
		default:
			res[v] = true
			break
		}
	}
	return res
}
