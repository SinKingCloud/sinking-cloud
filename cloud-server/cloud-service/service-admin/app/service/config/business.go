package config

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/model"
	"strconv"
	"strings"
)

const (
	systemIdentID = "Sys"
	userIdentID   = "User_"
)

// getConfig 统一获取配置（用户+系统通用）
func (s *service) getConfig(cacheKey string, queryIdent string, group string, ttl int64, refresh bool) map[string]string {
	config := make(map[string]string)
	if !refresh {
		data, exists := s.cache.Get(cacheKey)
		if exists {
			if err := json.Unmarshal([]byte(data), &config); err == nil {
				return config
			}
		}
	}
	config = s.selectWithIdentIDAndKey(queryIdent, group)
	data, err := json.Marshal(config)
	if err == nil {
		s.cache.SetWithExpire(cacheKey, string(data), ttl)
	}
	return config
}

// setConfig 统一设置配置（用户+系统通用）
func (s *service) setConfig(identId string, key string, value string) bool {
	lockKey := constant.LockConfigSet + identId + ":" + key
	if !s.cache.Lock(lockKey, constant.LockTimeConfigSet) {
		return false
	}
	defer s.cache.UnLock(lockKey)
	rowsAffected, err := s.repositoryConfig.UpdateByIdentIdAndKey(identId, key, value)
	if err != nil {
		return false
	}
	if rowsAffected <= 0 {
		err = s.repositoryConfig.Create(&model.Config{
			IdentId: identId,
			Key:     key,
			Value:   value,
		})
	}
	return err == nil
}

// GetUserConfigs 获取用户配置
func (s *service) GetUserConfigs(userId int64, group string, refresh bool) map[string]string {
	uidStr := strconv.FormatInt(userId, 10)
	cacheKey := constant.CacheNameWithUserConfig + uidStr + "_" + group
	return s.getConfig(cacheKey, userIdentID+uidStr, group, constant.CacheTimeWithUserConfig, refresh)
}

// clearUserConfigCache 清理用户配置缓存
func (s *service) clearUserConfigCache(userId int64, keysOrGroups ...string) bool {
	if len(keysOrGroups) == 0 {
		return true
	}
	uidStr := strconv.FormatInt(userId, 10)
	cacheKeyMap := make(map[string]struct{}, len(keysOrGroups))
	for _, key := range keysOrGroups {
		group := s.getGroupKey(key)
		cacheKeyMap[constant.CacheNameWithUserConfig+uidStr+"_"+group] = struct{}{}
	}
	cacheKeys := make([]string, 0, len(cacheKeyMap))
	for cacheKey := range cacheKeyMap {
		cacheKeys = append(cacheKeys, cacheKey)
	}
	return s.cache.Delete(cacheKeys...)
}

// SetUserConfig 修改用户配置
func (s *service) SetUserConfig(userId int64, key string, value string) bool {
	uidStr := strconv.FormatInt(userId, 10)
	defer s.clearUserConfigCache(userId, key, "")
	return s.setConfig(userIdentID+uidStr, key, value)
}

// SetUserConfigs 批量修改用户配置
func (s *service) SetUserConfigs(userId int64, data map[string]string) {
	uidStr := strconv.FormatInt(userId, 10)
	groups := make(map[string]struct{}, len(data)+1)
	groups[""] = struct{}{}
	for k, v := range data {
		s.setConfig(userIdentID+uidStr, k, v)
		groups[s.getGroupKey(k)] = struct{}{}
	}
	groupList := make([]string, 0, len(groups))
	for g := range groups {
		groupList = append(groupList, g)
	}
	s.clearUserConfigCache(userId, groupList...)
}

// ClearUserConfig 删除用户配置
func (s *service) ClearUserConfig(userId int64, group string) bool {
	if userId <= 0 {
		return false
	}
	uidStr := strconv.FormatInt(userId, 10)
	identId := userIdentID + uidStr
	configs, err := s.repositoryConfig.FindByIdentIdAndKey(identId, group)
	if err != nil {
		return false
	}
	if err = s.repositoryConfig.DeleteByIdentIdAndKey(identId, group); err != nil {
		return false
	}
	groups := s.getConfigGroups(configs, group)
	return s.clearUserConfigCache(userId, groups...)
}

// GetSysConfigs 获取系统配置
func (s *service) GetSysConfigs(group string, refresh bool) map[string]string {
	cacheKey := constant.CacheNameWithSysConfig + group
	return s.getConfig(cacheKey, systemIdentID, group, constant.CacheTimeWithSysConfig, refresh)
}

// clearSysConfigCache 清理系统配置缓存
func (s *service) clearSysConfigCache(keysOrGroups ...string) bool {
	if len(keysOrGroups) == 0 {
		return true
	}
	cacheKeyMap := make(map[string]struct{}, len(keysOrGroups))
	for _, key := range keysOrGroups {
		group := s.getGroupKey(key)
		cacheKeyMap[constant.CacheNameWithSysConfig+group] = struct{}{}
	}
	cacheKeys := make([]string, 0, len(cacheKeyMap))
	for cacheKey := range cacheKeyMap {
		cacheKeys = append(cacheKeys, cacheKey)
	}
	return s.cache.Delete(cacheKeys...)
}

// SetSysConfig 修改系统配置
func (s *service) SetSysConfig(key string, value string) bool {
	defer s.clearSysConfigCache(key, "")
	return s.setConfig(systemIdentID, key, value)
}

// SetSysConfigs 批量修改系统配置
func (s *service) SetSysConfigs(data map[string]string) {
	groups := make(map[string]struct{}, len(data)+1)
	groups[""] = struct{}{}
	for k, v := range data {
		s.setConfig(systemIdentID, k, v)
		groups[s.getGroupKey(k)] = struct{}{}
	}
	groupList := make([]string, 0, len(groups))
	for g := range groups {
		groupList = append(groupList, g)
	}
	s.clearSysConfigCache(groupList...)
}

// ClearSysConfig 删除系统配置
func (s *service) ClearSysConfig(group string) bool {
	configs, err := s.repositoryConfig.FindByIdentIdAndKey(systemIdentID, group)
	if err != nil {
		return false
	}
	if err = s.repositoryConfig.DeleteByIdentIdAndKey(systemIdentID, group); err != nil {
		return false
	}
	groups := s.getConfigGroups(configs, group)
	return s.clearSysConfigCache(groups...)
}

// getGroupKey 统一截取分组（theme.color → theme）
func (s *service) getGroupKey(key string) string {
	if key == "*" {
		return ""
	}
	if idx := strings.Index(key, "."); idx != -1 {
		return key[:idx]
	}
	return key
}

// getConfigGroups 获取需要清理的配置分组
func (s *service) getConfigGroups(configs []*model.Config, group string) []string {
	groupMap := make(map[string]struct{}, len(configs)+2)
	groupMap[""] = struct{}{}
	if group != "" {
		groupMap[s.getGroupKey(group)] = struct{}{}
	}
	for _, item := range configs {
		if item == nil || item.Key == "" {
			continue
		}
		groupMap[s.getGroupKey(item.Key)] = struct{}{}
	}
	groups := make([]string, 0, len(groupMap))
	for key := range groupMap {
		groups = append(groups, key)
	}
	return groups
}
