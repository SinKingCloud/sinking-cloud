package route

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	consul "github.com/SinKingCloud/sinking-go/sinking-consul/client"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadGateway(s)
}

// loadMiddleware 中间件
func loadMiddleware(s *sinking_web.Engine) {
	s.Use(context.HandleFunc(middleware.Cors), context.HandleFunc(middleware.Log))
}

// loadGateway 网关路由
func loadGateway(s *sinking_web.Engine) {
	handle := context.HandleFunc(func(c *context.Context) {
		path := c.Request.URL.Path
		trimmedPath := strings.TrimPrefix(path, "/")
		parts := strings.SplitN(trimmedPath, "/", 2)
		var serviceName, targetPath string
		switch len(parts) {
		case 1:
			serviceName += parts[0]
		case 2:
			serviceName += parts[0]
			targetPath = parts[1]
		}
		if serviceName == "" || serviceName == global.Config.GetString(constant.ConsulName) {
			c.Error("服务名称不能为空")
			c.Abort()
			return
		}
		clientIp := c.ClientIP(true)
		if clientIp == "" {
			clientIp = c.ClientIP(false)
		}
		node, err := global.Consul.GetService(serviceName, consul.Hash, clientIp)
		if err != nil {
			c.Error("获取服务节点失败: " + err.Error())
			c.Abort()
			return
		}
		var url, tempPath string
		tempPath = c.Request.URL.Path
		if strings.ToLower(c.Request.Header.Get("Connection")) == "upgrade" && strings.ToLower(c.Request.Header.Get("Upgrade")) == "websocket" {
			url = "ws://" + node.Address
		} else {
			url = "http://" + node.Address
		}
		_ = c.Proxy(url, nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
			r.URL.Path = targetPath
			r.URL.RawPath, r.RequestURI = targetPath, targetPath
			r.Header.Set("X-Forwarded-For", clientIp)
			r.Header.Set("X-Real-Ip", clientIp)
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Header.Set(constant.RpcTokenName, global.Config.GetString(constant.ConsulToken))
			r.Header.Set(constant.RequestId, c.GetRequestId())
		}, nil)
		c.Request.URL.Path = tempPath
	})
	s.ANY("/*", handle)
	s.ANY("/", handle)
}
