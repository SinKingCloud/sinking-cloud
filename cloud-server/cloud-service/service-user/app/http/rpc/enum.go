package rpc

import (
	"encoding/json"
	"service-user/app/enum"
)

type enumService struct {
}

// GetData 获取枚举配置
func (*enumService) GetData(params json.RawMessage) (interface{}, error) {
	var param string
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	if param == "" {
		return enum.Data, nil
	}
	return enum.Data[param], nil
}
