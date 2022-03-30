package web_uin

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/log"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_uin/uin"
	"gorm.io/gorm"
)

// GetServer 获取可用挂机server节点
func GetServer(key string) (data *model.AppUinServer, err error) {
	err = util.Database.Db.Where("`key`=? and `pid`>0 and `status`=0 ", key).Order("`uin_num` desc").First(&data).Error
	return
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
		data["uin_status"] = 1
		go func() {
			_ = ServerUinCreate(uinInfo.ServerId, uinInfo.Uin, uinInfo.Password)
			//写入QQ操作记录
			_ = log.GetIns().Create(&model.AppUinLog{
				UserId:  uinInfo.UserId,
				AppId:   uinInfo.AppId,
				UinId:   uinInfo.Id,
				Title:   "挂机地点变动",
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
