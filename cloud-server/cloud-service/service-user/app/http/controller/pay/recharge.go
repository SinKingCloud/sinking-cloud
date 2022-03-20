package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"strconv"
)

func (ControllerPay) Recharge(c *sinking_server.Context) {
	type Form struct {
		Money int `json:"money" default:"" validate:"required,numeric,min=1,max=99999" label:"充值金额"`
		Type  int `json:"type" default:"" validate:"oneof=0 1" label:"充值方式"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	min, _ := strconv.Atoi(configs[config.PayMinMoney])
	if form.Money < min {
		c.Error("系统限制最低充值" + strconv.Itoa(min) + "元")
		return
	}

	c.Success("下单成功")
}
