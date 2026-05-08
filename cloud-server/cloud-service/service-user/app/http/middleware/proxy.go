package middleware

import (
	"net/http"
	"net/http/httputil"
	"strings"

	constant2 "service-user/app/constant"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Proxy 转发到子服务
func Proxy(c *context.Context) {
	if rpc.IsRpcPath(c.Request.URL.Path) {
		c.Next()
		return
	}
	c.Next()
	if !c.GetBool(constant2.NeedProxy) {
		return
	}
	parts := strings.SplitN(strings.TrimPrefix(c.Request.URL.Path, "/"+constant.ServiceUser+"/"), "/", 2)
	var serviceName, targetPath string
	switch len(parts) {
	case 1:
		serviceName += parts[0]
	case 2:
		serviceName += parts[0]
		targetPath = "/" + parts[1]
	}
	serviceName = strings.ToLower(serviceName)
	targetPath = "/" + constant.ServiceUser + targetPath
	if serviceName == "" || serviceName == constant.ServiceUser {
		c.Error("请求资源不存在")
		c.Abort()
		return
	}
	node, err := global.App.Consul.GetService(serviceName, client.Hash, c.ClientIP(true))
	if err != nil {
		c.Error("请求资源不存在")
		c.Abort()
		return
	}
	var url string
	if strings.ToLower(c.Request.Header.Get("Connection")) == "upgrade" && strings.ToLower(c.Request.Header.Get("Upgrade")) == "websocket" {
		url = "ws://" + node.Address
	} else {
		url = "http://" + node.Address
	}
	_ = c.Proxy(url, nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
		r.URL.Path, r.URL.RawPath, r.RequestURI = targetPath, targetPath, targetPath
		r.Header.Set(constant.RpcTokenName, global.App.Config.GetString(constant.ConsulToken))
		r.Header.Set(constant.RpcOriginService, constant.ServiceUser)
		r.Header.Set(constant.RpcOriginAddress, global.App.Consul.GetAddress())
	}, nil)
}
