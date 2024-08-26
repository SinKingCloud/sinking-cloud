package web_uin

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	cookie2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/cookie"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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
	Data struct {
		MRes uinInfoResponseDataStruct `json:"mRes"`
	} `json:"data"`
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type uinInfoResponseDataStruct struct {
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

// GetGuestUinInfoBak 获取QQ信息(2号接口)
func GetGuestUinInfoBak(uin int64) *UinInfo {
	//获取cookie
	uc, err := cookie2.GetIns().FindByStatusDomainWithCache(cookie2.Normal, "vip.qq.com", false)
	if err != nil {
		return nil
	}
	u := "https://vip.qq.com/pk/index?param=" + strconv.FormatInt(uin, 10)
	req, err := http.NewRequest("GET", u, nil)
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
	var res *uinInfoResponseDataStruct
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
		if d, ok := levelList[key]; ok {
			useDgDays = useDgDays + 1.5*d
			useDgDays = math.Round(useDgDays*100) / 100
			s := strconv.FormatFloat(useDgDays, 'f', 2, 64)
			useDgDays, _ = strconv.ParseFloat(s, 64)
		}
	} else {
		useDgDays = useDgDays + 1.5
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

// GetGuestUinInfo 获取QQ信息（1号接口）
func GetGuestUinInfo(uin int64) *UinInfo {
	//获取cookie
	uc, err := cookie2.GetIns().FindByStatusDomainWithCache(cookie2.Normal, "vip.qq.com", false)
	if err != nil {
		return nil
	}
	sKey := ""
	pSkey := ""
	arr := strings.Split(uc.Cookie, ";")
	for _, v := range arr {
		temp := strings.Split(v, "=")
		if len(temp) == 2 && strings.Trim(temp[0], " ") == "p_skey" {
			pSkey = temp[1]
		}
		if len(temp) == 2 && strings.Trim(temp[0], " ") == "skey" {
			sKey = temp[1]
		}
	}
	gtk := utils.GetGTK(pSkey)
	reqBody := "{\"sClientIp\":\"127.0.0.1\",\"sSessionKey\":\"" + sKey + "\",\"iKeyType\":1,\"iAppId\":0,\"iUin\":\"" + strconv.FormatInt(uin, 10) + "\"}"
	u := "https://club.vip.qq.com/api/vip/getQQLevelInfo?g_tk=" + strconv.FormatInt(gtk, 10) + "&requestBody=" + url.QueryEscape(reqBody)
	req, err := http.NewRequest("GET", u, nil)
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
	var response *uinInfoResponseStruct
	err = json.Unmarshal(bytes, &response)
	if err != nil || response == nil {
		_, _ = cookie2.GetIns().FindByStatusDomainWithCache(cookie2.Normal, "vip.qq.com", true)
		return nil
	}
	if response.Ret != 0 {
		return nil
	}
	res := response.Data.MRes
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
		if d, ok := levelList[key]; ok {
			useDgDays = useDgDays + 1.5*d
			useDgDays = math.Round(useDgDays*100) / 100
			s := strconv.FormatFloat(useDgDays, 'f', 2, 64)
			useDgDays, _ = strconv.ParseFloat(s, 64)
		}
	} else {
		useDgDays = useDgDays + 1.5
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

// GetLevelInfo 获取等级信息
func GetLevelInfo(uin string) *UinInfo {
	u, err := strconv.ParseInt(uin, 10, 64)
	if err != nil || u <= 10000 || u >= 9999999999 {
		return nil
	}
	var info *UinInfo
	key := "LevelInfo_" + uin
	a, exists := util.Cache.Get(key)
	if exists {
		_ = json.Unmarshal([]byte(a), &info)
	}
	if info == nil {
		info = GetGuestUinInfo(u)
		if info == nil {
			info = GetGuestUinInfoBak(u)
		}
		if info != nil {
			bytes1, e := json.Marshal(info)
			if e == nil {
				util.Cache.SetWithExpire(key, string(bytes1), 3600*1000)
			}
		}
	}
	return info
}
