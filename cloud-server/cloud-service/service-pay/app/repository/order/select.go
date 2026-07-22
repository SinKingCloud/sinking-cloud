package order

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

func (r *Repository) Select(where *SelectOrder, queryPage *page.Query) (*page.Result[*Order], error) {
	query := r.Database.Db.Model(&model.Order{})
	if where != nil {
		if where.Id != "" {
			query = query.Where("id = ?", where.Id)
		}
		if where.UserId != "" {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.ChannelId != "" {
			query = query.Where("channel_id = ?", where.ChannelId)
		}
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.TradeNo != "" {
			query = query.Where("trade_no = ?", where.TradeNo)
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
	return r.Repository.SelectPage(query, queryPage)
}
