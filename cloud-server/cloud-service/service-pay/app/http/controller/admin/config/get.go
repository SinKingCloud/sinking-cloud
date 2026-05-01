package config

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Get 获取配置
func (Controller) Get(c *context.Context) {
	c.SuccessWithData("获取数据成功", service.Channel.GetConfig(false))
}
