package rbac

import (
	"errors"
	"service-admin/app/constant"
	"service-admin/app/model"
)

// RegisterPermission 注册权限（以 code 为唯一标识，不存在则插入，存在则更新）
func (s *service) RegisterPermission(name string, code string, path string) error {
	codeLockKey := constant.LockPermission + code
	if s.cache.Lock(codeLockKey, constant.LockTimePermission) {
		defer s.cache.UnLock(codeLockKey)
	} else {
		return errors.New("系统繁忙,请稍候重试")
	}
	perm := &model.Permission{
		Name: name,
		Code: code,
		Path: path,
	}
	existing, err := s.repositoryPermission.FindByCode(code)
	if err != nil {
		return err
	}
	if existing == nil {
		if err = s.repositoryPermission.Create(perm); err != nil {
			return err
		}
	} else {
		if err = s.repositoryPermission.UpdateByCode(code, perm); err != nil {
			return err
		}
	}
	// 刷新缓存
	s.clearCacheWithPermission(path)
	return nil
}

// CheckUserHasPermission 判断用户是否含有权限
func (s *service) CheckUserHasPermission(userId int64, path string) error {
	//根据path查询数据库，获取path对应的权限code
	permission, err := s.findPermissionByPath(path, false)
	if err != nil {
		return err
	}
	if permission == nil {
		return nil
	}
	//存在则根据用户ID查询用户角色，获取到所有角色ID
	roleIds, err := s.findRoleIdsByUserId(userId, false)
	if err != nil {
		return err
	}
	if len(roleIds) <= 0 {
		return errors.New("该账号还未分配权限")
	}
	//根据角色ID查询到所有权限code并是否含有path对应的code
	for _, roleId := range roleIds {
		permissionCodes, err2 := s.findPermissionCodesByRoleId(roleId, false)
		if err2 == nil && permissionCodes != nil {
			for _, permissionCode := range permissionCodes {
				if permissionCode == permission.Code {
					return nil
				}
			}
		}
	}
	return errors.New("无访问权限")
}
