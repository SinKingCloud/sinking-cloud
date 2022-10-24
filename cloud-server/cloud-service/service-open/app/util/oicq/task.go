package oicq

import (
	"errors"
	mobile "github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	oicqMap     = sync.Map{} // map[int64]*Oicq
	oicqMapTemp = sync.Map{} // map[int64]*Oicq
	heart       = atomic.Bool{}
)

// GetOicqInstance 获取QQ操作实例
func GetOicqInstance(uin int64, password [16]byte, device int, proxy *Proxy, loginToken []byte, temp bool) *Oicq {
	var data any
	var ok bool
	proto := mobile.IPad
	if temp {
		data, ok = oicqMapTemp.Load(uin)
	} else {
		data, ok = oicqMap.Load(uin)
	}
	if !ok {
		oicq := NewOicq(uin, password, device, proxy, loginToken, proto)
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
			oicq.Mobile.Release()
			oicq.Pc.Release()
			oicq.MacPc.Release()
			oicq.Mobile.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
			oicq.Pc.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
			oicq.MacPc.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
		}
		//检查是否更改密码
		if utils.Bin2HexTo(oicq.Mobile.PasswordMd5[:]) != utils.Bin2HexTo(password[:]) {
			oicq.Mobile.PasswordMd5 = password
			oicq.MacPc.PasswordMd5 = password
		}
		return oicq
	}
}

// DelOicqInstance 删除QQ操作实例
func DelOicqInstance(uin int64) {
	data, ok := oicqMap.Load(uin)
	if ok {
		data.(*Oicq).Mobile.Release()
		data.(*Oicq).Pc.Release()
		data.(*Oicq).MacPc.Release()
		oicqMap.Delete(uin)
	}
}

// DelTempOicqInstance 删除QQ操作实例
func DelTempOicqInstance(uin int64) {
	data, ok := oicqMapTemp.Load(uin)
	if ok {
		data.(*Oicq).Mobile.Release()
		data.(*Oicq).Pc.Release()
		data.(*Oicq).MacPc.Release()
		oicqMapTemp.Delete(uin)
	}
}

// GenToken 生成token
func GenToken(uin int64) ([]byte, error) {
	data, ok := oicqMap.Load(uin)
	if ok {
		return data.(*Oicq).Mobile.GenToken(), nil
	}
	return nil, errors.New("QQ不存在")
}

