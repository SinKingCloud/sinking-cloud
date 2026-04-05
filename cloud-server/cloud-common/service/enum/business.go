package enum

import (
	"encoding/json"
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Data 获取数据
func (*Service) Data(service string, key []string) (result map[string]interface{}) {
	err := rpc.InvokeFunc(service, constant.RpcGetEnumDataFunName, key, &result)
	if err != nil {
		return nil
	}
	return result
}

// GetDataHandler 获取枚举RPC处理器
func (*Service) GetDataHandler(enumData map[string]interface{}) client.RpcHandlerFunc {
	return func(params json.RawMessage) (interface{}, error) {
		var param []string
		err := json.Unmarshal(params, &param)
		if err != nil {
			return nil, err
		}
		if len(param) == 0 {
			return nil, errors.New("枚举名称不能为空")
		}
		data := make(map[string]interface{}, len(param))
		for _, v := range param {
			value, ok := enumData[v]
			if !ok {
				continue
			}
			if fn, ok := value.(func() interface{}); ok {
				data[v] = fn()
				continue
			}
			data[v] = value
		}
		return data, nil
	}
}
