package config

import (
	"encoding/json"
	"service-user/app/constant"
	"service-user/app/model"
	"strconv"
	"strings"
)

// getConfig 获取配置
func (s *service) getConfig(userId int64, key string, refresh bool) map[string]string {
	config := make(map[string]string)
	cacheKey := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10) + "_" + key
	ttl := constant.CacheTimeWithUserConfig
	if refresh {
		config = s.selectWithUserIDAndKey(userId, key)
		data, err := json.Marshal(config)
		if err == nil {
			s.cache.SetWithExpire(cacheKey, string(data), ttl)
		}
	} else {
		data, exists := s.cache.Get(cacheKey)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = s.selectWithUserIDAndKey(userId, key)
			data, err := json.Marshal(config)
			if err == nil {
				s.cache.SetWithExpire(cacheKey, string(data), ttl)
			}
		}
	}
	return config
}

// setConfig 设置配置
func (s *service) setConfig(userId int64, key string, value string) bool {
	lock := constant.LockConfigSet + strconv.FormatInt(userId, 10) + "_" + key
	if !s.cache.Lock(lock, constant.LockTimeConfigSet) {
		return false
	}
	defer s.cache.UnLock(lock)
	rowsAffected, err := s.repositoryConfig.UpdateByUserIdAndKey(userId, key, value)
	if err != nil {
		return false
	}
	if rowsAffected <= 0 {
		err = s.repositoryConfig.Create(&model.Config{
			UserId: userId,
			Key:    key,
			Value:  value,
		})
	}
	return err == nil
}

// GetUserConfigs 获取用户配置
func (s *service) GetUserConfigs(userId int64, group string, refresh bool) map[string]string {
	return s.getConfig(userId, group, refresh)
}

// clearUserConfigCache 清理用户配置缓存
func (s *service) clearUserConfigCache(userId int64, group string) bool {
	dotIndex := strings.Index(group, ".")
	if dotIndex != -1 {
		group = group[:dotIndex]
	}
	key := constant.CacheNameWithUserConfig + strconv.FormatInt(userId, 10) + "_" + group
	return s.cache.Delete(key)
}

// SetUserConfig 设置用户配置
func (s *service) SetUserConfig(userId int64, key string, value string) bool {
	defer s.clearUserConfigCache(userId, key)
	return s.setConfig(userId, key, value)
}

// SetUserConfigs 批量设置用户配置
func (s *service) SetUserConfigs(userId int64, data map[string]string) {
	groups := map[string]int{}
	for k, v := range data {
		s.setConfig(userId, k, v)
		dotIndex := strings.Index(k, ".")
		if dotIndex != -1 {
			groups[k[:dotIndex]] = 1
		} else {
			groups[k] = 1
		}
	}
	for k := range groups {
		s.clearUserConfigCache(userId, k)
	}
}
