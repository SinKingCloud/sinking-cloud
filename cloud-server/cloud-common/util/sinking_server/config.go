package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
)

// set 设置内容
func (c *Context) set(module string, config map[string]string) {
	c.Set(module, config)
}

// get 获取内容
func (c *Context) get(module string, key string) interface{} {
	config, exists := c.Get(module)
	if exists {
		configMap := config.(map[string]string)
		return configMap[key]
	}
	return nil
}

// SetSystemConfig 设置系统字典配置信息
func (c *Context) SetSystemConfig(config map[string]string) {
	c.set(util.SystemConfig, config)
}

// GetSystemConfig 获取系统字典配置信息
func (c *Context) GetSystemConfig(key string) interface{} {
	return c.get(util.SystemConfig, key)
}

// GetSystemConfigWithDefault 获取系统字典配置信息()
func (c *Context) GetSystemConfigWithDefault(key string) interface{} {
	return c.get(util.SystemConfig, key)
}

// SetWebUserConfig 设置网站用户配置信息
func (c *Context) SetWebUserConfig(config map[string]string) {
	c.set(util.WebUserConfig, config)
}

// GetWebUserConfig 获取网站用户配置信息
func (c *Context) GetWebUserConfig(key string) interface{} {
	return c.get(util.WebUserConfig, key)
}

// SetSystemUserConfig 设置系统用户配置信息
func (c *Context) SetSystemUserConfig(config map[string]string) {
	c.set(util.SystemUserConfig, config)
}

// GetSystemUserConfig 获取系统用户配置信息
func (c *Context) GetSystemUserConfig(key string) interface{} {
	return c.get(util.SystemUserConfig, key)
}
