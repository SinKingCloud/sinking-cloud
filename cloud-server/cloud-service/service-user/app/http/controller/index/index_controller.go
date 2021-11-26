package index

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/command/queues/producer"
)

func (ControllerIndex) Index(c *sinking_server.Context) {
	producer.NewWebSocketProducer().Build("1", "这是一个测试消息").Send()
	c.SuccessWithData("请求成功", c.GetWebUserInfo())
}
