package rbac

import (
	"service-admin/app/constant"
	"service-admin/app/model"
	"strconv"
)

// clearCacheWithPermissions 清理权限缓存
func (s *service) clearCacheWithPermissions(paths []string) {
	pathMap := make(map[string]struct{}, len(paths))
	keys := make([]string, 0, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}
		if _, ok := pathMap[path]; ok {
			continue
		}
		pathMap[path] = struct{}{}
		keys = append(keys, constant.CacheNameWithPermission+path)
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// clearCacheWithRolePermissions 清理角色权限缓存
func (s *service) clearCacheWithRolePermissions(roleIds []int64) {
	roleMap := make(map[int64]struct{}, len(roleIds))
	keys := make([]string, 0, len(roleIds))
	for _, roleId := range roleIds {
		if roleId <= 0 {
			continue
		}
		if _, ok := roleMap[roleId]; ok {
			continue
		}
		roleMap[roleId] = struct{}{}
		keys = append(keys, constant.CacheNameWithRolePermissions+strconv.FormatInt(roleId, 10))
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// clearCacheWithRoleMenus 清理角色菜单缓存
func (s *service) clearCacheWithRoleMenus(roleIds []int64) {
	roleMap := make(map[int64]struct{}, len(roleIds))
	keys := make([]string, 0, len(roleIds))
	for _, roleId := range roleIds {
		if roleId <= 0 {
			continue
		}
		if _, ok := roleMap[roleId]; ok {
			continue
		}
		roleMap[roleId] = struct{}{}
		keys = append(keys, constant.CacheNameWithRoleMenus+strconv.FormatInt(roleId, 10))
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// clearCacheWithUserRoles 清理用户角色缓存
func (s *service) clearCacheWithUserRoles(userIds []int64) {
	userMap := make(map[int64]struct{}, len(userIds))
	keys := make([]string, 0, len(userIds))
	for _, userId := range userIds {
		if userId <= 0 {
			continue
		}
		if _, ok := userMap[userId]; ok {
			continue
		}
		userMap[userId] = struct{}{}
		keys = append(keys, constant.CacheNameWithUserRoles+strconv.FormatInt(userId, 10))
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// clearCacheWithMenus 清理菜单缓存
func (s *service) clearCacheWithMenus(menus []*model.Menu) {
	keyMap := make(map[string]struct{}, len(menus)*3)
	keys := make([]string, 0, len(menus)*3)
	for _, item := range menus {
		if item == nil {
			continue
		}
		if item.Code != "" {
			cacheKey := constant.CacheNameWithMenus + item.Code
			if _, ok := keyMap[cacheKey]; !ok {
				keyMap[cacheKey] = struct{}{}
				keys = append(keys, cacheKey)
			}
		}
		if item.Id > 0 {
			cacheKey := constant.CacheNameWithMenus + "_list_" + strconv.FormatInt(item.Id, 10)
			if _, ok := keyMap[cacheKey]; !ok {
				keyMap[cacheKey] = struct{}{}
				keys = append(keys, cacheKey)
			}
		}
		if item.ParentId > 0 {
			cacheKey := constant.CacheNameWithMenus + "_list_" + strconv.FormatInt(item.ParentId, 10)
			if _, ok := keyMap[cacheKey]; !ok {
				keyMap[cacheKey] = struct{}{}
				keys = append(keys, cacheKey)
			}
		}
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// clearCacheWithEnums 清理枚举缓存
func (s *service) clearCacheWithEnums(names []string) {
	nameMap := make(map[string]struct{}, len(names))
	keys := make([]string, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		if _, ok := nameMap[name]; ok {
			continue
		}
		nameMap[name] = struct{}{}
		keys = append(keys, constant.CacheNameWithEnum+name)
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}
