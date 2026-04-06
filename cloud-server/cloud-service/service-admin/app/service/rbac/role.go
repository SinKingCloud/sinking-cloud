package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
)

// GetRoleNameEnum 获取角色名称枚举
func (s *service) GetRoleNameEnum(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithEnum + "role_names"
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result map[int64]string
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	roles, err := s.repositoryRole.SelectAll()
	if err != nil {
		return nil, err
	}
	result := map[int64]string{}
	for _, item := range roles {
		if item == nil || item.Id <= 0 {
			continue
		}
		result[item.Id] = item.Name
	}
	if data, err2 := json.Marshal(result); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithEnum)
	}
	return result, nil
}
