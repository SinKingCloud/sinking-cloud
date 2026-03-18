package config

import (
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// getSystemUserConfig 获取系统用户配置
func (*Service) getSystemUserConfig(userId int64, group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "GetUserConfig", &GetSystemUserConfigParam{
		UserId:  userId,
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return map[string]string{}
	}
	return config
}

// SetSystemUserConfig 设置用户配置
func (*Service) SetSystemUserConfig(userId int64, group string, key string, value string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetUserConfig", &SetSystemUserConfigParam{
		UserId: userId,
		Group:  group,
		Key:    key,
		Value:  value,
	}, &result)
	if err != nil {
		return false
	}
	// 清理缓存
	key = constant.SystemUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	global.App.Cache.Delete(key)
	return result
}

// SetSystemUserConfigs 批量设置用户配置
func (*Service) SetSystemUserConfigs(userId int64, group string, data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetUserConfigs", &SetSystemUserConfigsParam{
		UserId: userId,
		Group:  group,
		Data:   data,
	}, nil)
	if err != nil {
		return false
	}
	// 清理缓存
	key := constant.SystemUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	global.App.Cache.Delete(key)
	return true
}

// getSystemConfig 获取系统配置
func (*Service) getSystemConfig(group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "GetSystemConfig", &GetSystemConfigParam{
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return map[string]string{}
	}
	return config
}

// SetSystemConfig 设置系统配置
func (*Service) SetSystemConfig(key string, value string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetSystemConfig", &SetSystemConfigParam{
		Key:   key,
		Value: value,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// SetSystemConfigs 批量设置系统配置
func (*Service) SetSystemConfigs(data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "SetSystemConfigs", &SetSystemConfigsParam{
		Data: data,
	}, nil)
	if err != nil {
		return false
	}
	return true
}

// getWebUserConfig 获取系统用户配置
func (*Service) getWebUserConfig(userId int64, group string, refresh bool) (config map[string]string) {
	err := rpc.InvokeFunc(constant.ServiceUser, "GetUserConfig", &GetWebUserConfigParam{
		UserId:  userId,
		Group:   group,
		Refresh: refresh,
	}, &config)
	if err != nil {
		return map[string]string{}
	}
	return config
}

// SetWebUserConfig 设置用户配置
func (*Service) SetWebUserConfig(userId int64, group string, key string, value string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceUser, "SetUserConfig", &SetWebUserConfigParam{
		UserId: userId,
		Group:  group,
		Key:    key,
		Value:  value,
	}, &result)
	if err != nil {
		return false
	}
	// 清理缓存
	key = constant.WebUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	global.App.Cache.Delete(key)
	return result
}

// SetWebUserConfigs 批量设置用户配置
func (*Service) SetWebUserConfigs(userId int64, group string, data map[string]string) bool {
	err := rpc.InvokeFunc(constant.ServiceUser, "SetUserConfigs", &SetWebUserConfigsParam{
		UserId: userId,
		Group:  group,
		Data:   data,
	}, nil)
	if err != nil {
		return false
	}
	// 清理缓存
	key := constant.WebUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	global.App.Cache.Delete(key)
	return true
}
