package balance_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectBalanceLog, queryPage *page.Query) (*page.Result[*model.BalanceLog], error) {
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
	return r.Repository.SelectPage(query, queryPage)
}
