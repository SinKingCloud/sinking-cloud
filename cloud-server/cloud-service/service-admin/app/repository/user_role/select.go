package user_role

import "service-admin/app/model"

// SelectRoleIdsByUserId 根据用户ID查询角色ID列表
func (r *Repository) SelectRoleIdsByUserId(userId int64) ([]int64, error) {
	var roleIds []int64
	err := r.Database.Db.
		Model(&model.UserRole{}).
		Distinct("role_id").
		Where("user_id = ?", userId).
		Pluck("role_id", &roleIds).Error
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}
