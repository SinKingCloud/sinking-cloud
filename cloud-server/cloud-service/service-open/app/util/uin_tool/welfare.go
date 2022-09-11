package uin_tool

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var authWelfareIns *AutoWelfare

// GetAutoWelfareIns 获取单例
func GetAutoWelfareIns() *AutoWelfare {
	return authWelfareIns
}

type AutoWelfare struct {
	Leader  *WelfareUser //队长
	UserOne *WelfareUser //队员1
	UserTwo *WelfareUser //队员2
}

type WelfareUser struct {
	Uin   string
	Skey  string
	Pskey string
}

// NewAutoWelfare 实例化队伍
func NewAutoWelfare(Leader *WelfareUser, UserOne *WelfareUser, UserTwo *WelfareUser) *AutoWelfare {
	return &AutoWelfare{
		Leader:  Leader,
		UserOne: UserOne,
		UserTwo: UserTwo,
	}
}

// createTeam 创建队伍
func (a *AutoWelfare) createTeam(uin string, skey string, pskey string) string {
	url := "https://ti.qq.com/act/99welfare2022/trpc/CharityDayProxy/CreateTeam?bkn=1378514044"
	body := "{\"uin\":" + uin + "}"
	ret := a.request(url, body, a.getCookieString(uin, skey, pskey))
	if ret == "" {
		return ""
	}
	var res *createTeamResponse
	err := json.Unmarshal([]byte(ret), &res)
	if err != nil || res == nil || res.Response.TeamInfo.TeamId == "" {
		return ""
	}
	return res.Response.TeamInfo.TeamId
}

// joinTeam 加入队伍
func (a *AutoWelfare) joinTeam(uin string, skey string, pskey string, team string) bool {
	url := "https://ti.qq.com/act/99welfare2022/trpc/CharityDayProxy/JoinTeam?bkn=440346013"
	body := "{\"uin\":" + uin + ",\"team_id\":\"" + team + "\"}"
	ret := a.request(url, body, a.getCookieString(uin, skey, pskey))
	if ret == "" {
		return false
	}
	var res *joinTeamResponse
	err := json.Unmarshal([]byte(ret), &res)
	if err != nil || res == nil || res.Error.Code != 0 {
		return false
	}
	return true
}

// speed 领取加速
func (a *AutoWelfare) speed(uin string, skey string, pskey string, team string, types string) bool {
	url := "https://ti.qq.com/act/99welfare2022/trpc/CharityDayProxy/ReceiveAccelerate?bkn=1038960899"
	body := "{\"uin\":" + uin + ",\"team_id\":\"" + team + "\",\"role_type\":" + types + "}"
	ret := a.request(url, body, a.getCookieString(uin, skey, pskey))
	if ret == "" {
		return false
	}
	var res *speedResponse
	err := json.Unmarshal([]byte(ret), &res)
	if err != nil || res == nil {
		return false
	}
	return res.Error.Code == 0 || res.Error.Code == 303
}

// getMyTeam 获取我的组队个数
func (a *AutoWelfare) getMyTeamNum(uin string, skey string, pskey string) int {
	url := "https://ti.qq.com/act/99welfare2022/trpc/CharityDayProxy/GetMyTeamRecords?bkn=440346013"
	body := "{\"uin\":" + uin + "}"
	ret := a.request(url, body, a.getCookieString(uin, skey, pskey))
	if ret == "" {
		return -1
	}
	var res *myTeamResponse
	err := json.Unmarshal([]byte(ret), &res)
	if err != nil || res == nil {
		return -1
	}
	return len(res.Response.TeamInfos)
}

// CheckAllow 判断是否可以进入队伍
func (a *AutoWelfare) CheckAllow(uin string, skey string, pskey string) bool {
	return a.getMyTeamNum(uin, skey, pskey) == 0
}

// Run 执行
func (a *AutoWelfare) Run() bool {
	//创建队伍
	team := a.createTeam(a.Leader.Uin, a.Leader.Skey, a.Leader.Pskey)
	if team == "" {
		team = a.createTeam(a.Leader.Uin, a.Leader.Skey, a.Leader.Pskey)
	}
	if team != "" {
		//加入队伍
		join1 := a.joinTeam(a.UserOne.Uin, a.UserOne.Skey, a.UserOne.Pskey, team)
		if !join1 {
			time.Sleep(500 * time.Millisecond)
			join1 = a.joinTeam(a.UserOne.Uin, a.UserOne.Skey, a.UserOne.Pskey, team)
		}
		if join1 {
			join2 := a.joinTeam(a.UserTwo.Uin, a.UserTwo.Skey, a.UserTwo.Pskey, team)
			if !join2 {
				time.Sleep(500 * time.Millisecond)
				a.joinTeam(a.UserTwo.Uin, a.UserTwo.Skey, a.UserTwo.Pskey, team)
			}
			//队员2加入队伍
			if join2 {
				//领取加速
				if !a.speed(a.Leader.Uin, a.Leader.Skey, a.Leader.Pskey, team, "0") {
					time.Sleep(500 * time.Millisecond)
					a.speed(a.Leader.Uin, a.Leader.Skey, a.Leader.Pskey, team, "0")
				}
				if !a.speed(a.UserOne.Uin, a.UserOne.Skey, a.UserOne.Pskey, team, "1") {
					time.Sleep(500 * time.Millisecond)
					a.speed(a.UserOne.Uin, a.UserOne.Skey, a.UserOne.Pskey, team, "1")
				}
				if !a.speed(a.UserTwo.Uin, a.UserTwo.Skey, a.UserTwo.Pskey, team, "1") {
					time.Sleep(500 * time.Millisecond)
					a.speed(a.UserTwo.Uin, a.UserTwo.Skey, a.UserTwo.Pskey, team, "1")
				}
				return true
			}
		}
	}
	return false
}

