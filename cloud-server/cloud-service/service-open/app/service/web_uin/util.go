package web_uin

import (
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	cookie2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/cookie"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// GetLoginQrCode 获取设备锁二维码
func GetLoginQrCode(verifyUrl string) string {
	if !strings.Contains(verifyUrl, "uin=") || !strings.Contains(verifyUrl, "&sig=") {
		return ""
	}
	//解析参数
	uin := str.NewStringTool(verifyUrl).GetBetween("uin=", "&sig=")
	sig := str.NewStringTool(verifyUrl).GetBetween("&sig=", "-1")
	//构建url
	url := "https://ti.qq.com/safe/verify/proxy/domain/oidb.tim.qq.com/v3/oidbinterface/oidb_0xc9e_4?uid=" + uin + "&getqrcode=1&sdkappid=39998&actype=2&bkn="
	req, err := http.NewRequest("POST", url, strings.NewReader("{\"str_dev_auth_token\":\""+sig+"\",\"uint32_flag\":1}"))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Cookie", "domainId=338; uin=o"+uin+"; p_uin=o"+uin)
	req.Header.Set("Origin", "https://ti.qq.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("qname-service", "1104449:131072")
	req.Header.Set("qname-space", "Production")
	req.Header.Set("Referer", verifyUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.124 Safari/537.36 Edg/102.0.1245.41")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//解析返回值
	type result struct {
		Uint32GuaranteeStatus int    `json:"uint32_guarantee_status"`
		StrUrl                string `json:"str_url"`
		ActionStatus          string `json:"ActionStatus"`
		ErrorCode             int    `json:"ErrorCode"`
		ErrorInfo             string `json:"ErrorInfo"`
	}
	var res *result
	err = json.Unmarshal(bytes, &res)
	if err != nil || res == nil {
		return ""
	}
	if res.ErrorCode != 0 {
		return ""
	}
	d := "https://ti.qq.com/safe/scanresult?str_url=" + res.StrUrl + "&envfrom=&verify_id=&verify_scene=&uin=" + uin
	return d
}

type UinInfo struct {
	Uin                    int64    `json:"uin"`                         //账号
	NickName               string   `json:"nick_name"`                   //昵称
	Avatar                 string   `json:"avatar"`                      //头像
	Level                  int      `json:"level"`                       //等级
	TotalActiveDay         int      `json:"total_active_day"`            //活跃天
	NextLevelDay           int      `json:"next_level_day"`              //下次升级所需活跃天数
	NextCrownLevelDay      int      `json:"next_crown_level_day"`        //下次升级皇冠所需活跃天数
	IsVip                  bool     `json:"is_vip"`                      //是否会员
	IsSvip                 bool     `json:"is_svip"`                     //是否超级会员
	IsYearVip              bool     `json:"is_year_vip"`                 //是否年费会员
	IsBigVip               bool     `json:"is_big_vip"`                  //是否大会员
	IsBigYearVip           bool     `json:"is_big_year_vip"`             //是否年费大会员
	VipLevel               int      `json:"vip_level"`                   //会员等级
	BigVipLevel            int      `json:"big_vip_level"`               //大会员等级
	VipIcon                []string `json:"vip_icon"`                    //会员图标
	LevelIcon              []string `json:"level_icon"`                  //等级图标
	IfUseDayActiveDay      float64  `json:"if_use_day_active_day"`       //如果是使用代挂每日获取的活跃天
	IfUseNextLevelDay      int      `json:"if_use_next_level_day"`       //如果使用代挂下次升级所需天数
	IfUseNextCrownLevelDay int      `json:"if_use_next_crown_level_day"` //如果使用代挂下次升级皇冠所需天数
}
type uinInfoResponseStruct struct {
	Ret                 int    `json:"ret"`
	QzoneVisitor        string `json:"QzoneVisitor"`
	IContinueLogin      string `json:"iContinueLogin"`
	IMaxLvlRealDays     string `json:"iMaxLvlRealDays"`
	IMedal              string `json:"iMedal"`
	IMobileQQOnlineTime string `json:"iMobileQQOnlineTime"`
	INextLevelDay       string `json:"iNextLevelDay"`
	ISVip               string `json:"iSVip"`
	IYearVip            string `json:"iYearVip"`
	SFaceUrl            string `json:"sFaceUrl"`
	SNickName           string `json:"sNickName"`
	IBigClubLevel       string `json:"iBigClubLevel"`
	ICostMs             string `json:"iCostMs"`
	IDailySign          string `json:"iDailySign"`
	ISqqSpeedRate       string `json:"iSqqSpeedRate"`
	IVip                string `json:"iVip"`
	SpeedStarv3         string `json:"speedStarv3"`
	Lxby                string `json:"Lxby"`
	IBigClubVipFlag     string `json:"iBigClubVipFlag"`
	IYearBigClubFlag    string `json:"iYearBigClubFlag"`
	SpeedStar           string `json:"speedStar"`
	SpeedStarv2         string `json:"speedStarv2"`
	IBaseDays           string `json:"iBaseDays"`
	IBigClubGrowth      string `json:"iBigClubGrowth"`
	INoHideOnlineTime   string `json:"iNoHideOnlineTime"`
	IPCSafeOnline       string `json:"iPCSafeOnline"`
	IQzoneState         string `json:"iQzoneState"`
	ITotalActiveDay     string `json:"iTotalActiveDay"`
	SVIPStar            string `json:"SVIPStar"`
	WeishiVideoview     string `json:"WeishiVideoview"`
	IAddFriend          string `json:"iAddFriend"`
	IMaxLvlTotalDays    string `json:"iMaxLvlTotalDays"`
	IMobileGameOnline   string `json:"iMobileGameOnline"`
	IMobileQQOnline     string `json:"iMobileQQOnline"`
	IMqing              string `json:"iMqing"`
	INoHideOnline       string `json:"iNoHideOnline"`
	IPCQQOnlineTime     string `json:"iPCQQOnlineTime"`
	IQQLevel            string `json:"iQQLevel"`
	IRealDays           string `json:"iRealDays"`
	ISqq                string `json:"iSqq"`
	ISqqLevel           string `json:"iSqqLevel"`
	ISvrDays            string `json:"iSvrDays"`
	ITotalDays          string `json:"iTotalDays"`
	IVipLevel           string `json:"iVipLevel"`
	IBigClubSpeed       string `json:"iBigClubSpeed"`
	ICmshowActive       string `json:"iCmshowActive"`
	IPCQQOnline         string `json:"iPCQQOnline"`
	IQQSportStep        string `json:"iQQSportStep"`
	ISmallWorld         string `json:"iSmallWorld"`
	ISpeedType          string `json:"iSpeedType"`
	IVipSpeedRate       string `json:"iVipSpeedRate"`
}

// GetGuestUinInfo 获取QQ信息
func GetGuestUinInfo(uin int64) *UinInfo {
	//获取cookie
	uc, err := cookie2.GetIns().FindByStatusDomainWithCache(cookie2.Normal, "vip.qq.com", false)
	if err != nil {
		return nil
	}
	url := "https://vip.qq.com/pk/index?param=" + strconv.FormatInt(uin, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Cookie", uc.Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	body := string(bytes)
	if !strings.Contains(body, "我的等级pk") {
		_, _ = cookie2.GetIns().FindByStatusDomainWithCache(cookie2.Normal, "vip.qq.com", true)
		return nil
	}
	jsonString := str.NewStringTool(body).GetBetween("GUEST_LEVEL_INFO =", ";")
	var res *uinInfoResponseStruct
	err = json.Unmarshal([]byte(jsonString), &res)
	if err != nil || res == nil || res.Ret != 0 {
		return nil
	}
	var vipIcon []string
	if res.IVip != "1" && res.ISVip != "1" && res.IBigClubVipFlag != "1" {
		vipIcon = append(vipIcon, "https://imgcache.gtimg.cn/vipstyle/component/img/nameplate/vip/lv/vip1-l2-off.png")
	} else {
		if res.IVip == "1" || res.ISVip == "1" {
			name := "vip"
			if res.ISVip == "1" {
				name = "svip"
			}
			year := ""
			if res.IYearVip == "1" {
				year = "year-"
			}
			icon := ""
			level := res.IVipLevel
			if n, e := strconv.Atoi(level); e == nil && n >= 10 {
				icon = "https://qzonestyle.gtimg.cn/qzone/qzact/act/external/clubindex/qqvip/svip10/0-mp.png?nowebp&nosharpp&_bid=5065"
			} else {
				icon = "https://imgcache.gtimg.cn/vipstyle/component/img/nameplate/" + name + "/lv/" + year + name + level + "-l2-on.png"
			}
			vipIcon = append(vipIcon, icon)
		}
		if res.IBigClubVipFlag == "1" {
			vipIcon = append(vipIcon, "https://qzonestyle.gtimg.cn/vipstyle/component/img/nameplate/bvip/lv/LV"+res.IBigClubLevel+".png")
		}
	}
	//等级图标
	level, err := strconv.Atoi(res.IQQLevel)
	if err != nil {
		level = 0
	}
	var levelIcon []string
	iconArr := map[int]string{64: "https://imgcache.qq.com/xydata/qqVipLevel/item/3/crown.png",
		16: "https://imgcache.qq.com/xydata/qqVipLevel/item/3/sun.png",
		4:  "https://imgcache.qq.com/xydata/qqVipLevel/item/3/moon.png",
		1:  "https://imgcache.qq.com/xydata/qqVipLevel/item/3/star.png"}
	keys := []int{64, 16, 4, 1}
	tmp := level
	crownNum := 1
	for _, k := range keys {
		v := iconArr[k]
		num := tmp / k
		if k == 64 {
			crownNum = num + 1
		}
		tmp -= num * k
		for i := 0; i < num; i++ {
			levelIcon = append(levelIcon, v)
		}
	}
	totalActiveDay, err := strconv.Atoi(res.ITotalActiveDay)
	if err != nil {
		totalActiveDay = 0
	}
	nextCrownLevelDay := crownNum*64*crownNum*64 + 4*crownNum*64
	nextLevelDay, err := strconv.Atoi(res.INextLevelDay)
	if err != nil {
		nextLevelDay = 0
	}
	nextCrownLevelDay = nextCrownLevelDay - totalActiveDay
	vipLevel, err := strconv.Atoi(res.IVipLevel)
	if err != nil {
		vipLevel = 0
	}
	BigVipLevel, err := strconv.Atoi(res.IBigClubLevel)
	if err != nil {
		BigVipLevel = 0
	}
	levelList := map[string]float64{
		"vip1":          1.2,
		"vip2":          1.4,
		"vip3":          1.5,
		"vip4":          1.6,
		"vip5":          1.6,
		"vip6":          1.7,
		"vip7":          1.9,
		"vip8":          1.9,
		"vip9":          1.9,
		"vip10":         1.9,
		"vip10_1":       1.9,
		"vip10_2":       1.9,
		"vip10_3":       1.9,
		"year_vip1":     1.5,
		"year_vip2":     1.7,
		"year_vip3":     1.8,
		"year_vip4":     1.9,
		"year_vip5":     1.9,
		"year_vip6":     2.0,
		"year_vip7":     2.2,
		"year_vip8":     2.2,
		"year_vip9":     2.2,
		"year_vip10":    2.2,
		"year_vip10_1":  2.2,
		"year_vip10_2":  2.2,
		"year_vip10_3":  2.2,
		"svip1":         1.4,
		"svip2":         1.6,
		"svip3":         1.7,
		"svip4":         1.8,
		"svip5":         1.8,
		"svip6":         1.9,
		"svip7":         2.1,
		"svip8":         2.2,
		"svip9":         2.5,
		"svip10":        2.9,
		"svip10_1":      3.1,
		"svip10_2":      3.4,
		"svip10_3":      3.9,
		"year_svip1":    1.7,
		"year_svip2":    1.9,
		"year_svip3":    2.0,
		"year_svip4":    2.1,
		"year_svip5":    2.2,
		"year_svip6":    2.4,
		"year_svip7":    2.7,
		"year_svip8":    3.0,
		"year_svip9":    3.2,
		"year_svip10":   3.5,
		"year_svip10_1": 3.7,
		"year_svip10_2": 4.0,
		"year_svip10_3": 4.5,
	}
	useDgDays := 5.0
	key := ""
	if res.IVip == "1" {
		key = "vip"
	}
	if res.ISVip == "1" {
		key = "svip"
	}
	if key != "" && res.IYearVip == "1" {
		key = "year_" + key
	}
	if key != "" {
		key = key + strconv.Itoa(vipLevel)
		fmt.Println(key)
		if d, ok := levelList[key]; ok {
			useDgDays = useDgDays + 1.5*d
			useDgDays = math.Round(useDgDays*100) / 100
			s := strconv.FormatFloat(useDgDays, 'f', 2, 64)
			useDgDays, _ = strconv.ParseFloat(s, 64)
		}
	}
	ifUseNextLevelDay := int(math.Ceil(float64(nextLevelDay) / useDgDays))
	ifUseNextCrownLevelDay := int(math.Ceil(float64(nextCrownLevelDay) / useDgDays))
	return &UinInfo{
		Uin:                    uin,
		NickName:               res.SNickName,
		Avatar:                 res.SFaceUrl,
		Level:                  level,
		TotalActiveDay:         totalActiveDay,
		NextLevelDay:           nextLevelDay,
		NextCrownLevelDay:      nextCrownLevelDay,
		IsVip:                  res.IVip == "1",
		IsSvip:                 res.ISVip == "1",
		IsYearVip:              res.IYearVip == "1",
		IsBigVip:               res.IBigClubVipFlag == "1",
		IsBigYearVip:           res.IYearBigClubFlag == "1",
		VipLevel:               vipLevel,
		BigVipLevel:            BigVipLevel,
		VipIcon:                vipIcon,
		LevelIcon:              levelIcon,
		IfUseDayActiveDay:      useDgDays,
		IfUseNextLevelDay:      ifUseNextLevelDay,
		IfUseNextCrownLevelDay: ifUseNextCrownLevelDay,
	}
}
