package config

import (
	"encoding/json"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

// GetSystemConfigWithCache 获取系统配置（带缓存）
func (s *Service) GetSystemConfigWithCache(group string, refresh bool) map[string]string {
	key := constant.SystemConfigCacheName + group
	if !refresh {
		// 尝试从缓存获取
		s.systemConfigLock.RLock()
		if val, ok := global.App.Cache.Get(key); ok {
			s.systemConfigLock.RUnlock()
			var config map[string]string
			if err := json.Unmarshal([]byte(val), &config); err == nil {
				return config
			}
		}
		s.systemConfigLock.RUnlock()
	}
	// 缓存未命中或强制刷新，加写锁防止击穿
	s.systemConfigLock.Lock()
	defer s.systemConfigLock.Unlock()
	// 双重检查
	if !refresh {
		if val, ok := global.App.Cache.Get(key); ok {
			var config map[string]string
			if err := json.Unmarshal([]byte(val), &config); err == nil {
				return config
			}
		}
	}
	config := s.GetSystemConfig(group, false)
	if len(config) > 0 {
		b, _ := json.Marshal(config)
		global.App.Cache.SetWithExpire(key, string(b), int64(constant.SystemConfigCacheTime))
	}
	return config
}

// GetSystemUserConfigWithCache 获取系统用户配置（带缓存）
func (s *Service) GetSystemUserConfigWithCache(userId int64, group string, refresh bool) map[string]string {
	key := constant.SystemUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	if !refresh {
		if val, ok := global.App.Cache.Get(key); ok {
			var config map[string]string
			if err := json.Unmarshal([]byte(val), &config); err == nil {
				return config
			}
		}
	}
	config := s.GetSystemUserConfig(userId, group, false)
	if len(config) > 0 {
		b, _ := json.Marshal(config)
		global.App.Cache.SetWithExpire(key, string(b), int64(constant.SystemUserConfigCacheTime))
	}
	return config
}

// GetWebUserConfigWithCache 获取Web用户配置（带缓存）
func (s *Service) GetWebUserConfigWithCache(userId int64, group string, refresh bool) map[string]string {
	key := constant.WebUserConfigCacheName + strconv.FormatInt(userId, 10) + ":" + group
	if !refresh {
		if val, ok := global.App.Cache.Get(key); ok {
			var config map[string]string
			if err := json.Unmarshal([]byte(val), &config); err == nil {
				return config
			}
		}
	}
	config := s.GetWebUserConfig(userId, group, false)
	if len(config) > 0 {
		b, _ := json.Marshal(config)
		global.App.Cache.SetWithExpire(key, string(b), int64(constant.WebUserConfigCacheTime))
	}
	return config
}
