package log

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectLog, queryPage *page.Query) (*page.Result[*model.Log], error) {
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
	return r.Repository.SelectPage(query, queryPage)
}
