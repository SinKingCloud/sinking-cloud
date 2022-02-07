package index

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerIndex) Index(c *sinking_server.Context) {
	c.Success("服务正常")
}
