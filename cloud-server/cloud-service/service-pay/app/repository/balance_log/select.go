package balance_log

import "service-pay/app/model"

// Select 查询数据
func (r *Repository) Select(where *SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.BalanceLog, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.BalanceLog{})
	if where != nil {
		if where.UserId != "" {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
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
