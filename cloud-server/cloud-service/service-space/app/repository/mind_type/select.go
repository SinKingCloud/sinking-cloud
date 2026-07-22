package mind_type

import (
	"service-space/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectByUserId 根据用户ID查询自定义导图分类
func (r *Repository) SelectByUserId(userId int64, queryPage *page.Query) (*page.Result[*model.MindType], error) {
	if userId <= 0 || queryPage == nil {
		return page.New(queryPage, 0, make([]*model.MindType, 0)), nil
	}
	query := r.Database.Db.Model(&model.MindType{}).Where("user_id = ?", userId)
	return r.Repository.SelectPage(query, queryPage)
}
