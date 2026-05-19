package config

import (
	"service-pay/app/service"
	"service-pay/app/service/channel"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Pay 获取或更新支付配置
func (Controller) Pay(c *context.Context) {
	var form struct {
		Action string `json:"action" default:"get" validate:"required,oneof=get set" label:"操作类型"`
		Data   *struct {
			Min     float64 `json:"min" default:"0" validate:"min=0" label:"最小充值金额"`
			Max     float64 `json:"max" default:"0" validate:"min=0" label:"最大充值金额"`
			Balance int64   `json:"balance" default:"0" validate:"min=0" label:"余额支付通道"`
			Alipay  int64   `json:"alipay" default:"0" validate:"min=0" label:"支付宝支付通道"`
			Wxpay   int64   `json:"wxpay" default:"0" validate:"min=0" label:"微信支付通道"`
			Expire  int64   `json:"expire" default:"0" validate:"min=0" label:"订单有效期"`
		} `json:"data" label:"配置数据"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if form.Action == "set" {
		if form.Data == nil {
			c.Error("配置数据不能为空")
			return
		}
		err := service.Channel.SetConfig(&channel.Config{
			Min:     form.Data.Min,
			Max:     form.Data.Max,
			Balance: form.Data.Balance,
			Alipay:  form.Data.Alipay,
			Wxpay:   form.Data.Wxpay,
			Expire:  form.Data.Expire,
		})
		if err != nil {
			c.Error(err.Error())
			return
		}
		c.Success("更新配置成功")
		return
	}
	c.SuccessWithData("获取数据成功", service.Channel.GetConfig(false))
}
