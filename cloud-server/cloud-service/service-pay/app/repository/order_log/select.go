package order_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectOrderLog, queryPage *page.Query) (*page.Result[*OrderLog], error) {
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
	return r.Repository.SelectPage(query, queryPage)
}
