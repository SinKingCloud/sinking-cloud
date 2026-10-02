package order

import (
	"fmt"
	"service-pay/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
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
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "订单退款至余额", fmt.Sprintf("订单[%d]退款至用户余额，金额[%.4f]", form.Id, form.Money))
	c.Success("订单退款成功")
}
