package route

import (
	rpc2 "service-auth/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadRpc(s)
	loadApi(s)
}

// loadMiddleware 中间件
func loadMiddleware(s *sinking_web.Engine) {
}

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc(constant.RpcGetEnumDataFunName, rpc2.Enum.GetData) //获取枚举配置
}

// loadApi
func loadApi(s *sinking_web.Engine) {

}
