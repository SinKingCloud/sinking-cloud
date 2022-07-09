package uin_tool

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("错误：请求地址", url, "请求参数", s.toJson(data), "结果", string(bytes))
		return nil, err
	}
	body := str.NewStringTool(string(bytes)).ConvertToString("gbk", "utf-8")
	var res *response
	err = json.Unmarshal([]byte(body), &res)
	log.Println("请求地址", url, "请求参数", s.toJson(data), "结果", body)
	if err != nil {
		return nil, err
	}
	res.Body = body
	return res, nil
}

// response 响应结构
type response struct {
	Type     string `json:"Type"`
	State    string `json:"State"`
	Uin      int64  `json:"Uin"`
	Pskey    string `json:"pskey"`
	Skey     string `json:"skey"`
	Superkey string `json:"Superkey"`
	Body     string `json:"-"`
}

// SpeedStatus 任务状态
type SpeedStatus struct {
	CommonInfoGetQQLevelInfo struct {
		Lxby                string `json:"Lxby"`
		QzoneVisitor        string `json:"QzoneVisitor"` //访客加速
		SVIPStar            string `json:"SVIPStar"`
		WeishiVideoview     string `json:"WeishiVideoview"` //微视加速
		IAddFriend          string `json:"iAddFriend"`      //加好友加速
		IBaseDays           string `json:"iBaseDays"`       //基础服务加速
		IBigClubGrowth      string `json:"iBigClubGrowth"`
		IBigClubLevel       string `json:"iBigClubLevel"`
		IBigClubSpeed       string `json:"iBigClubSpeed"`
		IBigClubVipFlag     string `json:"iBigClubVipFlag"`
		ICmshowActive       string `json:"iCmshowActive"`       //厘米秀加速
		IContinueLogin      string `json:"iContinueLogin"`      //连续登陆加速
		ICostMs             string `json:"iCostMs"`             //本次请求毫秒数
		IDailySign          string `json:"iDailySign"`          //打卡加速
		IMaxLvlRealDays     string `json:"iMaxLvlRealDays"`     //开通会员基础服务加速活跃天数
		IMaxLvlTotalDays    string `json:"iMaxLvlTotalDays"`    //开通会员最高加速活跃天数
		IMedal              string `json:"iMedal"`              //勋章加速
		IMobileGameOnline   string `json:"iMobileGameOnline"`   //游戏加速
		IMobileQQOnline     string `json:"iMobileQQOnline"`     //手机加速
		IMobileQQOnlineTime string `json:"iMobileQQOnlineTime"` //手机在线时长
		IMqing              string `json:"iMqing"`
		INextLevelDay       string `json:"iNextLevelDay"`     //距离下次升级活跃天数
		INoHideOnline       string `json:"iNoHideOnline"`     //非隐身在线加速
		INoHideOnlineTime   string `json:"iNoHideOnlineTime"` //非隐身在线时长
		IPCQQOnline         string `json:"iPCQQOnline"`       //电脑加速
		IPCQQOnlineTime     string `json:"iPCQQOnlineTime"`   //电脑在线时长
		IPCSafeOnline       string `json:"iPCSafeOnline"`     //电脑管家加速
		IQQLevel            string `json:"iQQLevel"`          //QQ等级
		IQQSportStep        string `json:"iQQSportStep"`      //运动加速
		IQzoneState         string `json:"iQzoneState"`       //空间说说加速
		IRealDays           string `json:"iRealDays"`         //今日加速天数
		ISVip               string `json:"iSVip"`             //是否为超级会员
		ISpeedType          string `json:"iSpeedType"`
		ISqq                string `json:"iSqq"`
		ISqqLevel           string `json:"iSqqLevel"`
		ISqqSpeedRate       string `json:"iSqqSpeedRate"`
		ISvrDays            string `json:"iSvrDays"`        //服务加速
		ITotalActiveDay     string `json:"iTotalActiveDay"` //QQ活跃天数
		ITotalDays          string `json:"iTotalDays"`      //QQ最大活跃天数
		IVip                string `json:"iVip"`            //是否为会员
		IVipLevel           string `json:"iVipLevel"`       //会员等级
		IVipSpeedRate       string `json:"iVipSpeedRate"`   //会员成长值速度
		IYearBigClubFlag    string `json:"iYearBigClubFlag"`
		IYearVip            string `json:"iYearVip"`  //是否为年费会员
		SFaceUrl            string `json:"sFaceUrl"`  //头像地址
		SNickName           string `json:"sNickName"` //昵称
		SpeedStar           string `json:"speedStar"` //是否开通加速星
	} `json:"commonInfo.getQQLevelInfo"`
	CommonInfoGetQQInfo struct {
		Oidb0X99Ret           string      `json:"Oidb0x99Ret"`
		VipCareInfoRet        string      `json:"VipCareInfoRet"`
		IGrowthRet            string      `json:"iGrowthRet"`   //是否为会员
		IGrowthSpeed          string      `json:"iGrowthSpeed"` //会员成长值速度
		IGrowthValue          string      `json:"iGrowthValue"` //会员成长值
		IOpenSVip             string      `json:"iOpenSVip"`    //是否开通svip
		IOpenVip              string      `json:"iOpenVip"`     //是否开通vip
		IQQLevel              string      `json:"iQQLevel"`     //QQ等级
		IQQLevelRet           string      `json:"iQQLevelRet"`
		ISVipCareFlag         string      `json:"iSVipCareFlag"`
		IVipCareFlag          string      `json:"iVipCareFlag"`
		IVipLevel             string      `json:"iVipLevel"` //会员等级
		IVipType              string      `json:"iVipType"`
		SMiddlePicUrl         string      `json:"sMiddlePicUrl"`
		SNickName             string      `json:"sNickName"`
		SNickNameRet          string      `json:"sNickNameRet"`
		SPicUrlRet            string      `json:"sPicUrlRet"`
		SSVipExpireTime       string      `json:"sSVipExpireTime"`
		SVipExpireTime        string      `json:"sVipExpireTime"`
		SYearExpireTime       string      `json:"sYearExpireTime"`
		Uin                   string      `json:"uin"`
		IsMobileClub          string      `json:"isMobileClub"`
		ClubEndTimestamp      string      `json:"club_end_timestamp"`
		SuperclubEndTimestamp string      `json:"superclub_end_timestamp"`
		SuperclubYearPeriod   interface{} `json:"superclub_year_period"`
		ClubStatus            int         `json:"clubStatus"`
	} `json:"commonInfo.getQQInfo"`
	Pskey    string `json:"pskey"`
	Skey     string `json:"skey"`
	Superkey string `json:"Superkey"`
}
