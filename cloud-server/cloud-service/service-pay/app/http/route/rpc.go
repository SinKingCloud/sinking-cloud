package route

import (
	"service-pay/app/enum"
	rpc2 "service-pay/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc(constant.RpcGetEnumDataFunName, service.Enum.GetDataHandler(enum.Data)) //获取枚举配置
	rpc.RegisterFunc("RechargeOrderNotify", rpc2.Balance.RechargeOrderNotify)                //充值订单回调
	rpc.RegisterFunc("CleanExpiredUnPayOrder", rpc2.Order.CleanExpiredUnPayOrder)            //清理过期未支付订单
	rpc.RegisterFunc("CreateOrder", rpc2.Order.CreateOrder)                                  //生成订单
}
