package rpc

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

var rpcRequestPath = strings.ToLower(constant.RpcRequestPath)

// RegisterHandle 挂载RPC处理器到HTTP引擎
func RegisterHandle(engine *sinking_web.Engine) {
	engine.ANY(rpcRequestPath, func(c *sinking_web.Context) {
		c.StatusCode = 200
		global.App.Consul.RpcHandle().ServeHTTP(c.Writer, c.Request)
	})
}

// RegisterFunc 注册RPC方法
func RegisterFunc(action string, handler client.RpcHandlerFunc) {
	global.App.Consul.RpcRegister(action, handler)
}

// InvokeFunc 调用远程RPC方法
func InvokeFunc(serviceName string, action string, params interface{}, result interface{}, options ...interface{}) error {
	return global.App.Consul.RpcCall(serviceName, action, params, result, append([]interface{}{rpcRequestPath}, options...)...)
}

// IsRpcPath 判断是否为rpc请求路径
func IsRpcPath(path string) bool {
	return path == "/"+rpcRequestPath
}
