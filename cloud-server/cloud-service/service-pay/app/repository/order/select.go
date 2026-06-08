package order

import "service-pay/app/model"

func (r *Repository) Select(where *SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*Order, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Order{})
	if where != nil {
		if where.UserId != "" {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.ChannelId != "" {
			query = query.Where("channel_id = ?", where.ChannelId)
		}
		if where.TradeNo != "" {
			query = query.Where("trade_no = ?", where.TradeNo)
		}
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.OutTradeNo != "" {
			query = query.Where("out_trade_no = ?", where.OutTradeNo)
		}
		if where.PayStatus != "" {
			query = query.Where("pay_status = ?", where.PayStatus)
		}
		if where.PayType != "" {
			query = query.Where("pay_type = ?", where.PayType)
		}
		if where.PayTimeStart != "" {
			query = query.Where("pay_time >= ?", where.PayTimeStart)
		}
		if where.PayTimeEnd != "" {
			query = query.Where("pay_time <= ?", where.PayTimeEnd)
		}
		if where.CompleteTimeStart != "" {
			query = query.Where("complete_time >= ?", where.CompleteTimeStart)
		}
		if where.CompleteTimeEnd != "" {
			query = query.Where("complete_time <= ?", where.CompleteTimeEnd)
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
