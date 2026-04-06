package role_menu

import "service-admin/app/model"

// SelectMenuIdsByRoleId 根据角色ID查询菜单ID集合
func (r *Repository) SelectMenuIdsByRoleId(roleId int64) ([]int64, error) {
	var menuIds []int64
	err := r.Database.Db.
		Model(&model.RoleMenu{}).
		Distinct("menu_id").
		Where("role_id = ?", roleId).
		Pluck("menu_id", &menuIds).Error
	if err != nil {
		return nil, err
	}
	return menuIds, nil
}

// SelectRoleIdsByMenuIds 根据菜单ID集合查询角色ID集合
func (r *Repository) SelectRoleIdsByMenuIds(menuIds []int64) (roleIds []int64, err error) {
	if len(menuIds) == 0 {
		return nil, nil
	}
	roleIdMap := make(map[int64]struct{}, len(menuIds))
	err = r.Database.BatchExecute(menuIds, 1000, func(batch interface{}) error {
		var batchRoleIds []int64
		err := r.Database.Db.
			Model(&model.RoleMenu{}).
			Distinct("role_id").
			Where("menu_id IN ?", batch.([]int64)).
			Pluck("role_id", &batchRoleIds).Error
		if err != nil {
			return err
		}
		for _, roleId := range batchRoleIds {
			if _, ok := roleIdMap[roleId]; ok {
				continue
			}
			roleIdMap[roleId] = struct{}{}
			roleIds = append(roleIds, roleId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}
