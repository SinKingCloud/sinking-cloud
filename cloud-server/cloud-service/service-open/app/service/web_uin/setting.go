package web_uin

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	cookie2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/cookie"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin_info"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ServerUinCreate 挂机服务器创建账号
func ServerUinCreate(serverId int, uin string, pwd string, phone string, phoneArea string, loginType int) error {
	serverInfo, err := server2.GetIns().FindById(serverId)
	if err != nil || serverInfo == nil || serverInfo.Status != server2.Normal {
		return errors.New("未查询到挂机服务器信息,请更换挂机服务器")
	}
	s := uin_tool.NewUinTool(serverInfo.Addr, serverInfo.Proxy)
	if !s.Create(uin, pwd, phone, phoneArea, loginType) {
		return errors.New("更新挂机信息失败,请重试")
	}
	return nil
}

// ServerUinDelete 挂机服务器删除账号
func ServerUinDelete(serverId int, uin string) error {
	serverInfo, err := server2.GetIns().FindById(serverId)
	if err != nil || serverInfo == nil || serverInfo.Status != server2.Normal {
		return errors.New("未查询到挂机服务器信息")
	}
	s := uin_tool.NewUinTool(serverInfo.Addr, serverInfo.Proxy)
	s.MobileOffline(uin)
	s.PcOffline(uin)
	if !s.Delete(uin) {
		return errors.New("删除挂机失败")
	}
	return nil
}

// ChangeSetting 更改设置
func ChangeSetting(uinId int, newSetting map[string]string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	var set map[string]string
	err = json.Unmarshal([]byte(uinInfo.Setting), &set)
	if err != nil {
		set = make(map[string]string)
	}
	for k := range Setting {
		if newSetting[k] != "" && set[k] != newSetting[k] {
			if newSetting[k] == "1" {
				set[k] = "1"
			} else {
				set[k] = "0"
			}
		}
	}
	bytes, err := json.Marshal(set)
	if err != nil {
		return errors.New("修改配置错误")
	}
	err = uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{
		"setting": string(bytes),
	})
	if err != nil {
		return errors.New("修改配置失败")
	}
	return nil
}

// ChangePassword 更改挂机密码
func ChangePassword(uinId int, pwd string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	if pwd != uinInfo.Password || uinInfo.LoginType != 0 {
		data := make(map[string]interface{})
		data["password"] = pwd
		data["login_type"] = 0
		data["uin_status"] = 1
		go func() {
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Uin:     uinInfo.Uin,
				Type:    log.Password,
				Title:   "修改密码",
				Content: "修改挂机密码",
			})
		}()
		_ = ServerUinCreate(uinInfo.ServerId, uinInfo.Uin, pwd, uinInfo.Phone, uinInfo.PhoneArea, 0)
		err = uin.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			return errors.New("更改挂机密码失败")
		}
	}
	return nil
}

// ChangePhone 更改挂机手机
func ChangePhone(uinId int, phone string, phoneArea string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	if phone != uinInfo.Phone || uinInfo.LoginType != 1 {
		data := make(map[string]interface{})
		data["phone"] = phone
		data["phone_area"] = phoneArea
		data["login_type"] = 1
		data["uin_status"] = 1
		go func() {
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Uin:     uinInfo.Uin,
				Type:    log.Phone,
				Title:   "修改手机",
				Content: "修改挂机手机",
			})
		}()
		_ = ServerUinCreate(uinInfo.ServerId, uinInfo.Uin, uinInfo.Password, phone, phoneArea, 1)
		err = uin.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			return errors.New("更改挂机手机失败")
		}
	}
	return nil
}

// GetUinLevelInfo 获取等级信息
func GetUinLevelInfo(uinId int) (*model.AppUinInfo, error) {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return nil, errors.New("未查询到QQ信息")
	}
	refresh := false
	uinLevelInfo, e := uin_info.GetIns().FindByUin(uinInfo.Uin)
	//失效时刷新
	if (e != nil || uinLevelInfo == nil) && uinInfo.UinStatus == uin.Normal {
		refresh = true
	}
	//每隔三小时或者上次刷新在昨天时刷新一次数据
	if !refresh && e == nil && uinLevelInfo != nil {
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uinLevelInfo.UpdateTime.String(), time.Local)
		todayUnix := str.NewTimeTool().GetZeroTime(time.Now()).Unix()
		if stamp.Unix() < todayUnix || stamp.Unix()+60*60*3 < time.Now().Unix() {
			refresh = true
		}
	}
	//刷新QQ等级信息
	if refresh {
		server, e2 := server2.GetIns().FindById(uinInfo.ServerId)
		if e2 == nil && server != nil {
			_, err2, isFail := GetUinSpeedStatus(server.Addr, server.Proxy, uinInfo.Uin, false)
			if !isFail && err2 == nil {
				uinLevelInfo, e = uin_info.GetIns().FindByUin(uinInfo.Uin)
			}
		}
	}
	return uinLevelInfo, e
}

