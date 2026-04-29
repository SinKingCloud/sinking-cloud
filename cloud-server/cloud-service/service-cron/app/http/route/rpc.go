package route

import (
	"service-cron/app/enum"
	rpc2 "service-cron/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc(constant.RpcGetEnumDataFunName, service.Enum.GetDataHandler(enum.Data)) //获取枚举配置
	rpc.RegisterFunc("RegisterRpcTask", rpc2.Task.RegisterRpc)                               //注册rpc定时任务
	rpc.RegisterFunc("RegisterScriptTask", rpc2.Task.RegisterScript)                         //注册script定时任务
}
