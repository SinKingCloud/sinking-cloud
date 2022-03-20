package pay

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"

func (ControllerPay) LogList(c *sinking_server.Context) {
	c.Success("获取消费记录成功")
}
