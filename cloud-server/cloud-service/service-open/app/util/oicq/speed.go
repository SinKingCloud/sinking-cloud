package oicq

import (
	mobile "github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	pc "github.com/SinKingCloud/sinking-oicq/binary-pc/client"
	"log"
	"sync"
	"time"
)

// ActiveProject 执行项目
func (o *Oicq) ActiveProject(jobs []Job, cookies map[Job]string) (res map[Job]mobile.SpeedCode, resCookie map[Job]string) {
	//开始执行项目
	if cookies == nil {
		cookies = make(map[Job]string)
	}
	if res == nil {
		res = map[Job]mobile.SpeedCode{}
	}
	for _, v := range jobs {
		status := mobile.SpeedError
		cookie := ""
		switch v {
		case Ws: //微视加速
			_, cookie, status = o.Mobile.Speed.WeiShi(cookies[v])
			if status == mobile.SpeedCookieExpire {
				_, cookie, status = o.Mobile.Speed.WeiShi("")
			}
			break
		case Word: //小世界加速
			_, status = o.Mobile.Speed.FeedCloud()
			break
		case Game: //游戏加速
			_, cookie, status = o.Mobile.Speed.Game(cookies[v])
			if status == mobile.SpeedCookieExpire {
				_, cookie, status = o.Mobile.Speed.WeiShi("")
			}
			break
		case MobileClock: //打卡加速
			_, status = o.Mobile.Speed.DailySign()
			if status == mobile.SpeedCookieExpire {
				_, status = o.Mobile.Speed.DailySign()
			}
			break
		case ZoneClock: //说说加速
			_, status = o.Mobile.Speed.QZone()
			if status == mobile.SpeedCookieExpire {
				_, status = o.Mobile.Speed.QZone()
			}
			break
		case Friend: //加好友加速
			_, status = o.Mobile.Speed.Friend()
			break
		case Pc: //电脑加速
			if o.Pc.Online.Load() || o.MacPc.Online.Load() {
				status = mobile.SpeedSuccess
				break
			}
			o.Pc.Release()
			e, url := o.Pc.FetchQrCode()
			if e == nil {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					for i1 := 0; i1 < 30; i1++ {
						e2, state := o.Pc.FetchQrCodeStatus()
						if e2 != nil {
							break
						}
						if state == pc.ScanCodeSuccess {
							//登陆PC
							e3, l := o.Pc.Login(true)
							if e3 != nil {
								break
							}
							o.PcLoginTime = time.Now().Unix()
							if l.Status == pc.LoginSuccess {
								//关闭数据接收
								o.Pc.NetLoop(false)
								o.Pc.ForceOnline(true)  //强制在线,不间断发送心跳
								o.Pc.Online.Store(true) //设置为在线状态
								status = mobile.SpeedSuccess
							}
							break
						} else if state == pc.ScanCodeError {
							break
						}
						time.Sleep(1000 * time.Millisecond)
					}
					wg.Done()
				}(wg)
				//执行扫码操作
				time.Sleep(500 * time.Millisecond)
				_, e = o.Mobile.QrCodeAllow(utils.GetBetween(url, "p?k=", "&f="))
				if e == nil {
					//等待登陆结果
					wg.Wait()
					time.Sleep(500 * time.Millisecond)
					if o.Pc.Online.Load() {
						status = mobile.SpeedSuccess
					}
				} else {
					log.Println(o.Mobile.Uin, "pc授权登陆", e)
				}
			} else {
				log.Println(o.Mobile.Uin, "pc获取二维码失败", e, "开始使用mac协议")
			}
			if status == mobile.SpeedError {
				//mac协议挂机
				r := false
				o.MacPc.Release() //关闭连接
				e1 := o.MacPc.TokenLogin(o.LoginToken)
				if e1 == nil {
					log.Println(o.Mobile.Uin, "mac token登陆成功")
					r = true
				} else {
					log.Println(o.Mobile.Uin, "mac token登陆失败", e1, "开始尝试密码登陆")
					login, e2 := o.MacPc.Login()
					if e2 == nil && login.Success {
						log.Println(o.Mobile.Uin, "mac 密码登陆成功")
						r = true
					} else {
						log.Println(o.Mobile.Uin, "mac 密码登陆失败", e2, login)
					}
				}
				if r {
					o.MacPcLoginTime = time.Now().Unix()
					o.MacPc.NetLoop(false)     //关闭数据包接收
					o.MacPc.ForceOnline(true)  //强制在线,不间断发送心跳
					o.MacPc.Online.Store(true) //设置为在线状态
					status = mobile.SpeedSuccess
				} else {
					status = mobile.SpeedError
				}
			}
			break
		}
		res[v] = status
		if cookie != "" {
			cookies[v] = cookie
		}
	}
	return res, cookies
}
