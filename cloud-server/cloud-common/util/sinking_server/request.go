package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// GetRequestId 获取请求id
func (c *Context) GetRequestId() string {
	id := c.GetString(util.RequestId)
	if id == "" {
		id = str.GetSnowWorkIns().GetUuid()
		c.Set(util.RequestId, id)
	}
	return id
}
