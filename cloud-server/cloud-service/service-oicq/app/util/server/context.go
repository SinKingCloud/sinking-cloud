package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/user_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
)

// GetAppInfo 获取app信息
func GetAppInfo(c *sinking_server.Context) *user_app.UserApp {
	data, exists := c.Get(constant.AppInfo)
	if exists {
		return data.(*user_app.UserApp)
	} else {
		return nil
	}
}
