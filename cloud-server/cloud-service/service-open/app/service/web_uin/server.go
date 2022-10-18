package web_uin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/oicq"
	"github.com/SinKingCloud/sinking-oicq/binary-mobile/utils"
	"gorm.io/gorm"
	"strconv"
)

// GetServer 获取可用挂机server节点
func GetServer(key string) (data *model.AppUinServer, err error) {
	err = util.Database.Db.Where("`key`= ? and `pid`>0 and `status`=0 ", key).Order("`uin_num` asc").First(&data).Error
	if err != nil || data == nil {
		return nil, errors.New("无可用节点")
	}
	return data, nil
}

// GetOicqIns 获取QQ操作实例
func GetOicqIns(u string, newIns bool) (*oicq.Oicq, error) {
	uinInfo, err := uin.GetIns().FindByUin(u)
	if err != nil || uinInfo == nil {
		return nil, errors.New("未查询到QQ信息")
	}
	severInfo, err := server2.GetIns().FindByIdWithCache(uinInfo.ServerId, false)
	if err != nil || severInfo == nil {
		return nil, errors.New("未查询到挂机节点信息")
	}
	var proxy *oicq.Proxy
	_ = json.Unmarshal([]byte(severInfo.Proxy), &proxy)
	return GetOicqInsWithOption(uinInfo, proxy, newIns, uinInfo.Device)
}

// GetOicqInsById 获取QQ操作实例
func GetOicqInsById(id int, newIns bool) (*oicq.Oicq, error) {
	uinInfo, err := uin.GetIns().FindById(id)
	if err != nil || uinInfo == nil {
		return nil, errors.New("未查询到QQ信息")
	}
	severInfo, err := server2.GetIns().FindByIdWithCache(uinInfo.ServerId, false)
	if err != nil || severInfo == nil {
		return nil, errors.New("未查询到挂机节点信息")
	}
	var proxy *oicq.Proxy
	_ = json.Unmarshal([]byte(severInfo.Proxy), &proxy)
	return GetOicqInsWithOption(uinInfo, proxy, newIns, uinInfo.Device)
}

// GetOicqInsWithOption 获取QQ操作实例
func GetOicqInsWithOption(uinInfo *model.AppUin, proxy *oicq.Proxy, newIns bool, device int) (*oicq.Oicq, error) {
	q, err := strconv.ParseInt(uinInfo.Uin, 10, 64)
	if err != nil {
		return nil, errors.New("QQ格式错误")
	}
	return oicq.GetOicqInstance(q, utils.ByteToByte16(utils.Hex2Bin(uinInfo.Password)), device, proxy, utils.Hex2Bin(uinInfo.Token), newIns), nil
}

// ChangeServer 更改挂机地区
func ChangeServer(uinId int, key string) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	if uinInfo.ServerKey != key {
		server, err := GetServer(key)
		if err != nil || server == nil {
			return errors.New("未获取到可用节点,请更换挂机节点")
		}
		data := make(map[string]interface{})
		data["server_key"] = server.Key
		data["server_id"] = server.Id
		go func() {
			//写入QQ操作记录
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Uin:     uinInfo.Uin,
				Type:    log.Server,
				Title:   "地点变动",
				Content: "挂机地点变动为" + server.Name,
			})
			_ = server2.GetIns().UpdateById(uinInfo.ServerId, map[string]interface{}{
				"uin_num": gorm.Expr("`uin_num`-1"),
			})
			_ = server2.GetIns().UpdateById(server.Id, map[string]interface{}{
				"uin_num": gorm.Expr("`uin_num`+1"),
			})
		}()
		err = uin.GetIns().UpdateById(uinInfo.Id, data)
		if err != nil {
			return errors.New("更改挂机地区失败")
		}
	}
	return nil
}

// ChangeTiming 修改挂机时间
func ChangeTiming(uinId int, timing int) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	//修改挂机时间
	err = uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{
		"timing": timing,
	})
	if err != nil {
		return errors.New("更改挂机时间失败")
	}
	go func(timing int) {
		//写入QQ操作记录
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  uinInfo.UserId,
			AppId:   uinInfo.AppId,
			UinId:   uinInfo.Id,
			Uin:     uinInfo.Uin,
			Type:    log.Timing,
			Title:   "时间变动",
			Content: "挂机时间变动为" + fmt.Sprintf("%02d:%02d", timing/60, timing%60),
		})
	}(timing)
	return nil
}

// ChangeStatus 修改挂机状态
func ChangeStatus(uinId int, status int) error {
	uinInfo, err := uin.GetIns().FindById(uinId)
	if err != nil || uinInfo == nil {
		return errors.New("未查询到QQ信息")
	}
	//修改挂机时间
	err = uin.GetIns().UpdateById(uinInfo.Id, map[string]interface{}{
		"status": status,
	})
	if err != nil {
		return errors.New("更改挂机状态失败")
	}
	go func(status int) {
		//写入QQ操作记录
		s := "正常"
		if status == uin.Stop {
			s = "暂停"
		} else if status == uin.Error {
			s = "禁止"
		}
		_ = log.GetIns().Create(&model.AppUinLog{
			UserId:  uinInfo.UserId,
			AppId:   uinInfo.AppId,
			UinId:   uinInfo.Id,
			Uin:     uinInfo.Uin,
			Type:    log.Status,
			Title:   "状态变动",
			Content: "挂机状态变动为" + s,
		})
	}(status)
	return nil
}
