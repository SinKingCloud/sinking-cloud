package face

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/aliyun"
)

type ControllerFace struct {
}

func (ControllerFace) Check(c *sinking_server.Context) {
	data := c.GetRpcData()
	c.SuccessWithData("success", aliyun.FaceCheck(data["id"].(string)))
}
