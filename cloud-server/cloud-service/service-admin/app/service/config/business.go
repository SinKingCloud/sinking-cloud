package config

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/model"
	"strconv"
	"strings"
)

// getConfig 获取配置
func (s *service) getConfig(key string, group string, ttl int64, refresh bool) map[string]string {
	config := make(map[string]string)
	cacheKey := key + group
	if refresh {
		config = s.selectWithIdentIDAndKey(key, group)
		data, err := json.Marshal(config)
		if err == nil {
			s.cache.SetWithExpire(cacheKey, string(data), ttl)
		}
	} else {
		data, exists := s.cache.Get(cacheKey)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = s.selectWithIdentIDAndKey(key, group)
			data, err := json.Marshal(config)
			if err == nil {
				s.cache.SetWithExpire(cacheKey, string(data), ttl)
			}
		}
	}
	return config
}

// setConfig 设置配置
func (s *service) setConfig(identId string, key string, value string) bool {
	lock := constant.LockConfigSet + identId + key
	if !s.cache.Lock(lock, constant.LockTimeConfigSet) {
		return false
	}
	defer s.cache.UnLock(lock)
	rowsAffected, err := s.repo.UpdateByIdentIdAndKey(identId, key, value)
	if err != nil || rowsAffected <= 0 {
		err = s.repo.Create(&model.Config{
			IdentId: identId,
			Key:     key,
			Value:   value,
		})
	}
	return err == nil
}

// GetUserConfigs 获取用户配置
func (s *service) GetUserConfigs(userId int64, group string, refresh bool) map[string]string {
	key := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10)
	return s.getConfig(key, group, constant.CacheTimeWithUserConfig, refresh)
}

// clearUserConfigCache 清理用户配置缓存
func (s *service) clearUserConfigCache(userId int64, group string) bool {
	dotIndex := strings.Index(group, ".")
	if dotIndex != -1 {
		group = group[:dotIndex]
	}
	key := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10)
	return s.cache.Delete(key + group)
}

// SetUserConfig 设置用户配置
func (s *service) SetUserConfig(userId int64, group string, key string, value string) bool {
	identId := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10)
	defer s.clearUserConfigCache(userId, group)
	return s.setConfig(identId, key, value)
}

// SetUserConfigs 批量设置用户配置
func (s *service) SetUserConfigs(userId int64, group string, data map[string]string) {
	identId := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10)
	for k, v := range data {
		s.setConfig(identId, k, v)
	}
	s.clearUserConfigCache(userId, group)
}

// GetSysConfigs 获取配置
func (s *service) GetSysConfigs(group string, refresh bool) map[string]string {
	key := constant.CacheNameWithSysConfig
	return s.getConfig(key, group, constant.CacheTimeWithSysConfig, refresh)
}

// clearSysConfig 清理配置缓存
func (s *service) clearSysConfig(group string) bool {
	dotIndex := strings.Index(group, ".")
	if dotIndex != -1 {
		group = group[:dotIndex]
	}
	key := constant.CacheNameWithSysConfig + group
	return s.cache.Delete(key)
}

// SetSysConfig 设置配置
func (s *service) SetSysConfig(key string, value string) bool {
	identId := constant.CacheNameWithSysConfig
	defer s.clearSysConfig(key)
	return s.setConfig(identId, key, value)
}

// SetSysConfigs 批量设置app配置
func (s *service) SetSysConfigs(data map[string]string) {
	identId := constant.CacheNameWithSysConfig
	for k, v := range data {
		s.setConfig(identId, k, v)
		s.clearSysConfig(k)
	}
}
