package dashboard

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Service 服务数量
func (Controller) Service(c *context.Context) {
	data := service.Dashboard.GetServiceCount()
	c.SuccessWithData("获取数据成功", data)
}
