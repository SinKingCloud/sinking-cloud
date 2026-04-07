package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"strconv"
)

// getPermissionCodesByRoleId 根据roleId查询角色权限列表
func (s *service) getPermissionCodesByRoleId(roleId int64, refresh bool) (permissionCodes []string, err error) {
	cacheKey := constant.CacheNameWithRolePermissions + strconv.FormatInt(roleId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &permissionCodes); err == nil {
				return permissionCodes, nil
			}
		}
	}
	menuIds, err := s.getRoleMenuIdsByRoleId(roleId, false)
	if err != nil {
		return nil, err
	}
	if len(menuIds) > 0 {
		permissionCodes, err = s.repositoryMenuPermission.SelectPermissionCodesByMenuIds(menuIds)
		if err != nil {
			return nil, err
		}
	}
	if data, err2 := json.Marshal(permissionCodes); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithRolePermissions)
	}
	return permissionCodes, nil
}
