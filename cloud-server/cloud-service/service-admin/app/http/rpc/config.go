package rpc

import (
	"encoding/json"
	"service-admin/app/service"
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

// GetSystemConfig 获取系统配置
func (config) GetSystemConfig(params json.RawMessage) (interface{}, error) {
	var param struct {
		Group   string
		Refresh bool
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.GetSysConfigs(param.Group, param.Refresh), nil
}

// SetSystemConfig 设置系统配置
func (config) SetSystemConfig(params json.RawMessage) (interface{}, error) {
	var param struct {
		Key   string
		Value string
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.SetSysConfig(param.Key, param.Value), nil
}

// SetSystemConfigs 批量设置系统配置
func (config) SetSystemConfigs(params json.RawMessage) (interface{}, error) {
	var param struct {
		Data map[string]string
	}
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	service.Config.SetSysConfigs(param.Data)
	return true, nil
}
