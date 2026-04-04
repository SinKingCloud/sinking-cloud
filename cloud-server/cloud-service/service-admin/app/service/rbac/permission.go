package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/model"
	"service-admin/app/repository/permission"
)

// getPermissionByPath 根据path查询权限
func (s *service) getPermissionByPath(path string, refresh bool) (perm *model.Permission, err error) {
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

// SelectPermission 查询权限数据
func (s *service) SelectPermission(where *permission.SelectPermission, orderByField string, orderByType string, page int, pageSize int) ([]*model.Permission, int64, error) {
	return s.repositoryPermission.Select(where, orderByField, orderByType, page, pageSize)
}

// UpdatePermissionByIds 批量更新权限数据
func (s *service) UpdatePermissionByIds(ids []int64, data *permission.UpdatePermission) error {
	err := s.repositoryPermission.UpdateByIds(ids, data)
	if err == nil {

	}
	return err
}

// DeletePermissionByIds 批量删除权限数据
func (s *service) DeletePermissionByIds(ids []int64) error {
	err := s.repositoryPermission.DeleteByIds(ids)
	if err == nil {

	}
	return err
}
