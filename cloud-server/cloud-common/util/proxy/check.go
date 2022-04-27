package proxy

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type check struct {
}

var Check check

// Http 判断Http代理是否可用
func (check) Http(types string, ip string, port int, timeOut int) (int, error) {
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
