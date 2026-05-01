package config

import (
	"service-pay/app/service"
	"service-pay/app/service/channel"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Set 更新配置
func (Controller) Set(c *context.Context) {
	var form struct {
		Min     float64 `json:"min" default:"0" validate:"min=0" label:"最小充值金额"`
		Max     float64 `json:"max" default:"0" validate:"min=0" label:"最大充值金额"`
		Balance int64   `json:"balance" default:"0" validate:"min=0" label:"余额支付通道"`
		Alipay  int64   `json:"alipay" default:"0" validate:"min=0" label:"支付宝支付通道"`
		Wxpay   int64   `json:"wxpay" default:"0" validate:"min=0" label:"微信支付通道"`
		Expire  int64   `json:"expire" default:"0" validate:"min=0" label:"订单有效期"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.SetConfig(&channel.Config{
		Min:     form.Min,
		Max:     form.Max,
		Balance: form.Balance,
		Alipay:  form.Alipay,
		Wxpay:   form.Wxpay,
		Expire:  form.Expire,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新配置成功")
}
