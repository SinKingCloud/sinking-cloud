package context

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// GetRequestId 获取请求id
func (c *Context) GetRequestId() string {
	id := c.GetString(constant.RequestId)
	if id == "" {
		id = c.Request.Header.Get(constant.RequestId)
		if id == "" {
			id = str.GetSnowWorkIns().GetUuid()
			c.Set(constant.RequestId, id)
		}
	}
	return id
}

// GetRequestIp 获取请求ip
func (c *Context) GetRequestIp() string {
	ip := c.ClientIP(true)
	if ip == "" {
		ip = c.ClientIP(false)
	}
	if ip == "" {
		return "localhost"
	}
	return ip
}
