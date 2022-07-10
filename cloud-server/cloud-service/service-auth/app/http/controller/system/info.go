package system

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Info 网站详情
func (ControllerSystem) Info(c *sinking_server.Context) {
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.SystemGroup, false)
	c.SuccessWithData("获取成功", sinking_web.H{
		"name": configs[config.SystemAdminName],
	})
}
