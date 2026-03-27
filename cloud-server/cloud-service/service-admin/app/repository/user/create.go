package user

import "service-admin/app/model"

// Create 创建用户
func (r *Repository) Create(user *model.User) error {
	return r.Database.Db.Create(user).Error
}
