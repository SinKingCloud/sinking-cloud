package index

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerIndex) Index(c *sinking_server.Context) {
	c.SuccessWithData("请求成功", c.GetWebUserInfo())
}
