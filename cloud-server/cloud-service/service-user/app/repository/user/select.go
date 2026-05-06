package user

import "service-user/app/model"

// Select 查询用户数据
func (r *Repository) Select(where *SelectUser, orderByField string, orderByType string, page int, pageSize int) (list []*User, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.User{})
	if where != nil {
		if where.Account != "" {
			query = query.Where("account LIKE ?", "%"+where.Account+"%")
		}
		if where.Phone != "" {
			query = query.Where("phone LIKE ?", "%"+where.Phone+"%")
		}
		if where.Email != "" {
			query = query.Where("email LIKE ?", "%"+where.Email+"%")
		}
		if where.NickName != "" {
			query = query.Where("nick_name LIKE ?", "%"+where.NickName+"%")
		}
		if where.RealAuthStatus != "" {
			query = query.Where("real_auth_status = ?", where.RealAuthStatus)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
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
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
