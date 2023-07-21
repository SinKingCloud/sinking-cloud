package face

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/aliyun"
)

type ControllerFace struct {
}

func (ControllerFace) Check(c *sinking_server.Context) {
	data := c.GetRpcData()
	res, userId, cardNo, cardName := aliyun.FaceCheck(data["id"].(string))
	type Response struct {
		UserId   int    `json:"user_id"`
		Result   bool   `json:"result"`
		CardNo   string `json:"card_no"`
		CardName string `json:"card_name"`
	}
	c.SuccessWithData("success", &Response{
		UserId:   userId,
		Result:   res,
		CardNo:   cardNo,
		CardName: cardName,
	})
}
