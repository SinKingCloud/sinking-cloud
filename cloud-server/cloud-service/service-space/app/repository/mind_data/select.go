package mind_data

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"gorm.io/gorm/clause"
)

// Select 查询导图数据列表，不查询 value 字段
func (r *Repository) Select(where *SelectMindData, queryPage *page.Query) (*page.Result[*MindData], error) {
	if queryPage == nil {
		return page.New(queryPage, 0, make([]*MindData, 0)), nil
	}
	query := r.Database.Db.Model(&model.MindData{})
	if where != nil {
		if where.UserId != nil {
			query = query.Where("user_id = ?", *where.UserId)
		}
		if where.TypeId != nil {
			query = query.Where("type_id = ?", *where.TypeId)
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where(clause.Or(
				clause.Like{Column: clause.Column{Name: "title"}, Value: keyword},
				clause.Like{Column: clause.Column{Name: "describe"}, Value: keyword},
			))
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
