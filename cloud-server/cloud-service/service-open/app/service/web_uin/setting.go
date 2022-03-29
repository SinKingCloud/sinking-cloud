package web_uin

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
)

// UinSetting 功能设置
func UinSetting(uin *model.AppUin, key string, value string) {

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
		data["uin_status"] = 1
		go func() {
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
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
