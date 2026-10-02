package order

import (
	"fmt"
	"service-pay/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
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
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventDelete, "清理过期未支付订单", fmt.Sprintf("清理已过期%d天以上的未支付订单，共%d条", form.Day, total))
	c.Success(fmt.Sprintf("成功清理%d条订单", total))
}
