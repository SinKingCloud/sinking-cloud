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

// SelectByUserIds 根据用户ID集合查询用户角色关联
func (r *Repository) SelectByUserIds(userIds []int64) ([]*model.UserRole, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	list := make([]*model.UserRole, 0)
	err := r.Database.BatchExecute(userIds, 1000, func(batch interface{}) error {
		var items []*model.UserRole
		err := r.Database.Db.
			Model(&model.UserRole{}).
			Where("user_id IN ?", batch.([]int64)).
			Order("user_id ASC, role_id ASC").
			Find(&items).Error
		if err != nil {
			return err
		}
		list = append(list, items...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
