package web_uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// GetServer 获取可用挂机server节点
func GetServer(key string) (data *model.AppUinServer, err error) {
	err = util.Database.Db.Where("`key`=? and `pid`>0 and `status`=0 ", key).Order("`uin_num` desc").First(&data).Error
	return
}
