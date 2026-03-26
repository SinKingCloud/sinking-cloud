package rbac

import (
	"service-admin/app/constant"
)

// clearCacheWithPermission 清理缓存
func (s *service) clearCacheWithPermission(path string) {
	if path != "" {
		s.cache.Delete(constant.CacheNameWithPermission + path)
	}
}
