package task

import (
	"service-cron/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectAll 查询全部任务
func (r *Repository) SelectAll() (list []*model.Task, err error) {
	err = r.Database.Db.Model(&model.Task{}).Order("id asc").Find(&list).Error
	return
}

// Select 查询数据
func (r *Repository) Select(where *SelectTask, queryPage *page.Query) (*page.Result[*Task], error) {
	query := r.Database.Db.Model(&model.Task{})
	if where != nil {
		if where.Type != nil {
			query = query.Where("type = ?", *where.Type)
		}
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
		}
		if where.Concurrent != nil {
			query = query.Where("concurrent = ?", *where.Concurrent)
		}
		if where.Name != nil {
			query = query.Where("name LIKE ?", "%"+*where.Name+"%")
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
