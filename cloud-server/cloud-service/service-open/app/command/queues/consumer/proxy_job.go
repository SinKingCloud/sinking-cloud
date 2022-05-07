package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/server"
	uin2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// uinJobIns 消费者单例
var proxyJobIns queue.MqConsumer

func NewProxyJob() *proxyJob {
	return &proxyJob{}
}

// proxyJob 结构
type proxyJob struct {
}

func (l *proxyJob) Listen() {
	if proxyJobIns == nil {
		proxyJobIns = sinking.GetRocketConsumerServer(constant.ProxyJobGroup, false)
	}
	err := proxyJobIns.ListenReceiveMsgDo(constant.ProxyJobTopic, "", func(mqMsg queue.MqMsg) {
		s := &model.AppUinServer{}
		err := json.Unmarshal(mqMsg.Body, s)
		if err == nil {
			//判断代理
			_, e := l.test(s.Proxy, 30)
			status := 0
			if e != nil {
				_, e = l.test(s.Proxy, 30)
				if e != nil {
					status = 1
					go l.change(s) //故障迁移
				}
			}
			_ = server2.GetIns().UpdateById(s.Id, map[string]interface{}{
				"check_time": str.DateTime(time.Now()),
				"status":     status,
			})
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

// change 故障迁移
func (l *proxyJob) change(server *model.AppUinServer) {
	key := constant.LockUinServerChange + strconv.Itoa(server.Id)
	for true {
		//迁移锁
		if util.Cache.Lock(key, constant.LockTimeUinServerChange) {
			for true {
				list, er := uin2.GetIns().SelectByServerId(200, server.Id)
				if er == nil && len(list) > 0 {
					//更改账号挂机节点
					for _, v := range list {
						//挂机节点负载均衡
						newServer, e := web_uin.GetServer(v.ServerKey)
						if e == nil && newServer != nil {
							_ = uin2.GetIns().UpdateById(v.Id, map[string]interface{}{
								"server_id": newServer.Id,
							})
						} else {
							break
						}
					}
				} else {
					break
				}
			}
			util.Cache.UnLock(key)
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

// test 测试代理是否可用
func (l *proxyJob) test(addr string, timeOut int) (int, error) {
	proxyAddr := "http://" + addr
	testUrl := "http://icanhazip.com"
	proxy, err := url.Parse(proxyAddr)
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	httpClient := &http.Client{
		Timeout:   time.Second * time.Duration(timeOut),
		Transport: netTransport,
	}
	begin := time.Now()
	res, err := httpClient.Get(testUrl)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	return int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000), nil
}
