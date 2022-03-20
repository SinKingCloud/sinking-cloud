package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerPay) OrderList(c *sinking_server.Context) {
	c.Success("获取订单记录成功")
}
