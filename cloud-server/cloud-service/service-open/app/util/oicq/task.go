package oicq

import (
	"errors"
	mobile "github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	pc "github.com/SinKingCloud/sinking-oicq/binary-pc/client"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	oicqMap         = sync.Map{} // map[int64]*Oicq
	oicqMapLock     = sync.RWMutex{}
	oicqMapTemp     = sync.Map{} // map[int64]*Oicq
	oicqMapTempLock = sync.RWMutex{}
	heart           = atomic.Bool{}
	heartPc         = make(chan *Oicq, 100000)
	releasePc       = make(chan *Oicq, 100000)
	reConnPc        = make(chan *Oicq, 100000)
)

// GetOicqProto 获取协议类型
func GetOicqProto(proto int) mobile.ClientProtocol {
	protos := mobile.AndroidPad
	if proto == 0 {
		protos = mobile.IPad
	} else if proto == 1 {
		protos = mobile.AndroidPhone
	} else if proto == 2 {
		protos = mobile.AndroidPad
	} else if proto == 3 {
		protos = mobile.MacOS
	}
	return protos
}

// GetOicqProtoCode 获取协议类型标识
func GetOicqProtoCode(proto mobile.ClientProtocol) int {
	if proto == mobile.IPad {
		return 0
	} else if proto == mobile.AndroidPhone {
		return 1
	} else if proto == mobile.AndroidPad {
		return 2
	} else {
		return 3
	}
}

// GetOicqInstance 获取QQ操作实例
func GetOicqInstance(uin int64, password [16]byte, proxy *Proxy, proto int, device int, loginToken []byte, temp bool) *Oicq {
	var data any
	var ok bool
	if temp {
		data, ok = oicqMapTemp.Load(uin)
	} else {
		data, ok = oicqMap.Load(uin)
	}
	if !ok || data == nil {
		oicq := NewOicq(uin, password, proxy, device, loginToken, GetOicqProto(proto))
		if temp {
			oicqMapTemp.Store(uin, oicq)
		} else {
			oicqMap.Store(uin, oicq)
		}
		return oicq
	} else {
		oicq := data.(*Oicq)
		//检查是否更改代理
		if proxy != nil &&
			oicq.Mobile.Proxy != nil &&
			(proxy.Addr != oicq.Mobile.Proxy.Addr ||
				proxy.Port != oicq.Mobile.Proxy.Port ||
				proxy.User != oicq.Mobile.Proxy.User ||
				proxy.Pwd != oicq.Mobile.Proxy.Pwd) {
			oicq.Mobile.Disconnect()
			oicq.Pc.Disconnect()
			oicq.Mobile.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
			oicq.Pc.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
		}
		//检查是否更改密码
		if utils.Bin2HexTo(oicq.Mobile.PasswordMd5[:]) != utils.Bin2HexTo(password[:]) {
			oicq.Mobile.PasswordMd5 = password
		}
		return oicq
	}
}

// DelOicqInstance 删除QQ操作实例
func DelOicqInstance(uin int64) {
	oicqMapLock.RLock()
	defer oicqMapLock.RUnlock()
	data, ok := oicqMap.Load(uin)
	if ok {
		data.(*Oicq).Mobile.Release()
		data.(*Oicq).Mobile.Disconnect()
		data.(*Oicq).Pc.Release()
		data.(*Oicq).Pc.Disconnect()
		oicqMap.Store(uin, nil)
	}
}

// DelTempOicqInstance 删除QQ操作实例
func DelTempOicqInstance(uin int64) {
	oicqMapTempLock.RLock()
	defer oicqMapTempLock.RLock()
	data, ok := oicqMapTemp.Load(uin)
	if ok {
		data.(*Oicq).Mobile.Release()
		data.(*Oicq).Mobile.Disconnect()
		data.(*Oicq).Pc.Release()
		data.(*Oicq).Pc.Disconnect()
		oicqMapTemp.Delete(uin)
	}
}

