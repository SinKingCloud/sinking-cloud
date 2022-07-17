package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

type ControllerArea struct {
}

// List 获取支持的区号
func (ControllerArea) List(c *sinking_server.Context) {
	c.SuccessWithData("", str.Area.Get())
}
