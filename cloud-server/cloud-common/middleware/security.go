package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// InterceptIllegalRequest 拦截非法请求
func InterceptIllegalRequest(c *context.Context) {
	//RPC 请求直接放行
	if c.Request.URL.Path == "/"+constant.RpcRequestPath {
		c.Next()
		return
	}
	//其他请求需要验证 token
	token := c.Request.Header.Get(constant.RpcTokenName)
	if token == "" || token != global.Config.GetString(constant.ConsulToken) {
		c.Abort()
		c.NotAllow("非法请求", nil)
		return
	}
	c.Next()
}
