package role_permission

import "service-admin/app/model"

// FindPermissionCodesByRoleId 根据角色ID查询权限code列表（简化版）
func (r *Repository) FindPermissionCodesByRoleId(roleId int64) (permissionCodes []string, err error) {
	err = r.Database.Db.Model(&model.RolePermission{}).
		Distinct("permission_code").
		Where("role_id = ?", roleId).
		Pluck("permission_code", &permissionCodes).Error
	if err != nil {
		return nil, err
	}
	return permissionCodes, nil
}
