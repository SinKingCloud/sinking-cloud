package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/enum/role_status"
	"strconv"
)

// getRoleMenuIdsByRoleId 根据角色ID查询所有菜单ID
func (s *service) getRoleMenuIdsByRoleId(roleId int64, refresh bool) (menuIds []int64, err error) {
	cacheKey := constant.CacheNameWithRoleMenus + strconv.FormatInt(roleId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &menuIds); err == nil {
				return menuIds, nil
			}
		}
	}
	r, err := s.repositoryRole.FindById(roleId)
	if err == nil && r != nil && r.Id > 0 && r.Status == role_status.Enable {
		menuIds, err = s.repositoryRoleMenu.SelectMenuIdsByRoleId(r.Id)
		if err != nil {
			return nil, err
		}
	}
	if data, err2 := json.Marshal(menuIds); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithRoleMenus)
	}
	return menuIds, nil
}
