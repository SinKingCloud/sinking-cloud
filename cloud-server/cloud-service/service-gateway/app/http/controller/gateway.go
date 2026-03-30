package controller

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Gateway 网关转发
func Gateway(c *context.Context) {
	parts := strings.SplitN(strings.TrimPrefix(c.Request.URL.Path, "/"), "/", 2)
	var serviceName, targetPath string
	switch len(parts) {
	case 1:
		serviceName += parts[0]
	case 2:
		serviceName += parts[0]
		targetPath = parts[1]
	}
	serviceName = strings.ToLower(serviceName)
	if serviceName == "" || serviceName == constant.ServiceGateway {
		c.Error("获取服务节点失败: 服务不存在")
		c.Abort()
		return
	}
	clientIp := c.ClientIP(true)
	if clientIp == "" {
		clientIp = c.ClientIP(false)
	}
	node, err := global.App.Consul.GetService(serviceName, client.Hash, clientIp)
	if err != nil {
		c.Error("获取服务节点失败: " + err.Error())
		c.Abort()
		return
	}
	address := node.Address
	var url, tempPath string
	tempPath = c.Request.URL.Path
	if strings.ToLower(c.Request.Header.Get("Connection")) == "upgrade" && strings.ToLower(c.Request.Header.Get("Upgrade")) == "websocket" {
		url = "ws://" + address
	} else {
		url = "http://" + address
	}
	err = c.Proxy(url, nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
		r.URL.Path, r.URL.RawPath, r.RequestURI = targetPath, targetPath, targetPath
		r.Header.Set("X-Forwarded-For", clientIp)
		r.Header.Set("X-Real-Ip", clientIp)
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Header.Set(constant.RpcTokenName, global.App.Config.GetString(constant.ConsulToken))
		r.Header.Set(constant.RpcOriginService, constant.ServiceGateway)
		r.Header.Set(constant.RpcOriginAddress, global.App.Consul.GetAddress())
		r.Header.Set(constant.RequestId, c.GetRequestId())
	}, func(writer http.ResponseWriter, request *http.Request, err error) {
		c.Request = request
		c.Writer = writer
		c.Error("获取服务节点失败: 请求服务失败")
	})
	c.Request.URL.Path = tempPath
	if err != nil {
		c.Error("获取服务节点失败: 转发服务失败")
	}
}
