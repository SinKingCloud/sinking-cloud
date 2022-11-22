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
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin_info"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/oicq"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	"strconv"
	"time"
)

// ServerUinDelete 挂机服务器删除账号
func ServerUinDelete(u int) error {
	uinInfo, err := uin.GetIns().FindById(u)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	err = uin.GetIns().DeleteById(uinInfo.Id)
	if err != nil {
		return err
	}
	q, err := strconv.ParseInt(uinInfo.Uin, 10, 64)
	if err == nil {
		oicq.DelOicqInstance(q) //删除实例
	}
	go func(uinInfo *model.AppUin) {
		//写入QQ操作记录
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  uinInfo.UserId,
			AppId:   uinInfo.AppId,
			UinId:   uinInfo.Id,
			Uin:     uinInfo.Uin,
			Type:    log.Delete,
			Title:   "删除挂机",
			Content: "删除挂机账号",
		})
	}(uinInfo)
	return nil
}

// ChangeSetting 更改设置
func ChangeSetting(uinId int, newSetting map[JobName]string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	var set map[JobName]string
	err = json.Unmarshal([]byte(uinInfo.Setting), &set)
	if err != nil {
		set = make(map[JobName]string)
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
	if pwd != uinInfo.Password {
		data := make(map[string]interface{})
		data["password"] = pwd
		data["uin_status"] = uin.StatusUinError
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
		err = uin.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			return errors.New("更改挂机密码失败")
		}
	}
	return nil
}

// SaveLoginToken 保存登陆token
func SaveLoginToken(uinInfo *model.AppUin, token []byte) error {
	data := make(map[string]interface{})
	data["server_id"] = uinInfo.ServerId
	data["server_key"] = uinInfo.ServerKey
	data["timing"] = uinInfo.Timing
	data["proto"] = uinInfo.Proto
	data["token"] = utils.Bin2HexTo(token)
	data["uin_status"] = int(uin.StatusUinNormal)
	data["password"] = uinInfo.Password
	data["device"] = uinInfo.Device
	err := uin.GetIns().UpdateByUin(uinInfo.Uin, data)
	if err != nil {
		return errors.New("更改挂机密码失败")
	}
	return nil
}

// GetUinLevelInfo 获取等级信息
func GetUinLevelInfo(uinId int) (*model.AppUinInfo, error) {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return nil, errors.New("未查询到QQ信息")
	}
	uinLevelInfo, e := uin_info.GetIns().FindByUin(uinInfo.Uin)
	//失效时刷新
	refresh := false
	//重置任务状态
	reset := false
	if (e != nil || uinLevelInfo == nil) && uinInfo.UinStatus == uin.Normal {
		refresh = true
	}
	if uinLevelInfo == nil {
		reset = true
	}
	//上次刷新在昨天时或者大于120分钟时刷新一次数据
	if !refresh && e == nil && uinLevelInfo != nil {
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", uinLevelInfo.UpdateTime.String(), time.Local)
		todayUnix := str.NewTimeTool().GetZeroTime(time.Now()).Unix()
		if stamp.Unix() < todayUnix || stamp.Unix() <= time.Now().Unix()-7200 {
			refresh = true
		}
		if stamp.Unix() < todayUnix {
			reset = true
		}
	}
	//刷新QQ信息
	if refresh {
		_, _ = GetUinInfo(uinInfo.Uin)
	}
	//刷新任务信息
	if reset {
		var jobs []JobName
		for k := range Setting {
			jobs = append(jobs, k)
		}
		_ = SetUinJobStatus(uinInfo.Uin, jobs, 0)
	}
	//重新获取信息
	if refresh || reset {
		uinLevelInfo, e = uin_info.GetIns().FindByUin(uinInfo.Uin)
	}
	return uinLevelInfo, e
}

// GetUinInfo 获取QQ详情
func GetUinInfo(uin string) (status *client.SpeedStatus, err error) {
	cli, err := GetOicqIns(uin, false)
	if err != nil {
		return nil, err
	}
	err, status = cli.Mobile.Speed.Info()
	if err == nil && status != nil {
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
		_ = uin_info.GetIns().Save(info, nil)
	}
	return status, err
}

// SetUinSpeedStatus 设置加速任务状态
func SetUinSpeedStatus(uin string, status *client.SpeedStatus) error {
	if status == nil {
		return errors.New("加速状态不能为空")
	}
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
		"dress":        status.IDress,
	}
	return uin_info.GetIns().Save(info, speedInfo)
}

// SetUinJobStatus 设置任务状态
// status 0:未加速/1:已加速/2:已提交/3:已过期
func SetUinJobStatus(uin string, jobs []JobName, status int) error {
	info := &model.AppUinInfo{
		Uin: uin,
	}
	//加速状态
	speedInfo := make(map[string]string)
	for _, v := range jobs {
		speedInfo[string(v)] = strconv.Itoa(status)
	}
	return uin_info.GetIns().Save(info, speedInfo)
}

// GetUinCookies 获取QQ cookies
func GetUinCookies(uin string, domain string, refresh bool) (err error, cookies string) {
	key := constant.LockUinCookie + uin + "_" + domain
	//上锁
	if !util.Cache.Lock(key, constant.LockTimeUinCookie) {
		return errors.New("系统繁忙,请稍后重试"), ""
	}
	defer util.Cache.UnLock(key) //释放并发锁
	//查询数据库是否存在此cookie
	cookie, err := cookie2.GetIns().FindByUinDomainWithCache(uin, domain, refresh)
	if err != nil || cookie == nil {
		return errors.New("获取cookie失败"), ""
	} else {
		return nil, cookie.Cookie
	}
}

// SetUinCookies 设置QQ cookies
func SetUinCookies(uin string, domain string, cookies string) error {
	key := constant.LockUinCookie + uin + "_" + domain
	//上锁
	if !util.Cache.Lock(key, constant.LockTimeUinCookie) {
		return errors.New("系统繁忙,请稍后重试")
	}
	defer util.Cache.UnLock(key) //释放并发锁
	//查询数据库是否存在此cookie
	cookie, err := cookie2.GetIns().FindByUinDomainWithCache(uin, domain, false)
	if err != nil || cookie == nil {
		//写入cookie
		if cookie2.GetIns().Create(&model.AppUinCookie{
			Uin:    uin,
			Domain: domain,
			Cookie: cookies,
			Status: cookie2.Normal,
		}) != nil {
			return errors.New("设置cookie失败")
		}
	} else {
		//更新cookie
		if cookie2.GetIns().UpdateByUinDomain(uin, domain, map[string]interface{}{
			"cookie":     cookies,
			"status":     cookie2.Normal,
			"check_time": str.DateTime(time.Now()),
		}) != nil {
			return errors.New("设置cookie失败")
		}
	}
	//删除缓存
	_ = cookie2.GetIns().ClearByUinDomain(uin, domain)
	return nil
}
