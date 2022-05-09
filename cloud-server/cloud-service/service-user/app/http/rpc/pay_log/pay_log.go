package pay_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/pay_log"
)

type ControllerPayLog struct {
}

func (ControllerPayLog) Create(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	types := int(data["type"].(float64))
	logTypes := int(data["log_type"].(float64))
	money := data["money"].(float64)
	log := &model.WebUserPayLog{
		UserId:  userId,
		Type:    types,
		LogType: logTypes,
		Title:   data["title"].(string),
		Content: data["content"].(string),
		Money:   money,
	}
	err := pay_log.GetIns().Create(log)
	if err != nil {
		c.Error("fail")
	} else {
		c.Success("success")
	}
}
