package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// GetRequestId 获取请求id
func (c *Context) GetRequestId() string {
	id := c.Request.Header.Get(util.RequestId)
	if id == "" {
		id = str.GetSnowWorkIns().GetUuid()
		c.Set(util.RequestId, id)
	}
	return id
}

// GetRequestIp 获取请求ip
func (c *Context) GetRequestIp() string {
	ip := c.ClientIP(true)
	if ip == "" {
		ip = c.ClientIP(false)
	}
	return ip
}
