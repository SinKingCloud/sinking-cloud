package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
)

// SetSystemConfig 设置系统字典配置信息
func (c *Context) SetSystemConfig(config map[string]string) {
	c.Set(util.SystemConfig, config)
}

// GetSystemConfig 获取系统字典配置信息
func (c *Context) GetSystemConfig(key string) interface{} {
	config, exists := c.Get(util.SystemConfig)
	if exists {
		configMap := config.(map[string]string)
		return configMap[key]
	}
	return nil
}
