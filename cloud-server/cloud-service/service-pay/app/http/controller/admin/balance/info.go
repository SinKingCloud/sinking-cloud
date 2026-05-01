package balance

import (
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 用户余额详情
func (Controller) Info(c *context.Context) {
	var form struct {
		UserId int64 `json:"user_id" default:"" validate:"numeric,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Balance.FindByUserId(form.UserId, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}
