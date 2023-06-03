package oicq

import (
	"errors"
	mobile "github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	pc "github.com/SinKingCloud/sinking-oicq/binary-pc/client"
	"log"
	"runtime/debug"
	"strings"
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
		o.LoginTime = time.Now().Unix()
		status := mobile.SpeedError
		cookie := ""
		utils.Try(func() {
			switch v {
			case Vip: //会员签到
				_, status = o.Mobile.Speed.VipSign()
				break
			case BigVip: //大会员签到
				_, status = o.Mobile.Speed.BigVipSign()
				break
			case Lz: //黄钻签到
				_, status = o.Mobile.Speed.YellowSign()
				break
			case Ld: //乐斗签到
				_, status = o.Mobile.Speed.DaLeDouSign()
				break
			case WeiYun: //微云签到
				_, status = o.Mobile.Speed.WeiYunSign()
				break
			case Music: //音乐签到
				_, status = o.Mobile.Speed.QQMusicSign()
				break
			case GameCenter: //游戏中心签到
				_, status = o.Mobile.Speed.GameCenterSign()
				break
			case QQGame: //QQ游戏签到
				_, status = o.Mobile.Speed.QQGamSign()
				_, status = o.Mobile.Speed.DaoJuChengSign()
				_, status = o.Mobile.Speed.JingPinYeYouSign()
				break
			case Sweet: //情侣空间签到
				_, status = o.Mobile.Speed.SweetSign()
				break
			case Video: //腾讯视频
				_, status = o.Mobile.Speed.VideoSign()
				break
			case Ws: //微视加速
				_, cookie, status = o.Mobile.Speed.WeiShi(cookies[v])
				if status == mobile.SpeedCookieExpire {
					_, cookie, status = o.Mobile.Speed.WeiShi("")
				}
				break
			case Word: //小世界加速
				_, status = o.Mobile.Speed.WebFeedCloud()
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
				d := time.Now()
				d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
				if o.Pc.Online.Load() && o.Pc.LoginTime > 0 && o.PcLoginTime > d.Unix() {
					status = mobile.SpeedSuccess
					break
				}
				utils.Try(func() {
					o.Pc.NetLoop(true)
					o.Pc.ForceOnline(false)
					o.Pc.Release()
					o.Pc.Disconnect()
					var e error
					url := ""
					for i := 0; i < 3; i++ {
						e, url = o.Pc.FetchQrCode()
						if e == nil && strings.Contains(url, "&f=") {
							break
						}
					}
					if e == nil && len(url) > 32 {
						wg := &sync.WaitGroup{}
						wg.Add(1)
						go func(wg *sync.WaitGroup, o *Oicq, status mobile.SpeedCode) {
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
										o.Pc.Online.Store(true) //设置为在线状态
										o.Pc.ForceOnline(true)  //强制在线,不间断发送心跳
										status = mobile.SpeedSuccess
									}
									break
								} else if state == pc.ScanCodeError {
									break
								}
								time.Sleep(1000 * time.Millisecond)
							}
							o.Pc.NetLoop(false)
							if status != mobile.SpeedSuccess {
								o.Pc.Release()
								o.Pc.Disconnect()
							}
							wg.Done()
						}(wg, o, status)
						//执行扫码操作
						time.Sleep(500 * time.Millisecond)
						_, e = o.Mobile.QrCodeAllow(utils.GetBetween(url, "p?k=", "&f="))
						if e == nil {
							//等待登陆结果
							wg.Wait()
							time.Sleep(500 * time.Millisecond)
							if o.Pc.Online.Load() {
								status = mobile.SpeedSuccess
							} else {
								log.Println(o.Mobile.Uin, "pc登陆失败")
								e = errors.New("登陆失败")
							}
						} else {
							log.Println(o.Mobile.Uin, "pc授权登陆失败", e)
							e = errors.New("pc授权登陆失败")
						}
					} else {
						log.Println(o.Mobile.Uin, "pc获取二维码失败", e, url)
						e = errors.New("pc获取二维码失败")
					}
					if e != nil {
						o.Pc.Release()
						o.Pc.Disconnect()
					} else {
						o.Pc.NetLoop(false)
					}
				}, func(i interface{}) {
					log.Println(o.Mobile.Uin, "pc登陆错误", i)
					o.Pc.Release()
					o.Pc.Disconnect()
				})
				break
			}
		}, func(i interface{}) {
			log.Println("执行挂机项目失败", o.Mobile.Uin, v, i, string(debug.Stack()))
		})
		res[v] = status
		if cookie != "" {
			cookies[v] = cookie
		}
	}
	return res, cookies
}
