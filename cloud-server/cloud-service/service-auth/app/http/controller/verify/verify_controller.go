package verify

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerVerify) Get(c *sinking_server.Context) {
	c.Success(nil)
}

func (ControllerVerify) Check(c *sinking_server.Context) {
	c.Success(nil)
}
