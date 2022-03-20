package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerPerson) Log(c *sinking_server.Context) {
	c.Success("获取操作记录成功")
}
