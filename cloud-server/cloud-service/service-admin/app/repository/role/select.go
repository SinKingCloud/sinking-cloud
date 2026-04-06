package role

import "service-admin/app/model"

// SelectAll 查询全部角色
func (r *Repository) SelectAll() ([]*model.Role, error) {
	var roles []*model.Role
	err := r.Database.Db.
		Model(&model.Role{}).
		Order("id ASC").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
