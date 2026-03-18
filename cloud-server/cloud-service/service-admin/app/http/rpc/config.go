package rpc

import (
	"encoding/json"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
)

type configService struct {
}

// GetUserConfig 获取用户配置
func (*configService) GetUserConfig(params json.RawMessage) (interface{}, error) {
	var param *config.GetSystemUserConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.GetUserConfigs(param.UserId, param.Group, param.Refresh), nil
}

// SetUserConfig 设置用户配置
func (*configService) SetUserConfig(params json.RawMessage) (interface{}, error) {
	var param *config.SetSystemUserConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.SetUserConfig(param.UserId, param.Group, param.Key, param.Value), nil
}

// SetUserConfigs 批量设置用户配置
func (*configService) SetUserConfigs(params json.RawMessage) (interface{}, error) {
	var param *config.SetSystemUserConfigsParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	service.Config.SetUserConfigs(param.UserId, param.Group, param.Data)
	return true, nil
}

// GetSystemConfig 获取系统配置
func (*configService) GetSystemConfig(params json.RawMessage) (interface{}, error) {
	var param *config.GetSystemConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.GetSysConfigs(param.Group, param.Refresh), nil
}

// SetSystemConfig 设置系统配置
func (*configService) SetSystemConfig(params json.RawMessage) (interface{}, error) {
	var param *config.SetSystemConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.SetSysConfig(param.Key, param.Value), nil
}

// SetSystemConfigs 批量设置系统配置
func (*configService) SetSystemConfigs(params json.RawMessage) (interface{}, error) {
	var param *config.SetSystemConfigsParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	service.Config.SetSysConfigs(param.Data)
	return true, nil
}
