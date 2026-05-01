package order

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Deliver 订单发货
func (Controller) Deliver(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Order.Deliver(form.Id, c.GetUserInfo().Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("订单发货成功")
}
