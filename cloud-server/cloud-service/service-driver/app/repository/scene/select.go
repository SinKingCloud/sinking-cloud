package scene

import "service-driver/app/model"

// SelectIdNameMap 查询id name 映射
func (r *Repository) SelectIdNameMap() (map[int64]string, error) {
	var scene []struct {
		Id   int64
		Name string
	}
	err := r.Database.Db.
		Model(&model.Scene{}).
		Select("id", "name").
		Order("id ASC").
		Find(&scene).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(scene))
	for _, c := range scene {
		result[c.Id] = c.Name
	}
	return result, nil
}

// Select 查询数据
func (r *Repository) Select(where *SelectScene, orderByField string, orderByType string, page int, pageSize int) (list []*model.Scene, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Scene{})
	if where != nil {
		if where.ChannelId != "" {
			query = query.Where("channel_id = ?", where.ChannelId)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.Login != "" {
			query = query.Where("login = ?", where.Login)
		}
		if where.Name != "" {
			query = query.Where("name LIKE ?", "%"+where.Name+"%")
		}
		if where.Remark != "" {
			query = query.Where("remark LIKE ?", "%"+where.Remark+"%")
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
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(r.GetOrderBy(orderByField, orderByType)).Find(&list).Error
	return
}
