package rbac

import (
	"errors"
	"service-admin/app/constant"
	"service-admin/app/enum/menu"
	"service-admin/app/model"
	"service-admin/app/repository/permission"
	"sort"
)

// GetUserControl 获取用户页面控件
func (s *service) GetUserControl(userId int64, code string) []*UserControl {
	var controls []*UserControl
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
	controls = make([]*UserControl, 0, len(authControlMap))
	for _, item := range menus {
		if _, ok := authControlMap[item.Id]; ok {
			value := controlMap[item.Id]
			if value == nil {
				continue
			}
			controls = append(controls, &UserControl{
				Code: value.Code,
				Name: value.Name,
				Icon: value.Icon,
			})
		}
	}
	return controls
}

// GetUserMenus 获取用户菜单
func (s *service) GetUserMenus(userId int64) []*UserMenu {
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
	//对菜单数据进行树状组装
	menuMap := make(map[int64]*UserMenu, len(menus))
	var rootMenus []*UserMenu
	for _, item := range menus {
		menuMap[item.Id] = &UserMenu{
			Id:   item.Id,
			Code: item.Code,
			Name: item.Name,
			Type: item.Type,
			Icon: item.Icon,
			Sort: item.Sort,
		}
	}
	for _, item := range menus {
		current := menuMap[item.Id]
		if current == nil {
			continue
		}
		if item.ParentId == 0 {
			rootMenus = append(rootMenus, current)
			continue
		}
		if parentMenu, ok := menuMap[item.ParentId]; ok && parentMenu != nil {
			parentMenu.Children = append(parentMenu.Children, current)
			continue
		}
		rootMenus = append(rootMenus, current)
	}
	//菜单排序
	var sortMenus func(list []*UserMenu)
	sortMenus = func(list []*UserMenu) {
		if len(list) == 0 {
			return
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].Sort == list[j].Sort {
				return list[i].Id < list[j].Id
			}
			return list[i].Sort > list[j].Sort
		})
		for _, item := range list {
			if len(item.Children) > 0 {
				sortMenus(item.Children)
			}
		}
	}
	sortMenus(rootMenus)
	return rootMenus
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
