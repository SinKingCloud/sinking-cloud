package enum

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// Data 获取数据
func (*Service) Data(service string, key string) (result map[string]interface{}) {
	err := rpc.InvokeFunc(service, constant.RpcGetEnumDataFunName, key, &result)
	if err != nil {
		return nil
	}
	return result
}
