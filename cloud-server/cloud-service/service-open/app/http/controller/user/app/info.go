package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/server"
)

// Info 应用信息
func (ControllerApp) Info(c *sinking_server.Context) {
	c.SuccessWithData("获取数据成功", server.GetAppInfo(c))
}
