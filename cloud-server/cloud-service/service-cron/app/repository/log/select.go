package log

import "service-cron/app/model"

// Select 查询数据
func (r *Repository) Select(where *SelectLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Log{})
	if where != nil {
		if where.TaskId != "" {
			query = query.Where("task_id = ?", where.TaskId)
		}
		if where.CreateTimeStart != "" {
			query = query.Where("create_time >= ?", where.CreateTimeStart)
		}
		if where.CreateTimeEnd != "" {
			query = query.Where("create_time <= ?", where.CreateTimeEnd)
		}
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(r.GetOrderBy(orderByField, orderByType)).Find(&list).Error
	return
}
