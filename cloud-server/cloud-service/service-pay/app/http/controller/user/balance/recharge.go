package balance

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Recharge 充值
func (Controller) Recharge(c *context.Context) {
	var form struct {
		Money int    `json:"money" default:"" validate:"min=1" label:"充值金额"`
		Type  string `json:"type" default:"" label:"充值方式"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetUserInfo()
	orderId, err := service.Balance.CreateRechargeOrder(u.Id, form.Money, form.Type)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("创建订单成功", orderId)
}
