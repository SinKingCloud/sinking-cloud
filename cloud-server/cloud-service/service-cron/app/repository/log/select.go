package log

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectLog, page *page.Query) (list []*model.Log, total int64, err error) {
	offset := page.PageSize * (page.Page - 1)
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
	err = query.Count(&total).Limit(page.PageSize).Offset(offset).Order(r.GetOrderBy(page.OrderByField, page.OrderByType)).Find(&list).Error
	return
}
