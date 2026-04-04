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
	if len(ids) == 0 || data == nil {
		return nil
	}
	perms, err := s.repositoryPermission.SelectByIds(ids)
	if err != nil {
		return err
	}
	err = s.repositoryPermission.UpdateByIds(ids, data)
	if err == nil {
		paths := make([]string, 0, len(perms)+1)
		for _, item := range perms {
			if item == nil || item.Path == "" {
				continue
			}
			paths = append(paths, item.Path)
		}
		if data.Path != "" {
			paths = append(paths, data.Path)
		}
		go s.clearCacheWithPermissions(paths)
	}
	return err
}

// DeletePermissionByIds 批量删除权限数据
func (s *service) DeletePermissionByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	perms, err := s.repositoryPermission.SelectByIds(ids)
	if err != nil {
		return err
	}
	err = s.repositoryPermission.DeleteByIds(ids)
	if err == nil {
		paths := make([]string, 0, len(perms))
		codes := make([]string, 0, len(perms))
		for _, item := range perms {
			if item == nil {
				continue
			}
			if item.Path != "" {
				paths = append(paths, item.Path)
			}
			if item.Code != "" {
				codes = append(codes, item.Code)
			}
		}
		roleIds, err2 := s.repositoryRolePermission.SelectRoleIdsByPermissionCodes(codes)
		if err2 != nil {
			return err2
		}
		s.clearCacheWithPermissions(paths)
		s.clearCacheWithRolePermissions(roleIds)
	}
	return err
}
