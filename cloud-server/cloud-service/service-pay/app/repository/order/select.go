package order

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

func (r *Repository) Select(where *SelectOrder, queryPage *page.Query) (*page.Result[*Order], error) {
	query := r.Database.Db.Model(&model.Order{})
	if where != nil {
		if where.Id != nil {
			query = query.Where("id = ?", *where.Id)
		}
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.ChannelId != nil {
			query = query.Where("channel_id = ?", *where.ChannelId)
		}
		if where.Type != nil {
			query = query.Where("type = ?", *where.Type)
		}
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
		}
		if where.TradeNo != nil {
			query = query.Where("trade_no = ?", *where.TradeNo)
		}
		if where.PayStatus != nil {
			query = query.Where("pay_status = ?", *where.PayStatus)
		}
		if where.PayType != nil {
			query = query.Where("pay_type = ?", *where.PayType)
		}
		if where.PayTimeStart != nil {
			query = query.Where("pay_time >= ?", *where.PayTimeStart)
		}
		if where.PayTimeEnd != nil {
			query = query.Where("pay_time <= ?", *where.PayTimeEnd)
		}
		if where.CompleteTimeStart != nil {
			query = query.Where("complete_time >= ?", *where.CompleteTimeStart)
		}
		if where.CompleteTimeEnd != nil {
			query = query.Where("complete_time <= ?", *where.CompleteTimeEnd)
		}
		if where.CreateTimeStart != nil {
			query = query.Where("create_time >= ?", *where.CreateTimeStart)
		}
		if where.CreateTimeEnd != nil {
			query = query.Where("create_time <= ?", *where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != nil {
			query = query.Where("update_time >= ?", *where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != nil {
			query = query.Where("update_time <= ?", *where.UpdateTimeEnd)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
