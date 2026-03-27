package permission

import "service-admin/app/model"

// Create 创建权限
func (r *Repository) Create(perm *model.Permission) error {
	return r.Database.Db.Create(perm).Error
}
