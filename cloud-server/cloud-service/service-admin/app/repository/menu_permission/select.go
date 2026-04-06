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
