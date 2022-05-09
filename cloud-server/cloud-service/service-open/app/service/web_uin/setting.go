package web_uin

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/uin_tool"
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
