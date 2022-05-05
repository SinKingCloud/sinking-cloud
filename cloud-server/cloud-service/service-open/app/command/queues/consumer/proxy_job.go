package consumer

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/server"
	"io"
	"net/http"
	"net/url"
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
			_, e := l.test(s.Proxy, 10)
			status := 0
			if e != nil {
				status = 1
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
