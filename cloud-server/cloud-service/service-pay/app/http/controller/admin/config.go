package admin

import (
	"service-pay/app/service"
	"service-pay/app/service/channel"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type configController struct {
}

// Get 获取配置
func (configController) Get(c *context.Context) {
	c.SuccessWithData("获取数据成功", service.Channel.GetConfig(false))
}

// Set 更新配置
func (configController) Set(c *context.Context) {
	var form struct {
		Min    float64 `json:"min" default:"0" validate:"min=0" label:"最小充值金额(0为不限制)"`
		Max    float64 `json:"max" default:"0" validate:"min=0" label:"最大充值金额(0为不限制)"`
		Alipay int64   `json:"alipay" default:"0" validate:"min=0" label:"支付宝支付通道(0为关闭)"`
		Wxpay  int64   `json:"wxpay" default:"0" validate:"min=0" label:"微信支付通道(0为关闭)"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.SetConfig(&channel.Config{
		Min:    form.Min,
		Max:    form.Max,
		Alipay: form.Alipay,
		Wxpay:  form.Wxpay,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新配置成功")
}
