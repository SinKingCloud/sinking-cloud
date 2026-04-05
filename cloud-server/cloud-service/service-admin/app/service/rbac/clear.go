package rbac

import (
	"service-admin/app/constant"
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
