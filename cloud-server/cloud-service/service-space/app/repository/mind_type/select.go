package mind_type

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectByUserId 根据用户ID查询自定义导图分类
func (r *Repository) SelectByUserId(userId int64, queryPage *page.Query) (list []*model.MindType, total int64, err error) {
	if userId <= 0 || queryPage == nil {
		return []*model.MindType{}, 0, nil
	}
	offset := queryPage.PageSize * (queryPage.Page - 1)
	query := r.Database.Db.Model(&model.MindType{}).Where("user_id = ?", userId)
	err = query.Count(&total).Limit(queryPage.PageSize).Offset(offset).Order(r.GetOrderBy(queryPage.OrderByField, queryPage.OrderByType)).Find(&list).Error
	return
}
