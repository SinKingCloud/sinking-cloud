package middleware

import (
	"net/http"
	"net/http/httputil"
	constant2 "service-admin/app/constant"
	"service-admin/app/service"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Rbac rbac权限拦截
func Rbac(c *context.Context) {
	//RPC放行
	path := c.Request.URL.Path
	if rpc.IsRpcPath(path) {
		c.Next()
		return
	}
	//权限拦截
	user := c.GetUserInfo()
	if user != nil && user.Id > 0 {
		err := service.Rbac.CheckUserHasPermission(user.Id, "/"+constant.ServiceAdmin+path)
		if err != nil {
			c.NotAllow("暂无权限访问此内容")
			return
		}
	}
	//执行控制器
	c.Next()
	if !c.GetBool(constant2.NeedProxy) {
		return
	}
	//执行反向代理
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	var serviceName, targetPath string
	switch len(parts) {
	case 1:
		serviceName += parts[0]
	case 2:
		serviceName += parts[0]
		targetPath = "/" + parts[1]
	}
	serviceName = strings.ToLower(serviceName)
	targetPath = "/" + constant.ServiceAdmin + targetPath
	if serviceName == "" || serviceName == constant.ServiceAdmin {
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
		r.Header.Set(constant.RpcOriginService, constant.ServiceAdmin)
		r.Header.Set(constant.RpcOriginAddress, global.App.Consul.GetAddress())
	}, nil)
}
