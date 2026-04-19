package channel

import "service-driver/app/model"

// SelectIdNameMap 查询id name 映射
func (r *Repository) SelectIdNameMap() (map[int64]string, error) {
	var channel []struct {
		Id   int64
		Name string
	}
	err := r.Database.Db.
		Model(&model.Channel{}).
		Select("id", "name").
		Order("id ASC").
		Find(&channel).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(channel))
	for _, c := range channel {
		result[c.Id] = c.Name
	}
	return result, nil
}

// Select 查询数据
func (r *Repository) Select(where *SelectChannel, orderByField string, orderByType string, page int, pageSize int) (list []*Channel, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Channel{})
	if where != nil {
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
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
