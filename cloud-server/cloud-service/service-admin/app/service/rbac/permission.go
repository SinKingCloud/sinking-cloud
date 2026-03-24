package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/model"
)

// findPermissionByPath 根据path查询权限
func (s *service) findPermissionByPath(path string, refresh bool) (perm *model.Permission, err error) {
	cacheKey := constant.CacheNameWithPermission + path
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &perm); err == nil {
				return perm, nil
			}
		}
	}
	perm, err = s.repositoryPermission.FindByPath(path)
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(perm); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithPermission)
	}
	return perm, nil
}