// startDoHeartbeat 开启自动心跳
func startDoHeartbeat() {
	if !heart.Load() {
		heart.Store(true)
		go func() {
			for i := 0; i < 50; i++ {
				//下线
				go func(release chan *Oicq, id int) {
					for task := range release {
						isContinue := false
						utils.Try(func() {
							//判断是否加速完成,未加速完成重新挂机
							if task != nil && task.Mobile != nil {
								e, status := task.Mobile.Speed.Info()
								if e == nil && status != nil {
									if status.IPCQQOnline != "1" {
										task.LoginTime = time.Now().Unix()
										if task.TokenLogin() == nil {
											task.Pc.NetLoop(true)
											task.Pc.ForceOnline(false)
											var e2 error
											url := ""
											for i1 := 0; i1 < 3; i1++ {
												e2, url = task.Pc.FetchQrCode()
												if e == nil && strings.Contains(url, "&f=") {
													break
												}
											}
											if e2 == nil && len(url) > 32 {
												wg := &sync.WaitGroup{}
												wg.Add(1)
												go func(wg *sync.WaitGroup, task *Oicq) {
													succ := false
													for i1 := 0; i1 < 30; i1++ {
														e3, state := task.Pc.FetchQrCodeStatus()
														if e3 != nil {
															break
														}
														if state == pc.ScanCodeSuccess {
															//登陆PC
															e4, l := task.Pc.Login(true)
															if e4 != nil {
																break
															}
															task.PcLoginTime = time.Now().Unix()
															if l.Status == pc.LoginSuccess {
																succ = true
																//关闭数据接收
																task.Pc.NetLoop(false)
																task.Pc.Online.Store(true) //设置为在线状态
																task.Pc.ForceOnline(true)  //强制在线,不间断发送心跳
															}
															break
														} else if state == pc.ScanCodeError {
															break
														}
														time.Sleep(1000 * time.Millisecond)
													}
													task.Pc.NetLoop(false)
													if !succ {
														task.Pc.Release()
														task.Pc.Disconnect()
													}
													task.Mobile.Release()
													task.Mobile.Disconnect()
													wg.Done()
												}(wg, task)
												//执行扫码操作
												time.Sleep(500 * time.Millisecond)
												_, e = task.Mobile.QrCodeAllow(utils.GetBetween(url, "p?k=", "&f="))
												//等待登陆结果
												if e == nil {
													wg.Wait()
													time.Sleep(500 * time.Millisecond)
													if task.Pc.Online.Load() {
														log.Println(task.Mobile.Uin, "pc未加速重新登陆成功[id:", id, "]")
													} else {
														log.Println(task.Mobile.Uin, "pc未加速重新登陆失败[id:", id, "]")
														e = errors.New("pc未加速重新登陆失败")
													}
												} else {
													log.Println(task.Mobile.Uin, "pc未加速重新授权登陆失败[id:", id, "]", e)
													e = errors.New("pc未加速重新授权登陆失败")
												}
												task.Pc.NetLoop(false)
												task.Mobile.Release()
												task.Mobile.Disconnect()
												if e == nil {
													isContinue = true
													return
												}
											}
										}
										task.Mobile.Release()
										task.Mobile.Disconnect()
									}
								}
							} else {
								isContinue = true
							}
						}, func(i interface{}) {
							log.Println("pc下线出错[id:", id, "]", i)
							isContinue = false
						})
						if isContinue {
							continue
						}
						task.Pc.NetLoop(false)
						task.Pc.ForceOnline(false)
						task.Pc.Release()
						task.Pc.Disconnect()
						task.Pc.Online.Store(false)
						log.Println("pc下线[id:", id, "]", task.Pc.Uin, "登陆时间", task.PcLoginTime)
						task.PcLoginTime = 0
						task.PcHeartTime = 0
					}
				}(releasePc, i)
				//重连
				go func(tasks chan *Oicq, id int) {
					for task := range tasks {
						utils.Try(func() {
							if task == nil || task.Mobile == nil || task.Pc == nil {
								return
							}
							if task.PcReConnTime > 5 {
								log.Println("连接超限[id:", id, "]", task.Pc.Uin, "重试次数", task.PcReConnTime, "开始发送至下线队列")
								task.PcReConnTime = 0
								task.PcLoginTime = 0
								releasePc <- task
								return
							}
							r := false
							task.Pc.NetLoop(true)
							task.Pc.ForceOnline(false)
							for i1 := 0; i1 < 3; i1++ {
								if task.Pc.QuickReconnect() == nil && task.Pc.DoHeartbeat(true) == nil {
									r = true
									break
								} else {
									time.Sleep(500 * time.Millisecond)
								}
							}
							//重连成功则重新设置在线状态
							task.Pc.NetLoop(false)
							if r {
								log.Println("pc重连成功[id:", id, "]", task.Pc.Uin)
								task.Pc.Online.Store(true)
								task.Pc.ForceOnline(true)
								task.PcReConnTime++
							} else {
								log.Println("pc重连失败[id:", id, "]", task.Pc.Uin)
								task.Pc.Online.Store(false)
								task.Pc.ForceOnline(false)
								task.PcHeartTime = 0
								task.PcLoginTime = 0
							}
						}, func(i interface{}) {
							log.Println("pc重连出错[id:", id, "]", i)
						})
					}
				}(reConnPc, i)
				//心跳
				go func(tasks chan *Oicq, id int) {
					for task := range tasks {
						utils.Try(func() {
							//发送pc心跳
							if task != nil && task.Pc != nil && task.PcLoginTime > 0 && task.Pc.Online.Load() {
								//pc登陆满3小时下线
								if task.PcLoginTime > 0 && task.PcLoginTime+3600*3 < time.Now().Unix() {
									releasePc <- task
									task.PcLoginTime = time.Now().Unix()
								}
								//执行心跳
								e := task.Pc.DoHeartbeat(false)
								//心跳失败，重连(最大重试3次)
								if e != nil {
									log.Println("pc连接丢失[id:", id, "]", task.Pc.Uin, e, "开始重新连接")
									task.Pc.ForceOnline(false)
									task.Pc.Online.Store(false)
									reConnPc <- task
								} else {
									log.Println("发送心跳[id:", id, "]", task.Pc.Uin, e, "上次心跳时间", task.PcHeartTime)
									task.Pc.NetLoop(false)
									task.PcHeartTime = time.Now().Unix()
									task.PcReConnTime = 0
								}
							}
						}, func(i interface{}) {
							log.Println("发送心跳出错[id:", id, "]", i)
						})
					}
				}(heartPc, i)
			}
		}()
		//心跳
		go func() {
			for true {
				utils.Try(func() {
					log.Println("开始发送心跳")
					//发送数据到通道
					num := 0
					oicqMapLock.RLock()
					oicqMap.Range(func(key, value any) bool {
						utils.Try(func() {
							if value != nil {
								heartPc <- value.(*Oicq)
							}
						}, func(i interface{}) {
							log.Println("获取挂机实例出错", i)
						})
						num++
						return true
					})
					oicqMapLock.RUnlock()
					num2 := 0
					oicqMapTempLock.RLock()
					oicqMapTemp.Range(func(key, value any) bool {
						num2++
						return true
					})
					oicqMapTempLock.RUnlock()
					log.Println("当前挂机实例数量", num, "临时实例数量", num2)
				}, func(i interface{}) {
					log.Println("发送心跳出错", i)
					oicqMapLock.RUnlock()
					oicqMapTempLock.RUnlock()
				})
				time.Sleep(1 * time.Minute) //心跳周期1m
			}
		}()
		//删除临时实例
		go func() {
			for true {
				utils.Try(func() {
					log.Println("开始删除实例")
					oicqMapLock.RLock()
					oicqMapTemp.Range(func(key, value any) bool {
						utils.Try(func() {
							if value != nil {
								ins := value.(*Oicq)
								if ins.LoginTime+10*60 < time.Now().Unix() {
									//超时未操作，删除实例
									ins.Mobile.Release()
									ins.Mobile.Disconnect()
									ins.Pc.Release()
									ins.Pc.Disconnect()
									oicqMapTemp.Delete(key)
									log.Println("删除临时实例", key)
								}
							}
						}, func(i interface{}) {
							log.Println("删除临时实例失败", key)
						})
						return true
					})
					oicqMapLock.RUnlock()
					oicqMapTempLock.RLock()
					oicqMap.Range(func(key, value any) bool {
						utils.Try(func() {
							if value != nil {
								ins := value.(*Oicq)
								if ins.LoginTime > 0 && ins.LoginTime+10*60 < time.Now().Unix() {
									ins.LoginTime = 0
									ins.Mobile.Release()
									ins.Mobile.Disconnect()
									log.Println("断开实例", key)
								}
							}
						}, func(i interface{}) {
							log.Println("删除挂机实例失败", key)
						})
						return true
					})
					oicqMapTempLock.RUnlock()
				}, func(i interface{}) {
					log.Println("删除实例出错", i)
					oicqMapLock.RUnlock()
					oicqMapTempLock.RUnlock()
				})
				time.Sleep(1 * time.Minute)
			}
		}()
	}
}

// init 初始化加载
func init() {
	startDoHeartbeat()
}
