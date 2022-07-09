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
		if speedInfo != nil {
			bytes, e := json.Marshal(speedInfo)
			if e == nil {
				t := string(bytes)
				d["speed_info"] = t
				data.SpeedInfo = t
			}
		}
		if s.CountByUin(data.Uin) > 0 {
			err = s.UpdateByUin(data.Uin, d)
		} else {
			data.UpdateTime = str.DateTime(time.Now())
			err = s.Create(data)
		}
	} else {
		err = errors.New("系统繁忙")
	}
	return
}
