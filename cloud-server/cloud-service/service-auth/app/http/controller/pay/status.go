package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/order"
)

// Status 支付状态查询
func (p ControllerPay) Status(c *sinking_server.Context) {
	type Form struct {
		TradeNo string `form:"trade_no" json:"trade_no" default:"" validate:"required,numeric" label:"订单号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	userOrder, err := order.GetIns().FindByTradeNo(form.TradeNo)
	if err != nil || userOrder == nil {
		c.Error("未查询到有效订单")
		return
	}
	c.SuccessWithData("获取订单支付状态成功", userOrder.Status)
}
