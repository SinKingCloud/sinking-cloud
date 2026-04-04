package role_permission

import "service-admin/app/model"

// SelectPermissionCodesByRoleId 根据角色ID查询权限code列表（简化版）
func (r *Repository) SelectPermissionCodesByRoleId(roleId int64) (permissionCodes []string, err error) {
	err = r.Database.Db.Model(&model.RolePermission{}).
		Distinct("permission_code").
		Where("role_id = ?", roleId).
		Pluck("permission_code", &permissionCodes).Error
	if err != nil {
		return nil, err
	}
	return permissionCodes, nil
}

// SelectRoleIdsByPermissionCodes 根据权限code集合查询角色ID集合
func (r *Repository) SelectRoleIdsByPermissionCodes(permissionCodes []string) (roleIds []int64, err error) {
	if len(permissionCodes) == 0 {
		return nil, nil
	}
	roleIdMap := make(map[int64]struct{}, len(permissionCodes))
	err = r.Database.BatchExecute(permissionCodes, 1000, func(batch interface{}) error {
		var batchRoleIds []int64
		err = r.Database.Db.Model(&model.RolePermission{}).
			Distinct("role_id").
			Where("permission_code IN ?", batch.([]string)).
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
