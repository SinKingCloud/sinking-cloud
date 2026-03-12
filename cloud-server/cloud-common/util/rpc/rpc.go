package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// RegisterHandle 挂载RPC处理器到HTTP引擎
func RegisterHandle(engine *sinking_web.Engine) {
	engine.ANY(constant.RpcRequestPath, func(c *sinking_web.Context) {
		global.Consul.RpcHandle().ServeHTTP(c.Writer, c.Request)
	})
}

// RegisterFunc 注册RPC方法
func RegisterFunc(action string, handler client.RpcHandlerFunc) {
	global.Consul.RpcRegister(action, handler)
}

// InvokeFunc 调用远程RPC方法
func InvokeFunc(serviceName string, action string, params interface{}, result interface{}, options ...interface{}) error {
	return global.Consul.RpcCall(serviceName, action, params, result, append([]interface{}{constant.RpcRequestPath}, options...)...)
}
