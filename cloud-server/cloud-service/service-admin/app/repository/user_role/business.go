package user_role

// FindRoleIdsByUserId 根据用户ID查询角色ID列表
func (r *Repository) FindRoleIdsByUserId(userId int64) ([]int64, error) {
	var roleIds []int64
	err := r.Database.Db.
		Table("user_roles").
		Distinct("role_id").
		Where("user_id = ?", userId).
		Pluck("role_id", &roleIds).Error
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}
