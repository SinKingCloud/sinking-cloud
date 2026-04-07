package user_role

import "service-admin/app/model"

// CountByRoleIds 根据角色ID集合统计用户角色关联数量
func (r *Repository) CountByRoleIds(roleIds []int64) (int64, error) {
	if len(roleIds) == 0 {
		return 0, nil
	}
	var total int64
	err := r.Database.Db.
		Model(&model.UserRole{}).
		Where("role_id IN ?", roleIds).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
