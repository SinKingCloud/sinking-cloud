package mind_data

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"gorm.io/gorm/clause"
)

// Select 查询导图数据列表，不查询 value 字段
func (r *Repository) Select(where *SelectMindData, queryPage *page.Query) (list []*MindData, total int64, err error) {
	if queryPage == nil {
		return []*MindData{}, 0, nil
	}
	offset := queryPage.PageSize * (queryPage.Page - 1)
	query := r.Database.Db.Model(&model.MindData{})
	if where != nil {
		if where.UserId > 0 {
			query = query.Where("user_id = ?", where.UserId)
		}
		if where.TypeId > 0 {
			query = query.Where("type_id = ?", where.TypeId)
		}
		if where.Keyword != "" {
			keyword := "%" + where.Keyword + "%"
			query = query.Where(clause.Or(
				clause.Like{Column: clause.Column{Name: "title"}, Value: keyword},
				clause.Like{Column: clause.Column{Name: "describe"}, Value: keyword},
			))
		}
	}
	err = query.Count(&total).Limit(queryPage.PageSize).Offset(offset).Order(r.GetOrderBy(queryPage.OrderByField, queryPage.OrderByType)).Find(&list).Error
	return
}
