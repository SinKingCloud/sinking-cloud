package consumer

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"io"
	"log"
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
			log.Println("判断代理", s)
		}
	})
	if err != nil {
		util.Log.Println(err)
		return
	}
}

// Http 判断Http代理是否可用
func (l *proxyJob) Http(types string, ip string, port int, timeOut int) (int, error) {
	proxyAddr := types + "://" + ip + ":" + strconv.Itoa(port)
	var testUrl string
	if types == "https" {
		testUrl = "https://icanhazip.com"
	} else {
		testUrl = "http://icanhazip.com"
	}
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
	speed := int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000) //ms
	if res.StatusCode != http.StatusOK {
		return 0, errors.New("http status code error")
	}
	return speed, nil
}
