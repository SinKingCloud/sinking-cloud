package order_log

import "service-pay/app/model"

// Select 查询数据
func (r *Repository) Select(where *SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*OrderLog, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.OrderLog{})
	if where != nil {
		if where.OrderId != "" {
			query = query.Where("order_id = ?", where.OrderId)
		}
		if where.OperatorId != "" {
			query = query.Where("operator_id = ?", where.OperatorId)
		}
		if where.OperatorType != "" {
			query = query.Where("operator_type = ?", where.OperatorType)
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
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
