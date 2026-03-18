package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

type Object struct {
}

// GetSystemUserConfig 获取系统用户配置
func (*Object) GetSystemUserConfig(userId int64, group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "GetUserConfig", &GetSystemUserConfigParam{
		UserId:  userId,
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return nil
	}
	return config
}

// SetSystemUserConfig 设置用户配置
func (*Object) SetSystemUserConfig(userId int64, group string, key string, value string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetUserConfig", &SetSystemUserConfigParam{
		UserId: userId,
		Group:  group,
		Key:    key,
		Value:  value,
	}, nil)
	return err == nil
}

// SetSystemUserConfigs 批量设置用户配置
func (*Object) SetSystemUserConfigs(userId int64, group string, data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetUserConfigs", &SetSystemUserConfigsParam{
		UserId: userId,
		Group:  group,
		Data:   data,
	}, nil)
	return err == nil
}

// GetSystemConfig 获取系统配置
func (*Object) GetSystemConfig(group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "GetSystemConfig", &GetSystemConfigParam{
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return nil
	}
	return config
}

// SetSystemConfig 设置系统配置
func (*Object) SetSystemConfig(key string, value string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetSystemConfig", &SetSystemConfigParam{
		Key:   key,
		Value: value,
	}, nil)
	return err == nil
}

// SetSystemConfigs 批量设置系统配置
func (*Object) SetSystemConfigs(data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetSystemConfigs", &SetSystemConfigsParam{
		Data: data,
	}, nil)
	return err == nil
}

// GetWebUserConfig 获取系统用户配置
func (*Object) GetWebUserConfig(userId int64, group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceUser, "GetUserConfig", &GetWebUserConfigParam{
		UserId:  userId,
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return nil
	}
	return config
}

// SetWebUserConfig 设置用户配置
func (*Object) SetWebUserConfig(userId int64, group string, key string, value string) bool {
	err := rpc.InvokeFunc(constant.ServiceUser, "SetUserConfig", &SetWebUserConfigParam{
		UserId: userId,
		Group:  group,
		Key:    key,
		Value:  value,
	}, nil)
	return err == nil
}

// SetWebUserConfigs 批量设置用户配置
func (*Object) SetWebUserConfigs(userId int64, group string, data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceUser, "SetUserConfigs", &SetWebUserConfigsParam{
		UserId: userId,
		Group:  group,
		Data:   data,
	}, nil)
	return err == nil
}
