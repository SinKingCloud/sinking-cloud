package rpc

import (
	"encoding/json"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/rpc/config"
)

type configService struct {
}

// GetUserConfig 获取用户配置
func (*configService) GetUserConfig(params json.RawMessage) (interface{}, error) {
	var param *config.GetWebUserConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.GetUserConfigs(param.UserId, param.Group, param.Refresh), nil
}

// SetUserConfig 设置用户配置
func (*configService) SetUserConfig(params json.RawMessage) (interface{}, error) {
	var param *config.SetWebUserConfigParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	return service.Config.SetUserConfig(param.UserId, param.Group, param.Key, param.Value), nil
}

// SetUserConfigs 批量设置用户配置
func (*configService) SetUserConfigs(params json.RawMessage) (interface{}, error) {
	var param *config.SetWebUserConfigsParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	service.Config.SetUserConfigs(param.UserId, param.Group, param.Data)
	return true, nil
}
