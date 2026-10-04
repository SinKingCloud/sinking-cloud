package user_log

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询用户操作日志
func (r *Repository) Select(where *SelectUserLog, queryPage *page.Query) (*page.Result[*model.UserLog], error) {
	query := r.Database.Db.Model(&model.UserLog{})
	if where != nil {
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.RequestId != nil {
			query = query.Where("request_id = ?", *where.RequestId)
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
		if where.RequestIp != nil {
			query = query.Where("request_ip LIKE ?", "%"+*where.RequestIp+"%")
		}
		if where.Location != nil {
			query = query.Where("location LIKE ?", "%"+*where.Location+"%")
		}
		if where.Title != nil {
			query = query.Where("title LIKE ?", "%"+*where.Title+"%")
		}
		if where.Content != nil {
			query = query.Where("content LIKE ?", "%"+*where.Content+"%")
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(request_id = ? OR request_ip = ? OR location LIKE ? OR title LIKE ? OR content LIKE ?)", *where.Keyword, *where.Keyword, keyword, keyword, keyword)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
