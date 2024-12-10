package pay

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/order"
)

// Submit 支付跳转
func (p ControllerPay) Submit(c *sinking_server.Context) {
	type Form struct {
		TradeNo string `json:"trade_no" default:"" validate:"required,numeric" label:"订单号"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	userOrder, err := service.Order.FindByTradeNo(form.TradeNo)
	if err != nil || userOrder == nil || userOrder.Status != order.StatusUnPay {
		c.Error("未查询到有效订单")
		return
	}
	//构建支付请求
	var url string
	var err2 error
	switch userOrder.PayType {
	case order.PayTypeAliPay:
		url, err2 = p.Alipay(userOrder.TradeNo, fmt.Sprintf("%.2f", userOrder.Money), userOrder.Name)
		break
	default:
		err = errors.New("暂不支持该支付接口")
		break
	}
	if err2 != nil {
		c.Error(err2.Error())
		return
	}
	c.SuccessWithData("创建支付链接成功", url)
}
