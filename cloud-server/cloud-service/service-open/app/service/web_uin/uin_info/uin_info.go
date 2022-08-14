package uin_info

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"time"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// Save 更新或插入
func (s *service) Save(data *model.AppUinInfo, speedInfo map[string]string) (err error) {
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
		if data.TodayDay != "" {
			d["today_day"] = data.TodayDay
		}
		if data.Level > 0 {
			d["level"] = data.Level
		}
		//获取QQ信息
		info, e := s.FindByUin(data.Uin)
		if e != nil || info == nil {
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
			if e2 == nil && stamp.Unix() > str.NewTimeTool().GetZeroTime(time.Now()).Unix() {
				//今日执行过的就判断是否已经加速项目
				var status map[string]string
				e1 := json.Unmarshal([]byte(info.SpeedInfo), &status)
				if e1 == nil {
					for k, v := range status {
						if v != "1" && speedInfo[k] != "" {
							status[k] = speedInfo[k]
						}
					}
				}
				job = status
			} else {
				job = speedInfo
			}
			if job != nil {
				bytes, e1 := json.Marshal(job)
				if e1 == nil {
					t := string(bytes)
					d["speed_info"] = t
					data.SpeedInfo = t
				}
			}
			err = s.UpdateByUin(data.Uin, d)
		}
	} else {
		err = errors.New("系统繁忙")
	}
	return
}
