package middleware

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// InterceptIllegalRequest 拦截非法请求
func InterceptIllegalRequest(c *context.Context) {
	//RPC请求直接放行
	if c.Request.URL.Path == "/"+constant.RpcRequestPath {
		c.Next()
		return
	}
	//其他请求来源则禁止
	token := c.Request.Header.Get(constant.RpcTokenName)
	originService := c.Request.Header.Get(constant.RpcOriginService)
	originAddress := c.Request.Header.Get(constant.RpcOriginAddress)
	if token == "" || originService == "" || originAddress == "" || token != global.App.Config.GetString(constant.ConsulToken) {
		c.Abort()
		c.NotAllow("无权限", nil)
		return
	}
	//本地服务名称
	localService := global.App.Consul.GetName()
	//获取第一个请求路径
	parts := strings.SplitN(strings.TrimPrefix(c.Request.URL.Path, "/"), "/", 2)
	var firstPath string
	if len(parts) > 0 {
		firstPath += strings.ToLower(parts[0])
	}
	//如果是用户服务和后台服务请求自身或者网关则禁止
	if (localService == constant.ServiceUser || localService == constant.ServiceAdmin) && (firstPath == localService || firstPath == constant.ServiceGateway) {
		c.Abort()
		c.NotAllow("无权限", nil)
		return
	}
	//如果是其他服务请求用户服务和后台服务路径并且来源不是用户服务和后台服务则禁止
	if localService != constant.ServiceUser &&
		localService != constant.ServiceAdmin &&
		(firstPath == constant.ServiceUser || firstPath == constant.ServiceAdmin) &&
		(originService != constant.ServiceUser && originService != constant.ServiceAdmin) {
		c.Abort()
		c.NotAllow("无权限", nil)
		return
	}
	c.Next()
}
