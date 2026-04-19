package permission

import "service-admin/app/model"

// SelectCodeNameMap 查询code name 映射
func (r *Repository) SelectCodeNameMap() (map[string]string, error) {
	var data []struct {
		Code string `gorm:"column:code" json:"code"`
		Name string `gorm:"column:name" json:"name"`
	}
	err := r.Database.Db.
		Model(&model.Permission{}).
		Select("code", "name").
		Order("id ASC").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(data))
	for _, c := range data {
		result[c.Code] = c.Name
	}
	return result, nil
}

// SelectByCodes 根据标识集合查询权限
func (r *Repository) SelectByCodes(codes []string) ([]*model.Permission, error) {
	if len(codes) == 0 {
		return nil, nil
	}
	var permissionList []*model.Permission
	err := r.Database.BatchExecute(codes, 1000, func(batch interface{}) error {
		batchCodes := batch.([]string)
		var permissions []*model.Permission
		err := r.Database.Db.
			Model(&model.Permission{}).
			Where("code IN ?", batchCodes).
			Find(&permissions).Error
		if err != nil {
			return err
		}
		permissionList = append(permissionList, permissions...)
		return nil
	})
	return permissionList, err
}

// SelectByIds 根据ID集合查询权限
func (r *Repository) SelectByIds(ids []int64) ([]*model.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var permissionList []*model.Permission
	err := r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		batchIds := batch.([]int64)
		var permissions []*model.Permission
		err := r.Database.Db.
			Model(&model.Permission{}).
			Where("id IN ?", batchIds).
			Find(&permissions).Error
		if err != nil {
			return err
		}
		permissionList = append(permissionList, permissions...)
		return nil
	})
	return permissionList, err
}

// Select 查询数据
func (r *Repository) Select(where *SelectPermission, orderByField string, orderByType string, page int, pageSize int) (list []*model.Permission, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Permission{})
	if where != nil {
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
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
