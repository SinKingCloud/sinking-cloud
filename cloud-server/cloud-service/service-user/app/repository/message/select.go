package message

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询消息列表，不查询 content 字段。
func (r *Repository) Select(where *SelectMessage, queryPage *page.Query) (*page.Result[*Message], error) {
	query := r.Database.Db.Model(&model.Message{})
	if where != nil {
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(title LIKE ? OR content LIKE ?)", keyword, keyword)
		}
		if where.Title != nil {
			query = query.Where("title LIKE ?", "%"+*where.Title+"%")
		}
		if where.Content != nil {
			query = query.Where("content LIKE ?", "%"+*where.Content+"%")
		}
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
		}
		if where.ReadTimeStart != nil {
			query = query.Where("read_time >= ?", *where.ReadTimeStart)
		}
		if where.ReadTimeEnd != nil {
			query = query.Where("read_time <= ?", *where.ReadTimeEnd)
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
	query = query.Select("id", "title", "status", "read_time", "create_time", "update_time")
	return r.Repository.SelectPage(query, queryPage)
}