// GetLevelInfo 获取等级信息
func GetLevelInfo(uin string, skey string, pskey string) (*uin_tool.SpeedStatus, error, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://ti.qq.com/qqlevel/srf/getQQLevelAllInfo?bkn=683611393", nil)
	if err != nil {
		return nil, errors.New("发起请求失败"), false
	}
	req.Header.Set("Origin", "https://ti.qq.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://ti.qq.com/qqlevel/index?")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/104.0.5112.81")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "uin=o"+uin+"; skey="+skey+"; p_skey="+pskey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("请求失败"), false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("获取返回值失败"), false
	}
	type speedStatus struct {
		Data *uin_tool.SpeedStatus `json:"data"`
	}
	var t *speedStatus
	e := json.Unmarshal(bodyText, &t)
	if e != nil || t == nil || t.Data == nil {
		return nil, errors.New("cookies失效"), true
	}
	return t.Data, nil, false
}

// GetUinSpeedStatus 获取加速任务状态
func GetUinSpeedStatus(addr string, proxy string, uin string, isSync bool) (status *uin_tool.SpeedStatus, err error, isFail bool) {
	err, skey, pskey := GetUinCookies(addr, proxy, uin, "ti.qq.com", false)
	if err == nil && skey != "" && pskey != "" {
		status2, e, isExpire := GetLevelInfo(uin, skey, pskey)
		if e == nil && !isExpire && status2 != nil {
			status = status2
		}
	}
	if status == nil {
		s := uin_tool.NewUinTool(addr, proxy)
		status, err, isFail = s.Status(uin)
	} else {
		err = nil
		isFail = false
	}
	if err == nil && status != nil && !isFail {
		//更新QQ信息
		fun := func(status *uin_tool.SpeedStatus, uin string) {
			//qq信息
			level, e := strconv.Atoi(status.IQQLevel)
			if e != nil {
				level = 0
			}
			info := &model.AppUinInfo{
				Uin:         uin,
				NickName:    status.SNickName,
				Level:       level,
				ActivateDay: status.ITotalActiveDay,
				TodayDay:    status.IRealDays,
			}
			//加速状态
			speedInfo := map[string]string{
				"mobile":       status.IMobileQQOnline,
				"pc":           status.IPCQQOnline,
				"online":       status.INoHideOnline,
				"zone":         status.QzoneVisitor,
				"medal":        status.IMedal,
				"ws":           status.WeishiVideoview,
				"micro_show":   status.ICmshowActive,
				"safe":         status.IPCSafeOnline,
				"game":         status.IMobileGameOnline,
				"sport":        status.IQQSportStep,
				"mobile_clock": status.IDailySign,
				"zone_clock":   status.IQzoneState,
				"login":        status.IContinueLogin,
				"friend":       status.IAddFriend,
				"word":         status.ISmallWorld,
			}
			if time.Now().Day() > 8 {
				speedInfo["sport"] = "1"
			}
			for k := range Setting {
				if speedInfo[k] == "" {
					speedInfo[k] = "1"
				}
			}
			_ = uin_info.GetIns().Save(info, speedInfo)
		}
		if isSync {
			go fun(status, uin)
		} else {
			fun(status, uin)
		}
	}
	return status, err, isFail
}

// GetUinCookies 获取QQ cookies
func GetUinCookies(addr string, proxy string, uin string, domain string, refresh bool) (err error, skey string, pSkey string) {
	key := constant.LockUinCookie + uin + "_" + domain
	//上锁
	if !util.Cache.Lock(key, constant.LockTimeUinCookie) {
		return errors.New("系统繁忙,请稍后重试"), "", ""
	}
	defer util.Cache.UnLock(key) //释放并发锁
	//查询数据库是否存在此cookie
	cookie, err := cookie2.GetIns().FindByUinDomainWithCache(uin, domain, refresh)
	if err != nil || cookie == nil {
		//获取并写入cookie
		s := uin_tool.NewUinTool(addr, proxy)
		skey, pSkey = s.Cookies(uin, domain)
		if skey == "" || pSkey == "" {
			return errors.New("获取cookie失败"), "", ""
		}
		//写入数据库
		if cookie2.GetIns().Create(&model.AppUinCookie{
			Uin:    uin,
			Domain: domain,
			Cookie: cookie2.GetIns().CookieToStr(&cookie2.Cookie{
				Uin:   uin,
				Skey:  skey,
				Pskey: pSkey,
			}),
			Status: cookie2.Normal,
		}) != nil {
			return errors.New("获取cookie失败"), "", ""
		}
		return nil, skey, pSkey
	} else {
		//查看cookie是否到时间需要更新
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", cookie.CheckTime.String(), time.Local)
		//小于三小时或者cookie失效了则更新cookie
		if cookie.Status != cookie2.Normal || stamp.Unix() < time.Now().Unix()-3600*3 {
			s := uin_tool.NewUinTool(addr, proxy)
			skey, pSkey = s.Cookies(uin, domain)
			if skey == "" || pSkey == "" {
				return errors.New("获取cookie失败"), "", ""
			}
			//更新数据库cookie
			if cookie2.GetIns().UpdateByUinDomain(uin, domain, map[string]interface{}{
				"cookie": cookie2.GetIns().CookieToStr(&cookie2.Cookie{
					Uin:   uin,
					Skey:  skey,
					Pskey: pSkey,
				}),
				"status":     cookie2.Normal,
				"check_time": str.DateTime(time.Now()),
			}) == nil {
				return nil, skey, pSkey
			} else {
				return errors.New("获取cookie失败"), "", ""
			}
		}
		t := cookie2.GetIns().StrToCookie(cookie.Cookie)
		return nil, t.Skey, t.Pskey
	}
}
