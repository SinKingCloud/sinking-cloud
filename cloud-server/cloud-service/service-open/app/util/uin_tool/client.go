package uin_tool

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//请求client
var httpClient *http.Client

// getClient 获取http client
func getClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return httpClient
}

// sendRequest 发送请求
func (s *speed) sendRequest(url string, data interface{}) (*response, error) {
	url = "http://" + s.addr + url
	req, _ := http.NewRequest("POST", url, strings.NewReader(s.toJson(data)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res *response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// response 响应结构
type response struct {
	Type  string `json:"Type"`
	State string `json:"State"`
	Uin   int64  `json:"Uin"`
}
