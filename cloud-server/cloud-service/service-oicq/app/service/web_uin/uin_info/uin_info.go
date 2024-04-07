package uin_info

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"strconv"
	"sync"
	"time"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// Save 更新或插入
func (s *Service) Save(data *model.AppUinInfo, speedInfo map[string]string) (err error) {
	key := constant.LockUinInfo + data.Uin
	if util.Cache.Lock(key, constant.LockTimeUinInfo) {
		defer util.Cache.UnLock(key)
		d := map[string]interface{}{}
		if data.NickName != "" {
			d["nick_name"] = data.NickName
		}
		if data.ActivateDay != "" {
			d["activate_day"] = data.ActivateDay
		}
		if data.Level > 0 {
			d["level"] = data.Level
		}
		//获取QQ信息
		info, e := s.FindByUin(data.Uin)
		if e != nil || info == nil {
			if data.TodayDay != "" {
				d["today_day"] = data.TodayDay
			}
			if speedInfo != nil {
				bytes, e1 := json.Marshal(speedInfo)
				if e1 == nil {
					t := string(bytes)
					d["speed_info"] = t
					data.SpeedInfo = t
				}
			}
			//不存在
			data.UpdateTime = str.DateTime(time.Now())
			err = s.Create(data)
		} else {
			var job map[string]string
			//存在
			stamp, e2 := time.ParseInLocation("2006-01-02 15:04:05", info.UpdateTime.String(), time.Local)
			refresh := false
			if e2 == nil && stamp.Unix() > str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
				if data.TodayDay != "" {
					todayDay1, err1 := strconv.ParseFloat(info.TodayDay, 32)
					todayDay2, err2 := strconv.ParseFloat(data.TodayDay, 32)
					if err1 != nil || todayDay1 < 0 || (err2 == nil && todayDay1 <= todayDay2) {
						d["today_day"] = data.TodayDay
					}
				}
				refresh = true
			} else {
				d["today_day"] = data.TodayDay
			}
			if speedInfo != nil {
				if refresh {
					status := make(map[string]string)
					_ = json.Unmarshal([]byte(info.SpeedInfo), &status)
					if status != nil {
						for k := range speedInfo {
							status[k] = speedInfo[k]
						}
					}
					job = status
				} else {
					job = speedInfo
				}
			}
			if job != nil {
				bytes, e1 := json.Marshal(job)
				if e1 == nil {
					d["speed_info"] = string(bytes)
				}
			}
			err = s.UpdateByUin(data.Uin, d)
		}
	} else {
		err = errors.New("系统繁忙")
	}
	return
}
