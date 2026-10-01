package balance_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectBalanceLog, queryPage *page.Query) (*page.Result[*model.BalanceLog], error) {
	query := r.Database.Db.Model(&model.BalanceLog{})
	if where != nil {
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.Describe != nil {
			query = query.Where("cloud_balance_logs.describe LIKE ?", "%"+*where.Describe+"%")
		}
		if where.Type != nil {
			query = query.Where("type = ?", *where.Type)
		}
		if where.CreateTimeStart != nil {
			query = query.Where("create_time >= ?", *where.CreateTimeStart)
		}
		if where.CreateTimeEnd != nil {
			query = query.Where("create_time <= ?", *where.CreateTimeEnd)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
