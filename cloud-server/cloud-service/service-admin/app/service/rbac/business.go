package rbac

import (
	"errors"
	"service-admin/app/constant"
	"service-admin/app/enum/menu"
	"service-admin/app/model"
	"service-admin/app/repository/permission"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// GetPageControlsByUserId 根据用户ID获取页面控件
func (s *service) GetPageControlsByUserId(userId int64, code string) []*UserControl {
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

// GetUserMenuTree 获取用户菜单
func (s *service) GetUserMenuTree(userId int64) []*UserMenuTree {
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
	tree := s.buildMenuTree(menus)
	var formatMenus func(list []*MenuTree) []*UserMenuTree
	formatMenus = func(list []*MenuTree) []*UserMenuTree {
		if len(list) == 0 {
			return nil
		}
		result := make([]*UserMenuTree, 0, len(list))
		for _, item := range list {
			if item == nil {
				continue
			}
			result = append(result, &UserMenuTree{
				Code:     item.Code,
				Name:     item.Name,
				Type:     item.Type,
				Icon:     item.Icon,
				Sort:     item.Sort,
				Children: formatMenus(item.Children),
			})
		}
		return result
	}
	return formatMenus(tree)
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
	s.clearCacheWithEnums([]string{"permission"})
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
		s.clearCacheWithPermissions(paths)
		s.clearCacheWithEnums([]string{"permission"})
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
		s.clearCacheWithEnums([]string{"permission"})
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
	info, err = s.getPermissionByPath(data.Path, false)
	if err != nil {
		return err
	}
	if info != nil && info.Id > 0 {
		return errors.New("该路径已被使用")
	}
	data.UpdateTime = str.DateTime(time.Now())
	data.CreateTime = str.DateTime(time.Now())
	err = s.repositoryPermission.Create(data)
	if err == nil {
		s.clearCacheWithPermissions([]string{data.Path})
		s.clearCacheWithEnums([]string{"permission"})
	}
	return err
}

// GetAllMenuTree 获取所有菜单树状数据
func (s *service) GetAllMenuTree() []*MenuTree {
	menus, err := s.repositoryMenu.SelectWithType([]int{int(menu.TypeDir), int(menu.TypePage)})
	if err != nil {
		return nil
	}
	return s.buildMenuTree(menus)
}

// GetMenuOrControlInfoById 根据ID获取菜单或控件信息
func (s *service) GetMenuOrControlInfoById(id int64) (*MenuOrControlInfo, error) {
	info, err := s.repositoryMenu.FindById(id)
	if err != nil {
		return nil, errors.Join(errors.New("查询数据失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("该数据不存在")
	}
	menuOrControlInfo := &MenuOrControlInfo{
		Id:         info.Id,
		Code:       info.Code,
		Name:       info.Name,
		Icon:       info.Icon,
		Type:       info.Type,
		Sort:       info.Sort,
		Status:     info.Status,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}
	if menu.Type(info.Type) != menu.TypeDir {
		menuOrControlInfo.PermissionCodes, err = s.repositoryMenuPermission.SelectPermissionCodesByMenuId(info.Id)
		if err != nil {
			return nil, errors.Join(errors.New("查询权限数据失败"), err)
		}
	}
	return menuOrControlInfo, nil
}

// GetPageControlsByMenuId 获取页面控件
func (s *service) GetPageControlsByMenuId(id int64) []*Control {
	var controls []*Control
	if id <= 0 {
		return controls
	}
	menus, err := s.repositoryMenu.SelectByParentIdWithType(id, int(menu.TypeControl))
	if err != nil || len(menus) == 0 {
		return controls
	}
	n := len(menus)
	menuIds := make([]int64, 0, n)
	for _, item := range menus {
		if item == nil || item.Id <= 0 {
			continue
		}
		menuIds = append(menuIds, item.Id)
	}
	menuPermissions, err := s.repositoryMenuPermission.SelectByMenuIds(menuIds)
	if err != nil {
		return nil
	}
	permissionCodeMap := make(map[int64][]string, n)
	tempMap := make(map[int64]map[string]struct{}, n)
	for _, item := range menuPermissions {
		if item == nil || item.MenuId <= 0 || item.PermissionCode == "" {
			continue
		}
		if _, ok := tempMap[item.MenuId]; !ok {
			tempMap[item.MenuId] = make(map[string]struct{}, 5)
		}
		if _, ok := tempMap[item.MenuId][item.PermissionCode]; ok {
			continue
		}
		tempMap[item.MenuId][item.PermissionCode] = struct{}{}
		permissionCodeMap[item.MenuId] = append(permissionCodeMap[item.MenuId], item.PermissionCode)
	}
	controls = make([]*Control, 0, len(menus))
	for _, item := range menus {
		if item == nil || item.Id <= 0 {
			continue
		}
		controls = append(controls, &Control{
			Id:              item.Id,
			Code:            item.Code,
			Name:            item.Name,
			Icon:            item.Icon,
			PermissionCodes: permissionCodeMap[item.Id],
			Status:          item.Status,
			CreateTime:      item.CreateTime,
			UpdateTime:      item.UpdateTime,
		})
	}
	return controls
}

// DeleteMenuOrControlByIds 根据ID列表删除菜单或控件数据
func (s *service) DeleteMenuOrControlByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	menus, err := s.repositoryMenu.SelectAll()
	if err != nil {
		return errors.Join(errors.New("查询菜单数据失败"), err)
	}
	if len(menus) == 0 {
		return nil
	}
	menuMap := make(map[int64]*model.Menu, len(menus))
	childrenMap := make(map[int64][]int64, len(menus))
	for _, item := range menus {
		if item == nil || item.Id <= 0 {
			continue
		}
		menuMap[item.Id] = item
		childrenMap[item.ParentId] = append(childrenMap[item.ParentId], item.Id)
	}
	deleteMap := make(map[int64]*model.Menu, len(ids))
	var loadChildren func(menuId int64)
	loadChildren = func(menuId int64) {
		if _, ok := deleteMap[menuId]; ok {
			return
		}
		info, ok := menuMap[menuId]
		if !ok || info == nil {
			return
		}
		deleteMap[menuId] = info
		for _, childId := range childrenMap[menuId] {
			loadChildren(childId)
		}
	}
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		loadChildren(id)
	}
	if len(deleteMap) == 0 {
		return nil
	}
	deleteIds := make([]int64, 0, len(deleteMap))
	deleteMenus := make([]*model.Menu, 0, len(deleteMap))
	for _, item := range deleteMap {
		if item == nil || item.Id <= 0 {
			continue
		}
		deleteIds = append(deleteIds, item.Id)
		deleteMenus = append(deleteMenus, item)
	}
	if len(deleteIds) == 0 {
		return nil
	}
	roleIds, err := s.repositoryRoleMenu.SelectRoleIdsByMenuIds(deleteIds)
	if err != nil {
		return errors.Join(errors.New("查询角色菜单关联失败"), err)
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		err = s.repositoryMenuPermission.DeleteByMenuIds(deleteIds, tx)
		if err != nil {
			return errors.Join(errors.New("删除菜单权限关联失败"), err)
		}
		err = s.repositoryRoleMenu.DeleteByMenuIds(deleteIds, tx)
		if err != nil {
			return errors.Join(errors.New("删除角色菜单关联失败"), err)
		}
		err = s.repositoryMenu.DeleteByIds(deleteIds, tx)
		if err != nil {
			return errors.Join(errors.New("删除菜单数据失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.clearCacheWithRoleMenus(roleIds)
	s.clearCacheWithMenus(deleteMenus)
	return nil
}
