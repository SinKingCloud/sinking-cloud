package mind_version

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询导图版本列表，不查询 value 字段
func (r *Repository) Select(where *SelectMindVersion, queryPage *page.Query) (*page.Result[*MindVersion], error) {
	if queryPage == nil {
		return page.New(queryPage, 0, make([]*MindVersion, 0)), nil
	}
	query := r.Database.Db.Model(&model.MindVersion{})
	if where != nil {
		if where.DataId != nil {
			query = query.Where("data_id = ?", *where.DataId)
		}
		if where.Type != nil {
			query = query.Where("type = ?", *where.Type)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
