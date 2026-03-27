package config

import (
	"service-admin/app/repository/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	GetUserConfigs(userId int64, group string, refresh bool) map[string]string //获取用户配置
	SetUserConfig(userId int64, group string, key string, value string) bool   //设置用户配置
	SetUserConfigs(userId int64, group string, data map[string]string)         //批量修改用户配置
	GetSysConfigs(group string, refresh bool) map[string]string                //获取系统配置
	SetSysConfig(key string, value string) bool                                //设置系统配置
	SetSysConfigs(data map[string]string)                                      //批量设置系统配置
}

// service 注入结构
type service struct {
	repositoryConfig config.Interface
	cache            cache.Interface
}

// NewService 实例化service
func NewService(repository config.Interface, cache cache.Interface) *service {
	return &service{
		repositoryConfig: repository,
		cache:            cache,
	}
}
