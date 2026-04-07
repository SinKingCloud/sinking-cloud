package rbac

import (
	"errors"
	"fmt"
	"service-admin/app/constant"
	"service-admin/app/enum/menu_status"
	"service-admin/app/enum/menu_type"
	"service-admin/app/enum/role_status"
	"service-admin/app/model"
	menu2 "service-admin/app/repository/menu"
	"service-admin/app/repository/permission"
	"service-admin/app/repository/role"
	"strconv"
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
	menus, err := s.repositoryMenu.SelectByIdsWithStatusAndType(menuIds, menu_status.Enable, []int{menu_type.Dir, menu_type.Page})
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
	now := str.DateTime(time.Now())
	perm := &model.Permission{
		Name:       name,
		Code:       code,
		Path:       path,
		UpdateTime: now,
		CreateTime: now,
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
			if value, ok := data.Path.(string); ok {
				paths = append(paths, value)
			}
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
	roleIds, err := s.repositoryRolePermission.SelectRoleIdsByPermissionCodes(codes)
	if err != nil {
		return err
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		err = s.repositoryRolePermission.DeleteByPermissionCodes(codes, tx)
		if err != nil {
			return err
		}
		return s.repositoryPermission.DeleteByIds(ids, tx)
	})
	if err != nil {
		return err
	}
	s.clearCacheWithPermissions(paths)
	s.clearCacheWithRolePermissions(roleIds)
	s.clearCacheWithEnums([]string{"permission"})
	return nil
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

// GetAllMenuTreeByDirAndPage 获取所有菜单树状数据
func (s *service) GetAllMenuTreeByDirAndPage() []*MenuTree {
	menus, err := s.repositoryMenu.SelectWithType([]int{menu_type.Dir, menu_type.Page})
	if err != nil {
		return nil
	}
	return s.buildMenuTree(menus)
}

// GetAllMenuTreeByStatusEnabled 获取菜单树带控件
func (s *service) GetAllMenuTreeByStatusEnabled() []*MenuTree {
	menus, err := s.repositoryMenu.SelectWithStatus(menu_status.Enable)
	if err != nil {
		return nil
	}
	return s.buildMenuTree(menus)
}

// GetRoleInfoById 获取角色详情
func (s *service) GetRoleInfoById(id int64) (*Role, error) {
	info, err := s.repositoryRole.FindById(id)
	if err != nil {
		return nil, errors.Join(errors.New("查询角色失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("该角色不存在")
	}
	menuIds, err := s.repositoryRoleMenu.SelectMenuIdsByRoleId(info.Id)
	if err != nil {
		return nil, errors.Join(errors.New("查询角色菜单失败"), err)
	}
	return &Role{
		Id:         info.Id,
		Name:       info.Name,
		Status:     info.Status,
		MenuIds:    menuIds,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}, nil
}

// CreateRole 创建角色
func (s *service) CreateRole(name string, status int, menuIds []int64) error {
	if _, ok := role_status.Status()[status]; !ok {
		return errors.New("状态数据不合法")
	}
	queryMenuIds := make([]int64, 0, len(menuIds))
	menuIdMap := make(map[int64]struct{}, len(menuIds))
	for _, menuId := range menuIds {
		if menuId <= 0 {
			continue
		}
		if _, ok := menuIdMap[menuId]; ok {
			continue
		}
		menuIdMap[menuId] = struct{}{}
		queryMenuIds = append(queryMenuIds, menuId)
	}
	validMenuIds := make([]int64, 0, len(queryMenuIds))
	if len(queryMenuIds) > 0 {
		menus, err := s.repositoryMenu.SelectByIds(queryMenuIds)
		if err != nil {
			return errors.Join(errors.New("查询菜单失败"), err)
		}
		menuIdMap = make(map[int64]struct{}, len(menus))
		for _, item := range menus {
			if item == nil || item.Id <= 0 {
				continue
			}
			menuIdMap[item.Id] = struct{}{}
		}
		for _, menuId := range queryMenuIds {
			if _, ok := menuIdMap[menuId]; !ok {
				continue
			}
			validMenuIds = append(validMenuIds, menuId)
		}
	}
	permissionCodes := make([]string, 0)
	if len(validMenuIds) > 0 {
		menuPermissions, err := s.repositoryMenuPermission.SelectByMenuIds(validMenuIds)
		if err != nil {
			return errors.Join(errors.New("查询菜单权限失败"), err)
		}
		permissionCodeMap := make(map[string]struct{}, len(menuPermissions))
		permissionCodes = make([]string, 0, len(menuPermissions))
		for _, item := range menuPermissions {
			if item == nil || item.PermissionCode == "" {
				continue
			}
			if _, ok := permissionCodeMap[item.PermissionCode]; ok {
				continue
			}
			permissionCodeMap[item.PermissionCode] = struct{}{}
			permissionCodes = append(permissionCodes, item.PermissionCode)
		}
	}
	now := str.DateTime(time.Now())
	roleData := &model.Role{
		Name:       name,
		Status:     status,
		UpdateTime: now,
		CreateTime: now,
	}
	err := s.database.Transaction(func(tx *gorm.DB) error {
		err := s.repositoryRole.Create(roleData, tx)
		if err != nil {
			return errors.Join(errors.New("创建角色失败"), err)
		}
		if roleData.Id <= 0 {
			return errors.New("创建角色失败")
		}
		if len(validMenuIds) > 0 {
			roleMenus := make([]*model.RoleMenu, 0, len(validMenuIds))
			for _, menuId := range validMenuIds {
				roleMenus = append(roleMenus, &model.RoleMenu{
					RoleId:     roleData.Id,
					MenuId:     menuId,
					UpdateTime: now,
					CreateTime: now,
				})
			}
			err = s.repositoryRoleMenu.Creates(roleMenus, tx)
			if err != nil {
				return errors.Join(errors.New("创建角色菜单关联失败"), err)
			}
		}
		if len(permissionCodes) > 0 {
			rolePermissions := make([]*model.RolePermission, 0, len(permissionCodes))
			for _, permissionCode := range permissionCodes {
				rolePermissions = append(rolePermissions, &model.RolePermission{
					RoleId:         roleData.Id,
					PermissionCode: permissionCode,
					CreateTime:     now,
				})
			}
			err = s.repositoryRolePermission.Creates(rolePermissions, tx)
			if err != nil {
				return errors.Join(errors.New("创建角色权限关联失败"), err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.clearCacheWithEnums([]string{"role_names"})
	return nil
}

// AssignRoleToUser 分配用户角色
func (s *service) AssignRoleToUser(userId int64, roleIds []int64) error {
	if userId <= 0 {
		return errors.New("用户ID错误")
	}
	userLockKey := constant.LockUserRoleAssign + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(userLockKey, constant.LockTimeUserRoleAssign) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(userLockKey)

	currentRoleIds, err := s.repositoryUserRole.SelectRoleIdsByUserId(userId)
	if err != nil {
		return errors.Join(errors.New("查询用户角色失败"), err)
	}
	queryRoleIds := make([]int64, 0, len(roleIds))
	roleIdMap := make(map[int64]struct{}, len(roleIds))
	for _, roleId := range roleIds {
		if roleId <= 0 {
			return errors.New("角色ID错误")
		}
		if _, ok := roleIdMap[roleId]; ok {
			continue
		}
		roleIdMap[roleId] = struct{}{}
		queryRoleIds = append(queryRoleIds, roleId)
	}
	lockRoleIds := make([]int64, 0, len(currentRoleIds)+len(queryRoleIds))
	lockRoleIds = append(lockRoleIds, currentRoleIds...)
	lockRoleIds = append(lockRoleIds, queryRoleIds...)
	lockKeys, err := s.lockRoles(lockRoleIds)
	if err != nil {
		return err
	}
	defer func() {
		if len(lockKeys) > 0 {
			s.cache.UnLock(lockKeys...)
		}
	}()

	validRoleIds := make([]int64, 0, len(queryRoleIds))
	if len(queryRoleIds) > 0 {
		roles, err := s.repositoryRole.SelectByIds(queryRoleIds)
		if err != nil {
			return errors.Join(errors.New("查询角色失败"), err)
		}
		if len(roles) != len(queryRoleIds) {
			return errors.New("存在角色不存在")
		}
		roleIdMap = make(map[int64]struct{}, len(roles))
		for _, item := range roles {
			if item == nil || item.Id <= 0 {
				continue
			}
			roleIdMap[item.Id] = struct{}{}
		}
		for _, roleId := range queryRoleIds {
			if _, ok := roleIdMap[roleId]; !ok {
				continue
			}
			validRoleIds = append(validRoleIds, roleId)
		}
	}

	currentRoleMap := make(map[int64]struct{}, len(currentRoleIds))
	for _, roleId := range currentRoleIds {
		if roleId <= 0 {
			continue
		}
		currentRoleMap[roleId] = struct{}{}
	}
	if len(currentRoleMap) == len(validRoleIds) {
		same := true
		for _, roleId := range validRoleIds {
			if _, ok := currentRoleMap[roleId]; !ok {
				same = false
				break
			}
		}
		if same {
			return nil
		}
	}

	now := str.DateTime(time.Now())
	userRoles := make([]*model.UserRole, 0, len(validRoleIds))
	for _, roleId := range validRoleIds {
		userRoles = append(userRoles, &model.UserRole{
			UserId:     userId,
			RoleId:     roleId,
			UpdateTime: now,
			CreateTime: now,
		})
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		err = s.repositoryUserRole.DeleteByUserIds([]int64{userId}, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户角色关联失败"), err)
		}
		if len(userRoles) == 0 {
			return nil
		}
		err = s.repositoryUserRole.Creates(userRoles, tx)
		if err != nil {
			return errors.Join(errors.New("创建用户角色关联失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.clearCacheWithUserRoles([]int64{userId})
	return nil
}

// DeleteRoleByIds 根据ID集合删除角色
func (s *service) DeleteRoleByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	lockKeys, err := s.lockRoles(ids)
	if err != nil {
		return err
	}
	defer func() {
		if len(lockKeys) > 0 {
			s.cache.UnLock(lockKeys...)
		}
	}()
	total, err := s.repositoryUserRole.CountByRoleIds(ids)
	if err != nil {
		return errors.Join(errors.New("查询角色使用情况失败"), err)
	}
	if total > 0 {
		return errors.New("存在用户正在使用该角色，不允许删除")
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		err = s.repositoryRolePermission.DeleteByRoleIds(ids, tx)
		if err != nil {
			return errors.Join(errors.New("删除角色权限关联失败"), err)
		}
		err = s.repositoryRoleMenu.DeleteByRoleIds(ids, tx)
		if err != nil {
			return errors.Join(errors.New("删除角色菜单关联失败"), err)
		}
		err = s.repositoryRole.DeleteByIds(ids, tx)
		if err != nil {
			return errors.Join(errors.New("删除角色数据失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.clearCacheWithRolePermissions(ids)
	s.clearCacheWithRoleMenus(ids)
	s.clearCacheWithEnums([]string{"role_names"})
	return nil
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
	if info.Type != menu_type.Dir {
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
	menus, err := s.repositoryMenu.SelectByParentIdWithType(id, int(menu_type.Control))
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

// CreateMenuOrControl 创建菜单或控件数据
func (s *service) CreateMenuOrControl(parentId int64, code string, name string, icon string, Type int, sort int, statu int, permissionCodes []string) error {
	if parentId <= 0 {
		parentId = 0
	}
	types := menu_type.Types()
	status := menu_status.Status()
	if _, ok := status[statu]; !ok {
		return errors.New("状态数据不合法")
	}
	if _, ok := types[Type]; !ok {
		return errors.New("类型数据不合法")
	}
	if parentId == 0 && Type == menu_type.Control {
		return errors.New(fmt.Sprintf("%s类型必须添加在%s类型下", types[menu_type.Control], types[menu_type.Page]))
	}
	info, err := s.repositoryMenu.FindByCode(code)
	if err != nil {
		return errors.Join(errors.New("创建失败"), err)
	}
	if info != nil && info.Id > 0 {
		return errors.New("该数据标识已被占用")
	}
	if parentId > 0 {
		info, err = s.repositoryMenu.FindById(parentId)
		if err != nil {
			return errors.Join(errors.New("创建失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("该父级数据不存在")
		}
		if Type == menu_type.Control && info.Type != menu_type.Page {
			return errors.New(fmt.Sprintf("%s类型必须添加在%s类型下", types[menu_type.Control], types[menu_type.Page]))
		}
		if (Type == menu_type.Page || Type == menu_type.Dir) && info.Type != menu_type.Dir {
			return errors.New(fmt.Sprintf("%s或%s类型必须添加在%s类型下", types[menu_type.Dir], types[menu_type.Page], types[menu_type.Dir]))
		}
	}
	if Type == menu_type.Dir && len(permissionCodes) > 0 {
		return errors.New(fmt.Sprintf("%s类型不能绑定权限", types[menu_type.Dir]))
	}
	now := str.DateTime(time.Now())
	menuData := &model.Menu{
		ParentId:   parentId,
		Code:       code,
		Name:       name,
		Type:       Type,
		Icon:       icon,
		Sort:       sort,
		Status:     statu,
		CreateTime: now,
		UpdateTime: now,
	}
	permissionCodeMap := make(map[string]struct{}, len(permissionCodes))
	permissionList := make([]*model.MenuPermission, 0, len(permissionCodes))
	for _, permissionCode := range permissionCodes {
		if permissionCode == "" {
			continue
		}
		if _, ok := permissionCodeMap[permissionCode]; ok {
			continue
		}
		permissionCodeMap[permissionCode] = struct{}{}
		permissionList = append(permissionList, &model.MenuPermission{
			PermissionCode: permissionCode,
			CreateTime:     now,
		})
	}
	if len(permissionList) > 0 {
		codes := make([]string, 0, len(permissionList))
		for _, item := range permissionList {
			if item == nil || item.PermissionCode == "" {
				continue
			}
			codes = append(codes, item.PermissionCode)
		}
		permissions, err := s.repositoryPermission.SelectByCodes(codes)
		if err != nil {
			return errors.Join(errors.New("创建失败"), err)
		}
		permissionMap := make(map[string]struct{}, len(permissions))
		for _, item := range permissions {
			if item == nil || item.Code == "" {
				continue
			}
			permissionMap[item.Code] = struct{}{}
		}
		validPermissionList := make([]*model.MenuPermission, 0, len(permissionList))
		for _, item := range permissionList {
			if item == nil || item.PermissionCode == "" {
				continue
			}
			if _, ok := permissionMap[item.PermissionCode]; !ok {
				continue
			}
			validPermissionList = append(validPermissionList, item)
		}
		permissionList = validPermissionList
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		err = s.repositoryMenu.Create(menuData, tx)
		if err != nil {
			return err
		}
		if menuData.Id <= 0 || len(permissionList) == 0 {
			return nil
		}
		for _, item := range permissionList {
			if item == nil {
				continue
			}
			item.MenuId = menuData.Id
		}
		return s.repositoryMenuPermission.Creates(permissionList, tx)
	})
	if err != nil {
		return errors.Join(errors.New("创建失败"), err)
	}
	s.clearCacheWithMenus([]*model.Menu{menuData})
	return nil
}

// UpdateMenuOrControl 更新菜单或控件数据
func (s *service) UpdateMenuOrControl(id int64, data *menu2.UpdateMenu, permissionCodes []string) error {
	if id <= 0 || data == nil {
		return errors.New("参数错误")
	}
	info, err := s.repositoryMenu.FindById(id)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return errors.New("该数据不存在")
	}
	types := menu_type.Types()
	status := menu_status.Status()
	hasUpdateData := false
	newParentId := int64(-1)
	if data.ParentId != nil {
		parentId, ok := data.ParentId.(int64)
		if !ok {
			return errors.New("父级数据不合法")
		}
		if parentId <= 0 {
			parentId = 0
		}
		if parentId == info.Id {
			return errors.New("父级不能是自己")
		}
		if parentId > 0 {
			parentInfo, err := s.repositoryMenu.FindById(parentId)
			if err != nil {
				return errors.Join(errors.New("更新失败"), err)
			}
			if parentInfo == nil || parentInfo.Id <= 0 {
				return errors.New("该父级数据不存在")
			}
			if info.Type == menu_type.Control && parentInfo.Type != menu_type.Page {
				return errors.New(fmt.Sprintf("%s类型必须添加在%s类型下", types[menu_type.Control], types[menu_type.Page]))
			}
			if (info.Type == menu_type.Page || info.Type == menu_type.Dir) && parentInfo.Type != menu_type.Dir {
				return errors.New(fmt.Sprintf("%s或%s类型必须添加在%s类型下", types[menu_type.Dir], types[menu_type.Page], types[menu_type.Dir]))
			}
			menus, err := s.repositoryMenu.SelectAll()
			if err != nil {
				return errors.Join(errors.New("更新失败"), err)
			}
			childrenMap := make(map[int64][]int64, len(menus))
			for _, item := range menus {
				if item == nil || item.Id <= 0 {
					continue
				}
				childrenMap[item.ParentId] = append(childrenMap[item.ParentId], item.Id)
			}
			targetMap := make(map[int64]struct{}, len(menus))
			var loadChildren func(menuId int64)
			loadChildren = func(menuId int64) {
				if _, ok = targetMap[menuId]; ok {
					return
				}
				targetMap[menuId] = struct{}{}
				for _, childId := range childrenMap[menuId] {
					loadChildren(childId)
				}
			}
			loadChildren(info.Id)
			if _, ok = targetMap[parentId]; ok {
				return errors.New("父级不能是自己的子级")
			}
		} else if info.Type == menu_type.Control {
			return errors.New(fmt.Sprintf("%s类型必须添加在%s类型下", types[menu_type.Control], types[menu_type.Page]))
		}
		hasUpdateData = true
		newParentId = parentId
	}
	if data.Name != nil || data.Icon != nil || data.Sort != nil {
		hasUpdateData = true
	}
	if data.Status != nil {
		statusValue, ok := data.Status.(int)
		if !ok {
			return errors.New("状态数据不合法")
		}
		if _, ok = status[statusValue]; !ok {
			return errors.New("状态数据不合法")
		}
		hasUpdateData = true
	}
	permissionCodesChanged := permissionCodes != nil
	permissionList := make([]*model.MenuPermission, 0)
	newPermissionCodeMap := make(map[string]struct{})
	if permissionCodesChanged {
		now := str.DateTime(time.Now())
		codes := permissionCodes
		permissionCodeMap := make(map[string]struct{}, len(codes))
		for _, permissionCode := range codes {
			if permissionCode == "" {
				continue
			}
			if _, ok := permissionCodeMap[permissionCode]; ok {
				continue
			}
			permissionCodeMap[permissionCode] = struct{}{}
			permissionList = append(permissionList, &model.MenuPermission{
				PermissionCode: permissionCode,
				CreateTime:     now,
			})
		}
		if info.Type == menu_type.Dir && len(permissionList) > 0 {
			return errors.New(fmt.Sprintf("%s类型不能绑定权限", types[menu_type.Dir]))
		}
		if len(permissionList) > 0 {
			codes = make([]string, 0, len(permissionList))
			for _, item := range permissionList {
				if item == nil || item.PermissionCode == "" {
					continue
				}
				codes = append(codes, item.PermissionCode)
			}
			permissions, err := s.repositoryPermission.SelectByCodes(codes)
			if err != nil {
				return errors.Join(errors.New("更新失败"), err)
			}
			permissionMap := make(map[string]struct{}, len(permissions))
			for _, item := range permissions {
				if item == nil || item.Code == "" {
					continue
				}
				permissionMap[item.Code] = struct{}{}
			}
			validPermissionList := make([]*model.MenuPermission, 0, len(permissionList))
			for _, item := range permissionList {
				if item == nil || item.PermissionCode == "" {
					continue
				}
				if _, ok := permissionMap[item.PermissionCode]; !ok {
					continue
				}
				validPermissionList = append(validPermissionList, item)
			}
			permissionList = validPermissionList
		}
		for _, item := range permissionList {
			if item == nil || item.PermissionCode == "" {
				continue
			}
			newPermissionCodeMap[item.PermissionCode] = struct{}{}
		}
		oldPermissionCodes, err := s.repositoryMenuPermission.SelectPermissionCodesByMenuId(info.Id)
		if err != nil {
			return errors.Join(errors.New("更新失败"), err)
		}
		oldPermissionCodeMap := make(map[string]struct{}, len(oldPermissionCodes))
		for _, permissionCode := range oldPermissionCodes {
			if permissionCode == "" {
				continue
			}
			oldPermissionCodeMap[permissionCode] = struct{}{}
		}
		if len(oldPermissionCodeMap) == len(newPermissionCodeMap) {
			same := true
			for permissionCode := range newPermissionCodeMap {
				if _, ok := oldPermissionCodeMap[permissionCode]; !ok {
					same = false
					break
				}
			}
			if same {
				permissionCodesChanged = false
			}
		}
	}
	if !hasUpdateData && !permissionCodesChanged {
		return nil
	}
	//更新数据 并重构权限
	err = s.database.Transaction(func(tx *gorm.DB) error {
		if hasUpdateData {
			err = s.repositoryMenu.UpdateById(info.Id, data, tx)
			if err != nil {
				return err
			}
		}
		if permissionCodesChanged {
			err = s.repositoryMenuPermission.DeleteByMenuIds([]int64{info.Id}, tx)
			if err != nil {
				return err
			}
			if len(permissionList) > 0 {
				for _, item := range permissionList {
					if item == nil {
						continue
					}
					item.MenuId = info.Id
				}
				err = s.repositoryMenuPermission.Creates(permissionList, tx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	//清理缓存
	clearMenus := []*model.Menu{info}
	if newParentId >= 0 && newParentId != info.ParentId {
		clearMenus = append(clearMenus, &model.Menu{ParentId: newParentId})
	}
	s.clearCacheWithMenus(clearMenus)
	return nil
}

// SelectRole 查询角色数据
func (s *service) SelectRole(where *role.SelectRole, orderByField string, orderByType string, page int, pageSize int) ([]*model.Role, int64, error) {
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("状态参数错误")
		}
		if _, ok := role_status.Status()[status]; !ok {
			return nil, 0, errors.New("状态参数不合法")
		}
	}
	return s.repositoryRole.Select(where, orderByField, orderByType, page, pageSize)
}