// startDoHeartbeat 开启自动心跳
func startDoHeartbeat() {
	if !heart.Load() {
		heart.Store(true)
		go func() {
			for true {
				//根据实例数量分配协程,每100个实例分配一个协程
				threadNum := 0
				oicqMap.Range(func(key, value any) bool {
					threadNum++
					return true
				})
				threadNum = threadNum / 100
				if threadNum == 0 {
					threadNum = 1
				}
				if threadNum >= 300 {
					threadNum = 300
				}
				//新建协程
				tasks := make(chan *Oicq)
				log.Println("协程开始，数量", threadNum)
				for i := 0; i < threadNum; i++ {
					go func(tasks chan *Oicq) {
						for task := range tasks {
							//检测mac pc
							if task.MacPcLoginTime > 0 || task.MacPc.Online.Load() {
								//pc登陆满3小时下线
								if task.MacPcLoginTime > 0 && task.MacPcLoginTime+3600*2.5 < time.Now().Unix() {
									task.MacPc.Release()
									task.MacPc.Online.Store(false)
									log.Println("mac pc下线", task.MacPc.Uin, "登陆时间", task.MacPcLoginTime)
									task.MacPcLoginTime = 0
									task.MacPcHeartTime = 0
								} else {
									e := task.MacPc.DoHeartbeat(false)
									//心跳失败，重连(最大重试3次)
									if e != nil {
										log.Println("mac连接丢失", task.MacPc.Uin, e, "开始重新连接")
										r := false
										for i1 := 0; i1 < 3; i1++ {
											if task.MacPc.QuickReconnect() == nil {
												r = true
												task.MacPc.NetLoop(false) //关闭数据监听
												break
											} else {
												time.Sleep(500 * time.Millisecond)
											}
										}
										//重连成功则重新设置在线状态
										if r {
											log.Println("mac重连成功", task.MacPc.Uin)
											task.MacPc.Online.Store(true)
										} else {
											log.Println("mac重连失败", task.MacPc.Uin)
											task.MacPc.Online.Store(false)
											task.MacPcHeartTime = 0
											task.MacPcLoginTime = 0
										}
									} else {
										log.Println("mac 心跳", task.MacPc.Uin, e, "上次心跳时间", task.MacPcHeartTime)
										task.MacPcHeartTime = time.Now().Unix()
									}
								}
							} else {
								task.MacPcHeartTime = 0
								task.MacPcLoginTime = 0
							}
							//发送pc心跳
							if task.PcLoginTime > 0 || task.Pc.Online.Load() {
								//pc登陆满3小时下线
								if task.PcLoginTime > 0 && task.PcLoginTime+3600*3 < time.Now().Unix() {
									task.Pc.Release()
									task.Pc.Online.Store(false)
									log.Println("pc下线", task.Pc.Uin, "登陆时间", task.PcLoginTime)
									task.PcLoginTime = 0
									task.PcHeartTime = 0
								} else {
									//执行心跳
									e := task.Pc.DoHeartbeat(false)
									//心跳失败，重连(最大重试3次)
									if e != nil {
										log.Println("pc连接丢失", task.Pc.Uin, e, "开始重新连接")
										r := false
										for i1 := 0; i1 < 3; i1++ {
											if task.Pc.QuickReconnect() == nil {
												r = true
												task.Pc.NetLoop(false) //关闭数据监听
												break
											} else {
												time.Sleep(500 * time.Millisecond)
											}
										}
										//重连成功则重新设置在线状态
										if r {
											log.Println("pc重连成功", task.Pc.Uin)
											task.Pc.Online.Store(true)
										} else {
											log.Println("pc重连失败", task.Pc.Uin, "开始切换mac协议")
											//继续使用mac协议挂机
											task.MacPc.Release()     //关闭连接
											task.MacPc.NetLoop(true) //打开数据包接收
											e1 := task.MacPc.TokenLogin(task.LoginToken)
											r2 := false
											if e1 == nil {
												log.Println(task.Mobile.Uin, "mac token 重连成功")
												r2 = true
											} else {
												log.Println(task.Mobile.Uin, "mac token 重连失败", e1)
											}
											if r2 {
												task.MacPcLoginTime = time.Now().Unix()
												task.MacPc.NetLoop(false)     //关闭数据包接收
												task.MacPc.ForceOnline(true)  //强制在线,不间断发送心跳
												task.MacPc.Online.Store(true) //设置为在线状态
											} else {
												task.Pc.Online.Store(false)
												task.PcHeartTime = 0
												task.PcLoginTime = 0
											}
										}
									} else {
										log.Println("发送心跳", task.Pc.Uin, e, "上次心跳时间", task.PcHeartTime)
										task.PcHeartTime = time.Now().Unix()
									}
								}
							} else {
								task.PcHeartTime = 0
								task.PcLoginTime = 0
							}
						}
						log.Println("协程关闭", time.Now().Unix())
					}(tasks)
				}
				//发送数据到通道
				oicqMap.Range(func(key, value any) bool {
					tasks <- value.(*Oicq)
					return true
				})
				//关闭通道
				close(tasks)
				time.Sleep(1 * time.Minute) //心跳周期1m
				num := 0
				oicqMap.Range(func(key, value any) bool {
					num++
					return true
				})
				num2 := 0
				oicqMapTemp.Range(func(key, value any) bool {
					num2++
					return true
				})
				log.Println("当前挂机实例数量", num, "临时实例数量", num2)
			}
		}()
	}
}

// init 初始化加载
func init() {
	startDoHeartbeat()
}
