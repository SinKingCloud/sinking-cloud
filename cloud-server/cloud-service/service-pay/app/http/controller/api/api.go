package api

import (
	"service-pay/app/enum/pay_type"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type Controller struct {
}

// Config 支付配置
func (Controller) Config(c *context.Context) {
	config := service.Channel.GetConfig(false)
	c.SuccessWithData("获取数据成功", sinking_web.H{
		"min":            config.Min,
		"max":            config.Max,
		pay_type.Balance: config.Balance > 0,
		pay_type.Alipay:  config.Alipay > 0,
		pay_type.Wxpay:   config.Wxpay > 0,
	})
}

// Submit 获取支付链接
func (Controller) Submit(c *context.Context) {

}

// Callback 支付回调
func (Controller) Callback(c *context.Context) {

}
