package uin_tool

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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
			Timeout: 60 * time.Second,
		}
	}
	return httpClient
}

// getClient 获取http client(自定义超时)
func getClientWithTimout(timeout int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

// sendRequest 发送请求
func (s *speed) sendRequest(url string, data interface{}, timeout int) (*response, error) {
	url = "http://" + s.addr + url
	req, err := http.NewRequest("POST", url, strings.NewReader(s.toJson(data)))
	if err != nil {
		log.Println("错误：请求地址", url, "请求参数", s.toJson(data))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	var client *http.Client
	if timeout < 0 {
		client = getClient()
	} else {
		client = getClientWithTimout(timeout)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("错误：请求地址", url, "请求参数", s.toJson(data))
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
		log.Println("错误：请求地址", url, "请求参数", s.toJson(data), "结果", string(body))
		return nil, err
	}
	var res *response
	err = json.Unmarshal(body, &res)
	log.Println("请求地址", url, "请求参数", s.toJson(data), "结果", string(body))
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

// SpeedStatus 任务状态
type SpeedStatus struct {
	QzoneVisitor        string `json:"QzoneVisitor"`
	SVIPStar            string `json:"SVIPStar"`
	WeishiVideoview     string `json:"WeishiVideoview"`
	IAddFriend          string `json:"iAddFriend"`
	IBaseDays           string `json:"iBaseDays"`
	IBigClubGrowth      string `json:"iBigClubGrowth"`
	IBigClubLevel       string `json:"iBigClubLevel"`
	IBigClubSpeed       string `json:"iBigClubSpeed"`
	IBigClubVipFlag     string `json:"iBigClubVipFlag"`
	ICmshowActive       string `json:"iCmshowActive"`
	IContinueLogin      string `json:"iContinueLogin"`
	ICostMs             string `json:"iCostMs"`
	IDailySign          string `json:"iDailySign"`
	IMaxLvlRealDays     string `json:"iMaxLvlRealDays"`
	IMaxLvlTotalDays    string `json:"iMaxLvlTotalDays"`
	IMedal              string `json:"iMedal"`
	IMobileGameOnline   string `json:"iMobileGameOnline"`
	IMobileQQOnline     string `json:"iMobileQQOnline"`
	IMobileQQOnlineTime string `json:"iMobileQQOnlineTime"`
	IMqing              string `json:"iMqing"`
	INextLevelDay       string `json:"iNextLevelDay"`
	INoHideOnline       string `json:"iNoHideOnline"`
	INoHideOnlineTime   string `json:"iNoHideOnlineTime"`
	IPCQQOnline         string `json:"iPCQQOnline"`
	IPCQQOnlineTime     string `json:"iPCQQOnlineTime"`
	IPCSafeOnline       string `json:"iPCSafeOnline"`
	IQQLevel            string `json:"iQQLevel"`
	IQQSportStep        string `json:"iQQSportStep"`
	IQzoneState         string `json:"iQzoneState"`
	IRealDays           string `json:"iRealDays"`
	ISVip               string `json:"iSVip"`
	ISpeedType          string `json:"iSpeedType"`
	ISqq                string `json:"iSqq"`
	ISqqLevel           string `json:"iSqqLevel"`
	ISqqSpeedRate       string `json:"iSqqSpeedRate"`
	ISvrDays            string `json:"iSvrDays"`
	ITotalActiveDay     string `json:"iTotalActiveDay"`
	ITotalDays          string `json:"iTotalDays"`
	IVip                string `json:"iVip"`
	IVipLevel           string `json:"iVipLevel"`
	IVipSpeedRate       string `json:"iVipSpeedRate"`
	IYearBigClubFlag    string `json:"iYearBigClubFlag"`
	IYearVip            string `json:"iYearVip"`
	SFaceUrl            string `json:"sFaceUrl"`
	SNickName           string `json:"sNickName"`
	SpeedStar           string `json:"speedStar"`
}
