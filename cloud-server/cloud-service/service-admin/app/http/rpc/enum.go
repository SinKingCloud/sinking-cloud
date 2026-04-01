package rpc

import (
	"encoding/json"
	"errors"
	"service-admin/app/enum"
)

type enumService struct {
}

// GetData 获取枚举配置
func (*enumService) GetData(params json.RawMessage) (interface{}, error) {
	var param []string
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	if param == nil || len(param) == 0 {
		return nil, errors.New("枚举名称不能为空")
	}
	data := make(map[string]interface{}, 10)
	for _, v := range param {
		data[v] = enum.Data[v]
	}
	return data, nil
}
