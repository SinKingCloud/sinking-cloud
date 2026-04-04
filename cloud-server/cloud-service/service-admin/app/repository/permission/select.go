package permission

import "service-admin/app/model"

// Select 查询数据
func (r *Repository) Select(where *SelectPermission, orderByField string, orderByType string, page int, pageSize int) (list []*model.Permission, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Permission{})
	if where.Code != "" {
		query = query.Where("code LIKE ?", "%"+where.Code+"%")
	}
	if where.Path != "" {
		query = query.Where("path LIKE ?", "%"+where.Path+"%")
	}
	if where.Name != "" {
		query = query.Where("name LIKE ?", "%"+where.Name+"%")
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