// getCookieString 获取cookie字符串
func (*AutoWelfare) getCookieString(uin string, skey string, pskey string) string {
	return "uin=o" + uin + "; skey=" + skey + "; p_uin=o" + uin + "; p_skey=" + pskey
}

// request 发送请求
func (*AutoWelfare) request(url string, body string, cookie string) string {
	client := &http.Client{}
	method := "GET"
	if body != "" {
		method = "POST"
	}
	var data = strings.NewReader(body)
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return ""
	}
	req.Header.Set("Origin", "https://ti.qq.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://ti.qq.com/act/99welfare2022/organize?_wv=3&_wwv=128")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; MI 8 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/89.0.4389.72 MQQBrowser/6.2 TBS/046125 Mobile Safari/537.36 V1_AND_SQ_8.8.98_3002_YYB_D A_8089800 QQ/8.8.98.8410 NetType/WIFI WebP/0.3.0 Pixel/1080 StatusBarHeight/90 SimpleUISwitch/0 QQTheme/1000 InMagicWin/0 StudyMode/0 CurrentMode/0 CurrentFontScale/1.0 GlobalDensityScale/0.9818182 AppId/537124039")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	fmt.Println("请求地址", url, "请求内容", body, "Cookie:", cookie, "返回内容", string(bodyText))
	return string(bodyText)
}

// createTeamResponse 创建队伍响应
type createTeamResponse struct {
	Response struct {
		TeamInfo struct {
			Members []interface{} `json:"members"`
			TeamId  string        `json:"team_id"`
			Leader  struct {
				Uin          string `json:"uin"`
				Avatar       string `json:"avatar"`
				IsAccelerate bool   `json:"is_accelerate"`
				IsAwardBall  bool   `json:"is_award_ball"`
			} `json:"leader"`
			IsHasBall         bool   `json:"is_has_ball"`
			IsAccelerateAward bool   `json:"is_accelerate_award"`
			CreateTime        string `json:"create_time"`
			IsExpired         bool   `json:"is_expired"`
			Cas               string `json:"cas"`
		} `json:"team_info"`
	} `json:"response"`
	RetCode  int `json:"retCode"`
	CostTime int `json:"costTime"`
	Error    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// joinTeamResponse 加入组队响应
type joinTeamResponse struct {
	Response struct {
		TeamInfo struct {
			Members []struct {
				Uin          string `json:"uin"`
				Avatar       string `json:"avatar"`
				IsAccelerate bool   `json:"is_accelerate"`
				IsAwardBall  bool   `json:"is_award_ball"`
			} `json:"members"`
			TeamId string `json:"team_id"`
			Leader struct {
				Uin          string `json:"uin"`
				Avatar       string `json:"avatar"`
				IsAccelerate bool   `json:"is_accelerate"`
				IsAwardBall  bool   `json:"is_award_ball"`
			} `json:"leader"`
			IsHasBall         bool   `json:"is_has_ball"`
			IsAccelerateAward bool   `json:"is_accelerate_award"`
			CreateTime        string `json:"create_time"`
			IsExpired         bool   `json:"is_expired"`
			Cas               string `json:"cas"`
		} `json:"team_info"`
	} `json:"response"`
	RetCode  int `json:"retCode"`
	CostTime int `json:"costTime"`
	Error    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// myTeamResponse 我的组队响应
type myTeamResponse struct {
	Response struct {
		TeamInfos []struct {
			Members []struct {
				Uin          string `json:"uin"`
				Avatar       string `json:"avatar"`
				IsAccelerate bool   `json:"is_accelerate"`
				IsAwardBall  bool   `json:"is_award_ball"`
			} `json:"members"`
			TeamId string `json:"team_id"`
			Leader struct {
				Uin          string `json:"uin"`
				Avatar       string `json:"avatar"`
				IsAccelerate bool   `json:"is_accelerate"`
				IsAwardBall  bool   `json:"is_award_ball"`
			} `json:"leader"`
			IsHasBall         bool   `json:"is_has_ball"`
			IsAccelerateAward bool   `json:"is_accelerate_award"`
			CreateTime        string `json:"create_time"`
			IsExpired         bool   `json:"is_expired"`
			Cas               string `json:"cas"`
		} `json:"team_infos"`
	} `json:"response"`
	RetCode  int `json:"retCode"`
	CostTime int `json:"costTime"`
	Error    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// speedResponse 加速响应
type speedResponse struct {
	Response struct {
	} `json:"response"`
	RetCode  int `json:"retCode"`
	CostTime int `json:"costTime"`
	Error    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
