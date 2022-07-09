package web_uin

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin_info"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
	"strconv"
	"time"
)

// ServerUinCreate 挂机服务器创建账号
func ServerUinCreate(serverId int, uin string, pwd string, phone string, loginType int) error {
	serverInfo, err := server2.GetIns().FindById(serverId)
	if err != nil || serverInfo == nil || serverInfo.Status != server2.Normal {
		return errors.New("未查询到挂机服务器信息,请更换挂机服务器")
	}
	s := uin_tool.NewUinTool(serverInfo.Addr, serverInfo.Proxy)
	if !s.Create(uin, pwd, phone, loginType) {
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
	for k, _ := range Setting {
		if newSetting[k] != "" && set[k] != newSetting[k] {
			if newSetting[k] == "1" {
				set[k] = "1"
			} else {
				set[k] = "0"
			}
		}
	}
	str, err := json.Marshal(set)
	if err != nil {
		return errors.New("修改配置错误")
	}
	err = uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{
		"setting": string(str),
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
		_ = ServerUinCreate(uinInfo.ServerId, uinInfo.Uin, pwd, uinInfo.Phone, 0)
		err = uin.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			return errors.New("更改挂机密码失败")
		}
	}
	return nil
}

// ChangePhone 更改挂机手机
func ChangePhone(uinId int, phone string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	if phone != uinInfo.Phone || uinInfo.LoginType != 1 {
		data := make(map[string]interface{})
		data["phone"] = phone
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
		_ = ServerUinCreate(uinInfo.ServerId, uinInfo.Uin, uinInfo.Password, phone, 1)
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

// GetUinSpeedStatus 获取加速任务状态
func GetUinSpeedStatus(addr string, proxy string, uin string, isSync bool) (*uin_tool.SpeedStatus, error, bool) {
	s := uin_tool.NewUinTool(addr, proxy)
	status, err, isFail := s.Status(uin)
	if err == nil && status != nil && !isFail {
		//更新QQ信息
		fun := func(status *uin_tool.SpeedStatus, uin string) {
			levelInfo := status.CommonInfoGetQQLevelInfo
			uinInfo := status.CommonInfoGetQQInfo
			//qq信息
			level, e := strconv.Atoi(uinInfo.IQQLevel)
			if e != nil {
				level = 0
			}
			info := &model.AppUinInfo{
				Uin:         uin,
				NickName:    uinInfo.SNickName,
				Level:       level,
				ActivateDay: levelInfo.ITotalActiveDay,
				TodayDay:    levelInfo.IRealDays,
			}
			//加速状态
			speedInfo := map[string]string{
				"mobile":       levelInfo.IMobileQQOnline,
				"pc":           levelInfo.IPCQQOnline,
				"online":       levelInfo.INoHideOnline,
				"zone":         levelInfo.QzoneVisitor,
				"medal":        levelInfo.IMedal,
				"ws":           levelInfo.WeishiVideoview,
				"micro_show":   levelInfo.ICmshowActive,
				"safe":         levelInfo.IPCSafeOnline,
				"game":         levelInfo.IMobileGameOnline,
				"sport":        levelInfo.IQQSportStep,
				"mobile_clock": levelInfo.IDailySign,
				"zone_clock":   levelInfo.IQzoneState,
				"login":        levelInfo.IContinueLogin,
				"friend":       levelInfo.IAddFriend,
			}
			if time.Now().Day() > 8 {
				speedInfo["sport"] = "1"
			}
			for k, _ := range Setting {
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
