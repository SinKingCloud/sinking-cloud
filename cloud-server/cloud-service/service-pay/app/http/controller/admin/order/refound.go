package order

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Refound 订单退款
func (Controller) Refound(c *context.Context) {
	var form struct {
		Id    int64   `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
		Money float64 `json:"money" default:"0" validate:"omitempty,min=0" label:"退款金额"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Order.Refound(form.Id, form.Money, c.GetUserInfo().Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("订单退款成功")
}
