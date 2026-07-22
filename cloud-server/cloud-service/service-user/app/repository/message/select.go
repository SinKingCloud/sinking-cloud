package message

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询消息列表，不查询 content 字段。
func (r *Repository) Select(where *SelectMessage, queryPage *page.Query) (*page.Result[*Message], error) {
	query := r.Database.Db.Model(&model.Message{})
	if where != nil {
		if where.UserId > 0 {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.Title != "" {
			query = query.Where("title LIKE ?", "%"+where.Title+"%")
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.ReadTimeStart != "" {
			query = query.Where("read_time >= ?", where.ReadTimeStart)
		}
		if where.ReadTimeEnd != "" {
			query = query.Where("read_time <= ?", where.ReadTimeEnd)
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
	query = query.Select("id", "title", "status", "read_time", "create_time", "update_time")
	return r.Repository.SelectPage(query, queryPage)
}
