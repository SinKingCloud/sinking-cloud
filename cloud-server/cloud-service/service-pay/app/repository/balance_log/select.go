package balance_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectBalanceLog, page *page.Query) (list []*model.BalanceLog, total int64, err error) {
	offset := page.PageSize * (page.Page - 1)
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
	err = query.Count(&total).Limit(page.PageSize).Offset(offset).Order(r.GetOrderBy(page.OrderByField, page.OrderByType)).Find(&list).Error
	return
}
