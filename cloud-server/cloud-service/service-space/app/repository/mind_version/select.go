package mind_version

import (
	"service-space/app/model"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询导图版本列表，不查询 value 字段
func (r *Repository) Select(where *SelectMindVersion, queryPage *page.Query) (list []*MindVersion, total int64, err error) {
	if queryPage == nil {
		return []*MindVersion{}, 0, nil
	}
	offset := queryPage.PageSize * (queryPage.Page - 1)
	query := r.Database.Db.Model(&model.MindVersion{})
	if where != nil {
		if where.DataId > 0 {
			query = query.Where("data_id = ?", where.DataId)
		}
		if where.Type != "" {
			versionType, convertErr := strconv.Atoi(where.Type)
			if convertErr != nil {
				return []*MindVersion{}, 0, convertErr
			}
			query = query.Where(map[string]interface{}{"type": versionType})
		}
	}
	err = query.Count(&total).Limit(queryPage.PageSize).Offset(offset).Order(r.GetOrderBy(queryPage.OrderByField, queryPage.OrderByType)).Find(&list).Error
	return
}
