package menu_permission

import "service-admin/app/model"

// SelectPermissionCodesByMenuId 根据菜单ID查询权限code列表
func (r *Repository) SelectPermissionCodesByMenuId(menuId int64) (permissionCodes []string, err error) {
	if menuId <= 0 {
		return nil, nil
	}
	err = r.Database.Db.
		Model(&model.MenuPermission{}).
		Distinct("permission_code").
		Where("menu_id = ?", menuId).
		Pluck("permission_code", &permissionCodes).Error
	if err != nil {
		return nil, err
	}
	return permissionCodes, nil
}

// SelectPermissionCodesByMenuIds 根据菜单ID集合查询权限code列表
func (r *Repository) SelectPermissionCodesByMenuIds(menuIds []int64) (permissionCodes []string, err error) {
	if len(menuIds) == 0 {
		return nil, nil
	}
	permissionCodeMap := make(map[string]struct{}, len(menuIds))
	err = r.Database.BatchExecute(menuIds, 1000, func(batch interface{}) error {
		var batchCodes []string
		err := r.Database.Db.
			Model(&model.MenuPermission{}).
			Distinct("permission_code").
			Where("menu_id IN ?", batch.([]int64)).
			Pluck("permission_code", &batchCodes).Error
		if err != nil {
			return err
		}
		for _, permissionCode := range batchCodes {
			if permissionCode == "" {
				continue
			}
			if _, ok := permissionCodeMap[permissionCode]; ok {
				continue
			}
			permissionCodeMap[permissionCode] = struct{}{}
			permissionCodes = append(permissionCodes, permissionCode)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return permissionCodes, nil
}

// SelectMenuIdsByPermissionCodes 根据权限code集合查询菜单ID列表
func (r *Repository) SelectMenuIdsByPermissionCodes(permissionCodes []string) (menuIds []int64, err error) {
	if len(permissionCodes) == 0 {
		return nil, nil
	}
	menuIdMap := make(map[int64]struct{}, len(permissionCodes))
	err = r.Database.BatchExecute(permissionCodes, 1000, func(batch interface{}) error {
		var batchMenuIds []int64
		err := r.Database.Db.
			Model(&model.MenuPermission{}).
			Distinct("menu_id").
			Where("permission_code IN ?", batch.([]string)).
			Pluck("menu_id", &batchMenuIds).Error
		if err != nil {
			return err
		}
		for _, menuId := range batchMenuIds {
			if menuId <= 0 {
				continue
			}
			if _, ok := menuIdMap[menuId]; ok {
				continue
			}
			menuIdMap[menuId] = struct{}{}
			menuIds = append(menuIds, menuId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return menuIds, nil
}

// SelectByMenuIds 根据菜单ID集合查询菜单权限关联数据
func (r *Repository) SelectByMenuIds(menuIds []int64) ([]*model.MenuPermission, error) {
	if len(menuIds) == 0 {
		return nil, nil
	}
	var menuPermissions []*model.MenuPermission
	err := r.Database.BatchExecute(menuIds, 1000, func(batch interface{}) error {
		var items []*model.MenuPermission
		err := r.Database.Db.
			Model(&model.MenuPermission{}).
			Where("menu_id IN ?", batch.([]int64)).
			Find(&items).Error
		if err != nil {
			return err
		}
		menuPermissions = append(menuPermissions, items...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return menuPermissions, nil
}
