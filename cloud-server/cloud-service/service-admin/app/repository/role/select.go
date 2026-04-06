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

// Select 查询数据
func (r *Repository) Select(where *SelectRole, orderByField string, orderByType string, page int, pageSize int) (list []*model.Role, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Permission{})
	if where.Name != "" {
		query = query.Where("name LIKE ?", "%"+where.Name+"%")
	}
	if where.Status != "" {
		query = query.Where("status = ?", where.Status)
	}
	if where.CreateTimeStart != "" {
		query = query.Where("create_time >= ?", where.CreateTimeStart)
	}
	if where.CreateTimeEnd != "" {
		query = query.Where("create_time <= ?", where.CreateTimeEnd)
	}
	if where.UpdateTimeStart != "" {
		query = query.Where("update_time >= ?", where.UpdateTimeStart)
	}
	if where.UpdateTimeEnd != "" {
		query = query.Where("update_time <= ?", where.UpdateTimeEnd)
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
