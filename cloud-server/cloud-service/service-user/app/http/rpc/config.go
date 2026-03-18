package rpc

import (
	"encoding/json"
	"service-user/app/service"
)

type config struct {
}

// GetUserConfig 获取用户配置
func (config) GetUserConfig(params json.RawMessage) (interface{}, error) {
	var param struct {
		UserId  int64
		Group   string
		Refresh bool
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.GetUserConfigs(param.UserId, param.Group, param.Refresh), nil
}

// SetUserConfig 设置用户配置
func (config) SetUserConfig(params json.RawMessage) (interface{}, error) {
	var param struct {
		UserId int64
		Group  string
		Key    string
		Value  string
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.SetUserConfig(param.UserId, param.Group, param.Key, param.Value), nil
}

// SetUserConfigs 批量设置用户配置
func (config) SetUserConfigs(params json.RawMessage) (interface{}, error) {
	var param struct {
		UserId int64
		Group  string
		Data   map[string]string
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	service.Config.SetUserConfigs(param.UserId, param.Group, param.Data)
	return true, nil
}
