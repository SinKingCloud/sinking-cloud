package pay

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"

func (ControllerPay) Recharge(c *sinking_server.Context) {
	c.Success("下单成功")
}
