package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"strconv"
)

// getRoleIdsByUserId 根据userId查询用户角色列表
func (s *service) getRoleIdsByUserId(userId int64, refresh bool) (ids []int64, err error) {
	cacheKey := constant.CacheNameWithUserRoles + strconv.FormatInt(userId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &ids); err == nil {
				return ids, nil
			}
		}
	}
	ids, err = s.repositoryUserRole.SelectRoleIdsByUserId(userId)
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(ids); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithUserRoles)
	}
	return ids, nil
}
