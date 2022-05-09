package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/command/queues/producer"
)

type ControllerLog struct {
}

func (ControllerLog) Create(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	types := int(data["type"].(float64))
	log := &model.WebUserLog{
		UserId:    userId,
		RequestId: data["request_id"].(string),
		RequestIp: data["request_ip"].(string),
		Type:      types,
		Title:     data["title"].(string),
		Content:   data["content"].(string),
	}
	producer.NewWebUserLogProducer().Build(log).Send()
	c.Success("success")
}
