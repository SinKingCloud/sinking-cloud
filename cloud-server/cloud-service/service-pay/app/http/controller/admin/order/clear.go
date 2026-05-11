package order

import (
	"fmt"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Clear 清理订单
func (Controller) Clear(c *context.Context) {
	var form struct {
		Day int `json:"day" default:"" validate:"numeric,min=1" label:"保留天数"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	total, err := service.Order.CleanExpiredUnPayOrder(form.Day)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success(fmt.Sprintf("成功清理%d条订单", total))
}
