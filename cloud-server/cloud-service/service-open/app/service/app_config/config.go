package app_config

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"strconv"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// GetConfig 获取配置
func (c *service) GetConfig(key string, ttl int64, refresh bool) map[string]string {
	config := make(map[string]string)
	if refresh {
		config = c.SelectWithMap(key, "")
		data, err := json.Marshal(config)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), ttl)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = c.SelectWithMap(key, "")
			data, err := json.Marshal(config)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), ttl)
			}
		}
	}
	return config
}

// SetConfig 设置配置
func (s *service) SetConfig(identId string, key string, value string) bool {
	var total int64
	err := util.Database.Db.Where("`ident_id` = ? and `key` = ?", identId, key).Updates(&model.AppConfig{
		Value: value,
	}).Count(&total).Error
	if err != nil || total <= 0 {
		err = util.Database.Db.Create(&model.AppConfig{
			IdentId: identId,
			Key:     key,
			Value:   value,
		}).Error
	}
	return err == nil
}

// GetAppConfigs 获取app配置
func (c *service) GetAppConfigs(appId int, refresh bool) map[string]string {
	key := constant.CacheNameWithAppConfig + strconv.Itoa(appId)
	return c.GetConfig(key, constant.CacheTimeWithAppConfig, refresh)
}

// ClearAppConfig 清理app配置缓存
func (c *service) ClearAppConfig(appId int) bool {
	key := constant.CacheNameWithAppConfig + strconv.Itoa(appId)
	return util.Cache.Delete(key)
}

// SetAppConfig 设置app配置
func (c *service) SetAppConfig(appId int, key string, value string) bool {
	identId := constant.CacheNameWithAppConfig + strconv.Itoa(appId)
	return c.SetConfig(identId, key, value)
}

// SetAppConfigs 批量设置app配置
func (c *service) SetAppConfigs(appId int, data map[string]string) {
	identId := constant.CacheNameWithAppConfig + strconv.Itoa(appId)
	for k, v := range data {
		c.SetConfig(identId, k, v)
	}
}

// GetUserConfigs 获取用户配置
func (c *service) GetUserConfigs(userId int, refresh bool) map[string]string {
	key := constant.CacheNameWithUserAppConfig + strconv.Itoa(userId)
	return c.GetConfig(key, constant.CacheTimeWithUserAppConfig, refresh)
}

// ClearUserConfig 清理用户配置缓存
func (c *service) ClearUserConfig(userId int) bool {
	key := constant.CacheNameWithUserAppConfig + strconv.Itoa(userId)
	return util.Cache.Delete(key)
}

// SetUserConfig 设置用户配置
func (c *service) SetUserConfig(userId int, key string, value string) bool {
	identId := constant.CacheNameWithUserAppConfig + strconv.Itoa(userId)
	return c.SetConfig(identId, key, value)
}

// SetUserConfigs 批量设置用户配置
func (c *service) SetUserConfigs(userId int, data map[string]string) {
	identId := constant.CacheNameWithUserAppConfig + strconv.Itoa(userId)
	for k, v := range data {
		c.SetConfig(identId, k, v)
	}
}
