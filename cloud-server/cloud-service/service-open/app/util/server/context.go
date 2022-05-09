package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// GetAppInfo 获取app信息
func GetAppInfo(c *sinking_server.Context) *model.UserApp {
	data, exists := c.Get(constant.AppInfo)
	if exists {
		return data.(*model.UserApp)
	} else {
		return nil
	}
}
