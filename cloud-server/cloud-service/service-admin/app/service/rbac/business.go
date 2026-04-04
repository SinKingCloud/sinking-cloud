package rbac

import (
	"errors"
	"service-admin/app/constant"
	"service-admin/app/enum/menu"
	"service-admin/app/model"
	"service-admin/app/repository/permission"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// GetUserControl 获取用户页面控件
func (s *service) GetUserControl(userId int64, code string) []*Control {
	var controls []*Control
	info, err := s.getMenuByCode(code, false)
	if err != nil || info == nil || info.Id <= 0 {
		return controls
	}
	//拿到页面下的控件数据
	menus, err := s.getMenusByParentId(info.Id, false)
	if err != nil || len(menus) == 0 {
		return controls
	}
	tempMap := make(map[int64]struct{}, len(menus))
	controlMap := make(map[int64]*model.Menu, len(menus))
	for _, item := range menus {
		tempMap[item.Id] = struct{}{}
		controlMap[item.Id] = item
	}
	//根据用户ID查询用户角色，获取到所有角色ID
	roleIds, err := s.getRoleIdsByUserId(userId, false)
	if err != nil || len(roleIds) == 0 {
		return controls
	}
	//根据角色ID查询到所有菜单ID
	authControlMap := make(map[int64]struct{}, len(menus))
	for _, roleId := range roleIds {
		roleMenuIds, err2 := s.getRoleMenuIdsByRoleId(roleId, false)
		if err2 == nil && roleMenuIds != nil {
			for _, roleMenuId := range roleMenuIds {
				if _, ok := tempMap[roleMenuId]; ok {
					authControlMap[roleMenuId] = struct{}{}
				}
			}
		}
	}
	controls = make([]*Control, 0, len(authControlMap))
	for _, item := range menus {
		if _, ok := authControlMap[item.Id]; ok {
			value := controlMap[item.Id]
			if value == nil {
				continue
			}
			controls = append(controls, &Control{
				Code: value.Code,
				Name: value.Name,
				Icon: value.Icon,
			})
		}
	}
	return controls
}

// GetUserMenuTree 获取用户菜单
func (s *service) GetUserMenuTree(userId int64) []*MenuTree {
	//根据用户ID查询用户角色，获取到所有角色ID
	roleIds, err := s.getRoleIdsByUserId(userId, false)
	if err != nil || len(roleIds) == 0 {
		return nil
	}
	//根据角色ID查询到所有菜单ID
	var menuIds []int64
	tempMap := make(map[int64]struct{}, 200)
	for _, roleId := range roleIds {
		roleMenuIds, err2 := s.getRoleMenuIdsByRoleId(roleId, false)
		if err2 == nil && roleMenuIds != nil {
			for _, roleMenuId := range roleMenuIds {
				if _, ok := tempMap[roleMenuId]; !ok {
					tempMap[roleMenuId] = struct{}{}
					menuIds = append(menuIds, roleMenuId)
				}
			}
		}
	}
	//去重菜单ID后根据菜单ID查询菜单信息
	menus, err := s.repositoryMenu.SelectByIdsWithStatusAndType(menuIds, int(menu.StatusEnable), []int{int(menu.TypeDir), int(menu.TypePage)})
	if err != nil || len(menus) == 0 {
		return nil
	}
	//返回菜单树
	return s.buildMenuTree(menus)
}

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
	record, err := s.repositoryPermission.FindByCode(code)
	if err != nil {
		return err
	}
	if record == nil || record.Id <= 0 {
		if err = s.repositoryPermission.Create(perm); err != nil {
			return err
		}
	} else {
		if err = s.repositoryPermission.UpdateByCode(code, &permission.UpdatePermission{
			Name: perm.Name,
			Path: perm.Path,
		}); err != nil {
			return err
		}
	}
	// 刷新缓存
	s.clearCacheWithPermissions([]string{path})
	return nil
}

// CheckUserHasPermission 判断用户是否含有权限
func (s *service) CheckUserHasPermission(userId int64, path string) error {
	//根据path查询数据库，获取path对应的权限code
	perm, err := s.getPermissionByPath(path, false)
	if err != nil {
		return err
	}
	if perm == nil || perm.Id <= 0 {
		return nil
	}
	//存在则根据用户ID查询用户角色，获取到所有角色ID
	roleIds, err := s.getRoleIdsByUserId(userId, false)
	if err != nil {
		return err
	}
	if len(roleIds) <= 0 {
		return errors.New("该账号还未分配权限")
	}
	//根据角色ID查询到所有权限code并是否含有path对应的code
	for _, roleId := range roleIds {
		permissionCodes, err2 := s.getPermissionCodesByRoleId(roleId, false)
		if err2 == nil && permissionCodes != nil {
			for _, permissionCode := range permissionCodes {
				if permissionCode == perm.Code {
					return nil
				}
			}
		}
	}
	return errors.New("无访问权限")
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

// CreatePermission 创建权限
func (s *service) CreatePermission(data *model.Permission) error {
	if data.Code == "" {
		return errors.New("标识不能为空")
	}
	info, err := s.repositoryPermission.FindByCode(data.Code)
	if err != nil {
		return err
	}
	if info != nil && info.Id > 0 {
		return errors.New("该标识已被使用")
	}
	if data.Path == "" {
		return errors.New("路径不能为空")
	}
	info, err = s.repositoryPermission.FindByPath(data.Path)
	if err != nil {
		return err
	}
	if info != nil && info.Id > 0 {
		return errors.New("该路径已被使用")
	}
	data.UpdateTime = str.DateTime(time.Now())
	data.CreateTime = str.DateTime(time.Now())
	return s.repositoryPermission.Create(data)
}

// GetAllMenuTree 获取所有菜单树状数据
func (s *service) GetAllMenuTree() []*MenuTree {
	menus, err := s.repositoryMenu.SelectWithType([]int{int(menu.TypeDir), int(menu.TypePage)})
	if err != nil {
		return nil
	}
	return s.buildMenuTree(menus)
}
