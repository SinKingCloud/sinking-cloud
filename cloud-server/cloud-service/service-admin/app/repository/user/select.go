package user

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询用户数据
func (r *Repository) Select(where *SelectUser, queryPage *page.Query) (*page.Result[*User], error) {
	query := r.Database.Db.Model(&model.User{})
	if where != nil {
		if where.Status != nil {
			query = query.Where("status = ?", *where.Status)
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
		if where.Account != nil {
			query = query.Where("account LIKE ?", "%"+*where.Account+"%")
		}
		if where.Phone != nil {
			query = query.Where("phone LIKE ?", "%"+*where.Phone+"%")
		}
		if where.Email != nil {
			query = query.Where("email LIKE ?", "%"+*where.Email+"%")
		}
		if where.NickName != nil {
			query = query.Where("nick_name LIKE ?", "%"+*where.NickName+"%")
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
