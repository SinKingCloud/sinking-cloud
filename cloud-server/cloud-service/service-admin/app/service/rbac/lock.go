package rbac

import (
	"errors"
	"service-admin/app/constant"
	"sort"
	"strconv"
)

// lockRoles 对角色ID集合按顺序加锁，避免重叠角色操作并发执行
func (s *service) lockRoles(roleIds []int64) ([]string, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	idMap := make(map[int64]struct{}, len(roleIds))
	ids := make([]int64, 0, len(roleIds))
	for _, roleId := range roleIds {
		if roleId <= 0 {
			continue
		}
		if _, ok := idMap[roleId]; ok {
			continue
		}
		idMap[roleId] = struct{}{}
		ids = append(ids, roleId)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		key := constant.LockRole + strconv.FormatInt(id, 10)
		if !s.cache.Lock(key, constant.LockTimeRole) {
			if len(keys) > 0 {
				s.cache.UnLock(keys...)
			}
			return nil, errors.New("系统繁忙,请稍候重试")
		}
		keys = append(keys, key)
	}
	return keys, nil
}
