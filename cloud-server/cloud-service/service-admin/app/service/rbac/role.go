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
	maps, err := s.repositoryRole.SelectIdNameMap()
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(maps); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithEnum)
	}
	return maps, nil
}
